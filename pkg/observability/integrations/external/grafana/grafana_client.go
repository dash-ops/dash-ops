package grafana

import (
	"context"
	"net/http"
	"time"
)

type Config struct {
	URL      string
	APIToken string
	Timeout  int
}

type Client struct {
	config     *Config
	httpClient *http.Client
}

func NewClient(config *Config) *Client {
	timeout := time.Duration(config.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		config:     config,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) ListDashboards(ctx context.Context) ([]byte, error) {
	// TODO: chamada real ao Grafana
	return []byte(`[]`), nil
}

func (c *Client) GetDashboard(ctx context.Context, uid string) ([]byte, error) {
	// TODO: chamada real ao Grafana
	return []byte(`{}`), nil
}
