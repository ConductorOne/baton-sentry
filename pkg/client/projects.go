package client

import (
	"context"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

func (c *Client) ListProjects(ctx context.Context, orgID, cursor string) ([]Project, string, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(OrganizationProjectsUrl, orgID), nil)
	if err != nil {
		return nil, "", nil, err
	}

	if cursor != "" {
		q := req.URL.Query()
		q.Set("cursor", cursor)
		req.URL.RawQuery = q.Encode()
	}

	var target []Project
	res, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, "", rl, fmt.Errorf("failed to list projects: %w", err)
	}

	var nextCursor string
	if HasNextPage(res) {
		nextCursor = NextCursor(res)
	}
	return target, nextCursor, rl, nil
}

func (c *Client) ListTeamProjects(ctx context.Context, orgID, teamID, cursor string) ([]Project, string, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(TeamProjectsUrl, orgID, teamID), nil)
	if err != nil {
		return nil, "", nil, err
	}

	if cursor != "" {
		q := req.URL.Query()
		q.Set("cursor", cursor)
		req.URL.RawQuery = q.Encode()
	}

	var target []Project
	res, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, "", rl, fmt.Errorf("failed to list team projects: %w", err)
	}

	var nextCursor string
	if HasNextPage(res) {
		nextCursor = NextCursor(res)
	}
	return target, nextCursor, rl, nil
}

// https://docs.sentry.io/api/projects/list-a-projects-organization-members/
// Returns a list of active organization members that belong to any team assigned to the project.
func (c *Client) ListProjectMembers(ctx context.Context, orgID, projectID, cursor string) ([]ProjectMember, string, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(ProjectMembersUrl, orgID, projectID), nil)
	if err != nil {
		return nil, "", nil, err
	}

	if cursor != "" {
		q := req.URL.Query()
		q.Set("cursor", cursor)
		req.URL.RawQuery = q.Encode()
	}

	var target []ProjectMember
	res, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, "", rl, fmt.Errorf("failed to list project members: %w", err)
	}

	var nextCursor string
	if HasNextPage(res) {
		nextCursor = NextCursor(res)
	}
	return target, nextCursor, rl, nil
}

func (c *Client) AddTeamToProject(ctx context.Context, orgID, projectID, teamID string) (*v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(ProvisionProjectTeamUrl, orgID, projectID, teamID), nil)
	if err != nil {
		return nil, err
	}

	_, rl, err := c.doRequest(req, nil) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return rl, fmt.Errorf("failed to add team to project: %w", err)
	}

	return rl, nil
}

func (c *Client) DeleteTeamFromProject(ctx context.Context, orgID, projectID, teamID string) (*v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf(ProvisionProjectTeamUrl, orgID, projectID, teamID), nil)
	if err != nil {
		return nil, err
	}

	_, rl, err := c.doRequest(req, nil) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return rl, fmt.Errorf("failed to delete team from project: %w", err)
	}

	return rl, nil
}

func (c *Client) GetProject(ctx context.Context, orgID, projectID string) (*DetailedProject, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(ProjectsUrl, orgID, projectID), nil)
	if err != nil {
		return nil, nil, err
	}

	var target DetailedProject
	_, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, rl, fmt.Errorf("failed to get project: %w", err)
	}

	return &target, rl, nil
}
