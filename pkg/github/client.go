package github

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client interface
type Client interface {
	GetUserLogger(token *oauth2.Token) (*github.User, error)
	GetTeamsUserLogger(token *oauth2.Token) ([]*github.Team, error)
}

type client struct {
	oauthConfig *oauth2.Config
}

// NewClient Create a new github client
func NewClient(oauthConfig *oauth2.Config) (client, error) {
	return client{oauthConfig}, nil
}

func getClient(oauthConfig *oauth2.Config, token *oauth2.Token) *github.Client {
	return github.NewClient(oauthConfig.Client(context.Background(), token))
}

func (gh client) GetUserLogger(token *oauth2.Token) (*github.User, error) {
	client := getClient(gh.oauthConfig, token)
	user, _, err := client.Users.Get(context.Background(), "")
	return user, err
}

func (gh client) GetTeamsUserLogger(token *oauth2.Token) ([]*github.Team, error) {
	client := getClient(gh.oauthConfig, token)
	opt := github.ListOptions{}
	teams, _, err := client.Teams.ListUserTeams(context.Background(), &opt)
	return teams, err
}
