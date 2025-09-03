package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubUser_IsOrganization(t *testing.T) {
	tests := []struct {
		name     string
		user     GitHubUser
		expected bool
	}{
		{
			name: "user account",
			user: GitHubUser{
				Login: "testuser",
				Type:  "User",
			},
			expected: false,
		},
		{
			name: "organization account",
			user: GitHubUser{
				Login: "testorg",
				Type:  "Organization",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.IsOrganization()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGitHubUser_GetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		user     GitHubUser
		expected string
	}{
		{
			name: "with name",
			user: GitHubUser{
				Login: "testuser",
				Name:  "Test User",
			},
			expected: "Test User",
		},
		{
			name: "without name",
			user: GitHubUser{
				Login: "testuser",
				Name:  "",
			},
			expected: "testuser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.GetDisplayName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGitHubUser_Validate(t *testing.T) {
	tests := []struct {
		name        string
		user        GitHubUser
		expectError bool
	}{
		{
			name: "valid user",
			user: GitHubUser{
				ID:    12345,
				Login: "testuser",
			},
			expectError: false,
		},
		{
			name: "missing login",
			user: GitHubUser{
				ID: 12345,
			},
			expectError: true,
		},
		{
			name: "missing ID",
			user: GitHubUser{
				Login: "testuser",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGitHubTeam_HasMember(t *testing.T) {
	team := GitHubTeam{
		Name: "test-team",
		Members: []GitHubUser{
			{Login: "user1"},
			{Login: "user2"},
			{Login: "User3"}, // Test case sensitivity
		},
	}

	tests := []struct {
		name      string
		userLogin string
		expected  bool
	}{
		{
			name:      "existing member",
			userLogin: "user1",
			expected:  true,
		},
		{
			name:      "case insensitive match",
			userLogin: "USER3",
			expected:  true,
		},
		{
			name:      "non-existing member",
			userLogin: "user4",
			expected:  false,
		},
		{
			name:      "empty login",
			userLogin: "",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := team.HasMember(tt.userLogin)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGitHubTeam_GetFullName(t *testing.T) {
	tests := []struct {
		name     string
		team     GitHubTeam
		expected string
	}{
		{
			name: "with organization",
			team: GitHubTeam{
				Slug: "developers",
				Organization: &GitHubOrganization{
					Login: "myorg",
				},
			},
			expected: "myorg/developers",
		},
		{
			name: "without organization",
			team: GitHubTeam{
				Slug: "developers",
			},
			expected: "developers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.team.GetFullName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGitHubRepository_CanUserPush(t *testing.T) {
	tests := []struct {
		name     string
		repo     GitHubRepository
		expected bool
	}{
		{
			name: "user can push",
			repo: GitHubRepository{
				Permissions: &RepositoryPermissions{
					Push: true,
				},
			},
			expected: true,
		},
		{
			name: "user cannot push",
			repo: GitHubRepository{
				Permissions: &RepositoryPermissions{
					Push: false,
				},
			},
			expected: false,
		},
		{
			name: "no permissions set",
			repo: GitHubRepository{
				Permissions: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.repo.CanUserPush()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGitHubIssue_HasLabel(t *testing.T) {
	issue := GitHubIssue{
		Labels: []GitHubLabel{
			{Name: "bug"},
			{Name: "enhancement"},
			{Name: "High Priority"},
		},
	}

	tests := []struct {
		name      string
		labelName string
		expected  bool
	}{
		{
			name:      "existing label",
			labelName: "bug",
			expected:  true,
		},
		{
			name:      "case insensitive match",
			labelName: "HIGH PRIORITY",
			expected:  true,
		},
		{
			name:      "non-existing label",
			labelName: "documentation",
			expected:  false,
		},
		{
			name:      "empty label",
			labelName: "",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := issue.HasLabel(tt.labelName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGitHubPullRequest_GetChangeSize(t *testing.T) {
	pr := GitHubPullRequest{
		Additions: 150,
		Deletions: 75,
	}

	result := pr.GetChangeSize()
	assert.Equal(t, 225, result)
}

func TestGitHubPullRequest_IsLargePR(t *testing.T) {
	tests := []struct {
		name     string
		pr       GitHubPullRequest
		expected bool
	}{
		{
			name: "large PR",
			pr: GitHubPullRequest{
				Additions: 400,
				Deletions: 200, // Total: 600 > 500
			},
			expected: true,
		},
		{
			name: "small PR",
			pr: GitHubPullRequest{
				Additions: 50,
				Deletions: 25, // Total: 75 < 500
			},
			expected: false,
		},
		{
			name: "medium PR",
			pr: GitHubPullRequest{
				Additions: 300,
				Deletions: 200, // Total: 500 = 500
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pr.IsLargePR()
			assert.Equal(t, tt.expected, result)
		})
	}
}
