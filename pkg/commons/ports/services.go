package commons

import (
	"context"

	commonsModels "github.com/dash-ops/dash-ops/pkg/commons/models"
)

// AuthenticationService defines the interface for authentication services
type AuthenticationService interface {
	// ValidateToken validates a token and returns user data
	ValidateToken(ctx context.Context, token string) (*commonsModels.UserData, error)

	// RefreshToken refreshes an existing token
	RefreshToken(ctx context.Context, token string) (string, error)

	// RevokeToken revokes a token
	RevokeToken(ctx context.Context, token string) error
}

// AuthorizationService defines the interface for authorization services
type AuthorizationService interface {
	// CheckPermission checks if user has required permissions
	CheckPermission(ctx context.Context, userData *commonsModels.UserData, requiredPermissions []string) (bool, error)

	// GetUserPermissions returns all permissions for a user
	GetUserPermissions(ctx context.Context, userData *commonsModels.UserData) ([]string, error)
}

// LoggingService defines the interface for logging services
type LoggingService interface {
	// LogRequest logs an HTTP request
	LogRequest(ctx context.Context, method, path, userAgent, clientIP string)

	// LogError logs an error with context
	LogError(ctx context.Context, err error, message string, fields map[string]interface{})

	// LogInfo logs an info message
	LogInfo(ctx context.Context, message string, fields map[string]interface{})
}
