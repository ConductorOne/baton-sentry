package connector

import (
	"context"
	"fmt"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-sentry/pkg/client"
)

const teamMembership = "member"

type teamBuilder struct {
	client *client.Client
}

func (o *teamBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return teamResourceType
}

func newTeamResource(team client.Team, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"org_id": parentResourceID.Resource,
	}
	return resourceSdk.NewGroupResource(
		team.Name,
		teamResourceType,
		// <orgID>/<teamID>
		fmt.Sprintf("%s/%s", parentResourceID.Resource, team.ID),
		[]resourceSdk.GroupTraitOption{
			resourceSdk.WithGroupProfile(profile),
		},
		resourceSdk.WithParentResourceID(parentResourceID),
	)
}

func (o *teamBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	ann := annotations.New()
	var cursor string
	if pToken != nil {
		cursor = pToken.Token
	}

	orgID := parentResourceID.Resource
	teams, nextCursor, rl, err := o.client.ListTeams(ctx, orgID, cursor)
	if rl != nil {
		ann.WithRateLimiting(rl)
	}
	if err != nil {
		return nil, "", ann, err
	}

	ret := make([]*v2.Resource, 0, len(teams))
	for _, team := range teams {
		resource, err := newTeamResource(team, parentResourceID)
		if err != nil {
			return nil, "", ann, err
		}
		ret = append(ret, resource)
	}

	return ret, nextCursor, ann, nil
}

func (o *teamBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			resource,
			teamMembership,
			entitlement.WithDescription(fmt.Sprintf("Member of %s team", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Member of %s team", resource.DisplayName)),
			entitlement.WithGrantableTo(userResourceType),
		),
	}, "", nil, nil
}

func (o *teamBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	ann := annotations.New()
	var cursor string
	if pToken != nil {
		cursor = pToken.Token
	}

	orgID := resource.ParentResourceId.Resource
	teamID := strings.Split(resource.Id.Resource, "/")[1]
	members, nextCursor, rl, err := o.client.ListTeamMembers(ctx, orgID, teamID, cursor)
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
		ret = append(ret, grant.NewGrant(resource, teamMembership, resourceId))
	}

	return ret, nextCursor, ann, nil
}

func (o *teamBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	ann := annotations.New()
	split := strings.Split(entitlement.Resource.Id.Resource, "/")

	orgId := split[0]
	teamId := split[1]
	memberId := principal.Id.Resource
	teamName := entitlement.Resource.DisplayName

	member, rl, err := o.client.GetOrganizationMember(ctx, orgId, memberId)
	if rl != nil {
		ann.WithRateLimiting(rl)
	}
	if err != nil {
		return ann, fmt.Errorf("baton-sentry: failed to get organization member: %w", err)
	}

	for _, name := range member.Teams {
		if name == teamName {
			return annotations.New(&v2.GrantAlreadyExists{}), nil
		}
	}

	rl, err = o.client.AddOrgMemberToTeam(ctx, orgId, memberId, teamId)
	if rl != nil {
		ann.WithRateLimiting(rl)
	}
	if err != nil {
		return ann, fmt.Errorf("baton-sentry: failed to add organization member to team: %w", err)
	}

	return ann, nil
}

func (o *teamBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	ann := annotations.New()
	entitlement := grant.Entitlement
	split := strings.Split(entitlement.Resource.Id.Resource, "/")

	orgId := split[0]
	teamId := split[1]
	memberId := grant.Principal.Id.Resource
	teamName := entitlement.Resource.DisplayName

	member, rl, err := o.client.GetOrganizationMember(ctx, orgId, memberId)
	if rl != nil {
		ann.WithRateLimiting(rl)
	}
	if err != nil {
		return ann, fmt.Errorf("baton-sentry: failed to get organization member: %w", err)
	}

	exists := false
	for _, name := range member.Teams {
		if name == teamName {
			exists = true
			break
		}
	}

	if !exists {
		return annotations.New(&v2.GrantAlreadyRevoked{}), nil
	}

	rl, err = o.client.DeleteOrgMemberFromTeam(ctx, orgId, memberId, teamId)
	if rl != nil {
		ann.WithRateLimiting(rl)
	}
	if err != nil {
		return ann, fmt.Errorf("baton-sentry: failed to delete organization member from team: %w", err)
	}

	return ann, nil
}

func newTeamBuilder(client *client.Client) *teamBuilder {
	return &teamBuilder{
		client: client,
	}
}
