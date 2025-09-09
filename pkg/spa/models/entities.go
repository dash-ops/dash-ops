package spa

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// SPAConfig represents SPA server configuration
type SPAConfig struct {
	StaticPath    string            `yaml:"static_path" json:"static_path"`
	IndexPath     string            `yaml:"index_path" json:"index_path"`
	CacheControl  string            `yaml:"cache_control" json:"cache_control"`
	Compression   bool              `yaml:"compression" json:"compression"`
	CORS          CORSConfig        `yaml:"cors" json:"cors"`
	Security      SecurityConfig    `yaml:"security" json:"security"`
	CustomHeaders map[string]string `yaml:"custom_headers" json:"custom_headers"`
	Enabled       bool              `yaml:"enabled" json:"enabled"`
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	Enabled          bool     `yaml:"enabled" json:"enabled"`
	AllowedOrigins   []string `yaml:"allowed_origins" json:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods" json:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers" json:"allowed_headers"`
	ExposedHeaders   []string `yaml:"exposed_headers" json:"exposed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials" json:"allow_credentials"`
	MaxAge           int      `yaml:"max_age" json:"max_age"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	ContentSecurityPolicy string `yaml:"content_security_policy" json:"content_security_policy"`
	XFrameOptions         string `yaml:"x_frame_options" json:"x_frame_options"`
	XContentTypeOptions   string `yaml:"x_content_type_options" json:"x_content_type_options"`
	ReferrerPolicy        string `yaml:"referrer_policy" json:"referrer_policy"`
	PermissionsPolicy     string `yaml:"permissions_policy" json:"permissions_policy"`
}

// FileInfo represents information about a served file
type FileInfo struct {
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	ModTime      time.Time `json:"mod_time"`
	ContentType  string    `json:"content_type"`
	ETag         string    `json:"etag,omitempty"`
	Compressed   bool      `json:"compressed"`
	CacheControl string    `json:"cache_control"`
}

// RequestInfo represents information about a request
type RequestInfo struct {
	Method       string            `json:"method"`
	Path         string            `json:"path"`
	UserAgent    string            `json:"user_agent"`
	RemoteAddr   string            `json:"remote_addr"`
	Referer      string            `json:"referer,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	Timestamp    time.Time         `json:"timestamp"`
	ResponseCode int               `json:"response_code"`
	ResponseSize int64             `json:"response_size"`
	Duration     time.Duration     `json:"duration"`
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

// GetIndexFilePath returns full path to index file
func (sc *SPAConfig) GetIndexFilePath() string {
	return filepath.Join(sc.StaticPath, sc.IndexPath)
}

// GetCacheControl returns cache control header value
func (sc *SPAConfig) GetCacheControl() string {
	if sc.CacheControl != "" {
		return sc.CacheControl
	}
	return "public, max-age=3600" // Default 1 hour
}

// ShouldCompress checks if compression should be enabled
func (sc *SPAConfig) ShouldCompress() bool {
	return sc.Compression
}

// GetSecurityHeaders returns security headers map
func (sc *SPAConfig) GetSecurityHeaders() map[string]string {
	headers := make(map[string]string)

	if sc.Security.ContentSecurityPolicy != "" {
		headers["Content-Security-Policy"] = sc.Security.ContentSecurityPolicy
	}

	if sc.Security.XFrameOptions != "" {
		headers["X-Frame-Options"] = sc.Security.XFrameOptions
	} else {
		headers["X-Frame-Options"] = "DENY" // Default
	}

	if sc.Security.XContentTypeOptions != "" {
		headers["X-Content-Type-Options"] = sc.Security.XContentTypeOptions
	} else {
		headers["X-Content-Type-Options"] = "nosniff" // Default
	}

	if sc.Security.ReferrerPolicy != "" {
		headers["Referrer-Policy"] = sc.Security.ReferrerPolicy
	}

	if sc.Security.PermissionsPolicy != "" {
		headers["Permissions-Policy"] = sc.Security.PermissionsPolicy
	}

	return headers
}

// GetCORSHeaders returns CORS headers if enabled
func (sc *SPAConfig) GetCORSHeaders() map[string]string {
	headers := make(map[string]string)

	if !sc.CORS.Enabled {
		return headers
	}

	if len(sc.CORS.AllowedOrigins) > 0 {
		headers["Access-Control-Allow-Origin"] = strings.Join(sc.CORS.AllowedOrigins, ", ")
	}

	if len(sc.CORS.AllowedMethods) > 0 {
		headers["Access-Control-Allow-Methods"] = strings.Join(sc.CORS.AllowedMethods, ", ")
	}

	if len(sc.CORS.AllowedHeaders) > 0 {
		headers["Access-Control-Allow-Headers"] = strings.Join(sc.CORS.AllowedHeaders, ", ")
	}

	if len(sc.CORS.ExposedHeaders) > 0 {
		headers["Access-Control-Expose-Headers"] = strings.Join(sc.CORS.ExposedHeaders, ", ")
	}

	if sc.CORS.AllowCredentials {
		headers["Access-Control-Allow-Credentials"] = "true"
	}

	if sc.CORS.MaxAge > 0 {
		headers["Access-Control-Max-Age"] = fmt.Sprintf("%d", sc.CORS.MaxAge)
	}

	return headers
}

// Domain methods for FileInfo

// IsImage checks if file is an image
func (fi *FileInfo) IsImage() bool {
	imageTypes := []string{
		"image/jpeg", "image/jpg", "image/png", "image/gif",
		"image/webp", "image/svg+xml", "image/bmp", "image/ico",
	}

	for _, imageType := range imageTypes {
		if fi.ContentType == imageType {
			return true
		}
	}
	return false
}

// IsJavaScript checks if file is JavaScript
func (fi *FileInfo) IsJavaScript() bool {
	return fi.ContentType == "application/javascript" ||
		fi.ContentType == "text/javascript"
}

// IsCSS checks if file is CSS
func (fi *FileInfo) IsCSS() bool {
	return fi.ContentType == "text/css"
}

// IsHTML checks if file is HTML
func (fi *FileInfo) IsHTML() bool {
	return fi.ContentType == "text/html"
}

// ShouldCache checks if file should be cached based on type
func (fi *FileInfo) ShouldCache() bool {
	// Cache static assets but not HTML files
	return !fi.IsHTML()
}

// Domain methods for SPAStats

// CalculateSuccessRate returns success rate percentage
func (ss *SPAStats) CalculateSuccessRate() float64 {
	if ss.TotalRequests == 0 {
		return 0
	}
	return float64(ss.SuccessfulRequests) / float64(ss.TotalRequests) * 100
}

// CalculateErrorRate returns error rate percentage
func (ss *SPAStats) CalculateErrorRate() float64 {
	if ss.TotalRequests == 0 {
		return 0
	}
	return float64(ss.ErrorRequests) / float64(ss.TotalRequests) * 100
}

// UpdateUptime updates uptime based on start time
func (ss *SPAStats) UpdateUptime() {
	ss.Uptime = time.Since(ss.StartTime)
}
