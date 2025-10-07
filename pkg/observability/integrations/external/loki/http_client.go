package loki

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// queryRange queries logs from Loki within a time range (private method)
func (c *LokiClient) queryRange(ctx context.Context, params wire.LokiQueryParams) (*wire.LokiQueryResponse, error) {
	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("query", params.Query)
	queryParams.Set("start", strconv.FormatInt(params.Start.UnixNano(), 10))
	queryParams.Set("end", strconv.FormatInt(params.End.UnixNano(), 10))

	if params.Limit > 0 {
		queryParams.Set("limit", strconv.Itoa(params.Limit))
	}

	if params.Direction != "" {
		queryParams.Set("direction", params.Direction)
	}

	if params.Step != "" {
		queryParams.Set("step", params.Step)
	}

	endpoint := fmt.Sprintf("%s/loki/api/v1/query_range?%s", c.baseURL, queryParams.Encode())

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query range: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result wire.LokiQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// query performs an instant query at a single point in time (private method)
func (c *LokiClient) query(ctx context.Context, query string, ts time.Time, limit int) (*wire.LokiQueryResponse, error) {
	queryParams := url.Values{}
	queryParams.Set("query", query)
	queryParams.Set("time", strconv.FormatInt(ts.UnixNano(), 10))

	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}

	endpoint := fmt.Sprintf("%s/loki/api/v1/query?%s", c.baseURL, queryParams.Encode())

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result wire.LokiQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// listLabels retrieves all labels (private method)
func (c *LokiClient) listLabels(ctx context.Context, start, end time.Time) (*wire.LokiLabelsResponse, error) {
	queryParams := url.Values{}
	if !start.IsZero() {
		queryParams.Set("start", strconv.FormatInt(start.UnixNano(), 10))
	}
	if !end.IsZero() {
		queryParams.Set("end", strconv.FormatInt(end.UnixNano(), 10))
	}

	endpoint := fmt.Sprintf("%s/loki/api/v1/labels", c.baseURL)
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list labels: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result wire.LokiLabelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// getLabelValues retrieves all values for a specific label (private method)
func (c *LokiClient) getLabelValues(ctx context.Context, label string, start, end time.Time) (*wire.LokiLabelValuesResponse, error) {
	queryParams := url.Values{}
	if !start.IsZero() {
		queryParams.Set("start", strconv.FormatInt(start.UnixNano(), 10))
	}
	if !end.IsZero() {
		queryParams.Set("end", strconv.FormatInt(end.UnixNano(), 10))
	}

	endpoint := fmt.Sprintf("%s/loki/api/v1/label/%s/values", c.baseURL, url.PathEscape(label))
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get label values: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result wire.LokiLabelValuesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// getSeries retrieves the list of time series that match a log stream selector (private method)
func (c *LokiClient) getSeries(ctx context.Context, matchers []string, start, end time.Time) (*wire.LokiSeriesResponse, error) {
	queryParams := url.Values{}
	for _, match := range matchers {
		queryParams.Add("match[]", match)
	}
	if !start.IsZero() {
		queryParams.Set("start", strconv.FormatInt(start.UnixNano(), 10))
	}
	if !end.IsZero() {
		queryParams.Set("end", strconv.FormatInt(end.UnixNano(), 10))
	}

	endpoint := fmt.Sprintf("%s/loki/api/v1/series?%s", c.baseURL, queryParams.Encode())

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get series: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result wire.LokiSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// doRequest performs an HTTP request with authentication
func (c *LokiClient) doRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add authentication if configured
	if c.auth != nil {
		switch c.auth.Type {
		case "basic":
			req.SetBasicAuth(c.auth.Username, c.auth.Password)
		case "bearer":
			req.Header.Set("Authorization", "Bearer "+c.auth.Token)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// handleErrorResponse handles error responses from Loki
func (c *LokiClient) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %d: failed to read error response", resp.StatusCode)
	}

	var lokiErr wire.LokiError
	if err := json.Unmarshal(body, &lokiErr); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("loki error (%s): %s", lokiErr.ErrorType, lokiErr.Error)
}
