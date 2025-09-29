package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	obsIntegrationsAlertManager "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/alertmanager"
	obsIntegrationsLoki "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/loki"
	obsIntegrationsPrometheus "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/prometheus"
	obsIntegrationsTempo "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/tempo"
	obsRepositories "github.com/dash-ops/dash-ops/pkg/observability/repositories"
)

// HTTPHandler handles HTTP requests for the observability module
type HTTPHandler struct {
	// Repositories
	logsRepository    *obsRepositories.LogsRepository
	metricsRepository *obsRepositories.MetricsRepository
	tracesRepository  *obsRepositories.TracesRepository
	alertsRepository  *obsRepositories.AlertsRepository

	// Adapters
	responseAdapter *commonsHttp.ResponseAdapter
	requestAdapter  *commonsHttp.RequestAdapter
}

// NewHTTPHandler creates a new HTTP handler with DI
func NewHTTPHandler(
	lokiClient *obsIntegrationsLoki.LokiClient,
	prometheusClient *obsIntegrationsPrometheus.PrometheusClient,
	tempoClient *obsIntegrationsTempo.TempoClient,
	alertManagerClient *obsIntegrationsAlertManager.AlertManagerClient,
	responseAdapter *commonsHttp.ResponseAdapter,
	requestAdapter *commonsHttp.RequestAdapter,
) *HTTPHandler {
	// Create repositories with clients
	logsRepository := obsRepositories.NewLogsRepository(lokiClient)
	metricsRepository := obsRepositories.NewMetricsRepository(prometheusClient)
	tracesRepository := obsRepositories.NewTracesRepository(tempoClient)
	alertsRepository := obsRepositories.NewAlertsRepository(alertManagerClient)

	return &HTTPHandler{
		logsRepository:    logsRepository,
		metricsRepository: metricsRepository,
		tracesRepository:  tracesRepository,
		alertsRepository:  alertsRepository,
		responseAdapter:   responseAdapter,
		requestAdapter:    requestAdapter,
	}
}

// RegisterRoutes registers HTTP routes for the observability module
func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	if router == nil {
		return
	}

	// Create observability subrouter
	observabilityRouter := router.PathPrefix("/observability").Subrouter()

	// Logs endpoints
	observabilityRouter.HandleFunc("/logs", h.handleGetLogs).Methods("GET")
	observabilityRouter.HandleFunc("/logs/labels", h.handleGetLogLabels).Methods("GET")
	observabilityRouter.HandleFunc("/logs/levels", h.handleGetLogLevels).Methods("GET")

	// Metrics endpoints
	observabilityRouter.HandleFunc("/metrics", h.handleGetMetrics).Methods("GET")
	observabilityRouter.HandleFunc("/metrics/names", h.handleGetMetricNames).Methods("GET")

	// Traces endpoints
	observabilityRouter.HandleFunc("/traces", h.handleGetTraces).Methods("GET")
	observabilityRouter.HandleFunc("/traces/{traceId}", h.handleGetTraceDetail).Methods("GET")
	observabilityRouter.HandleFunc("/traces/services", h.handleGetTraceServices).Methods("GET")

	// Alerts endpoints
	observabilityRouter.HandleFunc("/alerts", h.handleGetAlerts).Methods("GET")
	observabilityRouter.HandleFunc("/alerts/silences", h.handleGetSilences).Methods("GET")
}

// handleGetLogs handles GET /observability/logs
func (h *HTTPHandler) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logs query
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    []interface{}{},
	})
}

// handleGetLogLabels handles GET /observability/logs/labels
func (h *HTTPHandler) handleGetLogLabels(w http.ResponseWriter, r *http.Request) {
	labels, err := h.logsRepository.GetLogLabels(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to get log labels: "+err.Error())
		return
	}

	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    labels,
	})
}

// handleGetLogLevels handles GET /observability/logs/levels
func (h *HTTPHandler) handleGetLogLevels(w http.ResponseWriter, r *http.Request) {
	levels, err := h.logsRepository.GetLogLevels(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to get log levels: "+err.Error())
		return
	}

	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    levels,
	})
}

// handleGetMetrics handles GET /observability/metrics
func (h *HTTPHandler) handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement metrics query
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    []interface{}{},
	})
}

// handleGetMetricNames handles GET /observability/metrics/names
func (h *HTTPHandler) handleGetMetricNames(w http.ResponseWriter, r *http.Request) {
	names, err := h.metricsRepository.GetMetricNames(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to get metric names: "+err.Error())
		return
	}

	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    names,
	})
}

// handleGetTraces handles GET /observability/traces
func (h *HTTPHandler) handleGetTraces(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement traces query
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    []interface{}{},
	})
}

// handleGetTraceDetail handles GET /observability/traces/{traceId}
func (h *HTTPHandler) handleGetTraceDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceID := vars["traceId"]

	// TODO: Implement trace detail retrieval
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"traceId": traceID,
		},
	})
}

// handleGetTraceServices handles GET /observability/traces/services
func (h *HTTPHandler) handleGetTraceServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.tracesRepository.GetServices(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to get trace services: "+err.Error())
		return
	}

	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    services,
	})
}

// handleGetAlerts handles GET /observability/alerts
func (h *HTTPHandler) handleGetAlerts(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement alerts query
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    []interface{}{},
	})
}

// handleGetSilences handles GET /observability/alerts/silences
func (h *HTTPHandler) handleGetSilences(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement silences retrieval
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    []interface{}{},
	})
}
