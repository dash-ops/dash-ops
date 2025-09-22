package spa

import (
	"fmt"
	"net/http"
	"time"

	spaLogic "github.com/dash-ops/dash-ops/pkg/spa/logic"
	spaModels "github.com/dash-ops/dash-ops/pkg/spa/models"
)

// Module represents the SPA module
type Module struct {
	Config        *spaModels.SPAConfig
	FileProcessor *spaLogic.FileProcessor
	Handler       http.Handler
	Stats         *spaModels.SPAStats
}

// NewModule creates and initializes a new SPA module
func NewModule(config *spaModels.SPAConfig) (*Module, error) {
	if config == nil {
		return nil, fmt.Errorf("SPA config cannot be nil")
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid SPA config: %w", err)
	}

	// Initialize file processor
	fileProcessor := spaLogic.NewFileProcessor()

	// Validate static path and index file
	if err := fileProcessor.ValidateStaticPath(config.StaticPath); err != nil {
		return nil, fmt.Errorf("static path validation failed: %w", err)
	}

	if err := fileProcessor.ValidateIndexFile(config.StaticPath, config.IndexPath); err != nil {
		return nil, fmt.Errorf("index file validation failed: %w", err)
	}

	// Initialize stats
	stats := &spaModels.SPAStats{
		StartTime: time.Now(),
	}

	// Create HTTP handler
	handler := &SPAHandler{
		config:        config,
		fileProcessor: fileProcessor,
		stats:         stats,
	}

	return &Module{
		Config:        config,
		FileProcessor: fileProcessor,
		Handler:       handler,
		Stats:         stats,
	}, nil
}

// SPAHandler implements http.Handler for serving SPA files
type SPAHandler struct {
	config        *spaModels.SPAConfig
	fileProcessor *spaLogic.FileProcessor
	stats         *spaModels.SPAStats
}

// ServeHTTP implements http.Handler interface
func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Update stats
	h.stats.TotalRequests++
	h.stats.LastRequest = startTime
	h.stats.UpdateUptime()

	// Normalize path
	requestPath := h.fileProcessor.NormalizePath(r.URL.Path)

	// Resolve file path
	filePath, isIndexFallback, err := h.fileProcessor.ResolveFilePath(requestPath, h.config.StaticPath, h.config.IndexPath)
	if err != nil {
		h.handleError(w, r, http.StatusBadRequest, err.Error(), startTime)
		return
	}

	// Get file info
	fileInfo, err := h.fileProcessor.GetFileInfo(filePath)
	if err != nil {
		h.handleError(w, r, http.StatusNotFound, "File not found", startTime)
		return
	}

	// Set headers
	h.setHeaders(w, fileInfo, isIndexFallback)

	// Serve file
	http.ServeFile(w, r, filePath)

	// Update success stats
	h.stats.SuccessfulRequests++
	h.stats.BytesServed += fileInfo.Size

	// Calculate average response time
	duration := time.Since(startTime)
	if h.stats.TotalRequests > 0 {
		h.stats.AverageResponseTime = time.Duration(
			(int64(h.stats.AverageResponseTime) + int64(duration)) / 2,
		)
	} else {
		h.stats.AverageResponseTime = duration
	}
}

// setHeaders sets appropriate headers for the response
func (h *SPAHandler) setHeaders(w http.ResponseWriter, fileInfo *spaModels.FileInfo, isIndexFallback bool) {
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
		w.Header().Set("Cache-Control", h.config.GetCacheControl())
	}

	// Set ETag for caching
	if fileInfo.ETag != "" {
		w.Header().Set("ETag", fileInfo.ETag)
	}

	// Set security headers
	securityHeaders := h.config.GetSecurityHeaders()
	for key, value := range securityHeaders {
		w.Header().Set(key, value)
	}

	// Set CORS headers
	corsHeaders := h.config.GetCORSHeaders()
	for key, value := range corsHeaders {
		w.Header().Set(key, value)
	}

	// Set custom headers
	for key, value := range h.config.CustomHeaders {
		w.Header().Set(key, value)
	}
}

// handleError handles error responses and updates stats
func (h *SPAHandler) handleError(w http.ResponseWriter, r *http.Request, statusCode int, message string, startTime time.Time) {
	h.stats.ErrorRequests++

	// Calculate duration
	duration := time.Since(startTime)

	// Log request info (in production, this would go to a proper logger)
	_ = &spaModels.RequestInfo{
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
