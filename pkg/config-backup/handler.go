package config

import (
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	"github.com/gorilla/mux"
)

func configPluginsHandler(plugins []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commons.RespondJSON(w, http.StatusOK, plugins)
	}
}

// MakeConfigHandlers Add config module endpoints
func MakeConfigHandlers(api *mux.Router, dashConfig DashYaml) {

	api.HandleFunc("/config/plugins", configPluginsHandler(dashConfig.Plugins)).
		Methods("GET", "OPTIONS").
		Name("configPlugins")
}
