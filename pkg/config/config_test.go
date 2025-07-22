package config

import (
	"testing"

	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Sentry
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Sentry{
				ApiToken: "asdfasdfaasdf",
			},
			wantErr: false,
		},
		{
			name: "invalid config - missing required fields",
			config: &Sentry{
				ApiToken: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := field.Validate(Config, tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
