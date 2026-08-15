package client

import (
	"context"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

// docs: https://docs.sentry.io/api/teams/

func (c *Client) GetTeam(ctx context.Context, orgID, teamID string) (*Team, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(TeamUrl, orgID, teamID), nil)
	if err != nil {
		return nil, nil, err
	}

	var target Team
	_, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, rl, fmt.Errorf("failed to get team: %w", err)
	}

	return &target, rl, nil
}

func (c *Client) ListTeams(ctx context.Context, orgID, cursor string) ([]Team, string, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.OrganizationTeamsURL(orgID), nil)
	if err != nil {
		return nil, "", nil, err
	}

	if cursor != "" {
		q := req.URL.Query()
		q.Set("cursor", cursor)
		req.URL.RawQuery = q.Encode()
	}

	var target []Team
	res, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, "", rl, fmt.Errorf("failed to list teams: %w", err)
	}

	var nextCursor string
	if HasNextPage(res) {
		nextCursor = NextCursor(res)
	}
	return target, nextCursor, rl, nil
}

func (c *Client) ListTeamMembers(ctx context.Context, orgID, teamID, cursor string) ([]TeamMember, string, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.TeamMembersURL(orgID, teamID), nil)
	if err != nil {
		return nil, "", nil, err
	}

	if cursor != "" {
		q := req.URL.Query()
		q.Set("cursor", cursor)
		req.URL.RawQuery = q.Encode()
	}

	var target []TeamMember
	res, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, "", rl, fmt.Errorf("failed to list team members: %w", err)
	}

	var nextCursor string
	if HasNextPage(res) {
		nextCursor = NextCursor(res)
	}
	return target, nextCursor, rl, nil
}

func (c *Client) AddOrgMemberToTeam(ctx context.Context, orgID, memberID, teamID string) (*v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.ProvisionTeamMemberURL(orgID, memberID, teamID), nil)
	if err != nil {
		return nil, err
	}

	_, rl, err := c.doRequest(req, nil) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return rl, fmt.Errorf("failed to add organization member to team: %w", err)
	}

	return rl, nil
}

func (c *Client) DeleteOrgMemberFromTeam(ctx context.Context, orgID, memberID, teamID string) (*v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.ProvisionTeamMemberURL(orgID, memberID, teamID), nil)
	if err != nil {
		return nil, err
	}

	_, rl, err := c.doRequest(req, nil) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return rl, fmt.Errorf("failed to delete organization member from team: %w", err)
	}

	return rl, nil
}
