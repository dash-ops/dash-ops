package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	settingsAdaptersHttp "github.com/dash-ops/dash-ops/pkg/settings/adapters/http"
	settingsControllers "github.com/dash-ops/dash-ops/pkg/settings/controllers"
	settingsPorts "github.com/dash-ops/dash-ops/pkg/settings/ports"
	settingsWire "github.com/dash-ops/dash-ops/pkg/settings/wire"
)

// HTTPHandler handles HTTP requests for settings module
type HTTPHandler struct {
	setupController    *settingsControllers.SetupController
	settingsController *settingsControllers.SettingsController
	configController   *settingsControllers.ConfigController
	statusRepo         settingsPorts.ConfigStatusRepository
	settingsAdapter    *settingsAdaptersHttp.SettingsAdapter
	configAdapter      *settingsAdaptersHttp.ConfigAdapter
	responseAdapter    *commonsHttp.ResponseAdapter
	requestAdapter     *commonsHttp.RequestAdapter
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(
	setupController *settingsControllers.SetupController,
	settingsController *settingsControllers.SettingsController,
	configController *settingsControllers.ConfigController,
	statusRepo settingsPorts.ConfigStatusRepository,
	settingsAdapter *settingsAdaptersHttp.SettingsAdapter,
	configAdapter *settingsAdaptersHttp.ConfigAdapter,
	responseAdapter *commonsHttp.ResponseAdapter,
	requestAdapter *commonsHttp.RequestAdapter,
) *HTTPHandler {
	return &HTTPHandler{
		setupController:    setupController,
		settingsController: settingsController,
		configController:   configController,
		statusRepo:         statusRepo,
		settingsAdapter:    settingsAdapter,
		configAdapter:      configAdapter,
		responseAdapter:    responseAdapter,
		requestAdapter:     requestAdapter,
	}
}

// RegisterRoutes registers all settings routes
func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	// Setup routes (no auth required, only when no plugins configured)
	setupRouter := router.PathPrefix("/settings/setup").Subrouter()
	setupRouter.HandleFunc("/status", h.handleGetSetupStatus).Methods("GET")
	setupRouter.HandleFunc("/configure", h.handleConfigureSetup).Methods("POST")

	// Settings routes (auth conditional)
	settingsRouter := router.PathPrefix("/settings").Subrouter()
	settingsRouter.HandleFunc("/config", h.handleGetSettingsConfig).Methods("GET")
	settingsRouter.HandleFunc("/config", h.handleUpdateSettingsConfig).Methods("PUT")

	// Config routes (legacy compatibility)
	router.HandleFunc("/config", h.handleGetConfig).Methods("GET")
	router.HandleFunc("/config/plugins", h.handleGetPlugins).Methods("GET")
	router.HandleFunc("/config/plugins/{name}", h.handleGetPluginStatus).Methods("GET")
	router.HandleFunc("/config/reload", h.handleReloadConfig).Methods("POST")
	router.HandleFunc("/config/validate", h.handleValidateConfig).Methods("GET")
	router.HandleFunc("/system/info", h.handleGetSystemInfo).Methods("GET")
}

// handleGetSetupStatus handles GET /settings/setup/status
func (h *HTTPHandler) handleGetSetupStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.setupController.GetSetupStatus(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := h.settingsAdapter.SetupStatusToResponse(status)
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// handleConfigureSetup handles POST /settings/setup/configure
func (h *HTTPHandler) handleConfigureSetup(w http.ResponseWriter, r *http.Request) {
	// Check if plugins are already configured
	hasPlugins, err := h.statusRepo.HasPluginsConfigured(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if hasPlugins {
		h.responseAdapter.WriteJSON(w, http.StatusForbidden, map[string]interface{}{
			"success": false,
			"error":   "plugins are already configured. Use settings endpoint to modify configuration",
		})
		return
	}

	// Parse request
	var req settingsWire.SetupConfigureRequest
	if err := h.requestAdapter.ParseJSON(r, &req); err != nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	// Convert to model
	setupConfig := h.settingsAdapter.RequestToSetupConfig(&req)
	if setupConfig == nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Invalid setup configuration")
		return
	}

	// Configure setup
	configPath, err := h.setupController.ConfigureSetup(r.Context(), setupConfig)
	if err != nil {
		if err.Error() == "plugins are already configured. Use settings endpoint to modify configuration" {
			h.responseAdapter.WriteJSON(w, http.StatusForbidden, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := &settingsWire.SetupConfigureResponse{
		Success:    true,
		Message:    "Configuration saved successfully",
		ConfigPath: configPath,
	}

	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// handleGetSettingsConfig handles GET /settings/config
func (h *HTTPHandler) handleGetSettingsConfig(w http.ResponseWriter, r *http.Request) {
	// TODO: Add auth check if OAuth2 is configured
	// For now, allow access without auth

	settings, err := h.settingsController.GetSettings(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := h.settingsAdapter.SettingsConfigToResponse(settings)
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// handleUpdateSettingsConfig handles PUT /settings/config
func (h *HTTPHandler) handleUpdateSettingsConfig(w http.ResponseWriter, r *http.Request) {
	// TODO: Add auth check if OAuth2 is configured
	// For now, allow access without auth

	// Parse request
	var req settingsWire.UpdateSettingsRequest
	if err := h.requestAdapter.ParseJSON(r, &req); err != nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	// Convert to model
	updateRequest := h.settingsAdapter.RequestToUpdateSettingsRequest(&req)
	if updateRequest == nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "Invalid update request")
		return
	}

	// Update settings
	response, err := h.settingsController.UpdateSettings(r.Context(), updateRequest)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	wireResponse := h.settingsAdapter.UpdateSettingsResponseToWire(response)
	h.responseAdapter.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    wireResponse,
	})
}

func (h *HTTPHandler) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := h.configController.GetConfig(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := h.configAdapter.ModelToConfigResponse(config)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

func (h *HTTPHandler) handleGetPlugins(w http.ResponseWriter, r *http.Request) {
	plugins, err := h.configController.GetPlugins(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := h.configAdapter.ModelToPluginsArray(plugins)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

func (h *HTTPHandler) handleGetPluginStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pluginName := vars["name"]

	if pluginName == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "plugin name is required")
		return
	}

	config, err := h.configController.GetConfig(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := h.configAdapter.ModelToPluginStatusResponse(pluginName, config)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

func (h *HTTPHandler) handleReloadConfig(w http.ResponseWriter, r *http.Request) {
	config, err := h.configController.ReloadConfig(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := h.configAdapter.ModelToConfigResponse(config)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

func (h *HTTPHandler) handleValidateConfig(w http.ResponseWriter, r *http.Request) {
	err := h.configController.ValidateConfig(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.responseAdapter.WriteSuccess(w, "Configuration is valid", nil)
}

func (h *HTTPHandler) handleGetSystemInfo(w http.ResponseWriter, r *http.Request) {
	version := "1.0.0"
	environment := "development"
	uptime := "24h"

	systemInfo, err := h.configController.GetSystemInfo(r.Context(), version, environment, uptime)
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := h.configAdapter.ModelToSystemInfoResponse(
		systemInfo.Config,
		systemInfo.Version,
		systemInfo.Environment,
		systemInfo.Uptime,
	)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}
