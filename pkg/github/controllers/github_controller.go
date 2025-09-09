package controllers

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	githubLogic "github.com/dash-ops/dash-ops/pkg/github/logic"
	githubModels "github.com/dash-ops/dash-ops/pkg/github/models"
	githubPorts "github.com/dash-ops/dash-ops/pkg/github/ports"
)

// GitHubController orchestrates GitHub business logic
type GitHubController struct {
	githubClient githubPorts.GitHubAPIClient
	teamResolver *githubLogic.TeamResolver
	oauthConfig  *oauth2.Config
}

// NewGitHubController creates a new GitHub controller
func NewGitHubController(
	githubClient githubPorts.GitHubAPIClient,
	teamResolver *githubLogic.TeamResolver,
	oauthConfig *oauth2.Config,
) *GitHubController {
	return &GitHubController{
		githubClient: githubClient,
		teamResolver: teamResolver,
		oauthConfig:  oauthConfig,
	}
}

// GetUser gets user information from GitHub API
func (gc *GitHubController) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
	if token == nil {
		return nil, fmt.Errorf("token is required")
	}

	if !token.Valid() {
		return nil, fmt.Errorf("token is invalid")
	}

	return gc.githubClient.GetUser(ctx, token)
}

// GetUserTeams gets user teams from GitHub API
func (gc *GitHubController) GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error) {
	if token == nil {
		return nil, fmt.Errorf("token is required")
	}

	if !token.Valid() {
		return nil, fmt.Errorf("token is invalid")
	}

	return gc.githubClient.GetUserTeams(ctx, token)
}

// GetUserProfile gets enhanced user profile using new models
func (gc *GitHubController) GetUserProfile(ctx context.Context, token *oauth2.Token) (*githubModels.UserProfile, error) {
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
	orgs, err := gc.githubClient.GetUserOrganizations(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	// Convert to UserProfile aggregate
	return gc.convertToUserProfile(user, teams, orgs), nil
}

// GetUserTeamsAdvanced gets user teams with enhanced data
func (gc *GitHubController) GetUserTeamsAdvanced(ctx context.Context, token *oauth2.Token, orgLogin string) ([]githubModels.GitHubTeam, error) {
	// Get user profile
	profile, err := gc.GetUserProfile(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Use team resolver to filter by organization
	teams := gc.teamResolver.ResolveUserTeams(profile, orgLogin)
	return teams, nil
}

// ValidateTeamMembership validates if user is member of specific team
func (gc *GitHubController) ValidateTeamMembership(ctx context.Context, token *oauth2.Token, orgLogin, teamSlug string) (bool, error) {
	profile, err := gc.GetUserProfile(ctx, token)
	if err != nil {
		return false, fmt.Errorf("failed to get user profile: %w", err)
	}

	return gc.teamResolver.ValidateTeamMembership(profile, orgLogin, teamSlug)
}

// ValidateOrganizationMembership validates if user is member of organization
func (gc *GitHubController) ValidateOrganizationMembership(ctx context.Context, token *oauth2.Token, orgLogin string) (bool, error) {
	profile, err := gc.GetUserProfile(ctx, token)
	if err != nil {
		return false, fmt.Errorf("failed to get user profile: %w", err)
	}

	return gc.teamResolver.ValidateOrganizationMembership(profile, orgLogin)
}

// convertToUserProfile converts GitHub API data to UserProfile aggregate
func (gc *GitHubController) convertToUserProfile(user *github.User, teams []*github.Team, orgs []*github.Organization) *githubModels.UserProfile {
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
