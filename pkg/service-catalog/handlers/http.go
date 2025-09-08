package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	scAdapters "github.com/dash-ops/dash-ops/pkg/service-catalog/adapters/http"
	servicecatalog "github.com/dash-ops/dash-ops/pkg/service-catalog/controllers"
	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scWire "github.com/dash-ops/dash-ops/pkg/service-catalog/wire"
)

// HTTPHandler handles HTTP requests for service catalog module
type HTTPHandler struct {
	controller      *servicecatalog.ServiceController
	serviceAdapter  *scAdapters.ServiceAdapter
	responseAdapter *commonsHttp.ResponseAdapter
	requestAdapter  *commonsHttp.RequestAdapter
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(
	controller *servicecatalog.ServiceController,
	serviceAdapter *scAdapters.ServiceAdapter,
	responseAdapter *commonsHttp.ResponseAdapter,
	requestAdapter *commonsHttp.RequestAdapter,
) *HTTPHandler {
	return &HTTPHandler{
		controller:      controller,
		serviceAdapter:  serviceAdapter,
		responseAdapter: responseAdapter,
		requestAdapter:  requestAdapter,
	}
}

// RegisterRoutes registers all service catalog routes
func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	// Service CRUD operations
	router.HandleFunc("/services", h.listServicesHandler).Methods("GET")
	router.HandleFunc("/services", h.createServiceHandler).Methods("POST")
	router.HandleFunc("/services/{name}", h.getServiceHandler).Methods("GET")
	router.HandleFunc("/services/{name}", h.updateServiceHandler).Methods("PUT")
	router.HandleFunc("/services/{name}", h.deleteServiceHandler).Methods("DELETE")

	// Service filtering and search (TODO: Implement missing handlers)
	// router.HandleFunc("/services/search", h.searchServicesHandler).Methods("GET")
	// router.HandleFunc("/services/by-team/{team}", h.listServicesByTeamHandler).Methods("GET")
	// router.HandleFunc("/services/by-tier/{tier}", h.listServicesByTierHandler).Methods("GET")

	// Service health and monitoring
	router.HandleFunc("/services/{name}/health", h.getServiceHealthHandler).Methods("GET")
	// router.HandleFunc("/services/health/batch", h.getBatchHealthHandler).Methods("POST")
	router.HandleFunc("/services/{name}/history", h.getServiceHistoryHandler).Methods("GET")

	// System information (TODO: Implement missing handlers)
	// router.HandleFunc("/system/history", h.getAllHistoryHandler).Methods("GET")
	// router.HandleFunc("/system/status", h.getSystemStatusHandler).Methods("GET")

	// Context resolution (TODO: Implement missing handlers)
	// router.HandleFunc("/context/deployment/{deployment}", h.resolveDeploymentServiceHandler).Methods("GET")
}

// createServiceHandler handles POST /services
func (h *HTTPHandler) createServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req scWire.CreateServiceRequest
	if err := h.requestAdapter.ParseJSON(r, &req); err != nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	// Get user context
	user, err := h.getUserContext(r)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusUnauthorized, "Authentication required: "+err.Error())
		return
	}

	// Transform to domain model
	service, err := h.serviceAdapter.RequestToModel(req)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Invalid service data: "+err.Error())
		return
	}

	// Call controller
	createdService, err := h.controller.CreateService(r.Context(), service, user)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to create service: "+err.Error())
		return
	}

	// Transform and respond
	response := h.serviceAdapter.ModelToResponse(createdService)
	h.responseAdapter.WriteCreated(w, "/services/"+createdService.Metadata.Name, response)
}

// getServiceHandler handles GET /services/{name}
func (h *HTTPHandler) getServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Service name is required")
		return
	}

	// Call controller
	service, err := h.controller.GetService(r.Context(), name)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusNotFound, "Service not found: "+err.Error())
		return
	}

	// Transform and respond
	response := h.serviceAdapter.ModelToResponse(service)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// updateServiceHandler handles PUT /services/{name}
func (h *HTTPHandler) updateServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Service name is required")
		return
	}

	// Parse request
	var req scWire.UpdateServiceRequest
	if err := h.requestAdapter.ParseJSON(r, &req); err != nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	// Get user context
	user, err := h.getUserContext(r)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusUnauthorized, "Authentication required: "+err.Error())
		return
	}

	// Get existing service
	existingService, err := h.controller.GetService(r.Context(), name)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusNotFound, "Service not found: "+err.Error())
		return
	}

	// Transform to domain model
	service, err := h.serviceAdapter.UpdateRequestToModel(req, existingService)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Invalid service data: "+err.Error())
		return
	}

	// Call controller
	updatedService, err := h.controller.UpdateService(r.Context(), service, user)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to update service: "+err.Error())
		return
	}

	// Transform and respond
	response := h.serviceAdapter.ModelToResponse(updatedService)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// deleteServiceHandler handles DELETE /services/{name}
func (h *HTTPHandler) deleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Service name is required")
		return
	}

	// Get user context
	user, err := h.getUserContext(r)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusUnauthorized, "Authentication required: "+err.Error())
		return
	}

	// Call controller
	err = h.controller.DeleteService(r.Context(), name, user)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to delete service: "+err.Error())
		return
	}

	h.responseAdapter.WriteNoContent(w)
}

// listServicesHandler handles GET /services
func (h *HTTPHandler) listServicesHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filter := h.parseServiceFilter(r)

	// Call controller
	serviceList, err := h.controller.ListServices(r.Context(), filter)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to list services: "+err.Error())
		return
	}

	// Transform and respond
	response := h.serviceAdapter.ModelListToResponse(serviceList)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getServiceHealthHandler handles GET /services/{name}/health
func (h *HTTPHandler) getServiceHealthHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Service name is required")
		return
	}

	// Call controller
	health, err := h.controller.GetServiceHealth(r.Context(), name)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to get service health: "+err.Error())
		return
	}

	// Transform and respond
	response := h.serviceAdapter.HealthModelToResponse(health)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getServiceHistoryHandler handles GET /services/{name}/history
func (h *HTTPHandler) getServiceHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Service name is required")
		return
	}

	// Call controller
	history, err := h.controller.GetServiceHistory(r.Context(), name)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, "Failed to get service history: "+err.Error())
		return
	}

	// Transform and respond
	response := h.serviceAdapter.HistoryModelToResponse(history)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// parseServiceFilter parses query parameters into ServiceFilter
func (h *HTTPHandler) parseServiceFilter(r *http.Request) *scModels.ServiceFilter {
	query := r.URL.Query()

	filter := &scModels.ServiceFilter{
		Team:   query.Get("team"),
		Search: query.Get("search"),
	}

	if tier := query.Get("tier"); tier != "" {
		filter.Tier = scModels.ServiceTier(tier)
	}

	if status := query.Get("status"); status != "" {
		filter.Status = scModels.ServiceStatus(status)
	}

	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	if offset := query.Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	return filter
}

// getUserContext extracts user context from request
func (h *HTTPHandler) getUserContext(r *http.Request) (*scModels.UserContext, error) {
	// TODO: Implement OAuth2 integration to get user context
	// For now, return a default user context for testing
	return &scModels.UserContext{
		Username: "test-user",
		Name:     "Test User",
		Email:    "test@example.com",
		Teams:    []string{"test-team"},
	}, nil
}

// Additional handlers would be implemented here...
// searchServicesHandler, listServicesByTeamHandler, etc.
