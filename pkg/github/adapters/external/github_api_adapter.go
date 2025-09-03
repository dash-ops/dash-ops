package external

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHubAPIAdapter handles GitHub API integration
type GitHubAPIAdapter struct {
	oauthConfig *oauth2.Config
}

// NewGitHubAPIAdapter creates a new GitHub API adapter
func NewGitHubAPIAdapter(oauthConfig *oauth2.Config) *GitHubAPIAdapter {
	return &GitHubAPIAdapter{
		oauthConfig: oauthConfig,
	}
}

// GetUser gets user information from GitHub API
func (gaa *GitHubAPIAdapter) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	client := gaa.createAuthenticatedClient(token)
	user, _, err := client.Users.Get(ctx, "")
	return user, err
}

// GetUserTeams gets user teams from GitHub API
func (gaa *GitHubAPIAdapter) GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
	client := gaa.createAuthenticatedClient(token)
	opt := github.ListOptions{}
	teams, _, err := client.Teams.ListUserTeams(ctx, &opt)
	return teams, err
}

// GetUserOrganizations gets user organizations from GitHub API
func (gaa *GitHubAPIAdapter) GetUserOrganizations(ctx context.Context, token *oauth2.Token) ([]*github.Organization, error) {
	client := gaa.createAuthenticatedClient(token)
	opt := github.ListOptions{}
	orgs, _, err := client.Organizations.List(ctx, "", &opt)
	return orgs, err
}

// GetUserRepositories gets user repositories from GitHub API
func (gaa *GitHubAPIAdapter) GetUserRepositories(ctx context.Context, token *oauth2.Token) ([]*github.Repository, error) {
	client := gaa.createAuthenticatedClient(token)
	opt := github.RepositoryListOptions{
		ListOptions: github.ListOptions{},
	}
	repos, _, err := client.Repositories.List(ctx, "", &opt)
	return repos, err
}

// GetOrganizationTeams gets teams for a specific organization
func (gaa *GitHubAPIAdapter) GetOrganizationTeams(ctx context.Context, token *oauth2.Token, orgLogin string) ([]*github.Team, error) {
	client := gaa.createAuthenticatedClient(token)
	opt := github.ListOptions{}
	teams, _, err := client.Teams.ListTeams(ctx, orgLogin, &opt)
	return teams, err
}

// createAuthenticatedClient creates an authenticated GitHub client
func (gaa *GitHubAPIAdapter) createAuthenticatedClient(token *oauth2.Token) *github.Client {
	return github.NewClient(gaa.oauthConfig.Client(context.Background(), token))
}
