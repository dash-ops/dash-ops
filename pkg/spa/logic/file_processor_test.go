package spa

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileProcessor_ResolveFilePath_WithDirectoryWithoutIndex_ReturnsIndex(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary directory
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	// Create index file in root
	indexFile := filepath.Join(staticPath, "index.html")
	err = os.WriteFile(indexFile, []byte("index content"), 0644)
	require.NoError(t, err)

	// Create subdirectory without index.html
	subDir := filepath.Join(staticPath, "subdir")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	requestPath := "/subdir"
	indexPath := "index.html"

	path, isIndex, err := processor.ResolveFilePath(requestPath, staticPath, indexPath)

	require.NoError(t, err)
	assert.True(t, isIndex)
	assert.Equal(t, indexFile, path)
}

func TestFileProcessor_ResolveFilePath_WithFileInSubdirectory_ReturnsFile(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary directory
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	// Create subdirectory with file
	subDir := filepath.Join(staticPath, "assets")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	filePath := filepath.Join(subDir, "style.css")
	err = os.WriteFile(filePath, []byte("css content"), 0644)
	require.NoError(t, err)

	requestPath := "/assets/style.css"
	indexPath := "index.html"
	expectedPath := filepath.Join(staticPath, "assets", "style.css")

	path, isIndex, err := processor.ResolveFilePath(requestPath, staticPath, indexPath)

	require.NoError(t, err)
	assert.False(t, isIndex)
	assert.Equal(t, expectedPath, path)
}

func TestFileProcessor_ValidateStaticPath_WithNonExistentDirectory_ReturnsError(t *testing.T) {
	processor := NewFileProcessor()

	err := processor.ValidateStaticPath("/nonexistent/directory")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "static directory does not exist")
}

func TestFileProcessor_ValidateIndexFile_WithNonExistentIndexFile_ReturnsError(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary directory
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	err = processor.ValidateIndexFile(staticPath, "nonexistent.html")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "index file does not exist")
}

func TestFileProcessor_ValidateIndexFile_WithValidIndexFile_ReturnsNoError(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary directory
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	// Create index file
	indexFile := filepath.Join(staticPath, "index.html")
	err = os.WriteFile(indexFile, []byte("<html>Index</html>"), 0644)
	require.NoError(t, err)

	err = processor.ValidateIndexFile(staticPath, "index.html")

	assert.NoError(t, err)
}

func TestFileProcessor_ValidateStaticPath_WithValidDirectory_ReturnsNoError(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary directory
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	err = processor.ValidateStaticPath(staticPath)

	assert.NoError(t, err)
}

func TestFileProcessor_ValidateStaticPath_WithFileInsteadOfDirectory_ReturnsError(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary file instead of directory
	tempFile := t.TempDir() + "/notadir"
	err := os.WriteFile(tempFile, []byte("content"), 0644)
	require.NoError(t, err)

	err = processor.ValidateStaticPath(tempFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "static path is not a directory")
}

func TestFileProcessor_IsPathSafe_WithUnsafePath_ReturnsFalse(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary directory
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	unsafePaths := []string{
		"../../../etc/passwd",
		"path/../../../etc/passwd",
		"path/with/../unsafe/..",
		"path/with/../../unsafe",
	}

	for _, path := range unsafePaths {
		t.Run(path, func(t *testing.T) {
			result := processor.isPathSafe(path, staticPath)
			assert.False(t, result, "Path %s should be unsafe", path)
		})
	}
}

func TestFileProcessor_IsPathSafe_WithSafePath_ReturnsTrue(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary directory
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	safePaths := []string{
		filepath.Join(staticPath, "index.html"),
		filepath.Join(staticPath, "assets", "style.css"),
		filepath.Join(staticPath, "js", "app.js"),
		filepath.Join(staticPath, "images", "logo.png"),
		filepath.Join(staticPath, "subdir", "page.html"),
	}

	for _, path := range safePaths {
		t.Run(path, func(t *testing.T) {
			result := processor.isPathSafe(path, staticPath)
			assert.True(t, result, "Path %s should be safe", path)
		})
	}
}

func TestFileProcessor_GetContentType_WithVariousExtensions_ReturnsCorrectTypes(t *testing.T) {
	processor := NewFileProcessor()

	testCases := []struct {
		filePath     string
		expectedType string
	}{
		{"style.css", "text/css; charset=utf-8"},
		{"app.js", "text/javascript; charset=utf-8"},
		{"script.js", "text/javascript; charset=utf-8"},
		{"data.json", "application/json"},
		{"image.png", "image/png"},
		{"image.jpg", "image/jpeg"},
		{"image.gif", "image/gif"},
		{"image.svg", "image/svg+xml"},
		{"document.pdf", "application/pdf"},
		{"archive.zip", "application/zip"},
		{"unknown.xyz", "chemical/x-xyz"},
	}

	for _, tc := range testCases {
		t.Run(tc.filePath, func(t *testing.T) {
			result := processor.getContentType(tc.filePath)
			assert.Equal(t, tc.expectedType, result)
		})
	}
}

func TestFileProcessor_GetContentType_WithEmptyExtension_ReturnsOctetStream(t *testing.T) {
	processor := NewFileProcessor()

	result := processor.getContentType("file")

	assert.Equal(t, "application/octet-stream", result)
}

func TestFileProcessor_GetFileInfo_WithValidFile_ReturnsFileInfo(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary file
	tempFile := t.TempDir() + "/test.html"
	err := os.WriteFile(tempFile, []byte("<html>Test</html>"), 0644)
	require.NoError(t, err)

	fileInfo, err := processor.GetFileInfo(tempFile)

	require.NoError(t, err)
	assert.NotNil(t, fileInfo)
	assert.Equal(t, tempFile, fileInfo.Path)
	assert.Equal(t, int64(17), fileInfo.Size) // 17 bytes
	assert.Equal(t, "text/html; charset=utf-8", fileInfo.ContentType)
	assert.NotEmpty(t, fileInfo.ETag)
}

func TestFileProcessor_GetFileInfo_WithNonExistentFile_ReturnsError(t *testing.T) {
	processor := NewFileProcessor()

	fileInfo, err := processor.GetFileInfo("/nonexistent/file.html")

	assert.Error(t, err)
	assert.Nil(t, fileInfo)
}

func TestFileProcessor_GenerateETag_WithValidFile_ReturnsETag(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary file
	tempFile := t.TempDir() + "/test.html"
	err := os.WriteFile(tempFile, []byte("<html>Test</html>"), 0644)
	require.NoError(t, err)

	// Get file info
	fileInfo, err := os.Stat(tempFile)
	require.NoError(t, err)

	etag := processor.generateETag(fileInfo)

	assert.NotEmpty(t, etag)
	assert.Contains(t, etag, "\"") // ETag should be quoted
	assert.Contains(t, etag, "17") // File size should be in ETag
}

func TestFileProcessor_ShouldServeIndex_WithIndexRequest_ReturnsTrue(t *testing.T) {
	processor := NewFileProcessor()

	requestPath := "/"
	result := processor.ShouldServeIndex(requestPath)

	assert.True(t, result)
}

func TestFileProcessor_ShouldServeIndex_WithNonIndexRequest_ReturnsTrue(t *testing.T) {
	processor := NewFileProcessor()

	requestPath := "/api/test"
	result := processor.ShouldServeIndex(requestPath)

	assert.True(t, result) // Paths without extensions should serve index
}

func TestFileProcessor_ShouldServeIndex_WithFileExtension_ReturnsFalse(t *testing.T) {
	processor := NewFileProcessor()

	requestPath := "/style.css"
	result := processor.ShouldServeIndex(requestPath)

	assert.False(t, result)
}

func TestFileProcessor_GetCompressionType_WithCompressibleFile_ReturnsGzip(t *testing.T) {
	processor := NewFileProcessor()

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

	for _, contentType := range compressibleTypes {
		t.Run(contentType, func(t *testing.T) {
			result := processor.GetCompressionType(contentType)
			assert.Equal(t, "gzip", result)
		})
	}
}

func TestFileProcessor_GetCompressionType_WithNonCompressibleFile_ReturnsEmpty(t *testing.T) {
	processor := NewFileProcessor()

	nonCompressibleTypes := []string{
		"image/png",
		"image/jpeg",
		"image/gif",
		"application/pdf",
		"application/zip",
	}

	for _, contentType := range nonCompressibleTypes {
		t.Run(contentType, func(t *testing.T) {
			result := processor.GetCompressionType(contentType)
			assert.Empty(t, result)
		})
	}
}

func TestFileProcessor_NormalizePath_WithValidPaths_ReturnsNormalizedPath(t *testing.T) {
	processor := NewFileProcessor()

	testCases := []struct {
		input    string
		expected string
	}{
		{"/", "/"},
		{"/index.html", "/index.html"},
		{"/assets/style.css", "/assets/style.css"},
		{"/api/test", "/api/test"},
		{"/subdir/page", "/subdir/page"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := processor.NormalizePath(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFileProcessor_NormalizePath_WithEmptyPath_ReturnsRoot(t *testing.T) {
	processor := NewFileProcessor()

	result := processor.NormalizePath("")

	assert.Equal(t, "/.", result)
}
