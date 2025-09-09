package controllers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	githubLogic "github.com/dash-ops/dash-ops/pkg/github/logic"
)

// MockGitHubAPIClient is a mock implementation of GitHubAPIClient for testing
type MockGitHubAPIClient struct {
	GetUserFunc              func(ctx context.Context, token *oauth2.Token) (*github.User, error)
	GetUserTeamsFunc         func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error)
	GetUserOrganizationsFunc func(ctx context.Context, token *oauth2.Token) ([]*github.Organization, error)
	GetUserRepositoriesFunc  func(ctx context.Context, token *oauth2.Token) ([]*github.Repository, error)
	GetOrganizationTeamsFunc func(ctx context.Context, token *oauth2.Token, orgLogin string) ([]*github.Team, error)
}

func (m *MockGitHubAPIClient) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, token)
	}
	return nil, nil
}

func (m *MockGitHubAPIClient) GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
	if m.GetUserTeamsFunc != nil {
		return m.GetUserTeamsFunc(ctx, token)
	}
	return nil, nil
}

func (m *MockGitHubAPIClient) GetUserOrganizations(ctx context.Context, token *oauth2.Token) ([]*github.Organization, error) {
	if m.GetUserOrganizationsFunc != nil {
		return m.GetUserOrganizationsFunc(ctx, token)
	}
	return nil, nil
}

func (m *MockGitHubAPIClient) GetUserRepositories(ctx context.Context, token *oauth2.Token) ([]*github.Repository, error) {
	if m.GetUserRepositoriesFunc != nil {
		return m.GetUserRepositoriesFunc(ctx, token)
	}
	return nil, nil
}

func (m *MockGitHubAPIClient) GetOrganizationTeams(ctx context.Context, token *oauth2.Token, orgLogin string) ([]*github.Team, error) {
	if m.GetOrganizationTeamsFunc != nil {
		return m.GetOrganizationTeamsFunc(ctx, token, orgLogin)
	}
	return nil, nil
}

func TestNewGitHubController_CreatesControllerWithDependencies(t *testing.T) {
	// Arrange
	mockClient := &MockGitHubAPIClient{}
	teamResolver := githubLogic.NewTeamResolver()
	oauthConfig := &oauth2.Config{}

	// Act
	controller := NewGitHubController(mockClient, teamResolver, oauthConfig)

	// Assert
	assert.NotNil(t, controller)
	assert.Equal(t, mockClient, controller.githubClient)
	assert.Equal(t, teamResolver, controller.teamResolver)
	assert.Equal(t, oauthConfig, controller.oauthConfig)
}

func TestGitHubController_GetUser_WithValidToken_ReturnsUser(t *testing.T) {
	// Arrange
	expectedUser := &github.User{
		Login: github.String("testuser"),
		ID:    github.Int64(123),
		Name:  github.String("Test User"),
		Email: github.String("test@example.com"),
	}

	mockClient := &MockGitHubAPIClient{
		GetUserFunc: func(ctx context.Context, token *oauth2.Token) (*github.User, error) {
			return expectedUser, nil
		},
	}

	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	user, err := controller.GetUser(context.Background(), validToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)
	assert.Equal(t, "testuser", user.GetLogin())
	assert.Equal(t, int64(123), user.GetID())
}

func TestGitHubController_GetUser_WithNilToken_ReturnsError(t *testing.T) {
	// Arrange
	mockClient := &MockGitHubAPIClient{}
	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})

	// Act
	user, err := controller.GetUser(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "token is required")
}

func TestGitHubController_GetUser_WithInvalidToken_ReturnsError(t *testing.T) {
	// Arrange
	mockClient := &MockGitHubAPIClient{}
	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})

	invalidToken := &oauth2.Token{
		AccessToken: "expired-token",
		Expiry:      time.Now().Add(-1 * time.Hour), // Expired
	}

	// Act
	user, err := controller.GetUser(context.Background(), invalidToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "token is invalid")
}

func TestGitHubController_GetUser_WithAPIError_ReturnsError(t *testing.T) {
	// Arrange
	apiError := errors.New("GitHub API error: rate limit exceeded")

	mockClient := &MockGitHubAPIClient{
		GetUserFunc: func(ctx context.Context, token *oauth2.Token) (*github.User, error) {
			return nil, apiError
		},
	}

	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	user, err := controller.GetUser(context.Background(), validToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, apiError, err)
}

func TestGitHubController_GetUserTeams_WithValidToken_ReturnsTeams(t *testing.T) {
	// Arrange
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

	mockClient := &MockGitHubAPIClient{
		GetUserTeamsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
			return expectedTeams, nil
		},
	}

	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	teams, err := controller.GetUserTeams(context.Background(), validToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, teams)
	assert.Len(t, teams, 2)
	assert.Equal(t, expectedTeams, teams)
	assert.Equal(t, "Engineering", teams[0].GetName())
	assert.Equal(t, "DevOps", teams[1].GetName())
}

func TestGitHubController_GetUserTeams_WithNilToken_ReturnsError(t *testing.T) {
	// Arrange
	mockClient := &MockGitHubAPIClient{}
	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})

	// Act
	teams, err := controller.GetUserTeams(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, teams)
	assert.Contains(t, err.Error(), "token is required")
}

func TestGitHubController_GetUserProfile_WithCompleteData_ReturnsFullProfile(t *testing.T) {
	// Arrange
	mockUser := &github.User{
		ID:        github.Int64(123),
		Login:     github.String("johndoe"),
		Name:      github.String("John Doe"),
		Email:     github.String("john@example.com"),
		AvatarURL: github.String("https://avatar.url"),
		Company:   github.String("TechCorp"),
	}

	mockTeams := []*github.Team{
		{
			ID:   github.Int64(1),
			Name: github.String("Backend Team"),
			Slug: github.String("backend-team"),
			Organization: &github.Organization{
				ID:    github.Int64(100),
				Login: github.String("techcorp"),
				Name:  github.String("TechCorp Inc"),
			},
			Permission: github.String("admin"),
		},
	}

	mockOrgs := []*github.Organization{
		{
			ID:          github.Int64(100),
			Login:       github.String("techcorp"),
			Name:        github.String("TechCorp Inc"),
			Description: github.String("Technology Company"),
			AvatarURL:   github.String("https://org.avatar.url"),
		},
	}

	mockClient := &MockGitHubAPIClient{
		GetUserFunc: func(ctx context.Context, token *oauth2.Token) (*github.User, error) {
			return mockUser, nil
		},
		GetUserTeamsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
			return mockTeams, nil
		},
		GetUserOrganizationsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Organization, error) {
			return mockOrgs, nil
		},
	}

	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	profile, err := controller.GetUserProfile(context.Background(), validToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, profile)

	// Assert User data
	assert.Equal(t, int64(123), profile.User.ID)
	assert.Equal(t, "johndoe", profile.User.Login)
	assert.Equal(t, "John Doe", profile.User.Name)
	assert.Equal(t, "john@example.com", profile.User.Email)
	assert.Equal(t, "TechCorp", profile.User.Company)

	// Assert Teams data
	assert.Len(t, profile.Teams, 1)
	assert.Equal(t, "Backend Team", profile.Teams[0].Name)
	assert.Equal(t, "backend-team", profile.Teams[0].Slug)
	assert.Equal(t, "admin", profile.Teams[0].Permission)
	assert.NotNil(t, profile.Teams[0].Organization)
	assert.Equal(t, "techcorp", profile.Teams[0].Organization.Login)

	// Assert Organizations data
	assert.Len(t, profile.Organizations, 1)
	assert.Equal(t, int64(100), profile.Organizations[0].ID)
	assert.Equal(t, "techcorp", profile.Organizations[0].Login)
	assert.Equal(t, "TechCorp Inc", profile.Organizations[0].Name)
	assert.Equal(t, "Technology Company", profile.Organizations[0].Description)
}

func TestGitHubController_GetUserProfile_WhenGetUserFails_ReturnsError(t *testing.T) {
	// Arrange
	userError := errors.New("failed to fetch user")

	mockClient := &MockGitHubAPIClient{
		GetUserFunc: func(ctx context.Context, token *oauth2.Token) (*github.User, error) {
			return nil, userError
		},
	}

	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	profile, err := controller.GetUserProfile(context.Background(), validToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "failed to get user")
}

func TestGitHubController_GetUserProfile_WhenGetTeamsFails_ReturnsError(t *testing.T) {
	// Arrange
	mockClient := &MockGitHubAPIClient{
		GetUserFunc: func(ctx context.Context, token *oauth2.Token) (*github.User, error) {
			return &github.User{Login: github.String("user")}, nil
		},
		GetUserTeamsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
			return nil, errors.New("teams API error")
		},
	}

	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	profile, err := controller.GetUserProfile(context.Background(), validToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "failed to get teams")
}

func TestGitHubController_GetUserTeamsAdvanced_WithOrgFilter_ReturnsFilteredTeams(t *testing.T) {
	// Arrange
	mockUser := &github.User{
		Login: github.String("testuser"),
	}

	mockTeams := []*github.Team{
		{
			Slug: github.String("team1"),
			Organization: &github.Organization{
				Login: github.String("org1"),
			},
		},
		{
			Slug: github.String("team2"),
			Organization: &github.Organization{
				Login: github.String("org2"),
			},
		},
	}

	mockClient := &MockGitHubAPIClient{
		GetUserFunc: func(ctx context.Context, token *oauth2.Token) (*github.User, error) {
			return mockUser, nil
		},
		GetUserTeamsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
			return mockTeams, nil
		},
		GetUserOrganizationsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Organization, error) {
			return []*github.Organization{}, nil
		},
	}

	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	teams, err := controller.GetUserTeamsAdvanced(context.Background(), validToken, "org1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, teams)
	// Team resolver will filter by org
}

func TestGitHubController_ValidateTeamMembership_UserIsMember_ReturnsTrue(t *testing.T) {
	// Arrange
	mockTeams := []*github.Team{
		{
			Slug: github.String("engineering"),
			Organization: &github.Organization{
				Login: github.String("techcorp"),
			},
		},
	}

	mockClient := &MockGitHubAPIClient{
		GetUserFunc: func(ctx context.Context, token *oauth2.Token) (*github.User, error) {
			return &github.User{Login: github.String("johndoe")}, nil
		},
		GetUserTeamsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
			return mockTeams, nil
		},
		GetUserOrganizationsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Organization, error) {
			return []*github.Organization{}, nil
		},
	}

	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	isMember, err := controller.ValidateTeamMembership(context.Background(), validToken, "techcorp", "engineering")

	// Assert
	assert.NoError(t, err)
	assert.True(t, isMember)
}

func TestGitHubController_ValidateOrganizationMembership_UserIsMember_ReturnsTrue(t *testing.T) {
	// Arrange
	mockOrgs := []*github.Organization{
		{
			Login: github.String("techcorp"),
		},
	}

	mockClient := &MockGitHubAPIClient{
		GetUserFunc: func(ctx context.Context, token *oauth2.Token) (*github.User, error) {
			return &github.User{Login: github.String("johndoe")}, nil
		},
		GetUserTeamsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
			return []*github.Team{}, nil
		},
		GetUserOrganizationsFunc: func(ctx context.Context, token *oauth2.Token) ([]*github.Organization, error) {
			return mockOrgs, nil
		},
	}

	controller := NewGitHubController(mockClient, githubLogic.NewTeamResolver(), &oauth2.Config{})
	validToken := &oauth2.Token{
		AccessToken: "valid-token",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	// Act
	isMember, err := controller.ValidateOrganizationMembership(context.Background(), validToken, "techcorp")

	// Assert
	assert.NoError(t, err)
	assert.True(t, isMember)
}

func TestGitHubController_convertToUserProfile_ConvertsAllFieldsCorrectly(t *testing.T) {
	// Arrange
	controller := NewGitHubController(&MockGitHubAPIClient{}, githubLogic.NewTeamResolver(), &oauth2.Config{})

	user := &github.User{
		ID:          github.Int64(456),
		Login:       github.String("janedoe"),
		Name:        github.String("Jane Doe"),
		Email:       github.String("jane@example.com"),
		AvatarURL:   github.String("https://avatar.jane.url"),
		HTMLURL:     github.String("https://github.com/janedoe"),
		Type:        github.String("User"),
		SiteAdmin:   github.Bool(false),
		Company:     github.String("StartupCo"),
		Location:    github.String("San Francisco"),
		Bio:         github.String("Software Engineer"),
		Blog:        github.String("https://janedoe.blog"),
		PublicRepos: github.Int(42),
		Followers:   github.Int(100),
		Following:   github.Int(50),
	}

	teams := []*github.Team{
		{
			ID:         github.Int64(10),
			Name:       github.String("Frontend Team"),
			Slug:       github.String("frontend-team"),
			Permission: github.String("push"),
			Organization: &github.Organization{
				ID:    github.Int64(200),
				Login: github.String("startupco"),
				Name:  github.String("StartupCo Inc"),
			},
		},
	}

	orgs := []*github.Organization{
		{
			ID:          github.Int64(200),
			Login:       github.String("startupco"),
			Name:        github.String("StartupCo Inc"),
			Description: github.String("Innovative Startup"),
			AvatarURL:   github.String("https://org.startup.url"),
			HTMLURL:     github.String("https://github.com/startupco"),
			Company:     github.String("StartupCo"),
			Location:    github.String("Silicon Valley"),
			Email:       github.String("contact@startupco.com"),
			Blog:        github.String("https://startupco.com"),
			PublicRepos: github.Int(25),
		},
	}

	// Act
	profile := controller.convertToUserProfile(user, teams, orgs)

	// Assert
	assert.NotNil(t, profile)

	// Assert User conversion
	assert.Equal(t, int64(456), profile.User.ID)
	assert.Equal(t, "janedoe", profile.User.Login)
	assert.Equal(t, "Jane Doe", profile.User.Name)
	assert.Equal(t, "jane@example.com", profile.User.Email)
	assert.Equal(t, "https://avatar.jane.url", profile.User.AvatarURL)
	assert.Equal(t, "https://github.com/janedoe", profile.User.HTMLURL)
	assert.Equal(t, "User", profile.User.Type)
	assert.False(t, profile.User.SiteAdmin)
	assert.Equal(t, "StartupCo", profile.User.Company)
	assert.Equal(t, "San Francisco", profile.User.Location)
	assert.Equal(t, "Software Engineer", profile.User.Bio)
	assert.Equal(t, "https://janedoe.blog", profile.User.Blog)
	assert.Equal(t, 42, profile.User.PublicRepos)
	assert.Equal(t, 100, profile.User.Followers)
	assert.Equal(t, 50, profile.User.Following)

	// Assert Teams conversion
	assert.Len(t, profile.Teams, 1)
	assert.Equal(t, int64(10), profile.Teams[0].ID)
	assert.Equal(t, "Frontend Team", profile.Teams[0].Name)
	assert.Equal(t, "frontend-team", profile.Teams[0].Slug)
	assert.Equal(t, "push", profile.Teams[0].Permission)
	assert.NotNil(t, profile.Teams[0].Organization)
	assert.Equal(t, int64(200), profile.Teams[0].Organization.ID)
	assert.Equal(t, "startupco", profile.Teams[0].Organization.Login)
	assert.Equal(t, "StartupCo Inc", profile.Teams[0].Organization.Name)

	// Assert Organizations conversion
	assert.Len(t, profile.Organizations, 1)
	assert.Equal(t, int64(200), profile.Organizations[0].ID)
	assert.Equal(t, "startupco", profile.Organizations[0].Login)
	assert.Equal(t, "StartupCo Inc", profile.Organizations[0].Name)
	assert.Equal(t, "Innovative Startup", profile.Organizations[0].Description)
	assert.Equal(t, "https://org.startup.url", profile.Organizations[0].AvatarURL)
	assert.Equal(t, "https://github.com/startupco", profile.Organizations[0].HTMLURL)
	assert.Equal(t, "StartupCo", profile.Organizations[0].Company)
	assert.Equal(t, "Silicon Valley", profile.Organizations[0].Location)
	assert.Equal(t, "contact@startupco.com", profile.Organizations[0].Email)
	assert.Equal(t, "https://startupco.com", profile.Organizations[0].Blog)
	assert.Equal(t, 25, profile.Organizations[0].PublicRepos)
}

func TestGitHubController_convertToUserProfile_WithNilOrganizationInTeam_SkipsTeam(t *testing.T) {
	// Arrange
	controller := NewGitHubController(&MockGitHubAPIClient{}, githubLogic.NewTeamResolver(), &oauth2.Config{})

	user := &github.User{
		Login: github.String("testuser"),
	}

	teams := []*github.Team{
		{
			ID:           github.Int64(1),
			Name:         github.String("Team Without Org"),
			Slug:         github.String("team-without-org"),
			Organization: nil, // No organization
		},
		{
			ID:   github.Int64(2),
			Name: github.String("Team With Org"),
			Slug: github.String("team-with-org"),
			Organization: &github.Organization{
				Login: github.String("validorg"),
			},
		},
	}

	// Act
	profile := controller.convertToUserProfile(user, teams, []*github.Organization{})

	// Assert
	assert.NotNil(t, profile)
	assert.Len(t, profile.Teams, 1) // Only team with organization should be included
	assert.Equal(t, "Team With Org", profile.Teams[0].Name)
}
