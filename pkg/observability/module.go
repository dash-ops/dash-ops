package observability

import (
	"fmt"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	obsAdaptersConfig "github.com/dash-ops/dash-ops/pkg/observability/adapters/config"
	"github.com/dash-ops/dash-ops/pkg/observability/handlers"
	obsIntegrationsLoki "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/loki"
	obsIntegrationsTempo "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/tempo"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
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

	// Create multiple log providers
	logsClients := make(map[string]ports.LogsClient)
	for _, provider := range obsConfig.Logs.Providers {
		if provider.Type == "loki" && provider.Enabled {
			lokiClient := obsIntegrationsLoki.NewLokiClient(&obsIntegrationsLoki.LokiConfig{
				URL:     provider.URL,
				Timeout: provider.Timeout,
				Auth:    &provider.Auth,
			})
			// LokiClient implements ports.LogsClient interface
			logsClients[provider.Name] = lokiClient
		}
	}

	// Create multiple trace providers
	tracesClients := make(map[string]ports.TracesClient)
	for _, provider := range obsConfig.Traces.Providers {
		if provider.Type == "tempo" && provider.Enabled {
			tempoClient := obsIntegrationsTempo.NewTempoClient(&obsIntegrationsTempo.TempoConfig{
				URL:     provider.URL,
				Timeout: provider.Timeout,
				Auth:    &provider.Auth,
			})
			// TempoClient implements ports.TracesClient interface
			tracesClients[provider.Name] = tempoClient
		}
	}

	// TODO: Create Prometheus client when implementing metrics
	// var prometheusClient *obsIntegrationsPrometheus.PrometheusClient
	// for _, provider := range obsConfig.Metrics.Providers {
	// 	if provider.Type == "prometheus" && provider.Enabled {
	// 		prometheusClient = obsIntegrationsPrometheus.NewPrometheusClient(&obsIntegrationsPrometheus.PrometheusConfig{
	// 			URL:     provider.URL,
	// 			Timeout: provider.Timeout,
	// 			Auth:    &provider.Auth,
	// 		})
	// 		break
	// 	}
	// }

	// TODO: Create AlertManager client when implementing alerts
	// var alertManagerClient *obsIntegrationsAlertManager.AlertManagerClient
	// if obsConfig.Alerts.Enabled {
	// 	for _, provider := range obsConfig.Alerts.Providers {
	// 		if provider.Type == "alertmanager" && provider.Enabled {
	// 			alertManagerClient = obsIntegrationsAlertManager.NewAlertManagerClient(&obsIntegrationsAlertManager.AlertManagerConfig{
	// 				URL:     provider.URL,
	// 				Timeout: provider.Timeout,
	// 				Auth:    &provider.Auth,
	// 			})
	// 			break
	// 		}
	// 	}
	// }

	// Initialize HTTP adapters
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize handler with clients (handler creates controllers internally)
	handler := handlers.NewHTTPHandler(
		logsClients,
		tracesClients,
		nil, // prometheusClients - will be implemented
		nil, // tempoClients (deprecated parameter) - will be removed
		nil, // alertManagerClients - will be implemented
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
