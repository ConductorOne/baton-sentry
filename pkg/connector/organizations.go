package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-sentry/pkg/client"
)

const organizationMembership = "member"

type organizationBuilder struct {
	client *client.Client
}

func (o *organizationBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
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

func (o *organizationBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	cursor := ""
	if pToken != nil {
		cursor = pToken.Token
	}

	orgs, res, ratelimitDescription, err := o.client.ListOrganizations(ctx, cursor)
	if err != nil {
		return nil, "", nil, fmt.Errorf("baton-sentry: failed to list organizations: %w", err)
	}
	var annotations annotations.Annotations
	annotations = *annotations.WithRateLimiting(ratelimitDescription)

	ret := make([]*v2.Resource, 0, len(orgs))
	for _, org := range orgs {
		resource, err := newOrgResource(org)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-sentry: failed to create resource for organization %s: %w", org.ID, err)
		}
		ret = append(ret, resource)
	}

	nextCursor := ""
	if client.HasNextPage(res) {
		nextCursor = client.NextCursor(res)
	}

	return ret, nextCursor, annotations, nil
}

func (o *organizationBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			resource,
			organizationMembership,
			entitlement.WithDescription(fmt.Sprintf("Member of %s organization", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Member of %s organization", resource.DisplayName)),
		),
	}, "", nil, nil
}

func (o *organizationBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	cursor := ""
	if pToken != nil {
		cursor = pToken.Token
	}
	members, res, ratelimitDescription, err := o.client.ListOrganizationMembers(ctx, resource.Id.Resource, cursor)
	if err != nil {
		return nil, "", nil, err
	}
	var annotations annotations.Annotations
	annotations = *annotations.WithRateLimiting(ratelimitDescription)

	ret := make([]*v2.Grant, 0, len(members))
	for _, member := range members {
		resourceId, err := resourceSdk.NewResourceID(userResourceType, member.ID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("baton-sentry: failed to create resource ID for user %s: %w", member.ID, err)
		}

		ret = append(ret, grant.NewGrant(resource, organizationMembership, resourceId))
	}

	nextCursor := ""
	if client.HasNextPage(res) {
		nextCursor = client.NextCursor(res)
	}

	return ret, nextCursor, annotations, nil
}

func newOrganizationBuilder(client *client.Client) *organizationBuilder {
	return &organizationBuilder{
		client: client,
	}
}
