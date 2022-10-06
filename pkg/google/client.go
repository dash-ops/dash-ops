// package google

// import (
// 	"context"

// 	"github.com/google/go-github/github"
// 	"golang.org/x/oauth2"
// )

// type User struct {
// 	ID   *int64  `json:"id,omitempty"`
// 	Name *string `json:"name,omitempty"`
// }

// type Team struct {
// 	ID           *int64        `json:"id,omitempty"`
// 	Name         *string       `json:"name,omitempty"`
// 	Slug         *string       `json:"slug,omitempty"`
// 	Organization *Organization `json:"organization,omitempty"`
// }

// type Organization struct {
// 	Login *string `json:"login,omitempty"`
// 	ID    *int64  `json:"id,omitempty"`
// 	Name  *string `json:"name,omitempty"`
// }

// // Client interface
// type Client interface {
// 	GetUserLogger(token *oauth2.Token) (*User, error)
// 	GetTeamsUserLogger(token *oauth2.Token) ([]*Team, error)
// }

// type client struct {
// 	oauthConfig *oauth2.Config
// }

// // NewClient Create a new google client
// func NewClient(oauthConfig *oauth2.Config) (client, error) {
// 	return client{oauthConfig}, nil
// }

// func getClient(oauthConfig *oauth2.Config, token *oauth2.Token) *google.Client {
// 	return github.NewClient(oauthConfig.Client(context.Background(), token))
// }

// func (gh client) GetUserLogger(token *oauth2.Token) (*github.User, error) {
// 	client := getClient(gh.oauthConfig, token)
// 	user, _, err := client.Users.Get(context.Background(), "")
// 	return user, err
// }

// func (gh client) GetTeamsUserLogger(token *oauth2.Token) ([]*github.Team, error) {
// 	client := getClient(gh.oauthConfig, token)
// 	opt := github.ListOptions{}
// 	teams, _, err := client.Teams.ListUserTeams(context.Background(), &opt)
// 	return teams, err
// }
