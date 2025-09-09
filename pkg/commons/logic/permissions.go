package commons

import "strings"

// PermissionChecker provides permission checking logic
type PermissionChecker struct{}

// NewPermissionChecker creates a new permission checker
func NewPermissionChecker() *PermissionChecker {
	return &PermissionChecker{}
}

// HasPermission checks if user has any of the required permissions
func (p *PermissionChecker) HasPermission(requiredPermissions []string, userPermissions []string) bool {
	if len(requiredPermissions) == 0 {
		return true // No permissions required
	}

	if len(userPermissions) == 0 {
		return false // User has no permissions
	}

	for _, required := range requiredPermissions {
		for _, userPerm := range userPermissions {
			if p.permissionMatches(required, userPerm) {
				return true
			}
		}
	}

	return false
}

// HasAllPermissions checks if user has all required permissions
func (p *PermissionChecker) HasAllPermissions(requiredPermissions []string, userPermissions []string) bool {
	if len(requiredPermissions) == 0 {
		return true
	}

	if len(userPermissions) == 0 {
		return false
	}

	for _, required := range requiredPermissions {
		found := false
		for _, userPerm := range userPermissions {
			if p.permissionMatches(required, userPerm) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// permissionMatches checks if two permissions match (case-insensitive)
func (p *PermissionChecker) permissionMatches(required, userPerm string) bool {
	return strings.EqualFold(strings.TrimSpace(required), strings.TrimSpace(userPerm))
}
