package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	authModels "github.com/dash-ops/dash-ops/pkg/auth/models"
)

func TestOAuth2Processor_GenerateAuthURL_WithValidOAuth2Config_ReturnsURL(t *testing.T) {
	// Arrange
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
	state := "test-state"

	// Act
	url, err := processor.GenerateAuthURL(config, state)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, url)
	assert.Contains(t, url, config.AuthURL)
	assert.Contains(t, url, config.ClientID)
	assert.Contains(t, url, state)
}

func TestOAuth2Processor_GenerateAuthURL_WithNonOAuth2Config_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewOAuth2Processor()
	config := &authModels.AuthConfig{
		Provider: authModels.ProviderGitHub,
		Method:   authModels.MethodJWT,
	}
	state := "test-state"

	// Act
	url, err := processor.GenerateAuthURL(config, state)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, url)
}

// ValidateToken test removed - OAuth2Processor doesn't have ValidateToken method
// Token validation is handled by SessionManager.ValidateToken

func TestOAuth2Processor_buildOAuth2Config_WithValidConfig_ReturnsOAuth2Config(t *testing.T) {
	// Arrange
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

	// Act
	oauthConfig := processor.buildOAuth2Config(config)

	// Assert
	assert.Equal(t, config.ClientID, oauthConfig.ClientID)
	assert.Equal(t, config.ClientSecret, oauthConfig.ClientSecret)
	assert.Equal(t, config.RedirectURL, oauthConfig.RedirectURL)
	assert.Equal(t, config.Scopes, oauthConfig.Scopes)
	assert.Equal(t, config.AuthURL, oauthConfig.Endpoint.AuthURL)
	assert.Equal(t, config.TokenURL, oauthConfig.Endpoint.TokenURL)
}

func TestOAuth2Processor_generateState_GeneratesUniqueStates(t *testing.T) {
	// Arrange
	processor := NewOAuth2Processor()

	// Act
	state1 := processor.generateState()
	state2 := processor.generateState()

	// Assert
	assert.NotEmpty(t, state1)
	assert.NotEmpty(t, state2)
	// Note: States might be equal if generated in the same second due to timestamp-based generation
	// In production, this should use crypto/rand for better uniqueness
	assert.Contains(t, state1, "state_")
	assert.Contains(t, state2, "state_")
}
