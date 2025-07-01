package client

import (
	"context"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

func (c *Client) ListOrganizations(ctx context.Context) ([]Organization, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, OrganizationsUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	var target []Organization
	var ratelimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&ratelimitData),
	)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, nil, fmt.Errorf("failed to list organizations: %s", res.Status)
	}

	return target, &ratelimitData, nil
}

// https://docs.sentry.io/api/guides/teams-tutorial/#list-an-organizations-teams-1
func (c *Client) ListOrganizationMembers(ctx context.Context, orgID string) ([]OrganizationMember, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(OrganizationMembersUrl, orgID), nil)
	if err != nil {
		return nil, nil, err
	}

	var target []OrganizationMember
	var ratelimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&ratelimitData),
	)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, nil, fmt.Errorf("failed to list organization members: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, nil, fmt.Errorf("failed to list organization members: %s", res.Status)
	}

	return target, &ratelimitData, nil
}
