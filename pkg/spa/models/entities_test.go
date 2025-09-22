package spa

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSPAConfig_Validate_WithValidConfig_ReturnsNoError(t *testing.T) {
	config := &SPAConfig{
		StaticPath:   "front/dist",
		IndexPath:    "index.html",
		CacheControl: "public, max-age=3600",
	}

	err := config.Validate()
	assert.NoError(t, err)
}

func TestSPAConfig_Validate_WithEmptyStaticPath_ReturnsError(t *testing.T) {
	config := &SPAConfig{
		StaticPath: "",
		IndexPath:  "index.html",
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "static path is required")
}

func TestSPAConfig_Validate_WithEmptyIndexPath_ReturnsError(t *testing.T) {
	config := &SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "",
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "index path is required")
}

func TestSPAConfig_Validate_WithAbsolutePath_ReturnsNoError(t *testing.T) {
	config := &SPAConfig{
		StaticPath: "/absolute/path",
		IndexPath:  "index.html",
	}

	err := config.Validate()
	assert.NoError(t, err)
}

func TestSPAConfig_Validate_WithRelativePath_ReturnsNoError(t *testing.T) {
	config := &SPAConfig{
		StaticPath: "./relative/path",
		IndexPath:  "index.html",
	}

	err := config.Validate()
	assert.NoError(t, err)
}

func TestSPAConfig_Validate_WithParentPath_ReturnsNoError(t *testing.T) {
	config := &SPAConfig{
		StaticPath: "../parent/path",
		IndexPath:  "index.html",
	}

	err := config.Validate()
	assert.NoError(t, err)
}

func TestSPAConfig_Validate_WithSimpleRelativePath_ReturnsNoError(t *testing.T) {
	config := &SPAConfig{
		StaticPath: "front/dist",
		IndexPath:  "index.html",
	}

	err := config.Validate()
	assert.NoError(t, err)
}

func TestSPAConfig_Validate_WithUnsafePath_ReturnsError(t *testing.T) {
	config := &SPAConfig{
		StaticPath: "path/with/../unsafe/..",
		IndexPath:  "index.html",
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "static path must be absolute or relative")
}

func TestSPAConfig_GetCacheControl_WithCustomValue_ReturnsCustomValue(t *testing.T) {
	config := &SPAConfig{
		CacheControl: "public, max-age=7200",
	}

	result := config.GetCacheControl()
	assert.Equal(t, "public, max-age=7200", result)
}

func TestSPAConfig_GetCacheControl_WithEmptyValue_ReturnsDefault(t *testing.T) {
	config := &SPAConfig{}

	result := config.GetCacheControl()
	assert.Equal(t, "public, max-age=3600", result)
}

func TestSPAConfig_GetSecurityHeaders_ReturnsDefaultHeaders(t *testing.T) {
	config := &SPAConfig{}

	headers := config.GetSecurityHeaders()
	assert.Equal(t, "DENY", headers["X-Frame-Options"])
	assert.Equal(t, "nosniff", headers["X-Content-Type-Options"])
}

func TestSPAConfig_GetCORSHeaders_ReturnsEmptyMap(t *testing.T) {
	config := &SPAConfig{}

	headers := config.GetCORSHeaders()
	assert.Empty(t, headers)
}

func TestFileInfo_IsHTML_WithHTMLContentType_ReturnsTrue(t *testing.T) {
	fileInfo := &FileInfo{
		ContentType: "text/html",
	}

	result := fileInfo.IsHTML()
	assert.True(t, result)
}

func TestFileInfo_IsHTML_WithNonHTMLContentType_ReturnsFalse(t *testing.T) {
	fileInfo := &FileInfo{
		ContentType: "text/css",
	}

	result := fileInfo.IsHTML()
	assert.False(t, result)
}

func TestSPAStats_UpdateUptime_UpdatesUptime(t *testing.T) {
	startTime := time.Now().Add(-5 * time.Minute)
	stats := &SPAStats{
		StartTime: startTime,
	}

	stats.UpdateUptime()

	// Uptime should be approximately 5 minutes
	assert.True(t, stats.Uptime > 4*time.Minute)
	assert.True(t, stats.Uptime < 6*time.Minute)
}
