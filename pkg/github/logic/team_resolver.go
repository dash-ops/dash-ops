package github

import (
	"fmt"
	"strings"
	"time"

	githubModels "github.com/dash-ops/dash-ops/pkg/github/models"
)

// TeamResolver provides team resolution and membership logic
type TeamResolver struct{}

// NewTeamResolver creates a new team resolver
func NewTeamResolver() *TeamResolver {
	return &TeamResolver{}
}

// ResolveUserTeams resolves all teams for a user across organizations
func (tr *TeamResolver) ResolveUserTeams(userProfile *githubModels.UserProfile, orgLogin string) []githubModels.GitHubTeam {
	if userProfile == nil {
		return []githubModels.GitHubTeam{}
	}

	if orgLogin == "" {
		return userProfile.Teams
	}

	return userProfile.GetTeamsInOrganization(orgLogin)
}

// ValidateTeamMembership validates if user is member of a specific team
func (tr *TeamResolver) ValidateTeamMembership(userProfile *githubModels.UserProfile, orgLogin, teamSlug string) (bool, error) {
	if userProfile == nil {
		return false, fmt.Errorf("user profile is required")
	}

	if orgLogin == "" {
		return false, fmt.Errorf("organization login is required")
	}

	if teamSlug == "" {
		return false, fmt.Errorf("team slug is required")
	}

	team := userProfile.GetTeamBySlug(orgLogin, teamSlug)
	return team != nil, nil
}

// ValidateOrganizationMembership validates if user is member of organization
func (tr *TeamResolver) ValidateOrganizationMembership(userProfile *githubModels.UserProfile, orgLogin string) (bool, error) {
	if userProfile == nil {
		return false, fmt.Errorf("user profile is required")
	}

	if orgLogin == "" {
		return false, fmt.Errorf("organization login is required")
	}

	org := userProfile.GetOrganizationByLogin(orgLogin)
	return org != nil, nil
}

// GetUserPermissionLevel determines user's highest permission level in organization
func (tr *TeamResolver) GetUserPermissionLevel(userProfile *githubModels.UserProfile, orgLogin string) PermissionLevel {
	if userProfile == nil {
		return PermissionLevelNone
	}

	// Check if user is org admin
	if userProfile.IsOrgAdmin(orgLogin) {
		return PermissionLevelAdmin
	}

	// Check team memberships for highest permission
	teams := userProfile.GetTeamsInOrganization(orgLogin)
	highestLevel := PermissionLevelNone

	for _, team := range teams {
		level := tr.getTeamPermissionLevel(team)
		if level > highestLevel {
			highestLevel = level
		}
	}

	// If user is in organization but no teams, they have member level
	if highestLevel == PermissionLevelNone {
		org := userProfile.GetOrganizationByLogin(orgLogin)
		if org != nil {
			return PermissionLevelMember
		}
	}

	return highestLevel
}

// FilterTeamsByPermission filters teams by minimum permission level
func (tr *TeamResolver) FilterTeamsByPermission(teams []githubModels.GitHubTeam, minLevel PermissionLevel) []githubModels.GitHubTeam {
	var filtered []githubModels.GitHubTeam

	for _, team := range teams {
		if tr.getTeamPermissionLevel(team) >= minLevel {
			filtered = append(filtered, team)
		}
	}

	return filtered
}

// GetTeamHierarchy builds team hierarchy for organization
func (tr *TeamResolver) GetTeamHierarchy(teams []githubModels.GitHubTeam, orgLogin string) *TeamHierarchy {
	hierarchy := &TeamHierarchy{
		Organization: orgLogin,
		Teams:        []TeamNode{},
	}

	// For now, create flat structure (GitHub teams don't have native hierarchy)
	for _, team := range teams {
		if team.Organization != nil && strings.EqualFold(team.Organization.Login, orgLogin) {
			hierarchy.Teams = append(hierarchy.Teams, TeamNode{
				Team:     team,
				Children: []TeamNode{}, // Flat structure
				Level:    0,
			})
		}
	}

	return hierarchy
}

// getTeamPermissionLevel determines permission level from team
func (tr *TeamResolver) getTeamPermissionLevel(team githubModels.GitHubTeam) PermissionLevel {
	switch team.Permission {
	case "admin":
		return PermissionLevelAdmin
	case "push":
		return PermissionLevelWrite
	case "pull":
		return PermissionLevelRead
	default:
		return PermissionLevelMember
	}
}

// PermissionLevel represents permission levels
type PermissionLevel int

const (
	PermissionLevelNone PermissionLevel = iota
	PermissionLevelRead
	PermissionLevelMember
	PermissionLevelWrite
	PermissionLevelAdmin
)

// String returns string representation of permission level
func (pl PermissionLevel) String() string {
	switch pl {
	case PermissionLevelNone:
		return "none"
	case PermissionLevelRead:
		return "read"
	case PermissionLevelMember:
		return "member"
	case PermissionLevelWrite:
		return "write"
	case PermissionLevelAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

// TeamHierarchy represents team hierarchy in an organization
type TeamHierarchy struct {
	Organization string     `json:"organization"`
	Teams        []TeamNode `json:"teams"`
	LastUpdated  time.Time  `json:"last_updated"`
}

// TeamNode represents a node in team hierarchy
type TeamNode struct {
	Team     githubModels.GitHubTeam `json:"team"`
	Children []TeamNode              `json:"children,omitempty"`
	Level    int                     `json:"level"`
}
