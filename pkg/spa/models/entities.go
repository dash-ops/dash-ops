package spa

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// SPAConfig represents SPA server configuration
type SPAConfig struct {
	StaticPath   string `yaml:"static_path" json:"static_path"`
	IndexPath    string `yaml:"index_path" json:"index_path"`
	CacheControl string `yaml:"cache_control" json:"cache_control"`
}

// FileInfo represents information about a served file
type FileInfo struct {
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	ETag        string `json:"etag,omitempty"`
}

// SPAStats represents SPA server statistics
type SPAStats struct {
	TotalRequests       int64         `json:"total_requests"`
	SuccessfulRequests  int64         `json:"successful_requests"`
	ErrorRequests       int64         `json:"error_requests"`
	BytesServed         int64         `json:"bytes_served"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	LastRequest         time.Time     `json:"last_request"`
	StartTime           time.Time     `json:"start_time"`
	Uptime              time.Duration `json:"uptime"`
}

// Domain methods for SPAConfig

// Validate validates SPA configuration
func (sc *SPAConfig) Validate() error {
	if sc.StaticPath == "" {
		return fmt.Errorf("static path is required")
	}

	if sc.IndexPath == "" {
		return fmt.Errorf("index path is required")
	}

	// Validate paths - accept absolute paths or relative paths (including those without ./ prefix)
	if filepath.IsAbs(sc.StaticPath) {
		// Absolute path is valid
	} else if strings.HasPrefix(sc.StaticPath, "./") || strings.HasPrefix(sc.StaticPath, "../") {
		// Explicit relative paths are valid
	} else if !strings.Contains(sc.StaticPath, "..") && sc.StaticPath != "" {
		// Simple relative paths without .. are valid (like "front/dist")
	} else {
		return fmt.Errorf("static path must be absolute or relative")
	}

	return nil
}

// GetCacheControl returns cache control header value
func (sc *SPAConfig) GetCacheControl() string {
	if sc.CacheControl != "" {
		return sc.CacheControl
	}
	return "public, max-age=3600" // Default 1 hour
}

// GetSecurityHeaders returns security headers map
func (sc *SPAConfig) GetSecurityHeaders() map[string]string {
	headers := make(map[string]string)
	headers["X-Frame-Options"] = "DENY"
	headers["X-Content-Type-Options"] = "nosniff"
	return headers
}

// GetCORSHeaders returns CORS headers if enabled
func (sc *SPAConfig) GetCORSHeaders() map[string]string {
	return make(map[string]string)
}

// Domain methods for FileInfo

// IsHTML checks if file is HTML
func (fi *FileInfo) IsHTML() bool {
	return fi.ContentType == "text/html"
}

// Domain methods for SPAStats

// UpdateUptime updates uptime based on start time
func (ss *SPAStats) UpdateUptime() {
	ss.Uptime = time.Since(ss.StartTime)
}
