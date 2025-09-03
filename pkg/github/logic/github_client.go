package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	githubModels "github.com/dash-ops/dash-ops/pkg/github/models"
	githubPorts "github.com/dash-ops/dash-ops/pkg/github/ports"
)

// GitHubClient handles GitHub business logic
type GitHubClient struct {
	apiClient    githubPorts.GitHubAPIClient
	teamResolver *TeamResolver
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient(apiClient githubPorts.GitHubAPIClient, teamResolver *TeamResolver) *GitHubClient {
	return &GitHubClient{
		apiClient:    apiClient,
		teamResolver: teamResolver,
	}
}

// GetUser gets user information (legacy compatibility)
func (gc *GitHubClient) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	if token == nil {
		return nil, fmt.Errorf("token is required")
	}

	if !token.Valid() {
		return nil, fmt.Errorf("token is invalid")
	}

	return gc.apiClient.GetUser(ctx, token)
}

// GetUserTeams gets user teams (legacy compatibility)
func (gc *GitHubClient) GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
	if token == nil {
		return nil, fmt.Errorf("token is required")
	}

	if !token.Valid() {
		return nil, fmt.Errorf("token is invalid")
	}

	return gc.apiClient.GetUserTeams(ctx, token)
}

// GetUserProfile gets enhanced user profile using new models
func (gc *GitHubClient) GetUserProfile(ctx context.Context, token *oauth2.Token) (*githubModels.UserProfile, error) {
	// Get basic user info
	user, err := gc.GetUser(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get teams
	teams, err := gc.GetUserTeams(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}

	// Get organizations
	orgs, err := gc.apiClient.GetUserOrganizations(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	// Convert to enhanced model
	return gc.convertToUserProfile(user, teams, orgs), nil
}

// GetUserTeamsAdvanced gets user teams with enhanced data
func (gc *GitHubClient) GetUserTeamsAdvanced(ctx context.Context, token *oauth2.Token, orgLogin string) ([]githubModels.GitHubTeam, error) {
	profile, err := gc.GetUserProfile(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Use team resolver to filter by organization
	teams := gc.teamResolver.ResolveUserTeams(profile, orgLogin)
	return teams, nil
}

// ValidateTeamMembership validates if user is member of specific team
func (gc *GitHubClient) ValidateTeamMembership(ctx context.Context, token *oauth2.Token, orgLogin, teamSlug string) (bool, error) {
	profile, err := gc.GetUserProfile(ctx, token)
	if err != nil {
		return false, fmt.Errorf("failed to get user profile: %w", err)
	}

	return gc.teamResolver.ValidateTeamMembership(profile, orgLogin, teamSlug)
}

// ValidateOrganizationMembership validates if user is member of organization
func (gc *GitHubClient) ValidateOrganizationMembership(ctx context.Context, token *oauth2.Token, orgLogin string) (bool, error) {
	profile, err := gc.GetUserProfile(ctx, token)
	if err != nil {
		return false, fmt.Errorf("failed to get user profile: %w", err)
	}

	return gc.teamResolver.ValidateOrganizationMembership(profile, orgLogin)
}

// convertToUserProfile converts GitHub API data to UserProfile aggregate
func (gc *GitHubClient) convertToUserProfile(user *github.User, teams []*github.Team, orgs []*github.Organization) *githubModels.UserProfile {
	// Convert basic user info
	enhancedUser := githubModels.GitHubUser{
		ID:          user.GetID(),
		Login:       user.GetLogin(),
		Name:        user.GetName(),
		Email:       user.GetEmail(),
		AvatarURL:   user.GetAvatarURL(),
		HTMLURL:     user.GetHTMLURL(),
		Type:        user.GetType(),
		SiteAdmin:   user.GetSiteAdmin(),
		Company:     user.GetCompany(),
		Location:    user.GetLocation(),
		Bio:         user.GetBio(),
		Blog:        user.GetBlog(),
		PublicRepos: user.GetPublicRepos(),
		Followers:   user.GetFollowers(),
		Following:   user.GetFollowing(),
	}

	// Convert teams
	var enhancedTeams []githubModels.GitHubTeam
	for _, team := range teams {
		if team.Organization != nil {
			enhancedTeam := githubModels.GitHubTeam{
				ID:   team.GetID(),
				Name: team.GetName(),
				Slug: team.GetSlug(),
				Organization: &githubModels.GitHubOrganization{
					ID:    team.Organization.GetID(),
					Login: team.Organization.GetLogin(),
					Name:  team.Organization.GetName(),
				},
				Permission: team.GetPermission(),
			}
			enhancedTeams = append(enhancedTeams, enhancedTeam)
		}
	}

	// Convert organizations
	var enhancedOrgs []githubModels.GitHubOrganization
	for _, org := range orgs {
		enhancedOrg := githubModels.GitHubOrganization{
			ID:          org.GetID(),
			Login:       org.GetLogin(),
			Name:        org.GetName(),
			Description: org.GetDescription(),
			AvatarURL:   org.GetAvatarURL(),
			HTMLURL:     org.GetHTMLURL(),
			Company:     org.GetCompany(),
			Location:    org.GetLocation(),
			Email:       org.GetEmail(),
			Blog:        org.GetBlog(),
			PublicRepos: org.GetPublicRepos(),
		}
		enhancedOrgs = append(enhancedOrgs, enhancedOrg)
	}

	return &githubModels.UserProfile{
		User:          enhancedUser,
		Organizations: enhancedOrgs,
		Teams:         enhancedTeams,
	}
}
