package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	authModels "github.com/dash-ops/dash-ops/pkg/auth-new/models"
)

// SessionManager handles authentication session logic
type SessionManager struct {
	sessionDuration time.Duration
}

// NewSessionManager creates a new session manager
func NewSessionManager(sessionDuration time.Duration) *SessionManager {
	if sessionDuration == 0 {
		sessionDuration = 24 * time.Hour // Default 24 hours
	}

	return &SessionManager{
		sessionDuration: sessionDuration,
	}
}

// CreateSession creates a new authentication session
func (sm *SessionManager) CreateSession(user *authModels.User, token *authModels.Token, ipAddress, userAgent string) (*authModels.AuthSession, error) {
	if user == nil {
		return nil, fmt.Errorf("user cannot be nil")
	}

	if token == nil {
		return nil, fmt.Errorf("token cannot be nil")
	}

	sessionID, err := sm.generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	now := time.Now()
	session := &authModels.AuthSession{
		SessionID: sessionID,
		UserID:    user.ID,
		Token:     token,
		CreatedAt: now,
		ExpiresAt: now.Add(sm.sessionDuration),
		LastUsed:  now,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	return session, nil
}

// ValidateSession validates an authentication session
func (sm *SessionManager) ValidateSession(session *authModels.AuthSession) error {
	if session == nil {
		return fmt.Errorf("session is nil")
	}

	if session.SessionID == "" {
		return fmt.Errorf("session ID is empty")
	}

	if session.UserID == "" {
		return fmt.Errorf("user ID is empty")
	}

	if !session.IsActive() {
		return fmt.Errorf("session is not active or expired")
	}

	return nil
}

// RefreshSession refreshes a session's expiry time
func (sm *SessionManager) RefreshSession(session *authModels.AuthSession) error {
	if err := sm.ValidateSession(session); err != nil {
		return fmt.Errorf("cannot refresh invalid session: %w", err)
	}

	now := time.Now()
	session.ExpiresAt = now.Add(sm.sessionDuration)
	session.LastUsed = now

	return nil
}

// IsSessionExpired checks if a session is expired
func (sm *SessionManager) IsSessionExpired(session *authModels.AuthSession) bool {
	if session == nil {
		return true
	}
	return time.Now().After(session.ExpiresAt)
}

// GetSessionTimeRemaining returns time remaining for a session
func (sm *SessionManager) GetSessionTimeRemaining(session *authModels.AuthSession) time.Duration {
	if session == nil {
		return 0
	}
	return time.Until(session.ExpiresAt)
}

// generateSessionID generates a cryptographically secure session ID
func (sm *SessionManager) generateSessionID() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// ExtendSessionDuration extends the session duration
func (sm *SessionManager) ExtendSessionDuration(session *authModels.AuthSession, extension time.Duration) error {
	if session == nil {
		return fmt.Errorf("session is nil")
	}

	if extension <= 0 {
		return fmt.Errorf("extension must be positive")
	}

	session.ExpiresAt = session.ExpiresAt.Add(extension)
	return nil
}
