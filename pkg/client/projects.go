package client

import (
	"context"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func (c *Client) ListProjects(ctx context.Context, orgID, cursor string) ([]Project, *http.Response, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(OrganizationProjectsUrl, orgID), nil)
	if err != nil {
		return nil, nil, nil, err
	}

	if cursor != "" {
		q := req.URL.Query()
		q.Set("cursor", cursor)
		req.URL.RawQuery = q.Encode()
	}

	var target []Project
	var ratelimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&ratelimitData),
	)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("failed to list projects: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("failed to list projects: %s", res.Status)
	}

	return target, res, &ratelimitData, nil
}

// https://docs.sentry.io/api/projects/list-a-projects-organization-members/
// Returns a list of active organization members that belong to any team assigned to the project.
func (c *Client) ListProjectMembers(ctx context.Context, orgID, projectID, cursor string) ([]ProjectMember, *http.Response, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(ProjectMembersUrl, orgID, projectID), nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var target []ProjectMember
	var ratelimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&ratelimitData),
	)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("failed to list project members: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("failed to list project members: %s", res.Status)
	}

	return target, res, &ratelimitData, nil
}

func (c *Client) AddTeamToProject(ctx context.Context, orgID, projectID, teamID string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(ProjectMembersUrl, orgID, projectID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, fmt.Errorf("failed to add team to project: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, fmt.Errorf("failed to add team to project: %s", res.Status)
	}

	return res, nil
}

func (c *Client) DeleteTeamFromProject(ctx context.Context, orgID, projectID, teamID string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf(ProjectMembersUrl, orgID, projectID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, fmt.Errorf("failed to delete team from project: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, fmt.Errorf("failed to delete team from project: %s", res.Status)
	}

	return res, nil
}

func (c *Client) GetProject(ctx context.Context, orgID, projectID string) (*DetailedProject, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(ProjectsUrl, orgID, projectID), nil)
	if err != nil {
		return nil, nil, err
	}

	var target DetailedProject
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
	)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, nil, fmt.Errorf("failed to get project: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, nil, fmt.Errorf("failed to get project: %s", res.Status)
	}

	return &target, res, nil
}
