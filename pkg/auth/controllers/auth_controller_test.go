package controllers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	authLogic "github.com/dash-ops/dash-ops/pkg/auth/logic"
	authModels "github.com/dash-ops/dash-ops/pkg/auth/models"
)

// MockGitHubService is a mock implementation of GitHubService for testing
type MockGitHubService struct {
	GetUserFunc      func(ctx context.Context, token *oauth2.Token) (*github.User, error)
	GetUserTeamsFunc func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error)
}

func (m *MockGitHubService) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, token)
	}
	return nil, nil
}

func (m *MockGitHubService) GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
	if m.GetUserTeamsFunc != nil {
		return m.GetUserTeamsFunc(ctx, token)
	}
	return nil, nil
}

func TestNewAuthController_CreatesControllerWithDependencies(t *testing.T) {
	// Arrange
	config := &authModels.AuthConfig{
		ClientID:        "test-client-id",
		ClientSecret:    "test-client-secret",
		URLLoginSuccess: "http://localhost:3000/success",
		OrgPermission:   "test-org",
	}
	
	oauth2Processor := authLogic.NewOAuth2Processor()
	sessionManager := authLogic.NewSessionManager(24 * time.Hour)
	mockGitHubService := &MockGitHubService{}

	// Act
	controller := NewAuthController(config, oauth2Processor, sessionManager, mockGitHubService)

	// Assert
	assert.NotNil(t, controller)
	assert.Equal(t, config, controller.config)
	assert.Equal(t, oauth2Processor, controller.oauth2Processor)
	assert.Equal(t, sessionManager, controller.sessionManager)
	assert.Equal(t, mockGitHubService, controller.githubService)
}

func TestAuthController_GenerateAuthURL_WithValidConfig_ReturnsAuthURL(t *testing.T) {
	// Arrange
	config := &authModels.AuthConfig{
		ClientID:        "test-client-id",
		ClientSecret:    "test-client-secret",
		URLLoginSuccess: "http://localhost:3000/success",
		Method:          authModels.MethodOAuth2,
		AuthURL:         "https://github.com/login/oauth/authorize",
		TokenURL:        "https://github.com/login/oauth/access_token",
		RedirectURL:     "http://localhost:8080/callback",
	}
	
	oauth2Processor := authLogic.NewOAuth2Processor()
	sessionManager := authLogic.NewSessionManager(24 * time.Hour)
	controller := NewAuthController(config, oauth2Processor, sessionManager, &MockGitHubService{})
	
	redirectURL := "http://localhost:8080/callback"

	// Act
	authURL, err := controller.GenerateAuthURL(context.Background(), redirectURL)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, authURL)
	assert.Contains(t, authURL, "https://github.com/login/oauth/authorize")
	assert.Contains(t, authURL, "client_id=test-client-id")
}

func TestAuthController_BuildRedirectURL_WithTokenAndState_ReturnsCorrectURL(t *testing.T) {
	// Arrange
	config := &authModels.AuthConfig{
		URLLoginSuccess: "http://localhost:3000/success",
	}
	
	controller := NewAuthController(config, nil, nil, nil)
	
	token := &oauth2.Token{
		AccessToken: "test-access-token",
	}
	
	state := "/dashboard"

	// Act
	redirectURL := controller.BuildRedirectURL(token, state)

	// Assert
	assert.Equal(t, "http://localhost:3000/success/dashboard?access_token=test-access-token", redirectURL)
}

func TestAuthController_BuildRedirectURL_WithTokenNoState_ReturnsBaseURL(t *testing.T) {
	// Arrange
	config := &authModels.AuthConfig{
		URLLoginSuccess: "http://localhost:3000/success",
	}
	
	controller := NewAuthController(config, nil, nil, nil)
	
	token := &oauth2.Token{
		AccessToken: "test-access-token",
	}

	// Act
	redirectURL := controller.BuildRedirectURL(token, "")

	// Assert
	assert.Equal(t, "http://localhost:3000/success?access_token=test-access-token", redirectURL)
}

func TestAuthController_GetUserProfile_WithValidToken_ReturnsUserProfile(t *testing.T) {
	// Arrange
	expectedUser := &github.User{
		Login: github.String("johndoe"),
		ID:    github.Int64(123),
		Name:  github.String("John Doe"),
		Email: github.String("john@example.com"),
	}
	
	mockGitHubService := &MockGitHubService{
		GetUserFunc: func(ctx context.Context, token *oauth2.Token) (*github.User, error) {
			return expectedUser, nil
		},
	}
	
	config := &authModels.AuthConfig{}
	controller := NewAuthController(config, nil, nil, mockGitHubService)
	
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	profile, err := controller.GetUserProfile(context.Background(), validToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	
	user, ok := profile.(*github.User)
	assert.True(t, ok)
	assert.Equal(t, expectedUser, user)
	assert.Equal(t, "johndoe", user.GetLogin())
	assert.Equal(t, int64(123), user.GetID())
}

func TestAuthController_GetUserProfile_WithNilToken_ReturnsError(t *testing.T) {
	// Arrange
	controller := NewAuthController(&authModels.AuthConfig{}, nil, nil, &MockGitHubService{})

	// Act
	profile, err := controller.GetUserProfile(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "token is required")
}

func TestAuthController_GetUserProfile_WithInvalidToken_ReturnsError(t *testing.T) {
	// Arrange
	controller := NewAuthController(&authModels.AuthConfig{}, nil, nil, &MockGitHubService{})
	
	invalidToken := &oauth2.Token{
		AccessToken: "expired-token",
		Expiry:      time.Now().Add(-1 * time.Hour), // Expired
	}

	// Act
	profile, err := controller.GetUserProfile(context.Background(), invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "token is invalid")
}

func TestAuthController_GetUserProfile_WhenServiceFails_ReturnsError(t *testing.T) {
	// Arrange
	mockGitHubService := &MockGitHubService{
		GetUserFunc: func(ctx context.Context, token *oauth2.Token) (*github.User, error) {
			return nil, errors.New("GitHub API error")
		},
	}
	
	controller := NewAuthController(&authModels.AuthConfig{}, nil, nil, mockGitHubService)
	
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	profile, err := controller.GetUserProfile(context.Background(), validToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "failed to get user profile")
}

func TestAuthController_GetUserPermissions_WithValidToken_ReturnsPermissions(t *testing.T) {
	// Arrange
	mockTeams := []*github.Team{
		{
			ID:   github.Int64(1),
			Name: github.String("Engineering"),
			Slug: github.String("engineering"),
			Organization: &github.Organization{
				Login: github.String("test-org"),
			},
		},
		{
			ID:   github.Int64(2),
			Name: github.String("DevOps"),
			Slug: github.String("devops"),
			Organization: &github.Organization{
				Login: github.String("test-org"),
			},
		},
	}
	
	mockGitHubService := &MockGitHubService{
		GetUserTeamsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
			return mockTeams, nil
		},
	}
	
	config := &authModels.AuthConfig{
		OrgPermission: "test-org",
	}
	
	oauth2Processor := authLogic.NewOAuth2Processor()
	controller := NewAuthController(config, oauth2Processor, nil, mockGitHubService)
	
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	permissions, err := controller.GetUserPermissions(context.Background(), validToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, permissions)
	assert.Equal(t, "test-org", permissions.Organization)
	assert.Len(t, permissions.Teams, 2)
	assert.Len(t, permissions.Groups, 2)
}

func TestAuthController_GetUserPermissions_WithNilToken_ReturnsError(t *testing.T) {
	// Arrange
	controller := NewAuthController(&authModels.AuthConfig{}, nil, nil, &MockGitHubService{})

	// Act
	permissions, err := controller.GetUserPermissions(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, permissions)
	assert.Contains(t, err.Error(), "token is required")
}

func TestAuthController_GetUserPermissions_WhenServiceFails_ReturnsError(t *testing.T) {
	// Arrange
	mockGitHubService := &MockGitHubService{
		GetUserTeamsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
			return nil, errors.New("teams API error")
		},
	}
	
	config := &authModels.AuthConfig{
		OrgPermission: "test-org",
	}
	
	oauth2Processor := authLogic.NewOAuth2Processor()
	controller := NewAuthController(config, oauth2Processor, nil, mockGitHubService)
	
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	permissions, err := controller.GetUserPermissions(context.Background(), validToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, permissions)
	assert.Contains(t, err.Error(), "failed to fetch user teams")
}

func TestAuthController_ValidateToken_WithValidToken_ReturnsNoError(t *testing.T) {
	// Arrange
	sessionManager := authLogic.NewSessionManager(24 * time.Hour)
	controller := NewAuthController(&authModels.AuthConfig{}, nil, sessionManager, nil)
	
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	err := controller.ValidateToken(context.Background(), validToken)

	// Assert
	assert.NoError(t, err)
}

func TestAuthController_ValidateToken_WithExpiredToken_ReturnsError(t *testing.T) {
	// Arrange
	sessionManager := authLogic.NewSessionManager(24 * time.Hour)
	controller := NewAuthController(&authModels.AuthConfig{}, nil, sessionManager, nil)
	
	expiredToken := &oauth2.Token{
		AccessToken: "expired-token",
		Expiry:      time.Now().Add(-1 * time.Hour),
	}

	// Act
	err := controller.ValidateToken(context.Background(), expiredToken)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is invalid or expired")
}

func TestAuthController_BuildUserData_WithValidToken_ReturnsUserData(t *testing.T) {
	// Arrange
	mockTeams := []*github.Team{
		{
			ID:   github.Int64(1),
			Name: github.String("Engineering"),
			Slug: github.String("engineering"),
			Organization: &github.Organization{
				Login: github.String("test-org"),
			},
		},
	}
	
	mockGitHubService := &MockGitHubService{
		GetUserTeamsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
			return mockTeams, nil
		},
	}
	
	config := &authModels.AuthConfig{
		OrgPermission: "test-org",
	}
	
	oauth2Processor := authLogic.NewOAuth2Processor()
	controller := NewAuthController(config, oauth2Processor, nil, mockGitHubService)
	
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	userData, err := controller.BuildUserData(context.Background(), validToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	assert.Equal(t, "test-org", userData.Org)
	assert.Len(t, userData.Groups, 1)
	assert.Contains(t, userData.Groups[0], "test-org")
}

func TestAuthController_BuildUserData_WithNilToken_ReturnsError(t *testing.T) {
	// Arrange
	controller := NewAuthController(&authModels.AuthConfig{}, nil, nil, &MockGitHubService{})

	// Act
	userData, err := controller.BuildUserData(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, userData)
	assert.Contains(t, err.Error(), "token is required")
}

func TestAuthController_BuildUserData_WithNoOrgAccess_ReturnsUserDataWithNoPermission(t *testing.T) {
	// Arrange
	mockTeams := []*github.Team{
		{
			ID:   github.Int64(1),
			Name: github.String("External Team"),
			Slug: github.String("external-team"),
			Organization: &github.Organization{
				Login: github.String("other-org"),
			},
		},
	}
	
	mockGitHubService := &MockGitHubService{
		GetUserTeamsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
			return mockTeams, nil
		},
	}
	
	config := &authModels.AuthConfig{
		OrgPermission: "test-org", // Different from team's org
	}
	
	oauth2Processor := authLogic.NewOAuth2Processor()
	controller := NewAuthController(config, oauth2Processor, nil, mockGitHubService)
	
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	userData, err := controller.BuildUserData(context.Background(), validToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	assert.Equal(t, "test-org", userData.Org)
	assert.Len(t, userData.Groups, 0) // No groups from test-org
}