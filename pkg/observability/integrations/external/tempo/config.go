package tempo

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// TempoConfig holds the configuration for Tempo client
type TempoConfig struct {
	URL     string
	Timeout string
	Auth    *models.AuthConfig
}

// TempoClient represents a client for interacting with Tempo
type TempoClient struct {
	baseURL    string
	httpClient *http.Client
	auth       *models.AuthConfig
}

// NewTempoClient creates a new Tempo client
func NewTempoClient(config *TempoConfig) *TempoClient {
	// Parse timeout
	timeout := 30 * time.Second
	if config.Timeout != "" {
		if t, err := time.ParseDuration(config.Timeout); err == nil {
			timeout = t
		}
	}

	return &TempoClient{
		baseURL: config.URL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		auth: config.Auth,
	}
}

// Validate validates the Tempo configuration
func (c *TempoConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("tempo URL is required")
	}
	return nil
}
