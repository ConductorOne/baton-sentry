package client

import (
	"context"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

// docs: https://docs.sentry.io/api/organizations/

func (c *Client) ListOrganizations(ctx context.Context, cursor string) ([]Organization, *http.Response, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, OrganizationsUrl, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	if cursor != "" {
		q := req.URL.Query()
		q.Set("cursor", cursor)
		req.URL.RawQuery = q.Encode()
	}

	var target []Organization
	var ratelimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&ratelimitData),
	)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("failed to list organizations: %s", res.Status)
	}

	return target, res, &ratelimitData, nil
}

// https://docs.sentry.io/api/guides/teams-tutorial/#list-an-organizations-teams-1
func (c *Client) ListOrganizationMembers(ctx context.Context, orgID, cursor string) ([]OrganizationMember, *http.Response, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(OrganizationMembersUrl, orgID), nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var target []OrganizationMember
	var ratelimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&ratelimitData),
	)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("failed to list organization members: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, nil, nil, fmt.Errorf("failed to list organization members: %s", res.Status)
	}

	return target, res, &ratelimitData, nil
}

func (c *Client) GetOrganizationMember(ctx context.Context, orgID, memberID string) (*DetailedMember, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(OrganizationOneMemberUrl, orgID, memberID), nil)
	if err != nil {
		return nil, nil, err
	}

	var target DetailedMember
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
	)

	if err != nil {
		logBody(ctx, res.Body)
		return nil, res, fmt.Errorf("failed to get detailed organization member: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, res, fmt.Errorf("failed to get detailed organization member: %s", res.Status)
	}

	return &target, res, nil
}
