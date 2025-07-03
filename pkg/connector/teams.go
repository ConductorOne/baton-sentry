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

const teamMembership = "member"

type teamBuilder struct {
	client *client.Client
}

func (o *teamBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return teamResourceType
}

func newTeamResource(team client.Team, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		// "status": org.Status.Name,
	}
	return resourceSdk.NewGroupResource(
		team.Name, teamResourceType,
		team.ID,
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

	cursor := ""
	if pToken != nil {
		cursor = pToken.Token
	}

	orgID := parentResourceID.Resource
	teams, res, ratelimitDescription, err := o.client.ListTeams(ctx, orgID, cursor)
	if err != nil {
		return nil, "", nil, err
	}
	var annotations annotations.Annotations
	annotations = *annotations.WithRateLimiting(ratelimitDescription)

	ret := make([]*v2.Resource, 0, len(teams))
	for _, team := range teams {
		resource, err := newTeamResource(team, parentResourceID)
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

func (o *teamBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			resource,
			teamMembership,
			entitlement.WithDescription(fmt.Sprintf("Member of %s team", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Member of %s team", resource.DisplayName)),
		),
	}, "", nil, nil
}

func (o *teamBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	cursor := ""
	if pToken != nil {
		cursor = pToken.Token
	}

	orgID := resource.ParentResourceId.Resource
	teamID := resource.Id.Resource
	members, res, ratelimitDescription, err := o.client.ListTeamMembers(ctx, orgID, teamID, cursor)
	if err != nil {
		return nil, "", nil, err
	}
	var annotations annotations.Annotations
	annotations = *annotations.WithRateLimiting(ratelimitDescription)

	ret := make([]*v2.Grant, 0, len(members))
	for _, member := range members {
		resourceId, err := resourceSdk.NewResourceID(userResourceType, member.ID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to create resource ID for user %s: %w", member.ID, err)
		}

		ret = append(ret, grant.NewGrant(resource, teamMembership, resourceId))
	}

	nextCursor := ""
	if client.HasNextPage(res) {
		nextCursor = client.NextCursor(res)
	}

	return ret, nextCursor, annotations, nil
}

func newTeamBuilder(client *client.Client) *teamBuilder {
	return &teamBuilder{
		client: client,
	}
}
