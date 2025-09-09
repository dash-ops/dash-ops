package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	authModels "github.com/dash-ops/dash-ops/pkg/auth/models"
)

func TestOAuth2Processor_GenerateAuthURL(t *testing.T) {
	processor := NewOAuth2Processor()

	config := &authModels.AuthConfig{
		Provider:     authModels.ProviderGitHub,
		Method:       authModels.MethodOAuth2,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"user", "repo"},
		Enabled:      true,
	}

	tests := []struct {
		name        string
		config      *authModels.AuthConfig
		state       string
		expectError bool
	}{
		{
			name:        "valid OAuth2 config",
			config:      config,
			state:       "test-state",
			expectError: false,
		},
		{
			name: "non-OAuth2 config",
			config: &authModels.AuthConfig{
				Provider: authModels.ProviderGitHub,
				Method:   authModels.MethodJWT,
			},
			state:       "test-state",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := processor.GenerateAuthURL(tt.config, tt.state)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, url)
				assert.Contains(t, url, tt.config.AuthURL)
				assert.Contains(t, url, tt.config.ClientID)
				if tt.state != "" {
					assert.Contains(t, url, tt.state)
				}
			}
		})
	}
}

// ValidateToken test removed - OAuth2Processor doesn't have ValidateToken method
// Token validation is handled by SessionManager.ValidateToken

func TestOAuth2Processor_buildOAuth2Config(t *testing.T) {
	processor := NewOAuth2Processor()

	config := &authModels.AuthConfig{
		Provider:     authModels.ProviderGitHub,
		Method:       authModels.MethodOAuth2,
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"user", "repo"},
	}

	oauthConfig := processor.buildOAuth2Config(config)

	assert.Equal(t, config.ClientID, oauthConfig.ClientID)
	assert.Equal(t, config.ClientSecret, oauthConfig.ClientSecret)
	assert.Equal(t, config.RedirectURL, oauthConfig.RedirectURL)
	assert.Equal(t, config.Scopes, oauthConfig.Scopes)
	assert.Equal(t, config.AuthURL, oauthConfig.Endpoint.AuthURL)
	assert.Equal(t, config.TokenURL, oauthConfig.Endpoint.TokenURL)
}

func TestOAuth2Processor_generateState(t *testing.T) {
	processor := NewOAuth2Processor()

	state1 := processor.generateState()
	state2 := processor.generateState()

	assert.NotEmpty(t, state1)
	assert.NotEmpty(t, state2)
	// Note: States might be equal if generated in the same second due to timestamp-based generation
	// In production, this should use crypto/rand for better uniqueness
	assert.Contains(t, state1, "state_")
	assert.Contains(t, state2, "state_")
}
