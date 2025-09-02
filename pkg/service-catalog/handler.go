package servicecatalog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/dash-ops/dash-ops/pkg/commons"
)

// Handler represents the HTTP handler for service catalog operations
type Handler struct {
	serviceCatalog *ServiceCatalog
}

// NewHandler creates a new service catalog HTTP handler
func NewHandler(serviceCatalog *ServiceCatalog) *Handler {
	return &Handler{
		serviceCatalog: serviceCatalog,
	}
}

// RegisterRoutes registers all service catalog routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Service CRUD operations
	router.HandleFunc("/services", h.listServicesHandler).Methods("GET")
	router.HandleFunc("/services", h.createServiceHandler).Methods("POST")
	router.HandleFunc("/services/{name}", h.getServiceHandler).Methods("GET")
	router.HandleFunc("/services/{name}", h.updateServiceHandler).Methods("PUT")
	router.HandleFunc("/services/{name}", h.deleteServiceHandler).Methods("DELETE")

	// Service filtering and search
	router.HandleFunc("/services/search", h.searchServicesHandler).Methods("GET")
	router.HandleFunc("/services/by-team/{team}", h.listServicesByTeamHandler).Methods("GET")
	router.HandleFunc("/services/by-tier/{tier}", h.listServicesByTierHandler).Methods("GET")

	// Service health and monitoring
	router.HandleFunc("/services/{name}/health", h.getServiceHealthHandler).Methods("GET")
	router.HandleFunc("/services/{name}/history", h.getServiceHistoryHandler).Methods("GET")

	// Repository and system information
	router.HandleFunc("/system/history", h.getAllHistoryHandler).Methods("GET")
	router.HandleFunc("/system/status", h.getSystemStatusHandler).Methods("GET")
}

// listServicesHandler handles GET /services
func (h *Handler) listServicesHandler(w http.ResponseWriter, r *http.Request) {
	serviceList, err := h.serviceCatalog.ListServices()
	if err != nil {
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list services: %v", err))
		return
	}

	commons.RespondJSON(w, http.StatusOK, serviceList)
}

// createServiceHandler handles POST /services
func (h *Handler) createServiceHandler(w http.ResponseWriter, r *http.Request) {
	var service Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		commons.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
		return
	}

	// Get user context from OAuth2 (if available)
	user := h.getUserContext(r)

	// Validate service
	if err := h.serviceCatalog.ValidateService(&service); err != nil {
		commons.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Service validation failed: %v", err))
		return
	}

	// Create service
	if err := h.serviceCatalog.CreateService(&service, user); err != nil {
		if err.Error() == fmt.Sprintf("service '%s' already exists", service.Metadata.Name) {
			commons.RespondError(w, http.StatusConflict, err.Error())
			return
		}
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create service: %v", err))
		return
	}

	// Return created service
	createdService, err := h.serviceCatalog.GetService(service.Metadata.Name)
	if err != nil {
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve created service: %v", err))
		return
	}

	commons.RespondJSON(w, http.StatusCreated, createdService)
}

// getServiceHandler handles GET /services/{name}
func (h *Handler) getServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	service, err := h.serviceCatalog.GetService(name)
	if err != nil {
		if err.Error() == fmt.Sprintf("service '%s' not found", name) {
			commons.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get service: %v", err))
		return
	}

	commons.RespondJSON(w, http.StatusOK, service)
}

// updateServiceHandler handles PUT /services/{name}
func (h *Handler) updateServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	var service Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		commons.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
		return
	}

	// Ensure the service name in URL matches the name in body
	if service.Metadata.Name != name {
		commons.RespondError(w, http.StatusBadRequest, "Service name in URL and body must match")
		return
	}

	// Get user context from OAuth2 (if available)
	user := h.getUserContext(r)

	// Validate service
	if err := h.serviceCatalog.ValidateService(&service); err != nil {
		commons.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Service validation failed: %v", err))
		return
	}

	// Update service
	if err := h.serviceCatalog.UpdateService(&service, user); err != nil {
		if err.Error() == fmt.Sprintf("service '%s' not found", name) {
			commons.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update service: %v", err))
		return
	}

	// Return updated service
	updatedService, err := h.serviceCatalog.GetService(name)
	if err != nil {
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve updated service: %v", err))
		return
	}

	commons.RespondJSON(w, http.StatusOK, updatedService)
}

// deleteServiceHandler handles DELETE /services/{name}
func (h *Handler) deleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Get user context from OAuth2 (if available)
	user := h.getUserContext(r)

	// Delete service
	if err := h.serviceCatalog.DeleteService(name, user); err != nil {
		if err.Error() == fmt.Sprintf("service '%s' not found", name) {
			commons.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete service: %v", err))
		return
	}

	commons.RespondJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Service '%s' deleted successfully", name),
	})
}

// searchServicesHandler handles GET /services/search?q=query
func (h *Handler) searchServicesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		commons.RespondError(w, http.StatusBadRequest, "Query parameter 'q' is required")
		return
	}

	serviceList, err := h.serviceCatalog.SearchServices(query)
	if err != nil {
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to search services: %v", err))
		return
	}

	commons.RespondJSON(w, http.StatusOK, serviceList)
}

// listServicesByTeamHandler handles GET /services/by-team/{team}
func (h *Handler) listServicesByTeamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	team := vars["team"]

	serviceList, err := h.serviceCatalog.ListServicesByTeam(team)
	if err != nil {
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list services by team: %v", err))
		return
	}

	commons.RespondJSON(w, http.StatusOK, serviceList)
}

// listServicesByTierHandler handles GET /services/by-tier/{tier}
func (h *Handler) listServicesByTierHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tier := vars["tier"]

	serviceList, err := h.serviceCatalog.ListServicesByTier(tier)
	if err != nil {
		commons.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Failed to list services by tier: %v", err))
		return
	}

	commons.RespondJSON(w, http.StatusOK, serviceList)
}

// getServiceHealthHandler handles GET /services/{name}/health
func (h *Handler) getServiceHealthHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Get authorization header to pass to Kubernetes API
	authHeader := r.Header.Get("Authorization")

	health, err := h.serviceCatalog.GetServiceHealth(name, authHeader)
	if err != nil {
		if err.Error() == fmt.Sprintf("service '%s' not found", name) {
			commons.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get service health: %v", err))
		return
	}

	commons.RespondJSON(w, http.StatusOK, health)
}

// getServiceHistoryHandler handles GET /services/{name}/history
func (h *Handler) getServiceHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	history, err := h.serviceCatalog.GetServiceHistory(name)
	if err != nil {
		if err.Error() == fmt.Sprintf("service '%s' not found", name) {
			commons.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "git versioning is not available" {
			commons.RespondError(w, http.StatusServiceUnavailable, "Git versioning is not available")
			return
		}
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get service history: %v", err))
		return
	}

	commons.RespondJSON(w, http.StatusOK, history)
}

// getAllHistoryHandler handles GET /system/history
func (h *Handler) getAllHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Parse optional limit parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // Default limit

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	allHistory, err := h.serviceCatalog.GetAllHistory()
	if err != nil {
		if err.Error() == "git versioning is not available" {
			commons.RespondError(w, http.StatusServiceUnavailable, "Git versioning is not available")
			return
		}
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get history: %v", err))
		return
	}

	// Apply limit
	if len(allHistory) > limit {
		allHistory = allHistory[:limit]
	}

	commons.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"history": allHistory,
		"total":   len(allHistory),
		"limit":   limit,
	})
}

// getSystemStatusHandler handles GET /system/status
func (h *Handler) getSystemStatusHandler(w http.ResponseWriter, r *http.Request) {
	repoStatus, err := h.serviceCatalog.GetRepositoryStatus()
	if err != nil {
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get repository status: %v", err))
		return
	}

	// Get service count
	serviceList, err := h.serviceCatalog.ListServices()
	if err != nil {
		commons.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get service count: %v", err))
		return
	}

	status := map[string]interface{}{
		"service_count":      serviceList.Total,
		"repository_status":  repoStatus,
		"storage_provider":   h.serviceCatalog.config.Storage.Provider,
		"versioning_enabled": h.serviceCatalog.versioning != nil && h.serviceCatalog.versioning.IsEnabled(),
		"versioning_provider": func() string {
			if h.serviceCatalog.config.Versioning.Enabled {
				return h.serviceCatalog.config.Versioning.Provider
			}
			return "none"
		}(),
	}

	commons.RespondJSON(w, http.StatusOK, status)
}

// getUserContext extracts user context from request (OAuth2 integration)
func (h *Handler) getUserContext(r *http.Request) *UserContext {
	// TODO: Implement OAuth2 integration to get user context
	// For now, return a default user context

	// Check if user data is available in context (from OAuth2 middleware)
	if userData := r.Context().Value(commons.UserDataKey); userData != nil {
		if user, ok := userData.(commons.UserData); ok {
			return &UserContext{
				Username: user.Org,                 // Use Org as username for now
				Name:     user.Org,                 // Use Org as name for now
				Email:    user.Org + "@github.com", // Generate email from org
				Teams:    user.Groups,              // GitHub teams will be in groups
			}
		}
	}

	// Default context for development/testing
	return &UserContext{
		Username: "anonymous",
		Name:     "Anonymous User",
		Email:    "anonymous@dash-ops.local",
		Teams:    []string{},
	}
}

// Middleware functions for future enhancements

// requireServiceOwnership middleware to check if user owns the service
func (h *Handler) requireServiceOwnership(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement team-based authorization
		// For now, allow all operations
		next(w, r)
	}
}

// validateServiceTier middleware to validate tier parameter
func (h *Handler) validateServiceTier(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if tier, exists := vars["tier"]; exists {
			validTiers := map[string]bool{
				"TIER-1": true,
				"TIER-2": true,
				"TIER-3": true,
			}
			if !validTiers[tier] {
				commons.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid tier '%s', must be TIER-1, TIER-2, or TIER-3", tier))
				return
			}
		}
		next(w, r)
	}
}
