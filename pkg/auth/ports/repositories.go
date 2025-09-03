package auth

import (
	"context"

	authModels "github.com/dash-ops/dash-ops/pkg/auth-new/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *authModels.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (*authModels.User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*authModels.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*authModels.User, error)

	// Update updates a user
	Update(ctx context.Context, user *authModels.User) error

	// Delete deletes a user
	Delete(ctx context.Context, id string) error

	// UpdateLastLogin updates user's last login time
	UpdateLastLogin(ctx context.Context, id string) error
}

// SessionRepository defines the interface for session data access
type SessionRepository interface {
	// Create creates a new session
	Create(ctx context.Context, session *authModels.AuthSession) error

	// GetByID retrieves a session by session ID
	GetByID(ctx context.Context, sessionID string) (*authModels.AuthSession, error)

	// GetByUserID retrieves all sessions for a user
	GetByUserID(ctx context.Context, userID string) ([]*authModels.AuthSession, error)

	// Update updates a session
	Update(ctx context.Context, session *authModels.AuthSession) error

	// Delete deletes a session
	Delete(ctx context.Context, sessionID string) error

	// DeleteByUserID deletes all sessions for a user
	DeleteByUserID(ctx context.Context, userID string) error

	// DeleteExpired deletes all expired sessions
	DeleteExpired(ctx context.Context) error

	// UpdateLastUsed updates session's last used time
	UpdateLastUsed(ctx context.Context, sessionID string) error
}

// TokenRepository defines the interface for token storage (optional, for token persistence)
type TokenRepository interface {
	// Store stores a token
	Store(ctx context.Context, userID string, token *authModels.Token) error

	// Get retrieves a token for a user
	Get(ctx context.Context, userID string) (*authModels.Token, error)

	// Delete deletes a token for a user
	Delete(ctx context.Context, userID string) error

	// DeleteExpired deletes all expired tokens
	DeleteExpired(ctx context.Context) error
}
