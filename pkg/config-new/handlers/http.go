package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons-new/adapters/http"
	httpAdapter "github.com/dash-ops/dash-ops/pkg/config-new/adapters/http"
	config "github.com/dash-ops/dash-ops/pkg/config-new/controllers"
)

// HTTPHandler handles HTTP requests for config module
type HTTPHandler struct {
	controller      *config.ConfigController
	configAdapter   *httpAdapter.ConfigAdapter
	responseAdapter *commonsHttp.ResponseAdapter
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(
	controller *config.ConfigController,
	configAdapter *httpAdapter.ConfigAdapter,
	responseAdapter *commonsHttp.ResponseAdapter,
) *HTTPHandler {
	return &HTTPHandler{
		controller:      controller,
		configAdapter:   configAdapter,
		responseAdapter: responseAdapter,
	}
}

// RegisterRoutes registers all config routes
func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/config", h.getConfigHandler).Methods("GET")
	router.HandleFunc("/config/plugins", h.getPluginsHandler).Methods("GET")
	router.HandleFunc("/config/plugins/{name}", h.getPluginStatusHandler).Methods("GET")
	router.HandleFunc("/config/reload", h.reloadConfigHandler).Methods("POST")
	router.HandleFunc("/config/validate", h.validateConfigHandler).Methods("GET")
	router.HandleFunc("/system/info", h.getSystemInfoHandler).Methods("GET")
}

// getConfigHandler handles GET /config
func (h *HTTPHandler) getConfigHandler(w http.ResponseWriter, r *http.Request) {
	config, err := h.controller.GetConfig(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := h.configAdapter.ModelToConfigResponse(config)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getPluginsHandler handles GET /config/plugins
func (h *HTTPHandler) getPluginsHandler(w http.ResponseWriter, r *http.Request) {
	plugins, err := h.controller.GetPlugins(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := h.configAdapter.ModelToPluginsResponse(plugins)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// getPluginStatusHandler handles GET /config/plugins/{name}
func (h *HTTPHandler) getPluginStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pluginName := vars["name"]

	if pluginName == "" {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, "plugin name is required")
		return
	}

	config, err := h.controller.GetConfig(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := h.configAdapter.ModelToPluginStatusResponse(pluginName, config)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// reloadConfigHandler handles POST /config/reload
func (h *HTTPHandler) reloadConfigHandler(w http.ResponseWriter, r *http.Request) {
	config, err := h.controller.ReloadConfig(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := h.configAdapter.ModelToConfigResponse(config)
	h.responseAdapter.WriteJSON(w, http.StatusOK, response)
}

// validateConfigHandler handles GET /config/validate
func (h *HTTPHandler) validateConfigHandler(w http.ResponseWriter, r *http.Request) {
	err := h.controller.ValidateConfig(r.Context())
	if err != nil {
		h.responseAdapter.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.responseAdapter.WriteSuccess(w, "Configuration is valid", nil)
}

// getSystemInfoHandler handles GET /system/info
func (h *HTTPHandler) getSystemInfoHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Get actual system information
	version := "1.0.0"
	environment := "development"
	uptime := "24h"

	systemInfo, err := h.controller.GetSystemInfo(r.Context(), version, environment, uptime)
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
