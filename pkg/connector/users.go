package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-sentry/pkg/client"
)

type userBuilder struct {
	client *client.Client
}

func (o *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return userResourceType
}

func newUserResource(member client.OrganizationMember, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"expired":       member.Expired,
		"invite_status": member.InviteStatus,
		"org_id":        parentResourceID.Resource,
	}

	return resourceSdk.NewUserResource(
		member.Name,
		userResourceType,
		member.ID,
		[]resourceSdk.UserTraitOption{
			resourceSdk.WithEmail(member.Email, true),
			resourceSdk.WithUserProfile(profile),
			resourceSdk.WithCreatedAt(member.DateCreated),
		},
		resourceSdk.WithParentResourceID(parentResourceID),
	)
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}
	var cursor string

	if pToken != nil {
		cursor = pToken.Token
	}

	members, res, ratelimitDescription, err := o.client.ListOrganizationMembers(ctx, parentResourceID.Resource, cursor)
	if err != nil {
		return nil, "", nil, err
	}
	var annotations annotations.Annotations
	annotations = *annotations.WithRateLimiting(ratelimitDescription)

	ret := make([]*v2.Resource, 0, len(members))
	for _, member := range members {
		resource, err := newUserResource(member, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		ret = append(ret, resource)
	}

	nextCursor := ""
	if client.HasNextPage(res) {
		nextCursor = client.NextCursor(res)
	}

	return ret, nextCursor, annotations, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *userBuilder) CreateAccountCapabilityDetails(_ context.Context) (*v2.CredentialDetailsAccountProvisioning, annotations.Annotations, error) {
	return &v2.CredentialDetailsAccountProvisioning{
		SupportedCredentialOptions: []v2.CapabilityDetailCredentialOption{
			v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
		},
		PreferredCredentialOption: v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_NO_PASSWORD,
	}, nil, nil
}

func (o *userBuilder) CreateAccount(ctx context.Context, accountInfo *v2.AccountInfo, credentialOptions *v2.CredentialOptions) (
	connectorbuilder.CreateAccountResponse,
	[]*v2.PlaintextData,
	annotations.Annotations,
	error,
) {
	pMap := accountInfo.Profile.AsMap()
	email, ok := pMap["email"].(string)
	if !ok {
		return nil, nil, nil, fmt.Errorf("baton-sentry: email not found in profile")
	}

	orgId, ok := pMap["orgID"].(string)
	if !ok {
		return nil, nil, nil, fmt.Errorf("baton-sentry: orgID not found in profile")
	}

	orgRole, _ := pMap["orgRole"].(string)
	err := o.client.AddMemberToOrganization(ctx, orgId, client.AddOrganizationMemberBody{
		Email:   email,
		OrgRole: orgRole,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("baton-sentry: failed to create account: %w", err)
	}

	return &v2.CreateAccountResponse_ActionRequiredResult{}, nil, nil, nil
}

func (o *userBuilder) Delete(ctx context.Context, resourceId *v2.ResourceId) (annotations.Annotations, error) {
	userID := resourceId.Resource
	orgID, err := client.FindUserOrgID(ctx, o.client, userID)
	if err != nil {
		return nil, fmt.Errorf("baton-sentry: failed to find organization for user %s: %w", resourceId.Resource, err)
	}

	err = o.client.DeleteMemberFromOrganization(ctx, orgID, userID)
	if err != nil {
		return nil, fmt.Errorf("baton-sentry: failed to delete user %s from organization %s: %w", userID, orgID, err)
	}

	return nil, nil
}

func newUserBuilder(client *client.Client) *userBuilder {
	return &userBuilder{
		client: client,
	}
}
