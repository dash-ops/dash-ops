package external

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"

	githubLogic "github.com/dash-ops/dash-ops/pkg/github/logic"
)

// MockGitHubAPIClient is a mock implementation of GitHubAPIClient
type MockGitHubAPIClient struct {
	mock.Mock
}

func (m *MockGitHubAPIClient) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*github.User), args.Error(1)
}

func (m *MockGitHubAPIClient) GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*github.Team), args.Error(1)
}

func (m *MockGitHubAPIClient) GetUserOrganizations(ctx context.Context, token *oauth2.Token) ([]*github.Organization, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*github.Organization), args.Error(1)
}

func (m *MockGitHubAPIClient) GetUserRepositories(ctx context.Context, token *oauth2.Token) ([]*github.Repository, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*github.Repository), args.Error(1)
}

func (m *MockGitHubAPIClient) GetOrganizationTeams(ctx context.Context, token *oauth2.Token, orgLogin string) ([]*github.Team, error) {
	args := m.Called(ctx, token, orgLogin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*github.Team), args.Error(1)
}

// Test GitHubClient.GetUser scenarios

func TestGitHubClient_GetUser_WithNilToken_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	// Act
	user, err := client.GetUser(ctx, nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "token is required")
}

func TestGitHubClient_GetUser_WithExpiredToken_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	expiredToken := &oauth2.Token{
		AccessToken: "expired-token",
		Expiry:      time.Now().Add(-time.Hour), // Expired
	}

	// Act
	user, err := client.GetUser(ctx, expiredToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "token is invalid")
}

func TestGitHubClient_GetUser_WithValidToken_ReturnsUser(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	expectedUser := &github.User{
		ID:    github.Int64(123),
		Login: github.String("testuser"),
		Name:  github.String("Test User"),
		Email: github.String("test@example.com"),
	}

	mockAPI.On("GetUser", ctx, validToken).Return(expectedUser, nil)

	// Act
	user, err := client.GetUser(ctx, validToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(123), user.GetID())
	assert.Equal(t, "testuser", user.GetLogin())
	assert.Equal(t, "Test User", user.GetName())
	assert.Equal(t, "test@example.com", user.GetEmail())
	mockAPI.AssertExpectations(t)
}

func TestGitHubClient_GetUser_WhenAPIFails_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	apiError := errors.New("GitHub API error: rate limit exceeded")
	mockAPI.On("GetUser", ctx, validToken).Return(nil, apiError)

	// Act
	user, err := client.GetUser(ctx, validToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "rate limit exceeded")
	mockAPI.AssertExpectations(t)
}

// Test GitHubClient.GetUserTeams scenarios

func TestGitHubClient_GetUserTeams_WithNilToken_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	// Act
	teams, err := client.GetUserTeams(ctx, nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, teams)
	assert.Contains(t, err.Error(), "token is required")
}

func TestGitHubClient_GetUserTeams_WithExpiredToken_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	expiredToken := &oauth2.Token{
		AccessToken: "expired-token",
		Expiry:      time.Now().Add(-time.Hour),
	}

	// Act
	teams, err := client.GetUserTeams(ctx, expiredToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, teams)
	assert.Contains(t, err.Error(), "token is invalid")
}

func TestGitHubClient_GetUserTeams_WithValidToken_ReturnsTeams(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	expectedTeams := []*github.Team{
		{
			ID:   github.Int64(1),
			Name: github.String("Engineering"),
			Slug: github.String("engineering"),
		},
		{
			ID:   github.Int64(2),
			Name: github.String("DevOps"),
			Slug: github.String("devops"),
		},
	}

	mockAPI.On("GetUserTeams", ctx, validToken).Return(expectedTeams, nil)

	// Act
	teams, err := client.GetUserTeams(ctx, validToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, teams)
	assert.Len(t, teams, 2)
	assert.Equal(t, "Engineering", teams[0].GetName())
	assert.Equal(t, "DevOps", teams[1].GetName())
	mockAPI.AssertExpectations(t)
}

func TestGitHubClient_GetUserTeams_WhenAPIFails_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	apiError := errors.New("GitHub API error: unauthorized")
	mockAPI.On("GetUserTeams", ctx, validToken).Return(nil, apiError)

	// Act
	teams, err := client.GetUserTeams(ctx, validToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, teams)
	assert.Contains(t, err.Error(), "unauthorized")
	mockAPI.AssertExpectations(t)
}

// Test GitHubClient.GetUserProfile scenarios

func TestGitHubClient_GetUserProfile_WithCompleteData_ReturnsFullProfile(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	expectedUser := &github.User{
		ID:        github.Int64(123),
		Login:     github.String("testuser"),
		Name:      github.String("Test User"),
		Email:     github.String("test@example.com"),
		AvatarURL: github.String("https://avatar.com/user"),
		HTMLURL:   github.String("https://github.com/testuser"),
		Type:      github.String("User"),
		SiteAdmin: github.Bool(false),
	}

	expectedTeams := []*github.Team{
		{
			ID:         github.Int64(1),
			Name:       github.String("team1"),
			Slug:       github.String("team1"),
			Permission: github.String("admin"),
			Organization: &github.Organization{
				ID:    github.Int64(1),
				Login: github.String("org1"),
				Name:  github.String("Organization 1"),
			},
		},
	}

	expectedOrgs := []*github.Organization{
		{
			ID:          github.Int64(1),
			Login:       github.String("org1"),
			Name:        github.String("Organization 1"),
			Description: github.String("Test organization"),
		},
	}

	mockAPI.On("GetUser", ctx, validToken).Return(expectedUser, nil)
	mockAPI.On("GetUserTeams", ctx, validToken).Return(expectedTeams, nil)
	mockAPI.On("GetUserOrganizations", ctx, validToken).Return(expectedOrgs, nil)

	// Act
	profile, err := client.GetUserProfile(ctx, validToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, int64(123), profile.User.ID)
	assert.Equal(t, "testuser", profile.User.Login)
	assert.Equal(t, "Test User", profile.User.Name)
	assert.Len(t, profile.Teams, 1)
	assert.Equal(t, "team1", profile.Teams[0].Name)
	assert.Len(t, profile.Organizations, 1)
	assert.Equal(t, "org1", profile.Organizations[0].Login)
	mockAPI.AssertExpectations(t)
}

func TestGitHubClient_GetUserProfile_WhenGetUserFails_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	userError := errors.New("failed to fetch user")
	mockAPI.On("GetUser", ctx, validToken).Return(nil, userError)

	// Act
	profile, err := client.GetUserProfile(ctx, validToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "failed to get user")
	mockAPI.AssertExpectations(t)
}

func TestGitHubClient_GetUserProfile_WhenGetTeamsFails_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	mockAPI.On("GetUser", ctx, validToken).Return(&github.User{
		ID:    github.Int64(123),
		Login: github.String("testuser"),
	}, nil)

	teamsError := errors.New("failed to fetch teams")
	mockAPI.On("GetUserTeams", ctx, validToken).Return(nil, teamsError)

	// Act
	profile, err := client.GetUserProfile(ctx, validToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "failed to get teams")
	mockAPI.AssertExpectations(t)
}

func TestGitHubClient_GetUserProfile_WhenGetOrganizationsFails_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	mockAPI.On("GetUser", ctx, validToken).Return(&github.User{
		ID:    github.Int64(123),
		Login: github.String("testuser"),
	}, nil)
	mockAPI.On("GetUserTeams", ctx, validToken).Return([]*github.Team{}, nil)

	orgsError := errors.New("failed to fetch organizations")
	mockAPI.On("GetUserOrganizations", ctx, validToken).Return(nil, orgsError)

	// Act
	profile, err := client.GetUserProfile(ctx, validToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "failed to get organizations")
	mockAPI.AssertExpectations(t)
}

// Test GitHubClient.GetUserTeamsAdvanced scenarios

func TestGitHubClient_GetUserTeamsAdvanced_WithOrgFilter_ReturnsFilteredTeams(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	mockAPI.On("GetUser", ctx, validToken).Return(&github.User{
		ID:    github.Int64(123),
		Login: github.String("testuser"),
	}, nil)
	mockAPI.On("GetUserTeams", ctx, validToken).Return([]*github.Team{}, nil)
	mockAPI.On("GetUserOrganizations", ctx, validToken).Return([]*github.Organization{}, nil)

	// Act
	teams, err := client.GetUserTeamsAdvanced(ctx, validToken, "org1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, teams)
	assert.Len(t, teams, 0) // Empty when no teams are mocked
	mockAPI.AssertExpectations(t)
}

func TestGitHubClient_GetUserTeamsAdvanced_WhenGetUserProfileFails_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	userError := errors.New("profile error")
	mockAPI.On("GetUser", ctx, validToken).Return(nil, userError)

	// Act
	teams, err := client.GetUserTeamsAdvanced(ctx, validToken, "org1")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, teams)
	assert.Contains(t, err.Error(), "failed to get user profile")
	mockAPI.AssertExpectations(t)
}

// Test GitHubClient.ValidateTeamMembership scenarios

func TestGitHubClient_ValidateTeamMembership_WithValidMembership_ReturnsResult(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	mockAPI.On("GetUser", ctx, validToken).Return(&github.User{
		ID:    github.Int64(123),
		Login: github.String("testuser"),
	}, nil)
	mockAPI.On("GetUserTeams", ctx, validToken).Return([]*github.Team{}, nil)
	mockAPI.On("GetUserOrganizations", ctx, validToken).Return([]*github.Organization{}, nil)

	// Act
	valid, err := client.ValidateTeamMembership(ctx, validToken, "org1", "team1")

	// Assert
	assert.NoError(t, err)
	assert.False(t, valid) // False since no teams with org1/team1 are mocked
	mockAPI.AssertExpectations(t)
}

func TestGitHubClient_ValidateTeamMembership_WhenGetUserProfileFails_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	userError := errors.New("profile error")
	mockAPI.On("GetUser", ctx, validToken).Return(nil, userError)

	// Act
	valid, err := client.ValidateTeamMembership(ctx, validToken, "org1", "team1")

	// Assert
	assert.Error(t, err)
	assert.False(t, valid)
	assert.Contains(t, err.Error(), "failed to get user profile")
	mockAPI.AssertExpectations(t)
}

// Test GitHubClient.ValidateOrganizationMembership scenarios

func TestGitHubClient_ValidateOrganizationMembership_WithValidMembership_ReturnsResult(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	mockAPI.On("GetUser", ctx, validToken).Return(&github.User{
		ID:    github.Int64(123),
		Login: github.String("testuser"),
	}, nil)
	mockAPI.On("GetUserTeams", ctx, validToken).Return([]*github.Team{}, nil)
	mockAPI.On("GetUserOrganizations", ctx, validToken).Return([]*github.Organization{}, nil)

	// Act
	valid, err := client.ValidateOrganizationMembership(ctx, validToken, "org1")

	// Assert
	assert.NoError(t, err)
	assert.False(t, valid) // False since no organizations are mocked
	mockAPI.AssertExpectations(t)
}

func TestGitHubClient_ValidateOrganizationMembership_WhenGetUserProfileFails_ReturnsError(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)
	ctx := context.Background()

	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(time.Hour),
	}

	userError := errors.New("profile error")
	mockAPI.On("GetUser", ctx, validToken).Return(nil, userError)

	// Act
	valid, err := client.ValidateOrganizationMembership(ctx, validToken, "org1")

	// Assert
	assert.Error(t, err)
	assert.False(t, valid)
	assert.Contains(t, err.Error(), "failed to get user profile")
	mockAPI.AssertExpectations(t)
}

// Test GitHubClient.convertToUserProfile

func TestGitHubClient_convertToUserProfile_ConvertsAllFieldsCorrectly(t *testing.T) {
	// Arrange
	mockAPI := new(MockGitHubAPIClient)
	teamResolver := githubLogic.NewTeamResolver()
	client := NewGitHubClient(mockAPI, teamResolver)

	user := &github.User{
		ID:          github.Int64(123),
		Login:       github.String("testuser"),
		Name:        github.String("Test User"),
		Email:       github.String("test@example.com"),
		AvatarURL:   github.String("https://avatar.com/user"),
		HTMLURL:     github.String("https://github.com/testuser"),
		Type:        github.String("User"),
		SiteAdmin:   github.Bool(false),
		Company:     github.String("Test Company"),
		Location:    github.String("Test Location"),
		Bio:         github.String("Test Bio"),
		Blog:        github.String("https://testuser.com"),
		PublicRepos: github.Int(10),
		Followers:   github.Int(5),
		Following:   github.Int(3),
	}

	teams := []*github.Team{
		{
			ID:         github.Int64(1),
			Name:       github.String("team1"),
			Slug:       github.String("team1"),
			Permission: github.String("admin"),
			Organization: &github.Organization{
				ID:    github.Int64(1),
				Login: github.String("org1"),
				Name:  github.String("Organization 1"),
			},
		},
	}

	orgs := []*github.Organization{
		{
			ID:          github.Int64(1),
			Login:       github.String("org1"),
			Name:        github.String("Organization 1"),
			Description: github.String("Test organization"),
			AvatarURL:   github.String("https://avatar.com/org"),
			HTMLURL:     github.String("https://github.com/org1"),
			Company:     github.String("Org Company"),
			Location:    github.String("Org Location"),
			Email:       github.String("org@example.com"),
			Blog:        github.String("https://org1.com"),
			PublicRepos: github.Int(20),
		},
	}

	// Act
	profile := client.convertToUserProfile(user, teams, orgs)

	// Assert
	assert.NotNil(t, profile)
	assert.Equal(t, int64(123), profile.User.ID)
	assert.Equal(t, "testuser", profile.User.Login)
	assert.Equal(t, "Test User", profile.User.Name)
	assert.Equal(t, "test@example.com", profile.User.Email)
	assert.Equal(t, "Test Company", profile.User.Company)
	assert.Equal(t, "Test Location", profile.User.Location)
	assert.Equal(t, 10, profile.User.PublicRepos)

	assert.Len(t, profile.Teams, 1)
	assert.Equal(t, "team1", profile.Teams[0].Name)
	assert.Equal(t, "team1", profile.Teams[0].Slug)
	assert.Equal(t, "admin", profile.Teams[0].Permission)

	assert.Len(t, profile.Organizations, 1)
	assert.Equal(t, "org1", profile.Organizations[0].Login)
	assert.Equal(t, "Organization 1", profile.Organizations[0].Name)
	assert.Equal(t, "Org Company", profile.Organizations[0].Company)
}
