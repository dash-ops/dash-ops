package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubUser_IsOrganization_WithUserAccount_ReturnsFalse(t *testing.T) {
	// Arrange
	user := GitHubUser{
		Login: "testuser",
		Type:  "User",
	}

	// Act
	result := user.IsOrganization()

	// Assert
	assert.False(t, result)
}

func TestGitHubUser_IsOrganization_WithOrganizationAccount_ReturnsTrue(t *testing.T) {
	// Arrange
	user := GitHubUser{
		Login: "testorg",
		Type:  "Organization",
	}

	// Act
	result := user.IsOrganization()

	// Assert
	assert.True(t, result)
}

func TestGitHubUser_GetDisplayName_WithName_ReturnsName(t *testing.T) {
	// Arrange
	user := GitHubUser{
		Login: "testuser",
		Name:  "Test User",
	}

	// Act
	result := user.GetDisplayName()

	// Assert
	assert.Equal(t, "Test User", result)
}

func TestGitHubUser_GetDisplayName_WithoutName_ReturnsLogin(t *testing.T) {
	// Arrange
	user := GitHubUser{
		Login: "testuser",
		Name:  "",
	}

	// Act
	result := user.GetDisplayName()

	// Assert
	assert.Equal(t, "testuser", result)
}

func TestGitHubUser_Validate_WithValidUser_ReturnsNoError(t *testing.T) {
	// Arrange
	user := GitHubUser{
		ID:    12345,
		Login: "testuser",
	}

	// Act
	err := user.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestGitHubUser_Validate_WithMissingLogin_ReturnsError(t *testing.T) {
	// Arrange
	user := GitHubUser{
		ID: 12345,
	}

	// Act
	err := user.Validate()

	// Assert
	assert.Error(t, err)
}

func TestGitHubUser_Validate_WithMissingID_ReturnsError(t *testing.T) {
	// Arrange
	user := GitHubUser{
		Login: "testuser",
	}

	// Act
	err := user.Validate()

	// Assert
	assert.Error(t, err)
}

func TestGitHubTeam_HasMember_WithExistingMember_ReturnsTrue(t *testing.T) {
	// Arrange
	team := GitHubTeam{
		Name: "test-team",
		Members: []GitHubUser{
			{Login: "user1"},
			{Login: "user2"},
			{Login: "User3"}, // Test case sensitivity
		},
	}
	userLogin := "user1"

	// Act
	result := team.HasMember(userLogin)

	// Assert
	assert.True(t, result)
}

func TestGitHubTeam_HasMember_WithCaseInsensitiveMatch_ReturnsTrue(t *testing.T) {
	// Arrange
	team := GitHubTeam{
		Name: "test-team",
		Members: []GitHubUser{
			{Login: "user1"},
			{Login: "user2"},
			{Login: "User3"}, // Test case sensitivity
		},
	}
	userLogin := "USER3"

	// Act
	result := team.HasMember(userLogin)

	// Assert
	assert.True(t, result)
}

func TestGitHubTeam_HasMember_WithNonExistingMember_ReturnsFalse(t *testing.T) {
	// Arrange
	team := GitHubTeam{
		Name: "test-team",
		Members: []GitHubUser{
			{Login: "user1"},
			{Login: "user2"},
			{Login: "User3"}, // Test case sensitivity
		},
	}
	userLogin := "user4"

	// Act
	result := team.HasMember(userLogin)

	// Assert
	assert.False(t, result)
}

func TestGitHubTeam_HasMember_WithEmptyLogin_ReturnsFalse(t *testing.T) {
	// Arrange
	team := GitHubTeam{
		Name: "test-team",
		Members: []GitHubUser{
			{Login: "user1"},
			{Login: "user2"},
			{Login: "User3"}, // Test case sensitivity
		},
	}
	userLogin := ""

	// Act
	result := team.HasMember(userLogin)

	// Assert
	assert.False(t, result)
}

func TestGitHubTeam_GetFullName_WithOrganization_ReturnsFullName(t *testing.T) {
	// Arrange
	team := GitHubTeam{
		Slug: "developers",
		Organization: &GitHubOrganization{
			Login: "myorg",
		},
	}

	// Act
	result := team.GetFullName()

	// Assert
	assert.Equal(t, "myorg/developers", result)
}

func TestGitHubTeam_GetFullName_WithoutOrganization_ReturnsSlug(t *testing.T) {
	// Arrange
	team := GitHubTeam{
		Slug: "developers",
	}

	// Act
	result := team.GetFullName()

	// Assert
	assert.Equal(t, "developers", result)
}

func TestGitHubRepository_CanUserPush_WithPushPermission_ReturnsTrue(t *testing.T) {
	// Arrange
	repo := GitHubRepository{
		Permissions: &RepositoryPermissions{
			Push: true,
		},
	}

	// Act
	result := repo.CanUserPush()

	// Assert
	assert.True(t, result)
}

func TestGitHubRepository_CanUserPush_WithoutPushPermission_ReturnsFalse(t *testing.T) {
	// Arrange
	repo := GitHubRepository{
		Permissions: &RepositoryPermissions{
			Push: false,
		},
	}

	// Act
	result := repo.CanUserPush()

	// Assert
	assert.False(t, result)
}

func TestGitHubRepository_CanUserPush_WithNoPermissionsSet_ReturnsFalse(t *testing.T) {
	// Arrange
	repo := GitHubRepository{
		Permissions: nil,
	}

	// Act
	result := repo.CanUserPush()

	// Assert
	assert.False(t, result)
}

func TestGitHubIssue_HasLabel_WithExistingLabel_ReturnsTrue(t *testing.T) {
	// Arrange
	issue := GitHubIssue{
		Labels: []GitHubLabel{
			{Name: "bug"},
			{Name: "enhancement"},
			{Name: "High Priority"},
		},
	}
	labelName := "bug"

	// Act
	result := issue.HasLabel(labelName)

	// Assert
	assert.True(t, result)
}

func TestGitHubIssue_HasLabel_WithCaseInsensitiveMatch_ReturnsTrue(t *testing.T) {
	// Arrange
	issue := GitHubIssue{
		Labels: []GitHubLabel{
			{Name: "bug"},
			{Name: "enhancement"},
			{Name: "High Priority"},
		},
	}
	labelName := "HIGH PRIORITY"

	// Act
	result := issue.HasLabel(labelName)

	// Assert
	assert.True(t, result)
}

func TestGitHubIssue_HasLabel_WithNonExistingLabel_ReturnsFalse(t *testing.T) {
	// Arrange
	issue := GitHubIssue{
		Labels: []GitHubLabel{
			{Name: "bug"},
			{Name: "enhancement"},
			{Name: "High Priority"},
		},
	}
	labelName := "documentation"

	// Act
	result := issue.HasLabel(labelName)

	// Assert
	assert.False(t, result)
}

func TestGitHubIssue_HasLabel_WithEmptyLabel_ReturnsFalse(t *testing.T) {
	// Arrange
	issue := GitHubIssue{
		Labels: []GitHubLabel{
			{Name: "bug"},
			{Name: "enhancement"},
			{Name: "High Priority"},
		},
	}
	labelName := ""

	// Act
	result := issue.HasLabel(labelName)

	// Assert
	assert.False(t, result)
}

func TestGitHubPullRequest_GetChangeSize_WithAdditionsAndDeletions_ReturnsTotal(t *testing.T) {
	// Arrange
	pr := GitHubPullRequest{
		Additions: 150,
		Deletions: 75,
	}

	// Act
	result := pr.GetChangeSize()

	// Assert
	assert.Equal(t, 225, result)
}

func TestGitHubPullRequest_IsLargePR_WithLargePR_ReturnsTrue(t *testing.T) {
	// Arrange
	pr := GitHubPullRequest{
		Additions: 400,
		Deletions: 200, // Total: 600 > 500
	}

	// Act
	result := pr.IsLargePR()

	// Assert
	assert.True(t, result)
}

func TestGitHubPullRequest_IsLargePR_WithSmallPR_ReturnsFalse(t *testing.T) {
	// Arrange
	pr := GitHubPullRequest{
		Additions: 50,
		Deletions: 25, // Total: 75 < 500
	}

	// Act
	result := pr.IsLargePR()

	// Assert
	assert.False(t, result)
}

func TestGitHubPullRequest_IsLargePR_WithMediumPR_ReturnsFalse(t *testing.T) {
	// Arrange
	pr := GitHubPullRequest{
		Additions: 300,
		Deletions: 200, // Total: 500 = 500
	}

	// Act
	result := pr.IsLargePR()

	// Assert
	assert.False(t, result)
}
