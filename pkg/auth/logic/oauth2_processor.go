package auth

import (
	"context"
	"fmt"
	"time"

	authModels "github.com/dash-ops/dash-ops/pkg/auth/models"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// OAuth2Processor handles OAuth2 authentication logic
type OAuth2Processor struct{}

// NewOAuth2Processor creates a new OAuth2 processor
func NewOAuth2Processor() *OAuth2Processor {
	return &OAuth2Processor{}
}

// GenerateAuthURL generates OAuth2 authorization URL
func (op *OAuth2Processor) GenerateAuthURL(config *authModels.AuthConfig, redirectURL string) (string, error) {
	if config.Method != authModels.MethodOAuth2 {
		return "", fmt.Errorf("GenerateAuthURL only supports OAuth2 method, got %s", config.Method)
	}
	oauthConfig := op.createOAuth2Config(config)
	return oauthConfig.AuthCodeURL(redirectURL), nil
}

// ExchangeCodeForToken exchanges authorization code for access token
func (op *OAuth2Processor) ExchangeCodeForToken(ctx context.Context, config *authModels.AuthConfig, code string) (*oauth2.Token, error) {
	oauthConfig := op.createOAuth2Config(config)
	return oauthConfig.Exchange(ctx, code)
}

// ProcessUserTeams processes GitHub teams into user permissions
func (op *OAuth2Processor) ProcessUserTeams(teams []*github.Team, orgPermission string) (*authModels.UserPermissions, error) {
	var userTeams []authModels.Team
	var groups []string

	for _, team := range teams {
		if team.Organization != nil && team.Organization.Login != nil && *team.Organization.Login == orgPermission {
			userTeam := authModels.Team{
				ID:   team.ID,
				Name: team.Name,
				Slug: team.Slug,
			}
			userTeams = append(userTeams, userTeam)

			if team.Slug != nil {
				groups = append(groups, fmt.Sprintf("%s*%s", orgPermission, *team.Slug))
			}
		}
	}

	return &authModels.UserPermissions{
		Organization: orgPermission,
		Teams:        userTeams,
		Groups:       groups,
	}, nil
}

// BuildUserData builds user data from GitHub teams
func (op *OAuth2Processor) BuildUserData(teams []*github.Team, orgPermission string) (*authModels.UserData, error) {
	userData := &authModels.UserData{
		Org:    orgPermission,
		Groups: []string{},
	}

	for _, team := range teams {
		if team.Organization != nil && team.Organization.Login != nil && *team.Organization.Login == orgPermission {
			if team.Slug != nil {
				userData.Groups = append(userData.Groups, fmt.Sprintf("%s%s%s", userData.Org, "*", *team.Slug))
			}
		}
	}

	return userData, nil
}

// createOAuth2Config creates OAuth2 config from auth config
func (op *OAuth2Processor) createOAuth2Config(config *authModels.AuthConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       config.Scopes,
		RedirectURL:  config.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
	}
}

// generateState generates a secure state parameter for OAuth2
func (op *OAuth2Processor) generateState() string {
	// In a real implementation, this should generate a cryptographically secure random string
	return fmt.Sprintf("state_%d", time.Now().Unix())
}

// Legacy methods for backward compatibility

// GenerateAuthURLLegacy generates OAuth2 authorization URL (legacy method)
func (op *OAuth2Processor) GenerateAuthURLLegacy(config *authModels.AuthConfig, state string) (string, error) {
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

// ExchangeCodeForTokenLegacy exchanges authorization code for access token (legacy method)
func (op *OAuth2Processor) ExchangeCodeForTokenLegacy(ctx context.Context, config *authModels.AuthConfig, code string) (*authModels.Token, error) {
	if !config.IsOAuth2() {
		return nil, fmt.Errorf("config is not for OAuth2 authentication")
	}

	oauthConfig := op.buildOAuth2Config(config)

	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return &authModels.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}, nil
}

// buildOAuth2Config builds OAuth2 config from auth config (legacy method)
func (op *OAuth2Processor) buildOAuth2Config(config *authModels.AuthConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       config.Scopes,
		RedirectURL:  config.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
	}
}
