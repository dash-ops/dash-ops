package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	spaAdapters "github.com/dash-ops/dash-ops/pkg/spa/adapters/http"
	spaLogic "github.com/dash-ops/dash-ops/pkg/spa/logic"
	spaModels "github.com/dash-ops/dash-ops/pkg/spa/models"
)

func TestNewHTTPHandler_CreatesHandler(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{}
	apiRouter := mux.NewRouter()

	spaAdapter := spaAdapters.NewSPAAdapter(config, fileProcessor, stats, apiRouter)
	handler := NewHTTPHandler(spaAdapter)

	assert.NotNil(t, handler)
	assert.Equal(t, spaAdapter, handler.spaAdapter)
}

func TestHTTPHandler_RegisterRoutes_RegistersSPARoutes(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{}
	apiRouter := mux.NewRouter()

	spaAdapter := spaAdapters.NewSPAAdapter(config, fileProcessor, stats, apiRouter)
	handler := NewHTTPHandler(spaAdapter)

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Test that routes are registered by making a request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should not panic and should handle the request
	assert.NotPanics(t, func() {
		router.ServeHTTP(w, req)
	})
}

func TestHTTPHandler_RegisterRoutes_WithAPIRoutes_HandlesCorrectly(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{}
	apiRouter := mux.NewRouter()

	// Add a test API route
	apiRouter.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API response"))
	}).Methods("GET")

	spaAdapter := spaAdapters.NewSPAAdapter(config, fileProcessor, stats, apiRouter)
	handler := NewHTTPHandler(spaAdapter)

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Test API route
	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should handle API route
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "API response", w.Body.String())
}

func TestHTTPHandler_RegisterRoutes_WithNonAPIRoutes_HandlesCorrectly(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{}
	apiRouter := mux.NewRouter()

	spaAdapter := spaAdapters.NewSPAAdapter(config, fileProcessor, stats, apiRouter)
	handler := NewHTTPHandler(spaAdapter)

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Test non-API route (should be handled by SPA)
	req := httptest.NewRequest("GET", "/some-page", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should not panic (SPA adapter will handle it)
	assert.NotPanics(t, func() {
		router.ServeHTTP(w, req)
	})
}
