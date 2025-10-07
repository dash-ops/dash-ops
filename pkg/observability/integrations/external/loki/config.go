package loki

import (
	"net/http"
	"strings"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// LokiConfig represents configuration for Loki client
type LokiConfig struct {
	URL     string             `json:"url"`
	Timeout string             `json:"timeout"`
	Auth    *models.AuthConfig `json:"auth,omitempty"`
}

// LokiClient handles direct communication with Loki API
type LokiClient struct {
	baseURL    string
	httpClient *http.Client
	auth       *models.AuthConfig
}

// NewLokiClient creates a new Loki client
func NewLokiClient(config *LokiConfig) *LokiClient {
	timeout := 30 * time.Second
	if config.Timeout != "" {
		if d, err := time.ParseDuration(config.Timeout); err == nil {
			timeout = d
		}
	}

	return &LokiClient{
		baseURL: strings.TrimSuffix(config.URL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
		auth: config.Auth,
	}
}
