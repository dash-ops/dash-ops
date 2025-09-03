package auth

import (
	"context"

	authModels "github.com/dash-ops/dash-ops/pkg/auth-new/models"
)

// ProviderService defines the interface for authentication provider services
type ProviderService interface {
	// GetUser retrieves user information from the provider
	GetUser(ctx context.Context, token *authModels.Token) (*authModels.User, error)

	// GetUserOrganizations retrieves user's organizations from the provider
	GetUserOrganizations(ctx context.Context, token *authModels.Token) ([]authModels.Organization, error)

	// GetUserTeams retrieves user's teams from the provider
	GetUserTeams(ctx context.Context, token *authModels.Token) ([]authModels.Team, error)

	// ValidateToken validates a token with the provider
	ValidateToken(ctx context.Context, token *authModels.Token) error

	// RevokeToken revokes a token with the provider
	RevokeToken(ctx context.Context, token *authModels.Token) error
}

// GitHubService defines the interface for GitHub-specific operations
type GitHubService interface {
	ProviderService

	// CheckOrgMembership checks if user is member of an organization
	CheckOrgMembership(ctx context.Context, token *authModels.Token, org string) (bool, error)

	// CheckTeamMembership checks if user is member of a team
	CheckTeamMembership(ctx context.Context, token *authModels.Token, org, team string) (bool, error)

	// GetUserEmails retrieves user's email addresses
	GetUserEmails(ctx context.Context, token *authModels.Token) ([]string, error)
}

// GoogleService defines the interface for Google-specific operations
type GoogleService interface {
	ProviderService

	// GetUserInfo retrieves user info from Google
	GetUserInfo(ctx context.Context, token *authModels.Token) (*GoogleUserInfo, error)
}

// GoogleUserInfo represents Google user information
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// NotificationService defines the interface for sending notifications
type NotificationService interface {
	// SendLoginNotification sends login notification
	SendLoginNotification(ctx context.Context, user *authModels.User, session *authModels.AuthSession) error

	// SendLogoutNotification sends logout notification
	SendLogoutNotification(ctx context.Context, user *authModels.User) error

	// SendSecurityAlert sends security alert
	SendSecurityAlert(ctx context.Context, user *authModels.User, message string) error
}

// AuditService defines the interface for audit logging
type AuditService interface {
	// LogLogin logs a login event
	LogLogin(ctx context.Context, user *authModels.User, session *authModels.AuthSession) error

	// LogLogout logs a logout event
	LogLogout(ctx context.Context, userID, sessionID string) error

	// LogTokenRefresh logs a token refresh event
	LogTokenRefresh(ctx context.Context, userID string) error

	// LogSecurityEvent logs a security event
	LogSecurityEvent(ctx context.Context, userID, event, details string) error
}
