package client

import (
	"context"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

// docs: https://docs.sentry.io/api/teams/

func (c *Client) ListTeams(ctx context.Context, orgID, cursor string) ([]Team, *http.Response, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(OrganizationTeamsUrl, orgID), nil)
	if err != nil {
		return nil, nil, nil, err
	}

	if cursor != "" {
		q := req.URL.Query()
		q.Set("cursor", cursor)
		req.URL.RawQuery = q.Encode()
	}

	var target []Team
	var ratelimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&ratelimitData),
	)

	if err != nil {
		if res != nil {
			logBody(ctx, res.Body)
		}
		return nil, nil, nil, fmt.Errorf("failed to list teams: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("failed to list teams: %s", res.Status)
	}

	return target, res, &ratelimitData, nil
}

func (c *Client) ListTeamMembers(ctx context.Context, orgID, teamID, cursor string) ([]TeamMember, *http.Response, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(TeamMembersUrl, orgID, teamID), nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var target []TeamMember
	var ratelimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&ratelimitData),
	)

	if err != nil {
		if res != nil {
			logBody(ctx, res.Body)
		}
		return nil, nil, nil, fmt.Errorf("failed to list teams members: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("failed to list teams members: %s", res.Status)
	}

	return target, res, &ratelimitData, nil
}

func (c *Client) AddOrgMemberToTeam(ctx context.Context, orgID, memberID, teamID string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(ProvisionTeamMemberUrl, orgID, memberID, teamID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		if res != nil {
			logBody(ctx, res.Body)
		}
		return nil, fmt.Errorf("failed to add organization member to team: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, fmt.Errorf("failed to add organization member to team: %s", res.Status)
	}

	return res, nil
}

func (c *Client) DeleteOrgMemberFromTeam(ctx context.Context, orgID, memberID, teamID string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf(ProvisionTeamMemberUrl, orgID, memberID, teamID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		if res != nil {
			logBody(ctx, res.Body)
		}
		return nil, fmt.Errorf("failed to delete organization member from team: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, fmt.Errorf("failed to delete organization member from team: %s", res.Status)
	}

	return res, nil
}
