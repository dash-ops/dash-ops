package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	obsAdapters "github.com/dash-ops/dash-ops/pkg/observability/adapters"
	"github.com/dash-ops/dash-ops/pkg/observability/controllers"
	obsIntegrationsAlertManager "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/alertmanager"
	obsIntegrationsPrometheus "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/prometheus"
	obsIntegrationsTempo "github.com/dash-ops/dash-ops/pkg/observability/integrations/external/tempo"
	"github.com/dash-ops/dash-ops/pkg/observability/logic"
	"github.com/dash-ops/dash-ops/pkg/observability/ports"
	"github.com/dash-ops/dash-ops/pkg/observability/repositories"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// HTTPHandler handles HTTP requests for the observability module
type HTTPHandler struct {
	// Controllers
	logsController           *controllers.LogsController
	metricsController        *controllers.MetricsController
	tracesController         *controllers.TracesController
	alertsController         *controllers.AlertsController
	serviceContextController *controllers.ServiceContextController
	explorerController       *controllers.ExplorerController

	// Service context for integration
	serviceContextRepository ports.ServiceContextRepository

	// Client maps for listing available providers
	logsClients   map[string]ports.LogsClient
	tracesClients map[string]ports.TracesClient

	// Adapters (wire <-> models)
	logsAdapter           *obsAdapters.LogsAdapter
	tracesAdapter         *obsAdapters.TracesAdapter
	serviceContextAdapter *obsAdapters.ServiceContextAdapter
	explorerAdapter       *obsAdapters.ExplorerAdapter
	responseAdapter       *commonsHttp.ResponseAdapter
	requestAdapter        *commonsHttp.RequestAdapter
}

// NewHTTPHandler creates a new HTTP handler with DI
func NewHTTPHandler(
	logsClients map[string]ports.LogsClient,
	tracesClients map[string]ports.TracesClient,
	prometheusClients map[string]*obsIntegrationsPrometheus.PrometheusClient,
	tempoClients map[string]*obsIntegrationsTempo.TempoClient,
	alertManagerClients map[string]*obsIntegrationsAlertManager.AlertManagerClient,
	responseAdapter *commonsHttp.ResponseAdapter,
	requestAdapter *commonsHttp.RequestAdapter,
) *HTTPHandler {
	// Initialize logic processors
	logProcessor := logic.NewLogProcessor()
	metricProcessor := logic.NewMetricProcessor()
	traceProcessor := logic.NewTraceProcessor()
	alertProcessor := logic.NewAlertProcessor()

	// Initialize data transformation adapters
	logsAdapter := obsAdapters.NewLogsAdapter()
	tracesAdapter := obsAdapters.NewTracesAdapter()
	serviceContextAdapter := obsAdapters.NewServiceContextAdapter()
	explorerAdapter := obsAdapters.NewExplorerAdapter(logsAdapter, tracesAdapter)

	// Initialize repositories with multiple clients
	logsRepo := repositories.NewLogsRepository(logsClients)
	tracesRepo := repositories.NewTracesRepository(tracesClients)

	// Initialize controllers with repositories
	logsController := controllers.NewLogsController(
		logsRepo,
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
		tracesRepo,
		traceProcessor,
	)

	alertsController := controllers.NewAlertsController(
		nil, // alertRepo - will be implemented
		nil, // notificationService
		nil, // cache
		alertProcessor,
	)

	serviceContextController := controllers.NewServiceContextController(nil) // Will be set via LoadDependencies
	explorerController := controllers.NewExplorerController(logsClients, tracesClients)

	return &HTTPHandler{
		logsController:           logsController,
		metricsController:        metricsController,
		tracesController:         tracesController,
		alertsController:         alertsController,
		serviceContextController: serviceContextController,
		explorerController:       explorerController,
		serviceContextRepository: nil, // Will be set via LoadDependencies
		logsClients:              logsClients,
		tracesClients:            tracesClients,
		logsAdapter:              logsAdapter,
		tracesAdapter:            tracesAdapter,
		serviceContextAdapter:    serviceContextAdapter,
		explorerAdapter:          explorerAdapter,
		responseAdapter:          responseAdapter,
		requestAdapter:           requestAdapter,
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
	observabilityRouter.HandleFunc("/traces/services", h.handleGetTraceServices).Methods("GET")
	observabilityRouter.HandleFunc("/traces/{traceId}", h.handleGetTraceDetail).Methods("GET")
	observabilityRouter.HandleFunc("/traces", h.handleGetTraces).Methods("GET")

	// Alerts endpoints
	observabilityRouter.HandleFunc("/alerts", h.handleGetAlerts).Methods("GET")
	observabilityRouter.HandleFunc("/alerts/silences", h.handleGetSilences).Methods("GET")

	// Service context endpoints
	observabilityRouter.HandleFunc("/services", h.handleGetServices).Methods("GET")

	// Explorer endpoint
	observabilityRouter.HandleFunc("/explorer", h.handleExplorerQuery).Methods("GET")

	// Providers endpoint
	observabilityRouter.HandleFunc("/providers", h.handleGetProviders).Methods("GET")
}

// handleGetLogs handles GET /observability/logs
func (h *HTTPHandler) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters to wire.LogsRequest
	queryParams := r.URL.Query()

	wireReq := &wire.LogsRequest{
		Provider: queryParams.Get("provider"),
		Service:  queryParams.Get("service"),
		Level:    queryParams.Get("level"),
		Query:    queryParams.Get("query"),
		Sort:     queryParams.Get("sort"),
		Order:    queryParams.Get("order"),
	}

	// Parse time range
	if startStr := queryParams.Get("start"); startStr != "" {
		if startUnix, err := strconv.ParseInt(startStr, 10, 64); err == nil {
			wireReq.StartTime = time.Unix(startUnix, 0)
		}
	}
	if endStr := queryParams.Get("end"); endStr != "" {
		if endUnix, err := strconv.ParseInt(endStr, 10, 64); err == nil {
			wireReq.EndTime = time.Unix(endUnix, 0)
		}
	}

	// Set defaults
	if wireReq.StartTime.IsZero() {
		wireReq.StartTime = time.Now().Add(-1 * time.Hour)
	}
	if wireReq.EndTime.IsZero() {
		wireReq.EndTime = time.Now()
	}

	// Parse limit
	if limitStr := queryParams.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			wireReq.Limit = limit
		}
	}
	if wireReq.Limit == 0 {
		wireReq.Limit = 100
	}

	// Parse offset
	if offsetStr := queryParams.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			wireReq.Offset = offset
		}
	}

	// Step 1: Transform wire -> models using adapter
	modelQuery := h.logsAdapter.WireRequestToModel(wireReq)

	// Step 2: Call controller with models and provider
	logs, err := h.logsController.QueryLogs(r.Context(), wireReq.Provider, modelQuery)
	if err != nil {
		// Step 3a: Transform error to wire response
		wireResp := h.logsAdapter.ErrorToWireResponse(err)
		h.responseAdapter.WriteJSON(w, http.StatusInternalServerError, wireResp)
		return
	}

	// Step 3b: Transform models -> wire using adapter with provider info
	wireResp := h.logsAdapter.ModelToWireResponse(logs, len(logs), len(logs) >= wireReq.Limit)

	// Add provider information to response
	if wireResp.Data.Provider.Name == "" {
		wireResp.Data.Provider = wire.ProviderInfo{
			Name:        wireReq.Provider,
			Type:        "loki", // TODO: Make this dynamic based on provider
			Description: "Log aggregation system",
		}
	}

	// Step 4: Return wire response
	h.responseAdapter.WriteJSON(w, http.StatusOK, wireResp)
}

// handleGetLogLabels handles GET /observability/logs/labels
func (h *HTTPHandler) handleGetLogLabels(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "provider parameter is required")
		return
	}

	labels, err := h.logsController.GetLogLabels(r.Context(), provider)
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
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "provider parameter is required")
		return
	}

	levels, err := h.logsController.GetLogLevels(r.Context(), provider)
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
	// TODO: Implement metrics query via controller
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Metrics endpoint - to be implemented",
		"data":    []interface{}{},
	})
}

// handleGetMetricNames handles GET /observability/metrics/names
func (h *HTTPHandler) handleGetMetricNames(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement via controller
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Metric names endpoint - to be implemented",
		"data":    []string{},
	})
}

// handleGetTraces handles GET /observability/traces
func (h *HTTPHandler) handleGetTraces(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters to wire.TracesRequest
	queryParams := r.URL.Query()

	wireReq := &wire.TracesRequest{
		Service:     queryParams.Get("service"),
		Operation:   queryParams.Get("operation"),
		TraceID:     queryParams.Get("trace_id"),
		Status:      queryParams.Get("status"),
		MinDuration: queryParams.Get("min_duration"),
		MaxDuration: queryParams.Get("max_duration"),
		Sort:        queryParams.Get("sort"),
		Order:       queryParams.Get("order"),
	}

	// Parse time range
	if startStr := queryParams.Get("start"); startStr != "" {
		if startUnix, err := strconv.ParseInt(startStr, 10, 64); err == nil {
			wireReq.StartTime = time.Unix(startUnix, 0)
		}
	}
	if endStr := queryParams.Get("end"); endStr != "" {
		if endUnix, err := strconv.ParseInt(endStr, 10, 64); err == nil {
			wireReq.EndTime = time.Unix(endUnix, 0)
		}
	}

	// Set defaults
	if wireReq.StartTime.IsZero() {
		wireReq.StartTime = time.Now().Add(-1 * time.Hour)
	}
	if wireReq.EndTime.IsZero() {
		wireReq.EndTime = time.Now()
	}

	// Parse limit
	if limitStr := queryParams.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			wireReq.Limit = limit
		}
	}
	if wireReq.Limit == 0 {
		wireReq.Limit = 20
	}

	// Extract provider parameter
	provider := queryParams.Get("provider")
	if provider == "" {
		h.responseAdapter.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "provider parameter is required",
		})
		return
	}

	// Step 1: Transform wire -> models using adapter
	query := h.tracesAdapter.WireRequestToModel(wireReq)

	// Step 2: Call controller with models
	traces, err := h.tracesController.QueryTraces(r.Context(), provider, query)
	if err != nil {
		response := h.tracesAdapter.ErrorToWireResponse(err)
		h.responseAdapter.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Step 3: Transform models -> wire using adapter
	response := h.tracesAdapter.ModelToWireResponse(traces, provider, "tempo")

	// Step 4: Return wire.TracesResponse as JSON
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// handleGetTraceDetail handles GET /observability/traces/{traceId}
func (h *HTTPHandler) handleGetTraceDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceID := vars["traceId"]

	// Extract provider parameter
	queryParams := r.URL.Query()
	provider := queryParams.Get("provider")
	if provider == "" {
		h.responseAdapter.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "provider parameter is required",
		})
		return
	}

	// Call controller to get trace detail
	trace, err := h.tracesController.GetTraceDetail(r.Context(), provider, traceID)
	if err != nil {
		response := h.tracesAdapter.ErrorToWireDetailResponse(err)
		h.responseAdapter.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Transform models -> wire using adapter
	response := h.tracesAdapter.ModelToWireDetailResponse(trace)

	// Return wire.TraceDetailResponse as JSON
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// handleGetTraceServices handles GET /observability/traces/services
func (h *HTTPHandler) handleGetTraceServices(w http.ResponseWriter, r *http.Request) {
	// Extract provider parameter
	queryParams := r.URL.Query()
	provider := queryParams.Get("provider")
	if provider == "" {
		h.responseAdapter.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "provider parameter is required",
		})
		return
	}

	// Call controller to get services
	services, err := h.tracesController.GetServices(r.Context(), provider)
	if err != nil {
		h.responseAdapter.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Return services list
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    services,
	})
}

// handleGetAlerts handles GET /observability/alerts
func (h *HTTPHandler) handleGetAlerts(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement alerts query via controller
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Alerts endpoint - to be implemented",
		"data":    []interface{}{},
	})
}

// handleGetSilences handles GET /observability/alerts/silences
func (h *HTTPHandler) handleGetSilences(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement silences retrieval via controller
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Silences endpoint - to be implemented",
		"data":    []interface{}{},
	})
}

// SetServiceContextRepository sets the service context repository for service-catalog integration
func (h *HTTPHandler) SetServiceContextRepository(repo ports.ServiceContextRepository) {
	h.serviceContextRepository = repo
	// Update service context controller with repository
	if h.serviceContextController != nil {
		h.serviceContextController = controllers.NewServiceContextController(repo)
	}
}

// handleGetServices handles GET /observability/services
func (h *HTTPHandler) handleGetServices(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters to wire.ServicesWithContextRequest
	queryParams := r.URL.Query()

	wireReq := &wire.ServicesWithContextRequest{
		Search: queryParams.Get("search"),
	}

	// Parse limit
	if limitStr := queryParams.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			wireReq.Limit = limit
		}
	}
	if wireReq.Limit == 0 {
		wireReq.Limit = 50 // default
	}

	// Parse offset
	if offsetStr := queryParams.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			wireReq.Offset = offset
		}
	}

	// Step 1: Transform wire -> model parameters using adapter
	search, limit, offset := h.serviceContextAdapter.WireRequestToModel(wireReq)

	// Step 2: Call controller with parameters
	services, total, err := h.serviceContextController.GetServicesWithContext(r.Context(), search, limit, offset)
	if err != nil {
		// Step 3a: Transform error to wire response
		wireResp := h.serviceContextAdapter.ErrorToWireResponse(err)
		h.responseAdapter.WriteJSON(w, http.StatusInternalServerError, wireResp)
		return
	}

	// Step 3b: Transform models -> wire using adapter
	hasMore := offset+limit < total
	wireResp := h.serviceContextAdapter.ModelToWireResponse(services, total, hasMore)

	// Step 4: Return wire response
	h.responseAdapter.WriteJSON(w, http.StatusOK, wireResp)
}

// handleExplorerQuery handles GET /observability/explorer
func (h *HTTPHandler) handleExplorerQuery(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters to wire.ExplorerQueryRequest
	queryParams := r.URL.Query()

	wireReq := &wire.ExplorerQueryRequest{
		Query:         queryParams.Get("query"),
		TimeRangeFrom: queryParams.Get("time_range_from"),
		TimeRangeTo:   queryParams.Get("time_range_to"),
		Provider:      queryParams.Get("provider"),
	}

	// Validate required parameters
	if wireReq.Query == "" {
		wireResp := h.explorerAdapter.ErrorToWireResponse(fmt.Errorf("query parameter is required"))
		h.responseAdapter.WriteJSON(w, http.StatusBadRequest, wireResp)
		return
	}

	// Step 1: Transform wire -> model parameters using adapter
	query, timeRangeFrom, timeRangeTo, provider := h.explorerAdapter.WireRequestToModel(wireReq)

	// Step 2: Call controller
	dataSource, results, total, executionTimeMs, err := h.explorerController.ExecuteQuery(
		r.Context(),
		query,
		timeRangeFrom,
		timeRangeTo,
		provider,
	)
	if err != nil {
		// Step 3a: Transform error to wire response
		wireResp := h.explorerAdapter.ErrorToWireResponse(err)
		h.responseAdapter.WriteJSON(w, http.StatusInternalServerError, wireResp)
		return
	}

	// Step 3b: Transform models -> wire using adapter
	wireResp := h.explorerAdapter.ModelToWireResponse(dataSource, results, total, query, executionTimeMs)

	// Step 4: Return wire response
	h.responseAdapter.WriteJSON(w, http.StatusOK, wireResp)
}

// handleGetProviders handles GET /observability/providers
func (h *HTTPHandler) handleGetProviders(w http.ResponseWriter, r *http.Request) {
	// Extract provider names from client maps
	logsProviders := make([]string, 0, len(h.logsClients))
	for name := range h.logsClients {
		logsProviders = append(logsProviders, name)
	}

	tracesProviders := make([]string, 0, len(h.tracesClients))
	for name := range h.tracesClients {
		tracesProviders = append(tracesProviders, name)
	}

	// TODO: Add metricsProviders when metrics clients are implemented
	metricsProviders := []string{}

	wireResp := &wire.ProvidersResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.ProvidersResponseData{
			LogsProviders:    logsProviders,
			TracesProviders:  tracesProviders,
			MetricsProviders: metricsProviders,
		},
	}

	h.responseAdapter.WriteJSON(w, http.StatusOK, wireResp)
}
