package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	githubAdapters "github.com/dash-ops/dash-ops/pkg/github/adapters/external"
	githubControllers "github.com/dash-ops/dash-ops/pkg/github/controllers"
	githubLogic "github.com/dash-ops/dash-ops/pkg/github/logic"
	githubModels "github.com/dash-ops/dash-ops/pkg/github/models"
)

// Module represents the github module - main entry point for the plugin
type Module struct {
	controller   *githubControllers.GitHubController
	githubClient *githubLogic.GitHubClient
}

// NewModule creates and initializes a new github module (main factory)
func NewModule(oauthConfig *oauth2.Config) (*Module, error) {
	if oauthConfig == nil {
		return nil, fmt.Errorf("oauth config cannot be nil")
	}

	// Initialize logic components
	teamResolver := githubLogic.NewTeamResolver()

	// Initialize external adapters
	apiAdapter := githubAdapters.NewGitHubAPIAdapter(oauthConfig)

	// Initialize controller
	controller := githubControllers.NewGitHubController(
		apiAdapter,
		teamResolver,
		oauthConfig,
	)

	return &Module{
		controller:   controller,
		githubClient: nil, // Not needed for now
	}, nil
}

// GetController returns the GitHub controller
func (m *Module) GetController() *githubControllers.GitHubController {
	return m.controller
}

// GetGitHubClient returns the GitHub client
func (m *Module) GetGitHubClient() *githubLogic.GitHubClient {
	return m.githubClient
}

// Client interface - for dependency injection into other modules
type Client interface {
	GetUserLogger(token *oauth2.Token) (*github.User, error)
	GetTeamsUserLogger(token *oauth2.Token) ([]*github.Team, error)
}

// NewClient creates a new github client - for dependency injection
func NewClient(oauthConfig *oauth2.Config) (Client, error) {
	module, err := NewModule(oauthConfig)
	if err != nil {
		return nil, err
	}

	return module, nil
}

// GetUserLogger gets user information - implements Client interface
func (m *Module) GetUserLogger(token *oauth2.Token) (*github.User, error) {
	return m.controller.GetUser(context.Background(), token)
}

// GetTeamsUserLogger gets user teams - implements Client interface
func (m *Module) GetTeamsUserLogger(token *oauth2.Token) ([]*github.Team, error) {
	return m.controller.GetUserTeams(context.Background(), token)
}

// Advanced methods using the new architecture (for future use)

// GetUserProfile gets enhanced user profile using new models
func (m *Module) GetUserProfile(ctx context.Context, token *oauth2.Token) (*githubModels.UserProfile, error) {
	return m.controller.GetUserProfile(ctx, token)
}

// GetUserTeamsAdvanced gets user teams with enhanced data
func (m *Module) GetUserTeamsAdvanced(ctx context.Context, token *oauth2.Token, orgLogin string) ([]githubModels.GitHubTeam, error) {
	return m.controller.GetUserTeamsAdvanced(ctx, token, orgLogin)
}

// Implementation of auth.GitHubService interface (for dependency injection)

// GetUser implements GitHubService interface for auth module
func (m *Module) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	return m.controller.GetUser(ctx, token)
}

// GetUserTeams implements GitHubService interface for auth module
func (m *Module) GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
	return m.controller.GetUserTeams(ctx, token)
}
