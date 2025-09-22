package handlers

import (
	"github.com/gorilla/mux"

	spaAdapters "github.com/dash-ops/dash-ops/pkg/spa/adapters/http"
)

// HTTPHandler handles HTTP requests for the SPA module
type HTTPHandler struct {
	spaAdapter *spaAdapters.SPAAdapter
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(spaAdapter *spaAdapters.SPAAdapter) *HTTPHandler {
	return &HTTPHandler{
		spaAdapter: spaAdapter,
	}
}

// RegisterRoutes registers HTTP routes for the SPA module
func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	// SPA routes with API middleware - serves static files and handles SPA routing
	router.PathPrefix("/").Handler(h.spaAdapter.CreateAPIMiddleware()(h.spaAdapter))
}
