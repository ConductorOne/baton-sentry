package connector

import (
	"context"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-sentry/pkg/client"
)

type userBuilder struct {
	client *client.Client
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

func newUserResource(member client.OrganizationMember) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"expired":       member.Expired,
		"invite_status": member.InviteStatus,
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
	)

}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	var annotations annotations.Annotations
	members, ratelimitDescription, err := o.client.ListOrganizationMembers(ctx, parentResourceID.Resource)
	annotations = *annotations.WithRateLimiting(ratelimitDescription)
	if err != nil {
		return nil, "", nil, err
	}

	ret := make([]*v2.Resource, 0, len(members))
	for _, member := range members {
		resource, err := newUserResource(member)
		if err != nil {
			return nil, "", nil, err
		}
		ret = append(ret, resource)
	}

	return ret, "", annotations, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client *client.Client) *userBuilder {
	return &userBuilder{
		client: client,
	}
}
