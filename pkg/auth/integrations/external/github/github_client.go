package github

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHubClient handles communication with GitHub API
type GitHubClient struct {
	oauthConfig *oauth2.Config
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient(oauthConfig *oauth2.Config) *GitHubClient {
	return &GitHubClient{
		oauthConfig: oauthConfig,
	}
}

// GetUser gets user information from GitHub API
func (c *GitHubClient) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	client := c.createAuthenticatedClient(token)
	user, _, err := client.Users.Get(ctx, "")
	return user, err
}

// GetUserTeams gets user teams from GitHub API
func (c *GitHubClient) GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
	client := c.createAuthenticatedClient(token)
	opt := github.ListOptions{}
	teams, _, err := client.Teams.ListUserTeams(ctx, &opt)
	return teams, err
}

// createAuthenticatedClient creates an authenticated GitHub client
func (c *GitHubClient) createAuthenticatedClient(token *oauth2.Token) *github.Client {
	return github.NewClient(c.oauthConfig.Client(context.Background(), token))
}
