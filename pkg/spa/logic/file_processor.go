package spa

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	spaModels "github.com/dash-ops/dash-ops/pkg/spa/models"
)

// FileProcessor handles file processing logic for SPA
type FileProcessor struct{}

// NewFileProcessor creates a new file processor
func NewFileProcessor() *FileProcessor {
	return &FileProcessor{}
}

// ResolveFilePath resolves the file path for a request
func (fp *FileProcessor) ResolveFilePath(requestPath, staticPath, indexPath string) (string, bool, error) {
	// Security check - detect path traversal attempts before cleaning
	if strings.Contains(requestPath, "..") {
		return "", false, fmt.Errorf("path traversal attempt detected")
	}

	// Clean and validate request path
	cleanPath := filepath.Clean(requestPath)

	// For root paths, always serve index.html for SPA routing
	if cleanPath == "." || cleanPath == "/" {
		indexFullPath := filepath.Join(staticPath, indexPath)
		return indexFullPath, true, nil
	}

	// Build full file path
	fullPath := filepath.Join(staticPath, cleanPath)

	// Additional security check - ensure path is within static directory
	if !fp.isPathSafe(fullPath, staticPath) {
		return "", false, fmt.Errorf("path traversal attempt detected")
	}

	// Check if file exists
	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		// File doesn't exist, serve index.html for SPA routing
		indexFullPath := filepath.Join(staticPath, indexPath)
		return indexFullPath, true, nil
	}

	if err != nil {
		return "", false, fmt.Errorf("failed to stat file: %w", err)
	}

	// If it's a directory, try to serve index.html from that directory
	if fileInfo.IsDir() {
		indexInDir := filepath.Join(fullPath, "index.html")
		if _, err := os.Stat(indexInDir); err == nil {
			return indexInDir, false, nil
		}
		// Directory without index, serve main index.html
		indexFullPath := filepath.Join(staticPath, indexPath)
		return indexFullPath, true, nil
	}

	return fullPath, false, nil
}

// GetFileInfo gets comprehensive file information
func (fp *FileProcessor) GetFileInfo(filePath string) (*spaModels.FileInfo, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Determine content type
	contentType := fp.getContentType(filePath)

	// Generate ETag (simplified)
	etag := fp.generateETag(fileInfo)

	return &spaModels.FileInfo{
		Path:        filePath,
		Size:        fileInfo.Size(),
		ContentType: contentType,
		ETag:        etag,
	}, nil
}

// ValidateStaticPath validates static directory path
func (fp *FileProcessor) ValidateStaticPath(staticPath string) error {
	if staticPath == "" {
		return fmt.Errorf("static path cannot be empty")
	}

	// Check if directory exists
	fileInfo, err := os.Stat(staticPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("static directory does not exist: %s", staticPath)
	}

	if err != nil {
		return fmt.Errorf("failed to access static directory: %w", err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("static path is not a directory: %s", staticPath)
	}

	return nil
}

// ValidateIndexFile validates index file exists
func (fp *FileProcessor) ValidateIndexFile(staticPath, indexPath string) error {
	if indexPath == "" {
		return fmt.Errorf("index path cannot be empty")
	}

	fullIndexPath := filepath.Join(staticPath, indexPath)
	fileInfo, err := os.Stat(fullIndexPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("index file does not exist: %s", fullIndexPath)
	}

	if err != nil {
		return fmt.Errorf("failed to access index file: %w", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("index path is a directory, not a file: %s", fullIndexPath)
	}

	return nil
}

// isPathSafe checks if the resolved path is safe (no path traversal)
func (fp *FileProcessor) isPathSafe(requestedPath, staticPath string) bool {
	// Get absolute paths
	absRequested, err := filepath.Abs(requestedPath)
	if err != nil {
		return false
	}

	absStatic, err := filepath.Abs(staticPath)
	if err != nil {
		return false
	}

	// Check if requested path is within static directory
	return strings.HasPrefix(absRequested, absStatic)
}

// getContentType determines content type from file extension
func (fp *FileProcessor) getContentType(filePath string) string {
	ext := filepath.Ext(filePath)
	contentType := mime.TypeByExtension(ext)

	if contentType == "" {
		// Default content types for common web files
		switch strings.ToLower(ext) {
		case ".html", ".htm":
			return "text/html; charset=utf-8"
		case ".css":
			return "text/css; charset=utf-8"
		case ".js", ".mjs":
			return "application/javascript; charset=utf-8"
		case ".json":
			return "application/json; charset=utf-8"
		case ".xml":
			return "application/xml; charset=utf-8"
		case ".svg":
			return "image/svg+xml"
		case ".woff":
			return "font/woff"
		case ".woff2":
			return "font/woff2"
		case ".ttf":
			return "font/ttf"
		case ".eot":
			return "application/vnd.ms-fontobject"
		default:
			return "application/octet-stream"
		}
	}

	return contentType
}

// generateETag generates a simple ETag based on file info
func (fp *FileProcessor) generateETag(fileInfo os.FileInfo) string {
	// Simple ETag based on size and modification time
	return fmt.Sprintf(`"%d-%d"`, fileInfo.Size(), fileInfo.ModTime().Unix())
}

// ShouldServeIndex determines if index.html should be served
func (fp *FileProcessor) ShouldServeIndex(requestPath string) bool {
	// Serve index.html for routes that don't have file extensions
	// This supports client-side routing in SPAs
	ext := filepath.Ext(requestPath)
	return ext == "" || requestPath == "/"
}

// GetCompressionType determines compression type for file
func (fp *FileProcessor) GetCompressionType(contentType string) string {
	// Determine what files should be compressed
	compressibleTypes := []string{
		"text/html",
		"text/css",
		"application/javascript",
		"application/json",
		"text/plain",
		"application/xml",
		"text/xml",
		"image/svg+xml",
	}

	for _, compressibleType := range compressibleTypes {
		if strings.HasPrefix(contentType, compressibleType) {
			return "gzip"
		}
	}

	return "" // No compression
}

// NormalizePath normalizes URL path for file serving
func (fp *FileProcessor) NormalizePath(path string) string {
	// Remove query parameters and fragments
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}
	if idx := strings.Index(path, "#"); idx != -1 {
		path = path[:idx]
	}

	// Clean path
	path = filepath.Clean(path)

	// Ensure it starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path
}
