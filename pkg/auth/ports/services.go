package ports

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHubService defines the interface for GitHub operations needed by auth
type GitHubService interface {
	// User operations
	GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error)
	GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error)
}

// Optional: Future provider interfaces
type AuthProviderService interface {
	GetUser(ctx context.Context, token *oauth2.Token) (interface{}, error)
	GetUserGroups(ctx context.Context, token *oauth2.Token) ([]string, error)
}
