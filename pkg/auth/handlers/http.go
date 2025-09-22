package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	authAdapters "github.com/dash-ops/dash-ops/pkg/auth/adapters/http"
	authControllers "github.com/dash-ops/dash-ops/pkg/auth/controllers"
	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	commonsModels "github.com/dash-ops/dash-ops/pkg/commons/models"
	"golang.org/x/oauth2"
)

// HTTPHandler handles HTTP requests for auth module
type HTTPHandler struct {
	controller      *authControllers.AuthController
	authAdapter     *authAdapters.AuthAdapter
	responseAdapter *commonsHttp.ResponseAdapter
	requestAdapter  *commonsHttp.RequestAdapter
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(
	controller *authControllers.AuthController,
	authAdapter *authAdapters.AuthAdapter,
	responseAdapter *commonsHttp.ResponseAdapter,
	requestAdapter *commonsHttp.RequestAdapter,
) *HTTPHandler {
	return &HTTPHandler{
		controller:      controller,
		authAdapter:     authAdapter,
		responseAdapter: responseAdapter,
		requestAdapter:  requestAdapter,
	}
}

// RegisterRoutes registers HTTP routes for the auth module
func (h *HTTPHandler) RegisterRoutes(apiRouter, internalRouter *mux.Router) {
	// OAuth2 routes (matching original oauth2 module)
	apiRouter.HandleFunc("/oauth", h.authorizeHandler).Methods("GET").Name("oauth")
	apiRouter.HandleFunc("/oauth/redirect", h.redirectHandler).Methods("GET").Name("oauthRedirect")

	// Add middleware to internal router (matching original)
	internalRouter.Use(h.oAuthMiddleware)

	// Internal routes (matching original provider handlers)
	internalRouter.HandleFunc("/me", h.meHandler).Methods("GET", "OPTIONS").Name("userLogger")
	internalRouter.HandleFunc("/me/permissions", h.mePermissionsHandler).Methods("GET", "OPTIONS").Name("userPermissions")

	// Add organization permission middleware if configured
	// This will be handled by the controller logic
}

// authorizeHandler handles OAuth2 authorization requests
func (h *HTTPHandler) authorizeHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := r.URL.Query().Get("redirect_url")

	// Generate authorization URL using controller
	authURL, err := h.controller.GenerateAuthURL(r.Context(), redirectURL)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to generate auth URL: "+err.Error())
		return
	}

	// Redirect to authorization URL
	http.Redirect(w, r, authURL, http.StatusPermanentRedirect)
}

// redirectHandler handles OAuth2 callback/redirect requests
func (h *HTTPHandler) redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Authorization code is required")
		return
	}

	// Exchange code for token using controller
	token, err := h.controller.ExchangeCodeForToken(r.Context(), code)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusUnauthorized, "Failed to exchange code for token: "+err.Error())
		return
	}

	if !token.Valid() {
		h.responseAdapter.WriteError(w, http.StatusUnauthorized, "Retrieved invalid token")
		return
	}

	// Build redirect URL using controller
	redirectURL := h.controller.BuildRedirectURL(token, state)

	// Redirect to success URL
	http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
}

// meHandler handles user profile requests
func (h *HTTPHandler) meHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from context
	token, ok := r.Context().Value(commonsModels.TokenKey).(*oauth2.Token)
	if !ok {
		h.responseAdapter.WriteError(w, http.StatusUnauthorized, "No token found in context")
		return
	}

	// Get user profile using controller
	user, err := h.controller.GetUserProfile(r.Context(), token)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return user profile (matching original contract)
	h.responseAdapter.WriteJSON(w, http.StatusOK, user)
}

// mePermissionsHandler handles user permissions requests
func (h *HTTPHandler) mePermissionsHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from context
	token, ok := r.Context().Value(commonsModels.TokenKey).(*oauth2.Token)
	if !ok {
		h.responseAdapter.WriteError(w, http.StatusUnauthorized, "No token found in context")
		return
	}

	// Get user permissions using controller
	permissions, err := h.controller.GetUserPermissions(r.Context(), token)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Transform to response format using adapter
	response := h.authAdapter.UserPermissionsToResponse(permissions)

	// Return permissions (matching original contract)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// oAuthMiddleware validates OAuth2 tokens
func (h *HTTPHandler) oAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const bearerSchema = "Bearer "
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			h.responseAdapter.WriteError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		if len(authHeader) <= len(bearerSchema) {
			h.responseAdapter.WriteError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		accessToken := authHeader[len(bearerSchema):]
		if accessToken == "" {
			h.responseAdapter.WriteError(w, http.StatusUnauthorized, "Access token is required")
			return
		}

		// Create token object
		token := &oauth2.Token{AccessToken: accessToken, TokenType: "Bearer"}

		// Validate token using controller
		if err := h.controller.ValidateToken(r.Context(), token); err != nil {
			h.responseAdapter.WriteError(w, http.StatusUnauthorized, "Invalid token: "+err.Error())
			return
		}

		// Add token to context
		ctx := context.WithValue(r.Context(), commonsModels.TokenKey, token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// OrgPermissionMiddleware validates organization permissions
func (h *HTTPHandler) OrgPermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from context
		token, ok := r.Context().Value(commonsModels.TokenKey).(*oauth2.Token)
		if !ok {
			h.responseAdapter.WriteError(w, http.StatusUnauthorized, "No token found in context")
			return
		}

		// Build user data using controller
		userData, err := h.controller.BuildUserData(r.Context(), token)
		if err != nil {
			h.responseAdapter.WriteError(w, http.StatusUnauthorized, "Failed to validate organization permissions: "+err.Error())
			return
		}

		// Add user data to context
		ctx := context.WithValue(r.Context(), commonsModels.UserDataKey, userData)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
