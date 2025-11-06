package tempo

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// search performs a search query on Tempo
func (c *TempoClient) search(ctx context.Context, params *wire.TempoQueryParams) (*wire.TempoSearchResponse, error) {
	// Build query URL
	u, err := url.Parse(fmt.Sprintf("%s/api/search", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add query parameters
	q := u.Query()
	if params.Query != "" {
		q.Set("q", params.Query)
	}
	if params.Start > 0 {
		q.Set("start", strconv.FormatInt(params.Start, 10))
	}
	if params.End > 0 {
		q.Set("end", strconv.FormatInt(params.End, 10))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.SpanLimit > 0 {
		q.Set("spss", strconv.Itoa(params.SpanLimit)) // spss = spans per span set
	}
	u.RawQuery = q.Encode()

	// Execute request
	var response wire.TempoSearchResponse
	if err := c.doRequest(ctx, u.String(), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// convertTraceIDToHex converts a base64 traceID to hexadecimal format (required by Tempo)
func convertTraceIDToHex(traceID string) (string, error) {
	// Remove URL encoding if present (e.g., %3D%3D -> ==)
	decoded, err := url.QueryUnescape(traceID)
	if err != nil {
		decoded = traceID
	}

	// Check if already in hex format (only contains hex characters)
	// If it contains non-hex characters, assume it's base64
	if isHexString(decoded) {
		return decoded, nil
	}

	// Decode from base64
	bytes, err := base64.StdEncoding.DecodeString(strings.TrimRight(decoded, "="))
	if err != nil {
		// If decoding fails, try with padding
		bytes, err = base64.StdEncoding.DecodeString(decoded)
		if err != nil {
			return "", fmt.Errorf("failed to decode traceID from base64: %w", err)
		}
	}

	// Convert to hex
	return hex.EncodeToString(bytes), nil
}

// isHexString checks if a string contains only hexadecimal characters
func isHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return len(s) > 0
}

// getTraceByID retrieves a trace by its ID
func (c *TempoClient) getTraceByID(ctx context.Context, traceID string) (*wire.TempoTraceByIDResponse, error) {
	// Convert traceID from base64 to hex (Tempo requires hex format)
	hexTraceID, err := convertTraceIDToHex(traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert traceID to hex: %w", err)
	}

	// Build query URL
	u := fmt.Sprintf("%s/api/traces/%s", c.baseURL, hexTraceID)

	// Execute request
	var response wire.TempoTraceByIDResponse
	if err := c.doRequest(ctx, u, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// getTags retrieves all available tags
func (c *TempoClient) getTags(ctx context.Context) (*wire.TempoSearchTagsResponse, error) {
	u := fmt.Sprintf("%s/api/search/tags", c.baseURL)

	var response wire.TempoSearchTagsResponse
	if err := c.doRequest(ctx, u, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// getTagValues retrieves all values for a specific tag
func (c *TempoClient) getTagValues(ctx context.Context, tagName string) (*wire.TempoSearchTagValuesResponse, error) {
	u := fmt.Sprintf("%s/api/search/tag/%s/values", c.baseURL, url.PathEscape(tagName))

	var response wire.TempoSearchTagValuesResponse
	if err := c.doRequest(ctx, u, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// doRequest executes an HTTP GET request and decodes the JSON response
func (c *TempoClient) doRequest(ctx context.Context, url string, result interface{}) error {
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication if configured
	if c.auth != nil {
		switch c.auth.Type {
		case "bearer":
			if c.auth.Token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.auth.Token))
			}
		case "basic":
			if c.auth.Username != "" && c.auth.Password != "" {
				req.SetBasicAuth(c.auth.Username, c.auth.Password)
			}
		}
	}

	// Set headers
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleErrorResponse(resp.StatusCode, body)
	}

	// Decode JSON response
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to decode response: %w (body: %s)", err, string(body))
	}

	return nil
}

// handleErrorResponse handles error responses from Tempo
func (c *TempoClient) handleErrorResponse(statusCode int, body []byte) error {
	// Try to parse as error response
	var errorResp struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	if err := json.Unmarshal(body, &errorResp); err == nil {
		if errorResp.Message != "" {
			return fmt.Errorf("tempo error (HTTP %d): %s", statusCode, errorResp.Message)
		}
		if errorResp.Error != "" {
			return fmt.Errorf("tempo error (HTTP %d): %s", statusCode, errorResp.Error)
		}
	}

	// Fallback to generic error
	return fmt.Errorf("tempo error (HTTP %d): %s", statusCode, string(body))
}
