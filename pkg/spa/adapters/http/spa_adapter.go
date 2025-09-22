package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	spaLogic "github.com/dash-ops/dash-ops/pkg/spa/logic"
	spaModels "github.com/dash-ops/dash-ops/pkg/spa/models"
	"github.com/gorilla/mux"
)

// SPAAdapter handles HTTP requests for SPA serving
type SPAAdapter struct {
	config        *spaModels.SPAConfig
	fileProcessor *spaLogic.FileProcessor
	stats         *spaModels.SPAStats
	apiRouter     *mux.Router
}

// NewSPAAdapter creates a new SPA adapter
func NewSPAAdapter(
	config *spaModels.SPAConfig,
	fileProcessor *spaLogic.FileProcessor,
	stats *spaModels.SPAStats,
	apiRouter *mux.Router,
) *SPAAdapter {
	return &SPAAdapter{
		config:        config,
		fileProcessor: fileProcessor,
		stats:         stats,
		apiRouter:     apiRouter,
	}
}

// responseRecorder captures the status code for middleware
type responseRecorder struct {
	http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func (r *responseRecorder) WriteHeader(code int) {
	if !r.headerWritten {
		r.statusCode = code
		r.headerWritten = true
		r.ResponseWriter.WriteHeader(code)
	}
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	// If Write is called before WriteHeader, Go automatically sets status to 200
	if !r.headerWritten {
		r.statusCode = http.StatusOK
		r.headerWritten = true
	}
	return r.ResponseWriter.Write(data)
}

// ServeHTTP implements http.Handler interface
func (a *SPAAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Update stats
	a.stats.TotalRequests++
	a.stats.LastRequest = startTime
	a.stats.UpdateUptime()

	// Normalize path
	requestPath := a.fileProcessor.NormalizePath(r.URL.Path)

	// Resolve file path
	filePath, isIndexFallback, err := a.fileProcessor.ResolveFilePath(requestPath, a.config.StaticPath, a.config.IndexPath)
	if err != nil {
		a.handleError(w, r, http.StatusBadRequest, err.Error(), startTime)
		return
	}

	// Get file info
	fileInfo, err := a.fileProcessor.GetFileInfo(filePath)
	if err != nil {
		a.handleError(w, r, http.StatusNotFound, "File not found", startTime)
		return
	}

	// Set headers
	a.setHeaders(w, fileInfo, isIndexFallback)

	// Serve file
	http.ServeFile(w, r, filePath)

	// Update success stats
	a.stats.SuccessfulRequests++
	a.stats.BytesServed += fileInfo.Size

	// Calculate average response time
	duration := time.Since(startTime)
	if a.stats.TotalRequests > 0 {
		a.stats.AverageResponseTime = time.Duration(
			(int64(a.stats.AverageResponseTime) + int64(duration)) / 2,
		)
	} else {
		a.stats.AverageResponseTime = duration
	}
}

// setHeaders sets appropriate headers for the response
func (a *SPAAdapter) setHeaders(w http.ResponseWriter, fileInfo *spaModels.FileInfo, isIndexFallback bool) {
	// Set content type
	w.Header().Set("Content-Type", fileInfo.ContentType)

	// Set cache control
	if isIndexFallback || fileInfo.IsHTML() {
		// Don't cache HTML files to ensure SPA routing works
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	} else {
		// Cache static assets
		w.Header().Set("Cache-Control", a.config.GetCacheControl())
	}

	// Set ETag for caching
	if fileInfo.ETag != "" {
		w.Header().Set("ETag", fileInfo.ETag)
	}

	// Set security headers
	securityHeaders := a.config.GetSecurityHeaders()
	for key, value := range securityHeaders {
		w.Header().Set(key, value)
	}

	// Set CORS headers
	corsHeaders := a.config.GetCORSHeaders()
	for key, value := range corsHeaders {
		w.Header().Set(key, value)
	}

}

// handleError handles error responses and updates stats
func (a *SPAAdapter) handleError(w http.ResponseWriter, r *http.Request, statusCode int, message string, startTime time.Time) {
	a.stats.ErrorRequests++

	// Calculate duration
	duration := time.Since(startTime)

	// Log request info (in production, this would go to a proper logger)
	_ = struct {
		Method       string
		Path         string
		UserAgent    string
		RemoteAddr   string
		Referer      string
		Timestamp    time.Time
		ResponseCode int
		Duration     time.Duration
	}{
		Method:       r.Method,
		Path:         r.URL.Path,
		UserAgent:    r.UserAgent(),
		RemoteAddr:   r.RemoteAddr,
		Referer:      r.Referer(),
		Timestamp:    startTime,
		ResponseCode: statusCode,
		Duration:     duration,
	}

	http.Error(w, message, statusCode)
}

// CreateAPIMiddleware creates middleware to handle API routes vs SPA routes
func (a *SPAAdapter) CreateAPIMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If it's an API route, let it be handled by the API router
			if strings.HasPrefix(r.URL.Path, "/api/") {
				// Create a response recorder to check if the route was handled
				recorder := &responseRecorder{ResponseWriter: w, statusCode: 0}
				a.apiRouter.ServeHTTP(recorder, r)

				// If no route was matched (status 0), return 404
				if recorder.statusCode == 0 {
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(map[string]string{"error": "API endpoint not found"})
					return
				}
				return
			}
			// Otherwise, serve the SPA
			next.ServeHTTP(w, r)
		})
	}
}
