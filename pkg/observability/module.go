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

	// Create external service clients
	lokiClient := obsIntegrationsLoki.NewLokiClient(&obsIntegrationsLoki.LokiConfig{
		URL:     obsConfig.Loki.URL,
		Timeout: obsConfig.Loki.Timeout,
	})

	prometheusClient := obsIntegrationsPrometheus.NewPrometheusClient(&obsIntegrationsPrometheus.PrometheusConfig{
		URL:     obsConfig.Prometheus.URL,
		Timeout: obsConfig.Prometheus.Timeout,
	})

	tempoClient := obsIntegrationsTempo.NewTempoClient(&obsIntegrationsTempo.TempoConfig{
		URL:     obsConfig.Tempo.URL,
		Timeout: obsConfig.Tempo.Timeout,
	})

	alertManagerClient := obsIntegrationsAlertManager.NewAlertManagerClient(&obsIntegrationsAlertManager.AlertManagerConfig{
		URL:     obsConfig.AlertManager.URL,
		Timeout: obsConfig.AlertManager.Timeout,
	})

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
