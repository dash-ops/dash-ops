package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	githubModels "github.com/dash-ops/dash-ops/pkg/github/models"
)

// Test TeamResolver.ResolveUserTeams scenarios

func TestTeamResolver_ResolveUserTeams_WithNilUserProfile_ReturnsEmptyList(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	
	// Act
	result := resolver.ResolveUserTeams(nil, "test-org")
	
	// Assert
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestTeamResolver_ResolveUserTeams_WithEmptyOrgLogin_ReturnsAllTeams(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Teams: []githubModels.GitHubTeam{
			{Name: "team1", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
			{Name: "team2", Organization: &githubModels.GitHubOrganization{Login: "org2"}},
		},
	}
	
	// Act
	result := resolver.ResolveUserTeams(userProfile, "")
	
	// Assert
	assert.Len(t, result, 2)
	assert.Equal(t, "team1", result[0].Name)
	assert.Equal(t, "team2", result[1].Name)
}

func TestTeamResolver_ResolveUserTeams_WithOrgFilter_ReturnsFilteredTeams(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Teams: []githubModels.GitHubTeam{
			{Name: "team1", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
			{Name: "team2", Organization: &githubModels.GitHubOrganization{Login: "org2"}},
			{Name: "team3", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
		},
	}
	
	// Act
	result := resolver.ResolveUserTeams(userProfile, "org1")
	
	// Assert
	assert.Len(t, result, 2)
	assert.Equal(t, "team1", result[0].Name)
	assert.Equal(t, "team3", result[1].Name)
	assert.Equal(t, "org1", result[0].Organization.Login)
	assert.Equal(t, "org1", result[1].Organization.Login)
}

// Test TeamResolver.ValidateTeamMembership scenarios

func TestTeamResolver_ValidateTeamMembership_WithNilUserProfile_ReturnsError(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	
	// Act
	result, err := resolver.ValidateTeamMembership(nil, "org1", "team1")
	
	// Assert
	assert.Error(t, err)
	assert.False(t, result)
	assert.Contains(t, err.Error(), "user profile is required")
}

func TestTeamResolver_ValidateTeamMembership_WithEmptyOrgLogin_ReturnsError(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Teams: []githubModels.GitHubTeam{
			{Name: "team1", Slug: "team1", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
		},
	}
	
	// Act
	result, err := resolver.ValidateTeamMembership(userProfile, "", "team1")
	
	// Assert
	assert.Error(t, err)
	assert.False(t, result)
	assert.Contains(t, err.Error(), "organization login is required")
}

func TestTeamResolver_ValidateTeamMembership_WithEmptyTeamSlug_ReturnsError(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Teams: []githubModels.GitHubTeam{
			{Name: "team1", Slug: "team1", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
		},
	}
	
	// Act
	result, err := resolver.ValidateTeamMembership(userProfile, "org1", "")
	
	// Assert
	assert.Error(t, err)
	assert.False(t, result)
	assert.Contains(t, err.Error(), "team slug is required")
}

func TestTeamResolver_ValidateTeamMembership_WithValidMembership_ReturnsTrue(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Teams: []githubModels.GitHubTeam{
			{Name: "team1", Slug: "team1", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
			{Name: "team2", Slug: "team2", Organization: &githubModels.GitHubOrganization{Login: "org2"}},
		},
	}
	
	// Act
	result, err := resolver.ValidateTeamMembership(userProfile, "org1", "team1")
	
	// Assert
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestTeamResolver_ValidateTeamMembership_WithInvalidMembership_ReturnsFalse(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Teams: []githubModels.GitHubTeam{
			{Name: "team1", Slug: "team1", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
			{Name: "team2", Slug: "team2", Organization: &githubModels.GitHubOrganization{Login: "org2"}},
		},
	}
	
	// Act
	result, err := resolver.ValidateTeamMembership(userProfile, "org1", "team2")
	
	// Assert
	assert.NoError(t, err)
	assert.False(t, result) // team2 belongs to org2, not org1
}

// Test TeamResolver.ValidateOrganizationMembership scenarios

func TestTeamResolver_ValidateOrganizationMembership_WithNilUserProfile_ReturnsError(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	
	// Act
	result, err := resolver.ValidateOrganizationMembership(nil, "org1")
	
	// Assert
	assert.Error(t, err)
	assert.False(t, result)
	assert.Contains(t, err.Error(), "user profile is required")
}

func TestTeamResolver_ValidateOrganizationMembership_WithEmptyOrgLogin_ReturnsError(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Organizations: []githubModels.GitHubOrganization{
			{Login: "org1"},
		},
	}
	
	// Act
	result, err := resolver.ValidateOrganizationMembership(userProfile, "")
	
	// Assert
	assert.Error(t, err)
	assert.False(t, result)
	assert.Contains(t, err.Error(), "organization login is required")
}

func TestTeamResolver_ValidateOrganizationMembership_WithValidMembership_ReturnsTrue(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Organizations: []githubModels.GitHubOrganization{
			{Login: "org1"},
			{Login: "org2"},
		},
	}
	
	// Act
	result, err := resolver.ValidateOrganizationMembership(userProfile, "org1")
	
	// Assert
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestTeamResolver_ValidateOrganizationMembership_WithInvalidMembership_ReturnsFalse(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Organizations: []githubModels.GitHubOrganization{
			{Login: "org1"},
			{Login: "org2"},
		},
	}
	
	// Act
	result, err := resolver.ValidateOrganizationMembership(userProfile, "org3")
	
	// Assert
	assert.NoError(t, err)
	assert.False(t, result)
}

// Test TeamResolver.GetUserPermissionLevel scenarios

func TestTeamResolver_GetUserPermissionLevel_WithNilUserProfile_ReturnsNone(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	
	// Act
	result := resolver.GetUserPermissionLevel(nil, "org1")
	
	// Assert
	assert.Equal(t, PermissionLevelNone, result)
}

func TestTeamResolver_GetUserPermissionLevel_WithAdminTeamPermission_ReturnsAdmin(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Teams: []githubModels.GitHubTeam{
			{
				Permission:   "admin",
				Organization: &githubModels.GitHubOrganization{Login: "org1"},
			},
		},
	}
	
	// Act
	result := resolver.GetUserPermissionLevel(userProfile, "org1")
	
	// Assert
	assert.Equal(t, PermissionLevelAdmin, result)
}

func TestTeamResolver_GetUserPermissionLevel_WithWriteTeamPermission_ReturnsWrite(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Teams: []githubModels.GitHubTeam{
			{
				Permission:   "push",
				Organization: &githubModels.GitHubOrganization{Login: "org1"},
			},
		},
	}
	
	// Act
	result := resolver.GetUserPermissionLevel(userProfile, "org1")
	
	// Assert
	assert.Equal(t, PermissionLevelWrite, result)
}

func TestTeamResolver_GetUserPermissionLevel_WithReadTeamPermission_ReturnsRead(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Teams: []githubModels.GitHubTeam{
			{
				Permission:   "pull",
				Organization: &githubModels.GitHubOrganization{Login: "org1"},
			},
		},
	}
	
	// Act
	result := resolver.GetUserPermissionLevel(userProfile, "org1")
	
	// Assert
	assert.Equal(t, PermissionLevelRead, result)
}

func TestTeamResolver_GetUserPermissionLevel_WithOrgMembershipNoTeams_ReturnsMember(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Organizations: []githubModels.GitHubOrganization{
			{Login: "org1"},
		},
	}
	
	// Act
	result := resolver.GetUserPermissionLevel(userProfile, "org1")
	
	// Assert
	assert.Equal(t, PermissionLevelMember, result)
}

func TestTeamResolver_GetUserPermissionLevel_WithNoPermission_ReturnsNone(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	userProfile := &githubModels.UserProfile{
		Organizations: []githubModels.GitHubOrganization{
			{Login: "org2"},
		},
	}
	
	// Act
	result := resolver.GetUserPermissionLevel(userProfile, "org1")
	
	// Assert
	assert.Equal(t, PermissionLevelNone, result)
}

// Test TeamResolver.FilterTeamsByPermission scenarios

func TestTeamResolver_FilterTeamsByPermission_FilterByAdminLevel_ReturnsOnlyAdmins(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	teams := []githubModels.GitHubTeam{
		{Permission: "admin", Name: "admin-team"},
		{Permission: "push", Name: "write-team"},
		{Permission: "pull", Name: "read-team"},
		{Permission: "member", Name: "member-team"},
	}
	
	// Act
	result := resolver.FilterTeamsByPermission(teams, PermissionLevelAdmin)
	
	// Assert
	assert.Len(t, result, 1)
	assert.Equal(t, "admin-team", result[0].Name)
}

func TestTeamResolver_FilterTeamsByPermission_FilterByWriteLevel_ReturnsWriteAndAbove(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	teams := []githubModels.GitHubTeam{
		{Permission: "admin", Name: "admin-team"},
		{Permission: "push", Name: "write-team"},
		{Permission: "pull", Name: "read-team"},
		{Permission: "member", Name: "member-team"},
	}
	
	// Act
	result := resolver.FilterTeamsByPermission(teams, PermissionLevelWrite)
	
	// Assert
	assert.Len(t, result, 2)
	assert.Equal(t, "admin-team", result[0].Name)
	assert.Equal(t, "write-team", result[1].Name)
}

func TestTeamResolver_FilterTeamsByPermission_FilterByReadLevel_ReturnsReadAndAbove(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	teams := []githubModels.GitHubTeam{
		{Permission: "admin", Name: "admin-team"},
		{Permission: "push", Name: "write-team"},
		{Permission: "pull", Name: "read-team"},
		{Permission: "member", Name: "member-team"},
	}
	
	// Act
	result := resolver.FilterTeamsByPermission(teams, PermissionLevelRead)
	
	// Assert
	assert.Len(t, result, 3)
	assert.Equal(t, "admin-team", result[0].Name)
	assert.Equal(t, "write-team", result[1].Name)
	assert.Equal(t, "read-team", result[2].Name)
}

func TestTeamResolver_FilterTeamsByPermission_FilterByMemberLevel_ReturnsAllTeams(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	teams := []githubModels.GitHubTeam{
		{Permission: "admin", Name: "admin-team"},
		{Permission: "push", Name: "write-team"},
		{Permission: "pull", Name: "read-team"},
		{Permission: "member", Name: "member-team"},
	}
	
	// Act
	result := resolver.FilterTeamsByPermission(teams, PermissionLevelMember)
	
	// Assert
	assert.Len(t, result, 4)
	assert.Equal(t, teams, result)
}

// Test TeamResolver.GetTeamHierarchy scenarios

func TestTeamResolver_GetTeamHierarchy_FilterByOrganization_ReturnsCorrectHierarchy(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	teams := []githubModels.GitHubTeam{
		{Name: "team1", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
		{Name: "team2", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
		{Name: "team3", Organization: &githubModels.GitHubOrganization{Login: "org2"}},
	}
	
	// Act
	result := resolver.GetTeamHierarchy(teams, "org1")
	
	// Assert
	require.NotNil(t, result)
	assert.Equal(t, "org1", result.Organization)
	assert.Len(t, result.Teams, 2)
	assert.Equal(t, "team1", result.Teams[0].Team.Name)
	assert.Equal(t, "team2", result.Teams[1].Team.Name)
}

func TestTeamResolver_GetTeamHierarchy_FilterByDifferentOrganization_ReturnsCorrectHierarchy(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	teams := []githubModels.GitHubTeam{
		{Name: "team1", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
		{Name: "team2", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
		{Name: "team3", Organization: &githubModels.GitHubOrganization{Login: "org2"}},
	}
	
	// Act
	result := resolver.GetTeamHierarchy(teams, "org2")
	
	// Assert
	require.NotNil(t, result)
	assert.Equal(t, "org2", result.Organization)
	assert.Len(t, result.Teams, 1)
	assert.Equal(t, "team3", result.Teams[0].Team.Name)
}

func TestTeamResolver_GetTeamHierarchy_NoTeamsForOrganization_ReturnsEmptyHierarchy(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	teams := []githubModels.GitHubTeam{
		{Name: "team1", Organization: &githubModels.GitHubOrganization{Login: "org1"}},
		{Name: "team2", Organization: &githubModels.GitHubOrganization{Login: "org2"}},
	}
	
	// Act
	result := resolver.GetTeamHierarchy(teams, "org3")
	
	// Assert
	require.NotNil(t, result)
	assert.Equal(t, "org3", result.Organization)
	assert.Empty(t, result.Teams)
}

// Test TeamResolver.getTeamPermissionLevel scenarios

func TestTeamResolver_getTeamPermissionLevel_WithAdminPermission_ReturnsAdmin(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	team := githubModels.GitHubTeam{Permission: "admin"}
	
	// Act
	result := resolver.getTeamPermissionLevel(team)
	
	// Assert
	assert.Equal(t, PermissionLevelAdmin, result)
}

func TestTeamResolver_getTeamPermissionLevel_WithPushPermission_ReturnsWrite(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	team := githubModels.GitHubTeam{Permission: "push"}
	
	// Act
	result := resolver.getTeamPermissionLevel(team)
	
	// Assert
	assert.Equal(t, PermissionLevelWrite, result)
}

func TestTeamResolver_getTeamPermissionLevel_WithPullPermission_ReturnsRead(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	team := githubModels.GitHubTeam{Permission: "pull"}
	
	// Act
	result := resolver.getTeamPermissionLevel(team)
	
	// Assert
	assert.Equal(t, PermissionLevelRead, result)
}

func TestTeamResolver_getTeamPermissionLevel_WithUnknownPermission_ReturnsMember(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	team := githubModels.GitHubTeam{Permission: "unknown"}
	
	// Act
	result := resolver.getTeamPermissionLevel(team)
	
	// Assert
	assert.Equal(t, PermissionLevelMember, result)
}

func TestTeamResolver_getTeamPermissionLevel_WithEmptyPermission_ReturnsMember(t *testing.T) {
	// Arrange
	resolver := NewTeamResolver()
	team := githubModels.GitHubTeam{Permission: ""}
	
	// Act
	result := resolver.getTeamPermissionLevel(team)
	
	// Assert
	assert.Equal(t, PermissionLevelMember, result)
}

// Test PermissionLevel.String scenarios

func TestPermissionLevel_String_None_ReturnsNone(t *testing.T) {
	// Arrange
	level := PermissionLevelNone
	
	// Act
	result := level.String()
	
	// Assert
	assert.Equal(t, "none", result)
}

func TestPermissionLevel_String_Read_ReturnsRead(t *testing.T) {
	// Arrange
	level := PermissionLevelRead
	
	// Act
	result := level.String()
	
	// Assert
	assert.Equal(t, "read", result)
}

func TestPermissionLevel_String_Member_ReturnsMember(t *testing.T) {
	// Arrange
	level := PermissionLevelMember
	
	// Act
	result := level.String()
	
	// Assert
	assert.Equal(t, "member", result)
}

func TestPermissionLevel_String_Write_ReturnsWrite(t *testing.T) {
	// Arrange
	level := PermissionLevelWrite
	
	// Act
	result := level.String()
	
	// Assert
	assert.Equal(t, "write", result)
}

func TestPermissionLevel_String_Admin_ReturnsAdmin(t *testing.T) {
	// Arrange
	level := PermissionLevelAdmin
	
	// Act
	result := level.String()
	
	// Assert
	assert.Equal(t, "admin", result)
}

func TestPermissionLevel_String_Unknown_ReturnsUnknown(t *testing.T) {
	// Arrange
	level := PermissionLevel(999)
	
	// Act
	result := level.String()
	
	// Assert
	assert.Equal(t, "unknown", result)
}