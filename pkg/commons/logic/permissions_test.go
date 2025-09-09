package commons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermissionChecker_HasPermission_WithUserHavingRequiredPermission_ReturnsTrue(t *testing.T) {
	// Arrange
	checker := NewPermissionChecker()
	requiredPermissions := []string{"read", "write"}
	userPermissions := []string{"read", "admin"}

	// Act
	result := checker.HasPermission(requiredPermissions, userPermissions)

	// Assert
	assert.True(t, result)
}

func TestPermissionChecker_HasPermission_WithUserNotHavingRequiredPermission_ReturnsFalse(t *testing.T) {
	// Arrange
	checker := NewPermissionChecker()
	requiredPermissions := []string{"admin"}
	userPermissions := []string{"read", "write"}

	// Act
	result := checker.HasPermission(requiredPermissions, userPermissions)

	// Assert
	assert.False(t, result)
}

func TestPermissionChecker_HasPermission_WithNoPermissionsRequired_ReturnsTrue(t *testing.T) {
	// Arrange
	checker := NewPermissionChecker()
	requiredPermissions := []string{}
	userPermissions := []string{"read"}

	// Act
	result := checker.HasPermission(requiredPermissions, userPermissions)

	// Assert
	assert.True(t, result)
}

func TestPermissionChecker_HasPermission_WithUserHavingNoPermissions_ReturnsFalse(t *testing.T) {
	// Arrange
	checker := NewPermissionChecker()
	requiredPermissions := []string{"read"}
	userPermissions := []string{}

	// Act
	result := checker.HasPermission(requiredPermissions, userPermissions)

	// Assert
	assert.False(t, result)
}

func TestPermissionChecker_HasPermission_WithCaseInsensitiveMatch_ReturnsTrue(t *testing.T) {
	// Arrange
	checker := NewPermissionChecker()
	requiredPermissions := []string{"READ"}
	userPermissions := []string{"read"}

	// Act
	result := checker.HasPermission(requiredPermissions, userPermissions)

	// Assert
	assert.True(t, result)
}

func TestPermissionChecker_HasPermission_WithWhitespaceHandling_ReturnsTrue(t *testing.T) {
	// Arrange
	checker := NewPermissionChecker()
	requiredPermissions := []string{" read "}
	userPermissions := []string{"read"}

	// Act
	result := checker.HasPermission(requiredPermissions, userPermissions)

	// Assert
	assert.True(t, result)
}

func TestPermissionChecker_HasAllPermissions_WithUserHavingAllRequiredPermissions_ReturnsTrue(t *testing.T) {
	// Arrange
	checker := NewPermissionChecker()
	requiredPermissions := []string{"read", "write"}
	userPermissions := []string{"read", "write", "admin"}

	// Act
	result := checker.HasAllPermissions(requiredPermissions, userPermissions)

	// Assert
	assert.True(t, result)
}

func TestPermissionChecker_HasAllPermissions_WithUserMissingOnePermission_ReturnsFalse(t *testing.T) {
	// Arrange
	checker := NewPermissionChecker()
	requiredPermissions := []string{"read", "write", "admin"}
	userPermissions := []string{"read", "write"}

	// Act
	result := checker.HasAllPermissions(requiredPermissions, userPermissions)

	// Assert
	assert.False(t, result)
}

func TestPermissionChecker_HasAllPermissions_WithNoPermissionsRequired_ReturnsTrue(t *testing.T) {
	// Arrange
	checker := NewPermissionChecker()
	requiredPermissions := []string{}
	userPermissions := []string{"read"}

	// Act
	result := checker.HasAllPermissions(requiredPermissions, userPermissions)

	// Assert
	assert.True(t, result)
}

func TestPermissionChecker_HasAllPermissions_WithUserHavingNoPermissions_ReturnsFalse(t *testing.T) {
	// Arrange
	checker := NewPermissionChecker()
	requiredPermissions := []string{"read"}
	userPermissions := []string{}

	// Act
	result := checker.HasAllPermissions(requiredPermissions, userPermissions)

	// Assert
	assert.False(t, result)
}
