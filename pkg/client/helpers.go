package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/peterhellberg/link"
)

// https://docs.sentry.io/api/pagination/
func HasNextPage(res *http.Response) bool {
	for _, l := range link.ParseResponse(res) {
		if l.Rel != "next" {
			continue
		}
		if v, ok := l.Extra["results"]; ok && v == "true" {
			return true
		}
	}
	return false
}

// https://docs.sentry.io/api/pagination/
func NextCursor(res *http.Response) string {
	for _, l := range link.ParseResponse(res) {
		if l.Rel == "next" {
			if v, ok := l.Extra["cursor"]; ok {
				return v
			}
		}
	}
	return ""
}

func FindUserOrgID(ctx context.Context, client *Client, userID string) (string, error) {
	allOrgs := []Organization{}
	cursor := ""
	for {
		organizations, res, _, err := client.ListOrganizations(ctx, cursor)
		if err != nil {
			return "", fmt.Errorf("failed to list organizations: %w", err)
		}
		allOrgs = append(allOrgs, organizations...)

		if !HasNextPage(res) {
			break
		}
		cursor = NextCursor(res)
	}

	userOrgID := ""
	for _, org := range allOrgs {
		_, _, err := client.GetOrganizationMember(ctx, org.ID, userID)
		if err != nil {
			continue
		}
		userOrgID = org.ID
		break
	}

	if userOrgID == "" {
		return "", fmt.Errorf("user with ID %s not found in any organization", userID)
	}
	return userOrgID, nil
}
