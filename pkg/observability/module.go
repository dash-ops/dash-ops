package observability

import (
	"fmt"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	obsAdaptersConfig "github.com/dash-ops/dash-ops/pkg/observability/adapters/config"
	"github.com/dash-ops/dash-ops/pkg/observability/controllers"
	"github.com/dash-ops/dash-ops/pkg/observability/handlers"
	obsIntegrationsLoki "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/loki"
	"github.com/dash-ops/dash-ops/pkg/observability/logic"
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

	// TODO: Create Tempo client when implementing traces
	// var tempoClient *obsIntegrationsTempo.TempoClient
	// for _, provider := range obsConfig.Traces.Providers {
	// 	if provider.Type == "tempo" && provider.Enabled {
	// 		tempoClient = obsIntegrationsTempo.NewTempoClient(&obsIntegrationsTempo.TempoConfig{
	// 			URL:     provider.URL,
	// 			Timeout: provider.Timeout,
	// 			Auth:    &provider.Auth,
	// 		})
	// 		break
	// 	}
	// }

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

	// Initialize logic processors
	logProcessor := logic.NewLogProcessor()
	metricProcessor := logic.NewMetricProcessor()
	traceProcessor := logic.NewTraceProcessor()
	alertProcessor := logic.NewAlertProcessor()

	// Initialize adapters (Loki adapter implements LogRepository interface)
	var lokiAdapter ports.LogRepository
	if lokiClient != nil {
		lokiAdapter = obsIntegrationsLoki.NewLokiAdapter(lokiClient)
	}

	// Initialize HTTP adapters
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize controllers
	logsController := controllers.NewLogsController(
		lokiAdapter,
		nil, // serviceRepo - will be wired later
		nil, // logService - will be wired later
		nil, // cache - will be wired later
		logProcessor,
	)

	metricsController := controllers.NewMetricsController(
		nil, // metricRepo - will be implemented
		nil, // serviceRepo
		nil, // metricService
		nil, // cache
		metricProcessor,
	)

	tracesController := controllers.NewTracesController(
		nil, // traceRepo - will be implemented
		nil, // serviceRepo
		nil, // traceService
		nil, // cache
		traceProcessor,
	)

	alertsController := controllers.NewAlertsController(
		nil, // alertRepo - will be implemented
		nil, // notificationService
		nil, // cache
		alertProcessor,
	)

	// Initialize handler with controllers
	handler := handlers.NewHTTPHandler(
		logsController,
		metricsController,
		tracesController,
		alertsController,
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
