package observability

import (
	"fmt"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	"github.com/dash-ops/dash-ops/pkg/observability/controllers"
	"github.com/dash-ops/dash-ops/pkg/observability/handlers"
	"github.com/dash-ops/dash-ops/pkg/observability/logic"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
)

// Module represents the observability module with all its components
type Module struct {
	// Core components
	LogsController    *controllers.LogsController
	MetricsController *controllers.MetricsController
	TracesController  *controllers.TracesController
	AlertsController  *controllers.AlertsController
	HealthController  *controllers.HealthController
	ConfigController  *controllers.ConfigController
	Handler           *handlers.HTTPHandler

	// Logic components
	LogProcessor       *logic.LogProcessor
	MetricProcessor    *logic.MetricProcessor
	TraceProcessor     *logic.TraceProcessor
	AlertProcessor     *logic.AlertProcessor
	DashboardProcessor *logic.DashboardProcessor

	// Adapters
	ResponseAdapter *commonsHttp.ResponseAdapter
	RequestAdapter  *commonsHttp.RequestAdapter

	// Repositories (interfaces - implementations injected)
	LogRepo       ports.LogRepository
	MetricRepo    ports.MetricRepository
	TraceRepo     ports.TraceRepository
	AlertRepo     ports.AlertRepository
	DashboardRepo ports.DashboardRepository
	ServiceRepo   ports.ServiceContextRepository

	// Services (interfaces - implementations injected)
	LogService           ports.LogService
	MetricService        ports.MetricService
	TraceService         ports.TraceService
	AlertService         ports.AlertService
	DashboardService     ports.DashboardService
	NotificationService  ports.NotificationService
	CacheService         ports.CacheService
	ConfigurationService ports.ConfigurationService
}

// ModuleConfig represents configuration for the observability module
type ModuleConfig struct {
	// Repository implementations
	LogRepo       ports.LogRepository
	MetricRepo    ports.MetricRepository
	TraceRepo     ports.TraceRepository
	AlertRepo     ports.AlertRepository
	DashboardRepo ports.DashboardRepository
	ServiceRepo   ports.ServiceContextRepository

	// Service implementations
	LogService           ports.LogService
	MetricService        ports.MetricService
	TraceService         ports.TraceService
	AlertService         ports.AlertService
	DashboardService     ports.DashboardService
	NotificationService  ports.NotificationService
	CacheService         ports.CacheService
	ConfigurationService ports.ConfigurationService
}

// NewModule creates and initializes a new observability module
func NewModule(config *ModuleConfig) (*Module, error) {
	if config == nil {
		return nil, fmt.Errorf("module config cannot be nil")
	}

	// Validate required dependencies
	if config.LogRepo == nil {
		return nil, fmt.Errorf("log repository is required")
	}
	if config.MetricRepo == nil {
		return nil, fmt.Errorf("metric repository is required")
	}
	if config.TraceRepo == nil {
		return nil, fmt.Errorf("trace repository is required")
	}
	if config.AlertRepo == nil {
		return nil, fmt.Errorf("alert repository is required")
	}
	if config.DashboardRepo == nil {
		return nil, fmt.Errorf("dashboard repository is required")
	}
	if config.ServiceRepo == nil {
		return nil, fmt.Errorf("service context repository is required")
	}

	// Initialize logic components
	logProcessor := logic.NewLogProcessor()
	metricProcessor := logic.NewMetricProcessor()
	traceProcessor := logic.NewTraceProcessor()
	alertProcessor := logic.NewAlertProcessor()
	dashboardProcessor := logic.NewDashboardProcessor()

	// Initialize adapters
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize per-domain controllers
	logsController := controllers.NewLogsController(
		config.LogRepo,
		config.ServiceRepo,
		config.LogService,
		config.CacheService,
		logProcessor,
	)

	metricsController := controllers.NewMetricsController(
		config.MetricRepo,
		config.ServiceRepo,
		config.MetricService,
		config.CacheService,
		metricProcessor,
	)

	tracesController := controllers.NewTracesController(
		config.TraceRepo,
		config.ServiceRepo,
		config.TraceService,
		config.CacheService,
		traceProcessor,
	)

	alertsController := controllers.NewAlertsController(
		config.AlertRepo,
		config.AlertService,
		config.CacheService,
		alertProcessor,
	)

	healthController := controllers.NewHealthController(
		config.CacheService,
	)

	configController := controllers.NewConfigController(
		config.ConfigurationService,
		config.NotificationService,
	)

	// Initialize handler
	handler := handlers.NewHTTPHandler(
		logsController,
		metricsController,
		tracesController,
		alertsController,
		healthController,
		configController,
		responseAdapter,
		requestAdapter,
	)

	return &Module{
		LogsController:       logsController,
		MetricsController:    metricsController,
		TracesController:     tracesController,
		AlertsController:     alertsController,
		HealthController:     healthController,
		ConfigController:     configController,
		Handler:              handler,
		LogProcessor:         logProcessor,
		MetricProcessor:      metricProcessor,
		TraceProcessor:       traceProcessor,
		AlertProcessor:       alertProcessor,
		DashboardProcessor:   dashboardProcessor,
		ResponseAdapter:      responseAdapter,
		RequestAdapter:       requestAdapter,
		LogRepo:              config.LogRepo,
		MetricRepo:           config.MetricRepo,
		TraceRepo:            config.TraceRepo,
		AlertRepo:            config.AlertRepo,
		DashboardRepo:        config.DashboardRepo,
		ServiceRepo:          config.ServiceRepo,
		LogService:           config.LogService,
		MetricService:        config.MetricService,
		TraceService:         config.TraceService,
		AlertService:         config.AlertService,
		DashboardService:     config.DashboardService,
		NotificationService:  config.NotificationService,
		CacheService:         config.CacheService,
		ConfigurationService: config.ConfigurationService,
	}, nil
}

// RegisterRoutes registers HTTP routes for the observability module
func (m *Module) RegisterRoutes(router *mux.Router) {
	if m.Handler == nil {
		return
	}
	m.Handler.RegisterRoutes(router)
}
