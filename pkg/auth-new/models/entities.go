package auth

import (
	"fmt"
	"time"
)

// AuthProvider represents different authentication providers
type AuthProvider string

const (
	ProviderGitHub AuthProvider = "github"
	ProviderGoogle AuthProvider = "google"
	ProviderOkta   AuthProvider = "okta"
	ProviderSAML   AuthProvider = "saml"
	ProviderLDAP   AuthProvider = "ldap"
	ProviderJWT    AuthProvider = "jwt"
)

// AuthMethod represents different authentication methods
type AuthMethod string

const (
	MethodOAuth2 AuthMethod = "oauth2"
	MethodSAML   AuthMethod = "saml"
	MethodLDAP   AuthMethod = "ldap"
	MethodJWT    AuthMethod = "jwt"
	MethodBasic  AuthMethod = "basic"
)

// User represents an authenticated user
type User struct {
	ID       string       `json:"id"`
	Username string       `json:"username"`
	Name     string       `json:"name"`
	Email    string       `json:"email"`
	Avatar   string       `json:"avatar,omitempty"`
	Provider AuthProvider `json:"provider"`

	// Organization/Team information
	Organizations []Organization `json:"organizations,omitempty"`
	Teams         []Team         `json:"teams,omitempty"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	LastLogin time.Time `json:"last_login"`
}

// Organization represents a user's organization
type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Role string `json:"role,omitempty"` // admin, member, etc.
}

// Team represents a user's team within an organization
type Team struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Organization string `json:"organization"`   // Organization slug
	Role         string `json:"role,omitempty"` // maintainer, member, etc.
}

// Token represents an authentication token
type Token struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	TokenType    string       `json:"token_type"`
	ExpiresAt    time.Time    `json:"expires_at"`
	Scopes       []string     `json:"scopes,omitempty"`
	Provider     AuthProvider `json:"provider"`
}

// AuthConfig represents authentication configuration for a provider
type AuthConfig struct {
	Provider        AuthProvider `yaml:"provider" json:"provider"`
	Method          AuthMethod   `yaml:"method" json:"method"`
	ClientID        string       `yaml:"clientId" json:"client_id"`
	ClientSecret    string       `yaml:"clientSecret" json:"client_secret"`
	AuthURL         string       `yaml:"authURL" json:"auth_url"`
	TokenURL        string       `yaml:"tokenURL" json:"token_url"`
	RedirectURL     string       `yaml:"redirectURL" json:"redirect_url"`
	URLLoginSuccess string       `yaml:"urlLoginSuccess" json:"url_login_success"`
	OrgPermission   string       `yaml:"orgPermission" json:"org_permission"`
	Scopes          []string     `yaml:"scopes" json:"scopes"`
	Enabled         bool         `yaml:"enabled" json:"enabled"`
}

// AuthSession represents an active authentication session
type AuthSession struct {
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	Token     *Token    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	LastUsed  time.Time `json:"last_used"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

// Methods for User entity

// HasOrganization checks if user belongs to a specific organization
func (u *User) HasOrganization(orgSlug string) bool {
	for _, org := range u.Organizations {
		if org.Slug == orgSlug {
			return true
		}
	}
	return false
}

// HasTeam checks if user belongs to a specific team
func (u *User) HasTeam(teamSlug string) bool {
	for _, team := range u.Teams {
		if team.Slug == teamSlug {
			return true
		}
	}
	return false
}

// HasTeamInOrg checks if user belongs to a specific team in an organization
func (u *User) HasTeamInOrg(teamSlug, orgSlug string) bool {
	for _, team := range u.Teams {
		if team.Slug == teamSlug && team.Organization == orgSlug {
			return true
		}
	}
	return false
}

// GetTeamsInOrg returns all teams for a user in a specific organization
func (u *User) GetTeamsInOrg(orgSlug string) []Team {
	var teams []Team
	for _, team := range u.Teams {
		if team.Organization == orgSlug {
			teams = append(teams, team)
		}
	}
	return teams
}

// Methods for Token entity

// IsValid checks if the token is valid and not expired
func (t *Token) IsValid() bool {
	if t.AccessToken == "" {
		return false
	}
	return time.Now().Before(t.ExpiresAt)
}

// IsExpired checks if the token is expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// TimeToExpiry returns time until token expires
func (t *Token) TimeToExpiry() time.Duration {
	return time.Until(t.ExpiresAt)
}

// Methods for AuthConfig entity

// Validate validates the authentication configuration
func (ac *AuthConfig) Validate() error {
	if ac.Provider == "" {
		return fmt.Errorf("provider is required")
	}

	if ac.Method == "" {
		return fmt.Errorf("method is required")
	}

	// OAuth2 specific validation
	if ac.Method == MethodOAuth2 {
		if ac.ClientID == "" {
			return fmt.Errorf("clientId is required for OAuth2")
		}
		if ac.ClientSecret == "" {
			return fmt.Errorf("clientSecret is required for OAuth2")
		}
		if ac.AuthURL == "" {
			return fmt.Errorf("authURL is required for OAuth2")
		}
		if ac.TokenURL == "" {
			return fmt.Errorf("tokenURL is required for OAuth2")
		}
		if ac.RedirectURL == "" {
			return fmt.Errorf("redirectURL is required for OAuth2")
		}
	}

	return nil
}

// IsOAuth2 checks if the config is for OAuth2 authentication
func (ac *AuthConfig) IsOAuth2() bool {
	return ac.Method == MethodOAuth2
}

// GetScopes returns scopes with defaults if empty
func (ac *AuthConfig) GetScopes() []string {
	if len(ac.Scopes) == 0 {
		// Default scopes based on provider
		switch ac.Provider {
		case ProviderGitHub:
			return []string{"user", "read:org"}
		case ProviderGoogle:
			return []string{"openid", "profile", "email"}
		default:
			return []string{}
		}
	}
	return ac.Scopes
}

// Methods for AuthSession entity

// IsActive checks if the session is active and not expired
func (as *AuthSession) IsActive() bool {
	if as.Token == nil || !as.Token.IsValid() {
		return false
	}
	return time.Now().Before(as.ExpiresAt)
}

// UpdateLastUsed updates the last used timestamp
func (as *AuthSession) UpdateLastUsed() {
	as.LastUsed = time.Now()
}
