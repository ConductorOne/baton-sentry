package connector

import (
	"context"
	"fmt"
	"io"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-sentry/pkg/client"
)

type Connector struct {
	client *client.Client
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newOrganizationBuilder(d.client),
		newUserBuilder(d.client),
		newTeamBuilder(d.client),
		newProjectBuilder(d.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Connector) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "My Baton Connector",
		Description: "The template implementation of a baton connector",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	return nil, nil
}

func getOrgId(resource *v2.Resource) (string, error) {
	groupTrait, err := resourceSdk.GetUserTrait(resource)
	if err != nil {
		return "", fmt.Errorf("baton-sentry: error getting traits: %w", err)
	}
	traits := groupTrait.GetProfile().AsMap()
	orgId, ok := traits["org_id"].(string)
	if !ok {
		return "", fmt.Errorf("baton-sentry: org_id not found in resource profile")
	}
	return orgId, nil
}

func getOrgIdForTeam(resource *v2.Resource) (string, error) {
	groupTrait, err := resourceSdk.GetGroupTrait(resource)
	if err != nil {
		return "", fmt.Errorf("baton-sentry: error getting traits: %w", err)
	}
	traits := groupTrait.GetProfile().AsMap()
	orgId, ok := traits["org_id"].(string)
	if !ok {
		return "", fmt.Errorf("baton-sentry: org_id not found in resource profile")
	}
	return orgId, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, apiToken string) (*Connector, error) {
	client, err := client.New(ctx, apiToken)
	if err != nil {
		return nil, err
	}
	return &Connector{
		client: client,
	}, nil
}
