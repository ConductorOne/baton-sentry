package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

// ErrMemberAlreadyExists is returned when the Sentry API indicates
// the member is already part of the organization.
var ErrMemberAlreadyExists = errors.New("member already exists")

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
		if res != nil {
			logBody(ctx, res.Body)
		}
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

	if cursor != "" {
		q := req.URL.Query()
		q.Set("cursor", cursor)
		req.URL.RawQuery = q.Encode()
	}

	var target []OrganizationMember
	var ratelimitData v2.RateLimitDescription
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
		uhttp.WithRatelimitData(&ratelimitData),
	)

	if err != nil {
		if res != nil {
			logBody(ctx, res.Body)
		}
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
		if res != nil {
			logBody(ctx, res.Body)
		}
		return nil, res, fmt.Errorf("failed to get detailed organization member: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logBody(ctx, res.Body)
		return nil, res, fmt.Errorf("failed to get detailed organization member: %s", res.Status)
	}

	return &target, res, nil
}

func (c *Client) AddMemberToOrganization(ctx context.Context, orgID string, member AddOrganizationMemberBody) error {
	v, err := json.Marshal(member)
	if err != nil {
		return fmt.Errorf("failed to marshal member: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(OrganizationMembersUrl, orgID), bytes.NewReader(v))
	if err != nil {
		return fmt.Errorf("failed to create request to add member to organization: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)

	if err != nil {
		if res != nil {
			body := readBody(ctx, res.Body)
			if isMemberAlreadyExistsError(body) {
				return ErrMemberAlreadyExists
			}
			if body != "" {
				return fmt.Errorf("failed to add member to organization: %w: %s", err, body)
			}
		}
		return fmt.Errorf("failed to add member to organization: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body := readBody(ctx, res.Body)
		if isMemberAlreadyExistsError(body) {
			return ErrMemberAlreadyExists
		}
		if body != "" {
			return fmt.Errorf("failed to add member to organization: %s: %s", res.Status, body)
		}
		return fmt.Errorf("failed to add member to organization: %s", res.Status)
	}

	return nil
}

// isMemberAlreadyExistsError checks if the Sentry API response body
// indicates the member is already part of the organization.
func isMemberAlreadyExistsError(body string) bool {
	lower := strings.ToLower(body)
	return strings.Contains(lower, "already a member") || strings.Contains(lower, "member already exists")
}

// UpdateOrganizationMemberRole updates a member's organization-level role.
// https://docs.sentry.io/api/organizations/update-an-organization-members-roles/
// NOTE: Changing organization roles is restricted to user auth tokens.
func (c *Client) UpdateOrganizationMemberRole(ctx context.Context, orgID, memberID, role string) (*DetailedMember, error) {
	body := UpdateOrganizationMemberBody{OrgRole: role}
	v, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update member role body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf(OrganizationOneMemberUrl, orgID, memberID), bytes.NewReader(v))
	if err != nil {
		return nil, fmt.Errorf("failed to create request to update member role: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var target DetailedMember
	res, err := c.Do(req,
		uhttp.WithJSONResponse(&target),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update organization member role: %w", err)
	}
	defer res.Body.Close()

	return &target, nil
}

func (c *Client) DeleteMemberFromOrganization(ctx context.Context, orgID, userID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf(OrganizationOneMemberUrl, orgID, userID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request to delete member: %w", err)
	}

	res, err := c.Do(req)
	if err != nil {
		if res != nil {
			body := readBody(ctx, res.Body)
			if body != "" {
				return fmt.Errorf("failed to delete member from organization: %w: %s", err, body)
			}
		}
		return fmt.Errorf("failed to delete member from organization: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body := readBody(ctx, res.Body)
		if body != "" {
			return fmt.Errorf("failed to delete member from organization: %s: %s", res.Status, body)
		}
		return fmt.Errorf("failed to delete member from organization: %s", res.Status)
	}

	return nil
}
