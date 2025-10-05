package observability

import (
	"fmt"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	obsAdaptersConfig "github.com/dash-ops/dash-ops/pkg/observability/adapters/config"
	"github.com/dash-ops/dash-ops/pkg/observability/handlers"
	obsIntegrationsAlertManager "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/alertmanager"
	obsIntegrationsLoki "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/loki"
	obsIntegrationsPrometheus "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/prometheus"
	obsIntegrationsTempo "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/tempo"
)

// Module represents the observability module with all its components
type Module struct {
	Handler *handlers.HTTPHandler
}

// NewModule creates and initializes a new observability module
func NewModule(fileConfig []byte) (*Module, error) {
	if fileConfig == nil {
		return nil, fmt.Errorf("module config cannot be nil")
	}

	// Parse observability configuration
	configAdapter := obsAdaptersConfig.NewConfigAdapter()
	obsConfig, err := configAdapter.ParseObservabilityConfigFromFileConfig(fileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse observability configuration: %w", err)
	}

	// Check if observability is enabled
	if !obsConfig.Enabled {
		return nil, fmt.Errorf("observability module is disabled in configuration")
	}

	// Create Loki client from first enabled provider
	var lokiClient *obsIntegrationsLoki.LokiClient
	for _, provider := range obsConfig.Logs.Providers {
		if provider.Type == "loki" && provider.Enabled {
			lokiClient = obsIntegrationsLoki.NewLokiClient(&obsIntegrationsLoki.LokiConfig{
				URL:     provider.URL,
				Timeout: provider.Timeout,
				Auth:    &provider.Auth,
			})
			break
		}
	}

	// Create Tempo client from first enabled provider
	var tempoClient *obsIntegrationsTempo.TempoClient
	for _, provider := range obsConfig.Traces.Providers {
		if provider.Type == "tempo" && provider.Enabled {
			tempoClient = obsIntegrationsTempo.NewTempoClient(&obsIntegrationsTempo.TempoConfig{
				URL:     provider.URL,
				Timeout: provider.Timeout,
				Auth:    &provider.Auth,
			})
			break
		}
	}

	// Create Prometheus client from first enabled provider
	var prometheusClient *obsIntegrationsPrometheus.PrometheusClient
	for _, provider := range obsConfig.Metrics.Providers {
		if provider.Type == "prometheus" && provider.Enabled {
			prometheusClient = obsIntegrationsPrometheus.NewPrometheusClient(&obsIntegrationsPrometheus.PrometheusConfig{
				URL:     provider.URL,
				Timeout: provider.Timeout,
				Auth:    &provider.Auth,
			})
			break
		}
	}

	// Create AlertManager client from first enabled provider (if alerts are enabled)
	var alertManagerClient *obsIntegrationsAlertManager.AlertManagerClient
	if obsConfig.Alerts.Enabled {
		for _, provider := range obsConfig.Alerts.Providers {
			if provider.Type == "alertmanager" && provider.Enabled {
				alertManagerClient = obsIntegrationsAlertManager.NewAlertManagerClient(&obsIntegrationsAlertManager.AlertManagerConfig{
					URL:     provider.URL,
					Timeout: provider.Timeout,
					Auth:    &provider.Auth,
				})
				break
			}
		}
	}

	// Initialize adapters
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize handler with DI
	handler := handlers.NewHTTPHandler(
		lokiClient,
		prometheusClient,
		tempoClient,
		alertManagerClient,
		responseAdapter,
		requestAdapter,
	)

	return &Module{
		Handler: handler,
	}, nil
}

// RegisterRoutes registers HTTP routes for the observability module
func (m *Module) RegisterRoutes(router *mux.Router) {
	if m.Handler == nil {
		return
	}
	m.Handler.RegisterRoutes(router)
}
