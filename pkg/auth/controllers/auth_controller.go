package controllers

import (
	"context"
	"fmt"

	authLogic "github.com/dash-ops/dash-ops/pkg/auth/logic"
	authModels "github.com/dash-ops/dash-ops/pkg/auth/models"
	authPorts "github.com/dash-ops/dash-ops/pkg/auth/ports"
	"golang.org/x/oauth2"
)

// AuthController orchestrates authentication business logic
type AuthController struct {
	config          *authModels.AuthConfig
	oauth2Processor *authLogic.OAuth2Processor
	sessionManager  *authLogic.SessionManager
	githubService   authPorts.GitHubService
}

// NewAuthController creates a new auth controller
func NewAuthController(
	config *authModels.AuthConfig,
	oauth2Processor *authLogic.OAuth2Processor,
	sessionManager *authLogic.SessionManager,
	githubService authPorts.GitHubService,
) *AuthController {
	return &AuthController{
		config:          config,
		oauth2Processor: oauth2Processor,
		sessionManager:  sessionManager,
		githubService:   githubService,
	}
}

// GenerateAuthURL generates OAuth2 authorization URL
func (ac *AuthController) GenerateAuthURL(ctx context.Context, redirectURL string) (string, error) {
	return ac.oauth2Processor.GenerateAuthURL(ac.config, redirectURL)
}

// ExchangeCodeForToken exchanges authorization code for access token
func (ac *AuthController) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return ac.oauth2Processor.ExchangeCodeForToken(ctx, ac.config, code)
}

// BuildRedirectURL builds the final redirect URL with token
func (ac *AuthController) BuildRedirectURL(token *oauth2.Token, state string) string {
	baseURL := ac.config.URLLoginSuccess
	if state != "" {
		baseURL += state
	}
	return baseURL + "?access_token=" + token.AccessToken
}

// GetUserProfile gets user profile from provider
func (ac *AuthController) GetUserProfile(ctx context.Context, token *oauth2.Token) (interface{}, error) {
	if token == nil {
		return nil, fmt.Errorf("token is required")
	}

	if !token.Valid() {
		return nil, fmt.Errorf("token is invalid")
	}

	user, err := ac.githubService.GetUser(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return user, nil
}

// GetUserPermissions gets user permissions from provider
func (ac *AuthController) GetUserPermissions(ctx context.Context, token *oauth2.Token) (*authModels.UserPermissions, error) {
	if token == nil {
		return nil, fmt.Errorf("token is required")
	}

	if !token.Valid() {
		return nil, fmt.Errorf("token is invalid")
	}

	orgPermission := ac.config.OrgPermission
	teams, err := ac.githubService.GetUserTeams(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user teams: %w", err)
	}

	// Process teams using business logic
	return ac.oauth2Processor.ProcessUserTeams(teams, orgPermission)
}

// ValidateToken validates the provided token
func (ac *AuthController) ValidateToken(ctx context.Context, token *oauth2.Token) error {
	return ac.sessionManager.ValidateToken(token)
}

// BuildUserData builds user data for context
func (ac *AuthController) BuildUserData(ctx context.Context, token *oauth2.Token) (*authModels.UserData, error) {
	if token == nil {
		return nil, fmt.Errorf("token is required")
	}

	orgPermission := ac.config.OrgPermission
	teams, err := ac.githubService.GetUserTeams(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate organization permissions: %w", err)
	}

	// Process teams to build user data
	return ac.oauth2Processor.BuildUserData(teams, orgPermission)
}
