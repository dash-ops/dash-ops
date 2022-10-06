package github

import (
	"context"
	"encoding/json"
	"io/ioutil"

	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

// Client interface
type Client interface {
	GetUserLogger(token *oauth2.Token) (*UserInfo, error)
	GetTeamsUserLogger(token *oauth2.Token) ([]*gh.Team, error)
}

type client struct {
	oauthConfig *oauth2.Config
}

// NewClient Create a new github client
func NewClient(oauthConfig *oauth2.Config) (client, error) {
	return client{oauthConfig}, nil
}

func getClient(oauthConfig *oauth2.Config, token *oauth2.Token) *gh.Client {
	return gh.NewClient(oauthConfig.Client(context.Background(), token))
}

func (auth client) GetUserLogger(token *oauth2.Token) (*UserInfo, error) {
	c := getClient(auth.oauthConfig, token)
	resp, err := c.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result UserInfo
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	user, _, err := c.Users.Get(context.Background(), "")
	return user, err
}

func (gh client) GetTeamsUserLogger(token *oauth2.Token) ([]*gh.Team, error) {
	client := getClient(gh.oauthConfig, token)
	opt := gh.ListOptions{}
	teams, _, err := client.Teams.ListUserTeams(context.Background(), &opt)
	return teams, err
}
