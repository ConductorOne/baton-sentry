package connector

import (
	"context"
	"fmt"
	"slices"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-sentry/pkg/client"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

// orgRole pairs a Sentry organization role slug with its description.
type orgRole struct {
	Slug        string
	Description string
}

// orgRoles lists the built-in Sentry organization roles.
// https://docs.sentry.io/organization/membership/#organization-level-roles
var orgRoles = []orgRole{
	{"billing", "Can manage subscription and billing details"},
	{"member", "Can view and act on events, as well as view most other data within the organization"},
	{"admin", "Can edit global integrations, manage projects, and add/remove teams (deprecated on Business/Enterprise plans)"},
	{"manager", "Has full management access to all teams and projects"},
	{"owner", "Has unrestricted access to the organization, its data, and its settings"},
}

// defaultOrgRole is the role assigned when revoking a non-default role.
const defaultOrgRole = "member"

type organizationBuilder struct {
	client *client.Client
}

func (o *organizationBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return organizationResourceType
}

func newOrgResource(org client.Organization) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"status": org.Status.Name,
	}
	return resourceSdk.NewGroupResource(
		org.Name, organizationResourceType,
		org.ID,
		[]resourceSdk.GroupTraitOption{
			resourceSdk.WithGroupProfile(profile),
		},
		resourceSdk.WithAnnotation(
			&v2.ChildResourceType{ResourceTypeId: userResourceType.Id},
			&v2.ChildResourceType{ResourceTypeId: teamResourceType.Id},
			&v2.ChildResourceType{ResourceTypeId: projectResourceType.Id},
		),
	)
}

func (o *organizationBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	ann := annotations.New()
	cursor := ""
	if pToken != nil {
		cursor = pToken.Token
	}

	orgs, nextCursor, rl, err := o.client.ListOrganizations(ctx, cursor)
	if rl != nil {
		ann.WithRateLimiting(rl)
	}
	if err != nil {
		return nil, "", ann, fmt.Errorf("baton-sentry: failed to list organizations: %w", err)
	}

	ret := make([]*v2.Resource, 0, len(orgs))
	for _, org := range orgs {
		resource, err := newOrgResource(org)
		if err != nil {
			return nil, "", ann, fmt.Errorf("baton-sentry: failed to create resource for organization %s: %w", org.ID, err)
		}
		ret = append(ret, resource)
	}

	return ret, nextCursor, ann, nil
}

// Entitlements returns one permission entitlement per Sentry organization role.
// Sentry org roles: billing, member, admin, manager, owner.
// https://docs.sentry.io/organization/membership/#organization-level-roles
func (o *organizationBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	entitlements := make([]*v2.Entitlement, 0, len(orgRoles))
	for _, role := range orgRoles {
		entitlements = append(entitlements, entitlement.NewPermissionEntitlement(
			resource,
			role.Slug,
			entitlement.WithDescription(fmt.Sprintf("%s organization %s role: %s", resource.DisplayName, role.Slug, role.Description)),
			entitlement.WithDisplayName(fmt.Sprintf("%s %s", resource.DisplayName, role.Slug)),
			entitlement.WithGrantableTo(userResourceType),
		))
	}

	return entitlements, "", nil, nil
}

// Grants returns a grant for each organization member with their actual org role.
// The orgRole field from the Sentry API response determines which role entitlement is granted.
func (o *organizationBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	ann := annotations.New()

	var cursor string
	if pToken != nil {
		cursor = pToken.Token
	}

	members, nextCursor, rl, err := o.client.ListOrganizationMembers(ctx, resource.Id.Resource, cursor)
	if rl != nil {
		ann.WithRateLimiting(rl)
	}
	if err != nil {
		return nil, "", ann, err
	}

	ret := make([]*v2.Grant, 0, len(members))
	for _, member := range members {
		resourceId, err := resourceSdk.NewResourceID(userResourceType, member.ID)
		if err != nil {
			return nil, "", ann, fmt.Errorf("baton-sentry: failed to create resource ID for user %s: %w", member.ID, err)
		}

		role := member.OrgRole
		if role == "" {
			role = defaultOrgRole
		}

		if !slices.ContainsFunc(orgRoles, func(r orgRole) bool { return r.Slug == role }) {
			l.Debug("unknown organization role, skipping grant",
				zap.String("role", role),
				zap.String("member_id", member.ID),
				zap.String("org_id", resource.Id.Resource),
			)
			continue
		}

		ret = append(ret, grant.NewGrant(resource, role, resourceId))
	}

	return ret, nextCursor, ann, nil
}

// Grant changes a member's organization role to the requested role.
// Uses PUT /api/0/organizations/{org}/members/{member_id}/ with orgRole body.
// NOTE: Changing org roles requires a user auth token (not org-level API tokens).
func (o *organizationBuilder) Grant(ctx context.Context, principal *v2.Resource, ent *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	ann := annotations.New()

	if principal.Id.ResourceType != userResourceType.Id {
		return nil, fmt.Errorf("baton-sentry: expected principal to be a user, got %s", principal.Id.ResourceType)
	}

	orgID := ent.Resource.Id.Resource
	memberID := principal.Id.Resource
	desiredRole := ent.Slug

	// Check if the member already has this role.
	member, rl, err := o.client.GetOrganizationMember(ctx, orgID, memberID)
	if rl != nil {
		ann.WithRateLimiting(rl)
	}
	if err != nil {
		return ann, fmt.Errorf("baton-sentry: failed to get organization member %s: %w", memberID, err)
	}

	if member.OrgRole == desiredRole {
		return annotations.New(&v2.GrantAlreadyExists{}), nil
	}

	l.Debug("granting organization role",
		zap.String("member_id", memberID),
		zap.String("org_id", orgID),
		zap.String("current_role", member.OrgRole),
		zap.String("desired_role", desiredRole),
	)

	_, rl, err = o.client.UpdateOrganizationMemberRole(ctx, orgID, memberID, desiredRole)
	if rl != nil {
		ann.WithRateLimiting(rl)
	}
	if err != nil {
		return ann, fmt.Errorf("baton-sentry: failed to update organization member role: %w", err)
	}

	return ann, nil
}

// Revoke changes a member's organization role back to the default "member" role
// when revoking a non-default role. If the revoked role is "member", this is a no-op
// since the member cannot be downgraded further without removing them from the organization.
func (o *organizationBuilder) Revoke(ctx context.Context, gnt *v2.Grant) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	ann := annotations.New()

	memberID := gnt.Principal.Id.Resource
	orgID := gnt.Entitlement.Resource.Id.Resource
	revokedRole := gnt.Entitlement.Slug

	// Check if the member still has this role.
	member, rl, err := o.client.GetOrganizationMember(ctx, orgID, memberID)
	if rl != nil {
		ann.WithRateLimiting(rl)
	}
	if err != nil {
		return ann, fmt.Errorf("baton-sentry: failed to get organization member %s: %w", memberID, err)
	}

	if member.OrgRole != revokedRole {
		return annotations.New(&v2.GrantAlreadyRevoked{}), nil
	}

	// If revoking the default role, there's nothing to downgrade to.
	if revokedRole == defaultOrgRole {
		l.Debug("cannot revoke default member role without removing from organization",
			zap.String("member_id", memberID),
			zap.String("org_id", orgID),
		)
		return nil, fmt.Errorf("baton-sentry: cannot revoke the default %q role; remove the user from the organization instead", defaultOrgRole)
	}

	l.Debug("revoking organization role, downgrading to default",
		zap.String("member_id", memberID),
		zap.String("org_id", orgID),
		zap.String("revoked_role", revokedRole),
		zap.String("new_role", defaultOrgRole),
	)

	_, rl, err = o.client.UpdateOrganizationMemberRole(ctx, orgID, memberID, defaultOrgRole)
	if rl != nil {
		ann.WithRateLimiting(rl)
	}
	if err != nil {
		return ann, fmt.Errorf("baton-sentry: failed to downgrade organization member role: %w", err)
	}

	return ann, nil
}

func newOrganizationBuilder(client *client.Client) *organizationBuilder {
	return &organizationBuilder{
		client: client,
	}
}
