package auth

import "time"

// AuthResponse represents authentication response
type AuthResponse struct {
	RedirectURL string `json:"redirect_url"`
	State       string `json:"state,omitempty"`
}

// TokenResponse represents token response
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"` // seconds
	ExpiresAt    time.Time `json:"expires_at"`
	Scopes       []string  `json:"scopes,omitempty"`
}

// UserResponse represents user information response
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar,omitempty"`
	Provider string `json:"provider"`

	Organizations []OrganizationResponse `json:"organizations,omitempty"`
	Teams         []TeamResponse         `json:"teams,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	LastLogin time.Time `json:"last_login"`
}

// OrganizationResponse represents organization response
type OrganizationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Role string `json:"role,omitempty"`
}

// TeamResponse represents team response
type TeamResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Organization string `json:"organization"`
	Role         string `json:"role,omitempty"`
}

// SessionResponse represents session information response
type SessionResponse struct {
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	LastUsed  time.Time `json:"last_used"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

// ProvidersResponse represents available providers response
type ProvidersResponse struct {
	Providers []ProviderInfo `json:"providers"`
}

// ProviderInfo represents provider information
type ProviderInfo struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Method      string   `json:"method"`
	Enabled     bool     `json:"enabled"`
	Scopes      []string `json:"scopes,omitempty"`
}

// LogoutResponse represents logout response
type LogoutResponse struct {
	Message       string `json:"message"`
	RedirectURL   string `json:"redirect_url,omitempty"`
	SessionsEnded int    `json:"sessions_ended,omitempty"`
}
