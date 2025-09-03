package commons

import (
	"testing"
)

func TestPermissionChecker_HasPermission(t *testing.T) {
	checker := NewPermissionChecker()

	tests := []struct {
		name                string
		requiredPermissions []string
		userPermissions     []string
		expected            bool
	}{
		{
			name:                "user has required permission",
			requiredPermissions: []string{"read", "write"},
			userPermissions:     []string{"read", "admin"},
			expected:            true,
		},
		{
			name:                "user doesn't have required permission",
			requiredPermissions: []string{"admin"},
			userPermissions:     []string{"read", "write"},
			expected:            false,
		},
		{
			name:                "no permissions required",
			requiredPermissions: []string{},
			userPermissions:     []string{"read"},
			expected:            true,
		},
		{
			name:                "user has no permissions",
			requiredPermissions: []string{"read"},
			userPermissions:     []string{},
			expected:            false,
		},
		{
			name:                "case insensitive match",
			requiredPermissions: []string{"READ"},
			userPermissions:     []string{"read"},
			expected:            true,
		},
		{
			name:                "whitespace handling",
			requiredPermissions: []string{" read "},
			userPermissions:     []string{"read"},
			expected:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checker.HasPermission(tt.requiredPermissions, tt.userPermissions)
			if result != tt.expected {
				t.Errorf("HasPermission() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestPermissionChecker_HasAllPermissions(t *testing.T) {
	checker := NewPermissionChecker()

	tests := []struct {
		name                string
		requiredPermissions []string
		userPermissions     []string
		expected            bool
	}{
		{
			name:                "user has all required permissions",
			requiredPermissions: []string{"read", "write"},
			userPermissions:     []string{"read", "write", "admin"},
			expected:            true,
		},
		{
			name:                "user missing one permission",
			requiredPermissions: []string{"read", "write", "admin"},
			userPermissions:     []string{"read", "write"},
			expected:            false,
		},
		{
			name:                "no permissions required",
			requiredPermissions: []string{},
			userPermissions:     []string{"read"},
			expected:            true,
		},
		{
			name:                "user has no permissions",
			requiredPermissions: []string{"read"},
			userPermissions:     []string{},
			expected:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checker.HasAllPermissions(tt.requiredPermissions, tt.userPermissions)
			if result != tt.expected {
				t.Errorf("HasAllPermissions() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
