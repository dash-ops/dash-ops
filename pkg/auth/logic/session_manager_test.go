package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authModels "github.com/dash-ops/dash-ops/pkg/auth/models"
)

func TestSessionManager_CreateSession(t *testing.T) {
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

	tests := []struct {
		name        string
		user        *authModels.User
		token       *authModels.Token
		ipAddress   string
		userAgent   string
		expectError bool
	}{
		{
			name:        "valid session creation",
			user:        user,
			token:       token,
			ipAddress:   "192.168.1.1",
			userAgent:   "Mozilla/5.0",
			expectError: false,
		},
		{
			name:        "nil user",
			user:        nil,
			token:       token,
			ipAddress:   "192.168.1.1",
			userAgent:   "Mozilla/5.0",
			expectError: true,
		},
		{
			name:        "nil token",
			user:        user,
			token:       nil,
			ipAddress:   "192.168.1.1",
			userAgent:   "Mozilla/5.0",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := manager.CreateSession(tt.user, tt.token, tt.ipAddress, tt.userAgent)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, session)
			} else {
				require.NoError(t, err)
				require.NotNil(t, session)

				assert.NotEmpty(t, session.SessionID)
				assert.Equal(t, tt.user.ID, session.UserID)
				assert.Equal(t, tt.token, session.Token)
				assert.Equal(t, tt.ipAddress, session.IPAddress)
				assert.Equal(t, tt.userAgent, session.UserAgent)
				assert.True(t, session.IsActive())
				assert.False(t, time.Now().After(session.ExpiresAt))
			}
		})
	}
}

func TestSessionManager_ValidateSession(t *testing.T) {
	manager := NewSessionManager(24 * time.Hour)

	validToken := &authModels.Token{
		AccessToken: "access-token",
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(time.Hour),
		Provider:    authModels.ProviderGitHub,
	}

	expiredToken := &authModels.Token{
		AccessToken: "expired-token",
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(-time.Hour), // Expired
		Provider:    authModels.ProviderGitHub,
	}

	tests := []struct {
		name        string
		session     *authModels.AuthSession
		expectError bool
	}{
		{
			name: "valid session",
			session: &authModels.AuthSession{
				SessionID: "session-123",
				UserID:    "user-123",
				Token:     validToken,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(24 * time.Hour),
				LastUsed:  time.Now(),
			},
			expectError: false,
		},
		{
			name:        "nil session",
			session:     nil,
			expectError: true,
		},
		{
			name: "empty session ID",
			session: &authModels.AuthSession{
				SessionID: "",
				UserID:    "user-123",
				Token:     validToken,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(24 * time.Hour),
				LastUsed:  time.Now(),
			},
			expectError: true,
		},
		{
			name: "empty user ID",
			session: &authModels.AuthSession{
				SessionID: "session-123",
				UserID:    "",
				Token:     validToken,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(24 * time.Hour),
				LastUsed:  time.Now(),
			},
			expectError: true,
		},
		{
			name: "expired session",
			session: &authModels.AuthSession{
				SessionID: "session-123",
				UserID:    "user-123",
				Token:     expiredToken,
				CreatedAt: time.Now().Add(-25 * time.Hour),
				ExpiresAt: time.Now().Add(-time.Hour), // Expired
				LastUsed:  time.Now().Add(-time.Hour),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ValidateSession(tt.session)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionManager_RefreshSession(t *testing.T) {
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

	err := manager.RefreshSession(session)
	require.NoError(t, err)

	// Check that expiry was extended
	assert.True(t, session.ExpiresAt.After(originalExpiresAt))
	assert.True(t, session.LastUsed.After(originalLastUsed))
}

func TestSessionManager_IsSessionExpired(t *testing.T) {
	manager := NewSessionManager(24 * time.Hour)

	tests := []struct {
		name     string
		session  *authModels.AuthSession
		expected bool
	}{
		{
			name:     "nil session",
			session:  nil,
			expected: true,
		},
		{
			name: "active session",
			session: &authModels.AuthSession{
				ExpiresAt: time.Now().Add(time.Hour),
			},
			expected: false,
		},
		{
			name: "expired session",
			session: &authModels.AuthSession{
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.IsSessionExpired(tt.session)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSessionManager_generateSessionID(t *testing.T) {
	manager := NewSessionManager(24 * time.Hour)

	id1, err1 := manager.generateSessionID()
	id2, err2 := manager.generateSessionID()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)  // Should generate different IDs
	assert.Equal(t, 64, len(id1)) // 32 bytes = 64 hex characters
	assert.Equal(t, 64, len(id2))
}
