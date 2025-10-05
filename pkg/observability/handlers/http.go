package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	"github.com/dash-ops/dash-ops/pkg/observability/controllers"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// HTTPHandler handles HTTP requests for the observability module
type HTTPHandler struct {
	// Controllers
	logsController    *controllers.LogsController
	metricsController *controllers.MetricsController
	tracesController  *controllers.TracesController
	alertsController  *controllers.AlertsController

	// Adapters
	responseAdapter *commonsHttp.ResponseAdapter
	requestAdapter  *commonsHttp.RequestAdapter
}

// NewHTTPHandler creates a new HTTP handler with DI
func NewHTTPHandler(
	logsController *controllers.LogsController,
	metricsController *controllers.MetricsController,
	tracesController *controllers.TracesController,
	alertsController *controllers.AlertsController,
	responseAdapter *commonsHttp.ResponseAdapter,
	requestAdapter *commonsHttp.RequestAdapter,
) *HTTPHandler {
	return &HTTPHandler{
		logsController:    logsController,
		metricsController: metricsController,
		tracesController:  tracesController,
		alertsController:  alertsController,
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
	// Parse query parameters
	query := r.URL.Query()

	// Build request from query parameters
	req := &wire.LogsRequest{
		Service: query.Get("service"),
		Level:   query.Get("level"),
		Query:   query.Get("query"),
		Sort:    query.Get("sort"),
		Order:   query.Get("order"),
	}

	// Parse time range
	if startStr := query.Get("start"); startStr != "" {
		if startUnix, err := strconv.ParseInt(startStr, 10, 64); err == nil {
			req.StartTime = time.Unix(startUnix, 0)
		}
	}
	if endStr := query.Get("end"); endStr != "" {
		if endUnix, err := strconv.ParseInt(endStr, 10, 64); err == nil {
			req.EndTime = time.Unix(endUnix, 0)
		}
	}

	// Set defaults
	if req.StartTime.IsZero() {
		req.StartTime = time.Now().Add(-1 * time.Hour)
	}
	if req.EndTime.IsZero() {
		req.EndTime = time.Now()
	}

	// Parse limit
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}
	if req.Limit == 0 {
		req.Limit = 100
	}

	// Parse offset
	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}

	// Call controller
	response, err := h.logsController.GetLogs(r.Context(), req)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to query logs: "+err.Error())
		return
	}

	// Return response
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// handleGetLogLabels handles GET /observability/logs/labels
func (h *HTTPHandler) handleGetLogLabels(w http.ResponseWriter, r *http.Request) {
	labels, err := h.logsController.GetLogLabels(r.Context())
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
	levels, err := h.logsController.GetLogLevels(r.Context())
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
	// TODO: Implement traces query via controller
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Traces endpoint - to be implemented",
		"data":    []interface{}{},
	})
}

// handleGetTraceDetail handles GET /observability/traces/{traceId}
func (h *HTTPHandler) handleGetTraceDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceID := vars["traceId"]

	// TODO: Implement trace detail retrieval via controller
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Trace detail endpoint - to be implemented",
		"data": map[string]interface{}{
			"traceId": traceID,
		},
	})
}

// handleGetTraceServices handles GET /observability/traces/services
func (h *HTTPHandler) handleGetTraceServices(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement via controller
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Trace services endpoint - to be implemented",
		"data":    []string{},
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
