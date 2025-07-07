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

const projectMembership = "member"

type projectBuilder struct {
	client *client.Client
}

func (o *projectBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return projectResourceType
}

func newProjectResource(project client.Project, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"org_id": parentResourceID.Resource,
	}
	return resourceSdk.NewGroupResource(
		project.Name,
		projectResourceType,
		project.ID,
		[]resourceSdk.GroupTraitOption{
			resourceSdk.WithGroupProfile(profile),
		},
		resourceSdk.WithParentResourceID(parentResourceID),
	)
}

func (o *projectBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	cursor := ""
	if pToken != nil {
		cursor = pToken.Token
	}

	orgID := parentResourceID.Resource
	projects, res, ratelimitDescription, err := o.client.ListProjects(ctx, orgID, cursor)
	if err != nil {
		return nil, "", nil, err
	}
	var annotations annotations.Annotations
	annotations = *annotations.WithRateLimiting(ratelimitDescription)

	ret := make([]*v2.Resource, 0, len(projects))
	for _, project := range projects {
		resource, err := newProjectResource(project, parentResourceID)
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

func (o *projectBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	// TODO: add grantable to teams
	// add grant expansion to include user memberships  in the grants function
	return []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			resource,
			projectMembership,
			entitlement.WithDescription(fmt.Sprintf("Member of %s project", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("Member of %s project", resource.DisplayName)),
			entitlement.WithGrantableTo(teamResourceType),
		),
	}, "", nil, nil
}

func (o *projectBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	cursor := ""
	if pToken != nil {
		cursor = pToken.Token
	}

	orgID := resource.ParentResourceId.Resource
	projectID := resource.Id.Resource
	members, res, ratelimitDescription, err := o.client.ListProjectMembers(ctx, orgID, projectID, cursor)
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

		ret = append(ret, grant.NewGrant(resource, projectMembership, resourceId))
	}

	nextCursor := ""
	if client.HasNextPage(res) {
		nextCursor = client.NextCursor(res)
	}

	return ret, nextCursor, annotations, nil
}

func (o *projectBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	orgId := principal.ParentResourceId.Resource
	teamId := principal.Id.Resource
	projectId := entitlement.Resource.Id.Resource

	project, _, err := o.client.GetProject(ctx, orgId, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	for _, team := range project.Teams {
		if team.ID == teamId {
			return annotations.New(&v2.GrantAlreadyExists{}), nil
		}
	}

	_, err = o.client.AddTeamToProject(ctx, orgId, projectId, teamId)
	if err != nil {
		return nil, fmt.Errorf("failed to add team to project: %w", err)
	}

	return nil, nil
}

func (o *projectBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	orgId := grant.Principal.ParentResourceId.Resource
	teamId := grant.Principal.Id.Resource
	projectId := grant.Entitlement.Resource.Id.Resource

	project, _, err := o.client.GetProject(ctx, orgId, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	exists := false
	for _, team := range project.Teams {
		if team.ID == teamId {
			exists = true
			break
		}

	}

	if !exists {
		return annotations.New(&v2.GrantAlreadyRevoked{}), nil
	}

	_, err = o.client.DeleteTeamFromProject(ctx, orgId, projectId, teamId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete team from project: %w", err)
	}

	return nil, nil
}

func newProjectBuilder(client *client.Client) *projectBuilder {
	return &projectBuilder{
		client: client,
	}
}
