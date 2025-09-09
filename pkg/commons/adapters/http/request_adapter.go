package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	commonsModels "github.com/dash-ops/dash-ops/pkg/commons/models"
)

// RequestAdapter handles HTTP request parsing and context management
type RequestAdapter struct{}

// NewRequestAdapter creates a new request adapter
func NewRequestAdapter() *RequestAdapter {
	return &RequestAdapter{}
}

// ParseJSON parses JSON request body into the given structure
func (r *RequestAdapter) ParseJSON(req *http.Request, v interface{}) error {
	if req.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields() // Strict parsing

	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}

// GetUserFromContext extracts user data from request context
func (r *RequestAdapter) GetUserFromContext(ctx context.Context) (*commonsModels.UserData, error) {
	userData, ok := ctx.Value(commonsModels.UserDataKey).(*commonsModels.UserData)
	if !ok {
		return nil, fmt.Errorf("user data not found in context")
	}

	if userData == nil {
		return nil, fmt.Errorf("user data is nil")
	}

	return userData, nil
}

// GetTokenFromContext extracts token from request context
func (r *RequestAdapter) GetTokenFromContext(ctx context.Context) (string, error) {
	token, ok := ctx.Value(commonsModels.TokenKey).(string)
	if !ok {
		return "", fmt.Errorf("token not found in context")
	}

	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}

// SetUserInContext sets user data in request context
func (r *RequestAdapter) SetUserInContext(ctx context.Context, userData *commonsModels.UserData) context.Context {
	return context.WithValue(ctx, commonsModels.UserDataKey, userData)
}

// SetTokenInContext sets token in request context
func (r *RequestAdapter) SetTokenInContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, commonsModels.TokenKey, token)
}

// GetAuthorizationHeader extracts the Authorization header from request
func (r *RequestAdapter) GetAuthorizationHeader(req *http.Request) string {
	return req.Header.Get("Authorization")
}

// GetBearerToken extracts bearer token from Authorization header
func (r *RequestAdapter) GetBearerToken(req *http.Request) (string, error) {
	auth := r.GetAuthorizationHeader(req)
	if auth == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	const bearerPrefix = "Bearer "
	if len(auth) < len(bearerPrefix) || auth[:len(bearerPrefix)] != bearerPrefix {
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := auth[len(bearerPrefix):]
	if token == "" {
		return "", fmt.Errorf("bearer token is empty")
	}

	return token, nil
}
