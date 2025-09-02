package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"

	authModels "github.com/dash-ops/dash-ops/pkg/auth-new/models"
)

// OAuth2Processor handles OAuth2 authentication logic
type OAuth2Processor struct{}

// NewOAuth2Processor creates a new OAuth2 processor
func NewOAuth2Processor() *OAuth2Processor {
	return &OAuth2Processor{}
}

// GenerateAuthURL generates OAuth2 authorization URL
func (op *OAuth2Processor) GenerateAuthURL(config *authModels.AuthConfig, state string) (string, error) {
	if !config.IsOAuth2() {
		return "", fmt.Errorf("config is not for OAuth2 authentication")
	}

	oauthConfig := op.buildOAuth2Config(config)

	// Add state parameter for security
	if state == "" {
		state = op.generateState()
	}

	url := oauthConfig.AuthCodeURL(state)
	return url, nil
}

// ExchangeCodeForToken exchanges authorization code for access token
func (op *OAuth2Processor) ExchangeCodeForToken(ctx context.Context, config *authModels.AuthConfig, code string) (*authModels.Token, error) {
	if !config.IsOAuth2() {
		return nil, fmt.Errorf("config is not for OAuth2 authentication")
	}

	oauthConfig := op.buildOAuth2Config(config)

	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	if !token.Valid() {
		return nil, fmt.Errorf("received invalid token")
	}

	return op.convertToAuthToken(token, config.Provider), nil
}

// RefreshToken refreshes an OAuth2 token
func (op *OAuth2Processor) RefreshToken(ctx context.Context, config *authModels.AuthConfig, refreshToken string) (*authModels.Token, error) {
	if !config.IsOAuth2() {
		return nil, fmt.Errorf("config is not for OAuth2 authentication")
	}

	oauthConfig := op.buildOAuth2Config(config)

	// Create token source with refresh token
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := oauthConfig.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return op.convertToAuthToken(newToken, config.Provider), nil
}

// ValidateToken validates an OAuth2 token
func (op *OAuth2Processor) ValidateToken(token *authModels.Token) error {
	if token == nil {
		return fmt.Errorf("token is nil")
	}

	if token.AccessToken == "" {
		return fmt.Errorf("access token is empty")
	}

	if token.IsExpired() {
		return fmt.Errorf("token is expired")
	}

	return nil
}

// buildOAuth2Config builds oauth2.Config from AuthConfig
func (op *OAuth2Processor) buildOAuth2Config(config *authModels.AuthConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes:       config.GetScopes(),
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
	}
}

// convertToAuthToken converts oauth2.Token to authModels.Token
func (op *OAuth2Processor) convertToAuthToken(token *oauth2.Token, provider authModels.AuthProvider) *authModels.Token {
	return &authModels.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresAt:    token.Expiry,
		Provider:     provider,
	}
}

// generateState generates a random state parameter for OAuth2
func (op *OAuth2Processor) generateState() string {
	// In a real implementation, this should generate a cryptographically secure random string
	return fmt.Sprintf("state_%d", time.Now().Unix())
}
