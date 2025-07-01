package connector

import (
	"context"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-sentry/pkg/client"
)

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
		),
	)
}

func (o *organizationBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var annotations annotations.Annotations
	orgs, ratelimitDescription, err := o.client.ListOrganizations(ctx)
	annotations = *annotations.WithRateLimiting(ratelimitDescription)
	if err != nil {
		return nil, "", nil, err
	}

	ret := make([]*v2.Resource, 0, len(orgs))
	for _, org := range orgs {
		resource, err := newOrgResource(org)
		if err != nil {
			return nil, "", nil, err
		}
		ret = append(ret, resource)

	}

	return ret, "", annotations, nil
}

func (o *organizationBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *organizationBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newOrganizationBuilder(client *client.Client) *organizationBuilder {
	return &organizationBuilder{
		client: client,
	}
}
