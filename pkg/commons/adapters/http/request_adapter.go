package http

import (
	"encoding/json"
	"fmt"
	"net/http"
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
