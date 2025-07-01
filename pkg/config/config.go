package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	// Add the SchemaFields for the Config.
	ApiToken = field.StringField(
		"api-token",
		field.WithDisplayName("API Token"),
		field.WithDescription("API Token for Sentry"),
	)

	ConfigurationFields = []field.SchemaField{ApiToken}

	// FieldRelationships defines relationships between the ConfigurationFields that can be automatically validated.
	// For example, a username and password can be required together, or an access token can be
	// marked as mutually exclusive from the username password pair.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run -tags=generate ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("Sentry"),
)
