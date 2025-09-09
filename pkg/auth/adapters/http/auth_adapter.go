package http

import (
	"fmt"

	authModels "github.com/dash-ops/dash-ops/pkg/auth/models"
	authWire "github.com/dash-ops/dash-ops/pkg/auth/wire"
)

// AuthAdapter handles transformation between models and wire formats
type AuthAdapter struct{}

// NewAuthAdapter creates a new auth adapter
func NewAuthAdapter() *AuthAdapter {
	return &AuthAdapter{}
}

// ModelToUserResponse converts User model to UserResponse
func (aa *AuthAdapter) ModelToUserResponse(user *authModels.User) authWire.UserResponse {
	var organizations []authWire.OrganizationResponse
	for _, org := range user.Organizations {
		organizations = append(organizations, authWire.OrganizationResponse{
			ID:   org.ID,
			Name: org.Name,
			Slug: org.Slug,
			Role: org.Role,
		})
	}

	var teams []authWire.TeamResponse
	for _, team := range user.Teams {
		var id string
		if team.ID != nil {
			id = fmt.Sprintf("%d", *team.ID)
		}
		var name string
		if team.Name != nil {
			name = *team.Name
		}
		var slug string
		if team.Slug != nil {
			slug = *team.Slug
		}

		teams = append(teams, authWire.TeamResponse{
			ID:           id,
			Name:         name,
			Slug:         slug,
			Organization: team.Organization,
			Role:         team.Role,
		})
	}

	return authWire.UserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Name:          user.Name,
		Email:         user.Email,
		Avatar:        user.Avatar,
		Provider:      string(user.Provider),
		Organizations: organizations,
		Teams:         teams,
		CreatedAt:     user.CreatedAt,
		LastLogin:     user.LastLogin,
	}
}

// ModelToTokenResponse converts Token model to TokenResponse
func (aa *AuthAdapter) ModelToTokenResponse(token *authModels.Token) authWire.TokenResponse {
	expiresIn := int(token.TimeToExpiry().Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}

	return authWire.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresIn:    expiresIn,
		ExpiresAt:    token.ExpiresAt,
		Scopes:       token.Scopes,
	}
}

// ModelToSessionResponse converts AuthSession model to SessionResponse
func (aa *AuthAdapter) ModelToSessionResponse(session *authModels.AuthSession) authWire.SessionResponse {
	return authWire.SessionResponse{
		SessionID: session.SessionID,
		UserID:    session.UserID,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
		LastUsed:  session.LastUsed,
		IPAddress: session.IPAddress,
		UserAgent: session.UserAgent,
	}
}

// ConfigsToProvidersResponse converts auth configs to providers response
func (aa *AuthAdapter) ConfigsToProvidersResponse(configs []*authModels.AuthConfig) authWire.ProvidersResponse {
	var providers []authWire.ProviderInfo

	for _, config := range configs {
		if config.Enabled {
			providers = append(providers, authWire.ProviderInfo{
				Name:        string(config.Provider),
				DisplayName: aa.getProviderDisplayName(config.Provider),
				Method:      string(config.Method),
				Enabled:     config.Enabled,
				Scopes:      config.GetScopes(),
			})
		}
	}

	return authWire.ProvidersResponse{
		Providers: providers,
	}
}

// RequestToLoginRequest converts wire request to domain request
func (aa *AuthAdapter) RequestToLoginRequest(req authWire.LoginRequest) (*LoginRequest, error) {
	return &LoginRequest{
		Provider:    authModels.AuthProvider(req.Provider),
		RedirectURL: req.RedirectURL,
		State:       req.State,
	}, nil
}

// RequestToTokenExchangeRequest converts wire request to domain request
func (aa *AuthAdapter) RequestToTokenExchangeRequest(req authWire.TokenExchangeRequest) (*TokenExchangeRequest, error) {
	return &TokenExchangeRequest{
		Code:        req.Code,
		State:       req.State,
		RedirectURL: req.RedirectURL,
	}, nil
}

// getProviderDisplayName returns human-readable provider name
func (aa *AuthAdapter) getProviderDisplayName(provider authModels.AuthProvider) string {
	switch provider {
	case authModels.ProviderGitHub:
		return "GitHub"
	case authModels.ProviderGoogle:
		return "Google"
	case authModels.ProviderOkta:
		return "Okta"
	case authModels.ProviderSAML:
		return "SAML"
	case authModels.ProviderLDAP:
		return "LDAP"
	case authModels.ProviderJWT:
		return "JWT"
	default:
		return string(provider)
	}
}

// Domain request types (internal to adapter)
type LoginRequest struct {
	Provider    authModels.AuthProvider
	RedirectURL string
	State       string
}

type TokenExchangeRequest struct {
	Code        string
	State       string
	RedirectURL string
}

// UserPermissionsToResponse converts user permissions to response format
func (aa *AuthAdapter) UserPermissionsToResponse(permissions *authModels.UserPermissions) map[string]interface{} {
	// Convert teams to interface{} format for JSON compatibility
	var teams []map[string]interface{}
	for _, team := range permissions.Teams {
		teams = append(teams, map[string]interface{}{
			"id":   team.ID,
			"name": team.Name,
			"slug": team.Slug,
		})
	}

	return map[string]interface{}{
		"organization": permissions.Organization,
		"teams":        teams,
		"groups":       permissions.Groups,
	}
}
