package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authModels "github.com/dash-ops/dash-ops/pkg/auth/models"
)

func TestSessionManager_CreateSession_WithValidUserAndToken_ReturnsSession(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)
	user := &authModels.User{
		ID:       "user-123",
		Username: "testuser",
		Name:     "Test User",
		Email:    "test@example.com",
		Provider: authModels.ProviderGitHub,
	}
	token := &authModels.Token{
		AccessToken: "access-token",
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Hour),
		Provider:    authModels.ProviderGitHub,
	}
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	// Act
	session, err := manager.CreateSession(user, token, ipAddress, userAgent)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, session)
	assert.NotEmpty(t, session.SessionID)
	assert.Equal(t, user.ID, session.UserID)
	assert.Equal(t, token, session.Token)
	assert.Equal(t, ipAddress, session.IPAddress)
	assert.Equal(t, userAgent, session.UserAgent)
	assert.True(t, session.IsActive())
	assert.False(t, time.Now().After(session.ExpiresAt))
}

func TestSessionManager_CreateSession_WithNilUser_ReturnsError(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)
	token := &authModels.Token{
		AccessToken: "access-token",
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Hour),
		Provider:    authModels.ProviderGitHub,
	}
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	// Act
	session, err := manager.CreateSession(nil, token, ipAddress, userAgent)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, session)
}

func TestSessionManager_CreateSession_WithNilToken_ReturnsError(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)
	user := &authModels.User{
		ID:       "user-123",
		Username: "testuser",
		Name:     "Test User",
		Email:    "test@example.com",
		Provider: authModels.ProviderGitHub,
	}
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	// Act
	session, err := manager.CreateSession(user, nil, ipAddress, userAgent)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, session)
}

func TestSessionManager_ValidateSession_WithValidSession_ReturnsNoError(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)
	validToken := &authModels.Token{
		AccessToken: "access-token",
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Hour),
		Provider:    authModels.ProviderGitHub,
	}
	session := &authModels.AuthSession{
		SessionID: "session-123",
		UserID:    "user-123",
		Token:     validToken,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		LastUsed:  time.Now(),
	}

	// Act
	err := manager.ValidateSession(session)

	// Assert
	assert.NoError(t, err)
}

func TestSessionManager_ValidateSession_WithNilSession_ReturnsError(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)

	// Act
	err := manager.ValidateSession(nil)

	// Assert
	assert.Error(t, err)
}

func TestSessionManager_ValidateSession_WithEmptySessionID_ReturnsError(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)
	validToken := &authModels.Token{
		AccessToken: "access-token",
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Hour),
		Provider:    authModels.ProviderGitHub,
	}
	session := &authModels.AuthSession{
		SessionID: "",
		UserID:    "user-123",
		Token:     validToken,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		LastUsed:  time.Now(),
	}

	// Act
	err := manager.ValidateSession(session)

	// Assert
	assert.Error(t, err)
}

func TestSessionManager_ValidateSession_WithEmptyUserID_ReturnsError(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)
	validToken := &authModels.Token{
		AccessToken: "access-token",
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Hour),
		Provider:    authModels.ProviderGitHub,
	}
	session := &authModels.AuthSession{
		SessionID: "session-123",
		UserID:    "",
		Token:     validToken,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		LastUsed:  time.Now(),
	}

	// Act
	err := manager.ValidateSession(session)

	// Assert
	assert.Error(t, err)
}

func TestSessionManager_ValidateSession_WithExpiredSession_ReturnsError(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)
	expiredToken := &authModels.Token{
		AccessToken: "expired-token",
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(-time.Hour), // Expired
		Provider:    authModels.ProviderGitHub,
	}
	session := &authModels.AuthSession{
		SessionID: "session-123",
		UserID:    "user-123",
		Token:     expiredToken,
		CreatedAt: time.Now().Add(-25 * time.Hour),
		ExpiresAt: time.Now().Add(-time.Hour), // Expired
		LastUsed:  time.Now().Add(-time.Hour),
	}

	// Act
	err := manager.ValidateSession(session)

	// Assert
	assert.Error(t, err)
}

func TestSessionManager_RefreshSession_WithValidSession_ExtendsExpiry(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)
	validToken := &authModels.Token{
		AccessToken: "access-token",
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Hour),
		Provider:    authModels.ProviderGitHub,
	}
	session := &authModels.AuthSession{
		SessionID: "session-123",
		UserID:    "user-123",
		Token:     validToken,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour), // Will expire soon
		LastUsed:  time.Now(),
	}
	originalExpiresAt := session.ExpiresAt
	originalLastUsed := session.LastUsed

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	// Act
	err := manager.RefreshSession(session)

	// Assert
	require.NoError(t, err)
	assert.True(t, session.ExpiresAt.After(originalExpiresAt))
	assert.True(t, session.LastUsed.After(originalLastUsed))
}

func TestSessionManager_IsSessionExpired_WithNilSession_ReturnsTrue(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)

	// Act
	result := manager.IsSessionExpired(nil)

	// Assert
	assert.True(t, result)
}

func TestSessionManager_IsSessionExpired_WithActiveSession_ReturnsFalse(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)
	session := &authModels.AuthSession{
		ExpiresAt: time.Now().Add(time.Hour),
	}

	// Act
	result := manager.IsSessionExpired(session)

	// Assert
	assert.False(t, result)
}

func TestSessionManager_IsSessionExpired_WithExpiredSession_ReturnsTrue(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)
	session := &authModels.AuthSession{
		ExpiresAt: time.Now().Add(-time.Hour),
	}

	// Act
	result := manager.IsSessionExpired(session)

	// Assert
	assert.True(t, result)
}

func TestSessionManager_generateSessionID_GeneratesUniqueIDs(t *testing.T) {
	// Arrange
	manager := NewSessionManager(24 * time.Hour)

	// Act
	id1, err1 := manager.generateSessionID()
	id2, err2 := manager.generateSessionID()

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)  // Should generate different IDs
	assert.Equal(t, 64, len(id1)) // 32 bytes = 64 hex characters
	assert.Equal(t, 64, len(id2))
}
