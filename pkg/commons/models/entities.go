package commons

import "time"

// UserData represents user information from authentication
type UserData struct {
	Org       string    `json:"org"`
	Groups    []string  `json:"groups"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Token     string    `json:"-"` // Don't serialize token
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// HasGroup checks if user belongs to a specific group
func (u *UserData) HasGroup(group string) bool {
	for _, g := range u.Groups {
		if g == group {
			return true
		}
	}
	return false
}

// HasAnyGroup checks if user belongs to any of the specified groups
func (u *UserData) HasAnyGroup(groups []string) bool {
	for _, group := range groups {
		if u.HasGroup(group) {
			return true
		}
	}
	return false
}

// IsValid checks if the user data is valid and not expired
func (u *UserData) IsValid() bool {
	if u.Org == "" || u.Token == "" {
		return false
	}
	if !u.ExpiresAt.IsZero() && u.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}
