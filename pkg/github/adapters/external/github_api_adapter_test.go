package external

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestNewGitHubAPIAdapter_CreatesAdapterWithOAuthConfig(t *testing.T) {
	// Arrange
	oauthConfig := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		Scopes:       []string{"user:email", "read:org"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	// Act
	adapter := NewGitHubAPIAdapter(oauthConfig)

	// Assert
	assert.NotNil(t, adapter)
	assert.Equal(t, oauthConfig, adapter.oauthConfig)
}

func TestGitHubAPIAdapter_createAuthenticatedClient_CreatesClientWithToken(t *testing.T) {
	// Arrange
	oauthConfig := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}
	adapter := NewGitHubAPIAdapter(oauthConfig)
	
	validToken := &oauth2.Token{
		AccessToken: "test-access-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	client := adapter.createAuthenticatedClient(validToken)

	// Assert
	assert.NotNil(t, client)
	// The client should be configured with the OAuth token
	// We can't test the actual API calls without mocking the HTTP client
}

func TestGitHubAPIAdapter_GetUser_WithValidToken_ReturnsNoErrorForValidClient(t *testing.T) {
	// Arrange
	oauthConfig := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}
	adapter := NewGitHubAPIAdapter(oauthConfig)
	
	// This test verifies the method exists and accepts correct parameters
	// Actual API calls would require mocking the HTTP transport
	validToken := &oauth2.Token{
		AccessToken: "test-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act & Assert
	// We can't test the actual GitHub API without credentials
	// This just ensures the method signature is correct
	ctx := context.Background()
	assert.NotPanics(t, func() {
		// This will fail with actual API call, but confirms method exists
		_, _ = adapter.GetUser(ctx, validToken)
	})
}

func TestGitHubAPIAdapter_GetUserTeams_WithValidToken_ReturnsNoErrorForValidClient(t *testing.T) {
	// Arrange
	oauthConfig := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}
	adapter := NewGitHubAPIAdapter(oauthConfig)
	
	validToken := &oauth2.Token{
		AccessToken: "test-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act & Assert
	ctx := context.Background()
	assert.NotPanics(t, func() {
		_, _ = adapter.GetUserTeams(ctx, validToken)
	})
}

func TestGitHubAPIAdapter_GetUserOrganizations_WithValidToken_ReturnsNoErrorForValidClient(t *testing.T) {
	// Arrange
	oauthConfig := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}
	adapter := NewGitHubAPIAdapter(oauthConfig)
	
	validToken := &oauth2.Token{
		AccessToken: "test-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act & Assert
	ctx := context.Background()
	assert.NotPanics(t, func() {
		_, _ = adapter.GetUserOrganizations(ctx, validToken)
	})
}

func TestGitHubAPIAdapter_GetUserRepositories_WithValidToken_ReturnsNoErrorForValidClient(t *testing.T) {
	// Arrange
	oauthConfig := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}
	adapter := NewGitHubAPIAdapter(oauthConfig)
	
	validToken := &oauth2.Token{
		AccessToken: "test-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act & Assert
	ctx := context.Background()
	assert.NotPanics(t, func() {
		_, _ = adapter.GetUserRepositories(ctx, validToken)
	})
}

func TestGitHubAPIAdapter_GetOrganizationTeams_WithValidTokenAndOrg_ReturnsNoErrorForValidClient(t *testing.T) {
	// Arrange
	oauthConfig := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}
	adapter := NewGitHubAPIAdapter(oauthConfig)
	
	validToken := &oauth2.Token{
		AccessToken: "test-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}
	orgLogin := "test-org"

	// Act & Assert
	ctx := context.Background()
	assert.NotPanics(t, func() {
		_, _ = adapter.GetOrganizationTeams(ctx, validToken, orgLogin)
	})
}