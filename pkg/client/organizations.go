package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

// docs: https://docs.sentry.io/api/organizations/

func (c *Client) ListOrganizations(ctx context.Context, cursor string) ([]Organization, string, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.OrganizationsURL(), nil)
	if err != nil {
		return nil, "", nil, err
	}

	if cursor != "" {
		q := req.URL.Query()
		q.Set("cursor", cursor)
		req.URL.RawQuery = q.Encode()
	}

	var target []Organization
	res, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, "", rl, fmt.Errorf("failed to list organizations: %w", err)
	}

	var nextCursor string
	if HasNextPage(res) {
		nextCursor = NextCursor(res)
	}
	return target, nextCursor, rl, nil
}

// https://docs.sentry.io/api/guides/teams-tutorial/#list-an-organizations-teams-1
func (c *Client) ListOrganizationMembers(ctx context.Context, orgID, cursor string) ([]OrganizationMember, string, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.OrganizationMembersURL(orgID), nil)
	if err != nil {
		return nil, "", nil, err
	}

	if cursor != "" {
		q := req.URL.Query()
		q.Set("cursor", cursor)
		req.URL.RawQuery = q.Encode()
	}

	var target []OrganizationMember
	res, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, "", rl, fmt.Errorf("failed to list organization members: %w", err)
	}

	var nextCursor string
	if HasNextPage(res) {
		nextCursor = NextCursor(res)
	}
	return target, nextCursor, rl, nil
}

func (c *Client) GetOrganizationMember(ctx context.Context, orgID, memberID string) (*DetailedMember, *v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.OrganizationOneMemberURL(orgID, memberID), nil)
	if err != nil {
		return nil, nil, err
	}

	var target DetailedMember
	_, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, rl, fmt.Errorf("failed to get organization member: %w", err)
	}

	return &target, rl, nil
}

func (c *Client) AddMemberToOrganization(ctx context.Context, orgID string, member AddOrganizationMemberBody) (*OrganizationMember, *v2.RateLimitDescription, error) {
	v, err := json.Marshal(member)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal member: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.OrganizationMembersURL(orgID), bytes.NewReader(v))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var target OrganizationMember
	_, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, rl, fmt.Errorf("failed to add member to organization: %w", err)
	}

	return &target, rl, nil
}

// UpdateOrganizationMemberRole updates a member's organization-level role.
// https://docs.sentry.io/api/organizations/update-an-organization-members-roles/
// NOTE: Changing organization roles is restricted to user auth tokens.
func (c *Client) UpdateOrganizationMemberRole(ctx context.Context, orgID, memberID, role string) (*DetailedMember, *v2.RateLimitDescription, error) {
	body := UpdateOrganizationMemberBody{OrgRole: role}
	v, err := json.Marshal(body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal update member role body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.OrganizationOneMemberURL(orgID, memberID), bytes.NewReader(v))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var target DetailedMember
	_, rl, err := c.doRequest(req, &target) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return nil, rl, fmt.Errorf("failed to update organization member role: %w", err)
	}

	return &target, rl, nil
}

func (c *Client) DeleteMemberFromOrganization(ctx context.Context, orgID, userID string) (*v2.RateLimitDescription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.OrganizationOneMemberURL(orgID, userID), nil)
	if err != nil {
		return nil, err
	}

	_, rl, err := c.doRequest(req, nil) //nolint:bodyclose // body closed by doRequest
	if err != nil {
		return rl, fmt.Errorf("failed to delete member from organization: %w", err)
	}

	return rl, nil
}
