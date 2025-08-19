package servicecatalog

import (
	"fmt"
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	"github.com/gorilla/mux"
)

// ServiceCatalogService handles business logic
type ServiceCatalogService struct {
	storage Storage
}

// NewServiceCatalogService creates a new service catalog service
func NewServiceCatalogService(storage Storage) *ServiceCatalogService {
	return &ServiceCatalogService{storage: storage}
}

// listServicesHandler handles GET /services
func (s *ServiceCatalogService) listServicesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters for filtering
		filter := ServiceFilter{
			Tier:   r.URL.Query().Get("tier"),
			Team:   r.URL.Query().Get("team"),
			Status: r.URL.Query().Get("status"),
			Search: r.URL.Query().Get("search"),
		}

		services, err := s.storage.ListServices(filter)
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, "Failed to list services: "+err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"services": services,
			"total":    len(services),
		})
	}
}

// getServiceHandler handles GET /services/{id}
func (s *ServiceCatalogService) getServiceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		service, err := s.storage.GetService(id)
		if err != nil {
			commons.RespondError(w, http.StatusNotFound, "Service not found")
			return
		}

		commons.RespondJSON(w, http.StatusOK, service)
	}
}

// createServiceHandler handles POST /services
func (s *ServiceCatalogService) createServiceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateServiceRequest
		if err := commons.DecodeJSON(r.Body, &req); err != nil {
			commons.RespondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}

		// Basic validation
		if req.Name == "" {
			commons.RespondError(w, http.StatusBadRequest, "Service name is required")
			return
		}
		if req.Description == "" {
			commons.RespondError(w, http.StatusBadRequest, "Service description is required")
			return
		}
		if req.Tier == "" {
			commons.RespondError(w, http.StatusBadRequest, "Service tier is required")
			return
		}
		if req.Team == "" {
			commons.RespondError(w, http.StatusBadRequest, "Service team is required")
			return
		}

		service, err := s.storage.CreateService(req)
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, "Failed to create service: "+err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusCreated, service)
	}
}

// updateServiceHandler handles PUT /services/{id}
func (s *ServiceCatalogService) updateServiceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var service Service
		if err := commons.DecodeJSON(r.Body, &service); err != nil {
			commons.RespondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
			return
		}

		service.ID = id // Ensure ID matches URL
		if err := s.storage.UpdateService(id, &service); err != nil {
			commons.RespondError(w, http.StatusInternalServerError, "Failed to update service: "+err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, service)
	}
}

// deleteServiceHandler handles DELETE /services/{id}
func (s *ServiceCatalogService) deleteServiceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		if err := s.storage.DeleteService(id); err != nil {
			commons.RespondError(w, http.StatusInternalServerError, "Failed to delete service: "+err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, map[string]string{"message": "Service deleted successfully"})
	}
}

// statsHandler handles GET /stats
func (s *ServiceCatalogService) statsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get all services
		services, err := s.storage.ListServices(ServiceFilter{})
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, "Failed to get stats: "+err.Error())
			return
		}

		// Calculate stats
		stats := map[string]interface{}{
			"total": len(services),
			"by_tier": map[string]int{
				"tier-1": 0,
				"tier-2": 0,
				"tier-3": 0,
			},
			"by_status": map[string]int{
				"active":     0,
				"inactive":   0,
				"deprecated": 0,
			},
		}

		for _, service := range services {
			// Count by tier
			if tierCount, ok := stats["by_tier"].(map[string]int); ok {
				tierCount[service.Tier]++
			}

			// Count by status
			if statusCount, ok := stats["by_status"].(map[string]int); ok {
				statusCount[service.Status]++
			}
		}

		commons.RespondJSON(w, http.StatusOK, stats)
	}
}

// MakeServiceCatalogHandlers creates and registers all service catalog routes
func MakeServiceCatalogHandlers(r *mux.Router, fileConfig []byte) {
	dashConfig, err := loadConfig(fileConfig)
	if err != nil {
		fmt.Printf("Failed to load service catalog config: %v\n", err)
		return
	}

	if len(dashConfig.ServiceCatalog) == 0 {
		fmt.Println("No service catalog configuration found")
		return
	}

	config := dashConfig.ServiceCatalog[0]

	// Initialize storage
	var storage Storage
	switch config.Storage {
	case "file":
		catalogPath := config.CatalogPath
		if catalogPath == "" {
			catalogPath = "./catalog/services"
		}
		storage = NewFileStorage(catalogPath)
	default:
		fmt.Printf("Unsupported storage type: %s\n", config.Storage)
		return
	}

	// Create service
	service := NewServiceCatalogService(storage)

	// Register routes
	r.HandleFunc("/servicecatalog/services", service.listServicesHandler()).
		Methods("GET", "OPTIONS").
		Name("serviceCatalogList")

	r.HandleFunc("/servicecatalog/services", service.createServiceHandler()).
		Methods("POST", "OPTIONS").
		Name("serviceCatalogCreate")

	r.HandleFunc("/servicecatalog/services/{id}", service.getServiceHandler()).
		Methods("GET", "OPTIONS").
		Name("serviceCatalogGet")

	r.HandleFunc("/servicecatalog/services/{id}", service.updateServiceHandler()).
		Methods("PUT", "OPTIONS").
		Name("serviceCatalogUpdate")

	r.HandleFunc("/servicecatalog/services/{id}", service.deleteServiceHandler()).
		Methods("DELETE", "OPTIONS").
		Name("serviceCatalogDelete")

	r.HandleFunc("/servicecatalog/stats", service.statsHandler()).
		Methods("GET", "OPTIONS").
		Name("serviceCatalogStats")

	fmt.Println("Service Catalog handlers registered")
}
