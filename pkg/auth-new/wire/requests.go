package auth

// LoginRequest represents a login request
type LoginRequest struct {
	Provider    string `json:"provider" validate:"required"`
	RedirectURL string `json:"redirect_url,omitempty"`
	State       string `json:"state,omitempty"`
}

// TokenExchangeRequest represents OAuth2 token exchange request
type TokenExchangeRequest struct {
	Code        string `json:"code" validate:"required"`
	State       string `json:"state,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutRequest represents logout request
type LogoutRequest struct {
	SessionID   string `json:"session_id,omitempty"`
	AllSessions bool   `json:"all_sessions,omitempty"`
}
