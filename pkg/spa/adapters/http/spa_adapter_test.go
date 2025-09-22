package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	spaLogic "github.com/dash-ops/dash-ops/pkg/spa/logic"
	spaModels "github.com/dash-ops/dash-ops/pkg/spa/models"
)

func TestNewSPAAdapter_CreatesAdapter(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{StartTime: time.Now()}
	apiRouter := mux.NewRouter()

	adapter := NewSPAAdapter(config, fileProcessor, stats, apiRouter)

	assert.NotNil(t, adapter)
	assert.Equal(t, config, adapter.config)
	assert.Equal(t, fileProcessor, adapter.fileProcessor)
	assert.Equal(t, stats, adapter.stats)
	assert.Equal(t, apiRouter, adapter.apiRouter)
}

func TestResponseRecorder_WriteHeader_SetsStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := &responseRecorder{ResponseWriter: w, statusCode: 0}

	recorder.WriteHeader(http.StatusNotFound)

	assert.Equal(t, http.StatusNotFound, recorder.statusCode)
	assert.True(t, recorder.headerWritten)
}

func TestResponseRecorder_WriteHeader_MultipleCalls_OnlyFirstOneCounts(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := &responseRecorder{ResponseWriter: w, statusCode: 0}

	recorder.WriteHeader(http.StatusNotFound)
	recorder.WriteHeader(http.StatusOK) // Should be ignored

	assert.Equal(t, http.StatusNotFound, recorder.statusCode)
}

func TestResponseRecorder_Write_SetsStatusCodeToOK(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := &responseRecorder{ResponseWriter: w, statusCode: 0}

	recorder.Write([]byte("test"))

	assert.Equal(t, http.StatusOK, recorder.statusCode)
	assert.True(t, recorder.headerWritten)
}

func TestResponseRecorder_Write_AfterWriteHeader_PreservesStatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	recorder := &responseRecorder{ResponseWriter: w, statusCode: 0}

	recorder.WriteHeader(http.StatusNotFound)
	recorder.Write([]byte("test"))

	assert.Equal(t, http.StatusNotFound, recorder.statusCode)
}

func TestSPAAdapter_CreateAPIMiddleware_ReturnsMiddleware(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{StartTime: time.Now()}
	apiRouter := mux.NewRouter()

	adapter := NewSPAAdapter(config, fileProcessor, stats, apiRouter)
	middleware := adapter.CreateAPIMiddleware()

	assert.NotNil(t, middleware)
}

func TestSPAAdapter_CreateAPIMiddleware_WithAPIRoute_HandlesAPIRequest(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{StartTime: time.Now()}
	apiRouter := mux.NewRouter()

	// Add a test API route
	apiRouter.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API response"))
	}).Methods("GET")

	adapter := NewSPAAdapter(config, fileProcessor, stats, apiRouter)
	middleware := adapter.CreateAPIMiddleware()

	// Create a mock handler for the middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("SPA response"))
	})

	handler := middleware(nextHandler)

	// Test API route
	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "API response", w.Body.String())
}

func TestSPAAdapter_CreateAPIMiddleware_WithNonAPIRoute_CallsNextHandler(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{StartTime: time.Now()}
	apiRouter := mux.NewRouter()

	adapter := NewSPAAdapter(config, fileProcessor, stats, apiRouter)
	middleware := adapter.CreateAPIMiddleware()

	// Create a mock handler for the middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("SPA response"))
	})

	handler := middleware(nextHandler)

	// Test non-API route
	req := httptest.NewRequest("GET", "/some-page", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "SPA response", w.Body.String())
}

func TestSPAAdapter_CreateAPIMiddleware_WithUnhandledAPIRoute_Returns404(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{StartTime: time.Now()}
	apiRouter := mux.NewRouter()

	adapter := NewSPAAdapter(config, fileProcessor, stats, apiRouter)
	middleware := adapter.CreateAPIMiddleware()

	// Create a mock handler for the middleware
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("SPA response"))
	})

	handler := middleware(nextHandler)

	// Test unhandled API route
	req := httptest.NewRequest("GET", "/api/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	// The response should contain "API endpoint not found" or be a 404 page
	body := w.Body.String()
	assert.True(t, body == "API endpoint not found\n" || body == "404 page not found\n")
}

func TestSPAAdapter_ServeHTTP_UpdatesStats(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{StartTime: time.Now()}
	apiRouter := mux.NewRouter()

	adapter := NewSPAAdapter(config, fileProcessor, stats, apiRouter)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// This will fail because the file doesn't exist, but stats should be updated
	adapter.ServeHTTP(w, req)

	assert.Equal(t, int64(1), stats.TotalRequests)
	assert.Equal(t, int64(1), stats.ErrorRequests)
	assert.Equal(t, int64(0), stats.SuccessfulRequests)
}

func TestSPAAdapter_setHeaders_SetsContentType(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{StartTime: time.Now()}
	apiRouter := mux.NewRouter()

	adapter := NewSPAAdapter(config, fileProcessor, stats, apiRouter)

	fileInfo := &spaModels.FileInfo{
		ContentType: "text/html",
		ETag:        "test-etag",
	}

	w := httptest.NewRecorder()
	adapter.setHeaders(w, fileInfo, false)

	assert.Equal(t, "text/html", w.Header().Get("Content-Type"))
	assert.Equal(t, "test-etag", w.Header().Get("ETag"))
}

func TestSPAAdapter_setHeaders_WithHTMLFile_SetsNoCacheHeaders(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{StartTime: time.Now()}
	apiRouter := mux.NewRouter()

	adapter := NewSPAAdapter(config, fileProcessor, stats, apiRouter)

	fileInfo := &spaModels.FileInfo{
		ContentType: "text/html",
	}

	w := httptest.NewRecorder()
	adapter.setHeaders(w, fileInfo, true) // isIndexFallback = true

	assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	assert.Equal(t, "0", w.Header().Get("Expires"))
}

func TestSPAAdapter_setHeaders_WithNonHTMLFile_SetsCacheHeaders(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath:   "front/dist",
		IndexPath:    "index.html",
		CacheControl: "public, max-age=7200",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{StartTime: time.Now()}
	apiRouter := mux.NewRouter()

	adapter := NewSPAAdapter(config, fileProcessor, stats, apiRouter)

	fileInfo := &spaModels.FileInfo{
		ContentType: "text/css",
	}

	w := httptest.NewRecorder()
	adapter.setHeaders(w, fileInfo, false)

	assert.Equal(t, "public, max-age=7200", w.Header().Get("Cache-Control"))
}

func TestSPAAdapter_handleError_UpdatesErrorStats(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}
	fileProcessor := spaLogic.NewFileProcessor()
	stats := &spaModels.SPAStats{StartTime: time.Now()}
	apiRouter := mux.NewRouter()

	adapter := NewSPAAdapter(config, fileProcessor, stats, apiRouter)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	startTime := time.Now()

	adapter.handleError(w, req, http.StatusNotFound, "Not found", startTime)

	assert.Equal(t, int64(1), stats.ErrorRequests)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Not found")
}
