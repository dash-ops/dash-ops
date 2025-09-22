package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	githubAdapters "github.com/dash-ops/dash-ops/pkg/github/adapters/external"
	githubControllers "github.com/dash-ops/dash-ops/pkg/github/controllers"
	githubLogic "github.com/dash-ops/dash-ops/pkg/github/logic"
)

// Module represents the github module - main entry point for the plugin
type Module struct {
	controller *githubControllers.GitHubController
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
		controller: controller,
	}, nil
}

// GetUser implements GitHubService interface for auth module
func (m *Module) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	return m.controller.GetUser(ctx, token)
}

// GetUserTeams implements GitHubService interface for auth module
func (m *Module) GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
	return m.controller.GetUserTeams(ctx, token)
}
