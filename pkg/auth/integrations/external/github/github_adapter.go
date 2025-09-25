package github

import (
	"context"

	"github.com/dash-ops/dash-ops/pkg/auth/ports"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHubAdapter transforms data between GitHub API and auth domain
type GitHubAdapter struct {
	client *GitHubClient
}

// NewGitHubAdapter creates a new GitHub adapter
func NewGitHubAdapter(oauthConfig *oauth2.Config) ports.GitHubService {
	return &GitHubAdapter{
		client: NewGitHubClient(oauthConfig),
	}
}

// GetUser gets user information from GitHub API
func (a *GitHubAdapter) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	return a.client.GetUser(ctx, token)
}

// GetUserTeams gets user teams from GitHub API
func (a *GitHubAdapter) GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
	return a.client.GetUserTeams(ctx, token)
}
