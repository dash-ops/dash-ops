package oauth2

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GithubClient interface
type GithubClient interface {
	GetUserLogger(token *oauth2.Token) (*github.User, error)
	GetOrgsUserLogger(token *oauth2.Token) ([]*github.Organization, error)
}

type githubClient struct {
	oauthConfig *oauth2.Config
}

// NewGithubClient Create a new github client
func NewGithubClient(oauthConfig *oauth2.Config) (GithubClient, error) {
	return githubClient{oauthConfig}, nil
}

func getClient(oauthConfig *oauth2.Config, token *oauth2.Token) *github.Client {
	return github.NewClient(oauthConfig.Client(context.Background(), token))
}

func (gh githubClient) GetUserLogger(token *oauth2.Token) (*github.User, error) {
	client := getClient(gh.oauthConfig, token)
	user, _, err := client.Users.Get(context.Background(), "")
	return user, err
}

func (gh githubClient) GetOrgsUserLogger(token *oauth2.Token) ([]*github.Organization, error) {
	client := getClient(gh.oauthConfig, token)
	opt := github.ListOptions{}
	orgs, _, err := client.Organizations.List(context.Background(), "", &opt)
	return orgs, err
}
