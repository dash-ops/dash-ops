package ports

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHubAPIClient defines the interface for GitHub API operations
type GitHubAPIClient interface {
	// User operations
	GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error)
	GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error)
	GetUserOrganizations(ctx context.Context, token *oauth2.Token) ([]*github.Organization, error)
	GetUserRepositories(ctx context.Context, token *oauth2.Token) ([]*github.Repository, error)

	// Organization operations
	GetOrganizationTeams(ctx context.Context, token *oauth2.Token, orgLogin string) ([]*github.Team, error)
}

// GitHubRepository defines repository operations for GitHub data
type GitHubRepository interface {
	// Cache operations (for future implementation)
	CacheUser(ctx context.Context, userID int64, user *github.User) error
	GetCachedUser(ctx context.Context, userID int64) (*github.User, error)

	CacheTeams(ctx context.Context, userID int64, teams []*github.Team) error
	GetCachedTeams(ctx context.Context, userID int64) ([]*github.Team, error)

	// Session/token operations
	StoreUserSession(ctx context.Context, userID int64, token *oauth2.Token) error
	GetUserSession(ctx context.Context, userID int64) (*oauth2.Token, error)
	InvalidateUserSession(ctx context.Context, userID int64) error
}
