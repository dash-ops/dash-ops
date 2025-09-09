package spa

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileProcessor_ResolveFilePath_WithRootPath_ReturnsIndexFile(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	indexPath := "index.html"

	// Create static directory
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	// Create test files
	err = os.WriteFile(filepath.Join(staticPath, "index.html"), []byte("index content"), 0644)
	require.NoError(t, err)

	requestPath := "/"
	expectedPath := filepath.Join(staticPath, indexPath)
	expectedIsIndex := true
	expectError := false

	// Act
	path, isIndex, err := processor.ResolveFilePath(requestPath, staticPath, indexPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, path)
		assert.Equal(t, expectedIsIndex, isIndex)
	}
}

func TestFileProcessor_ResolveFilePath_WithEmptyPath_ReturnsIndexFile(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	indexPath := "index.html"

	// Create static directory
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	// Create test files
	err = os.WriteFile(filepath.Join(staticPath, "index.html"), []byte("index content"), 0644)
	require.NoError(t, err)

	requestPath := ""
	expectedPath := filepath.Join(staticPath, indexPath)
	expectedIsIndex := true
	expectError := false

	// Act
	path, isIndex, err := processor.ResolveFilePath(requestPath, staticPath, indexPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, path)
		assert.Equal(t, expectedIsIndex, isIndex)
	}
}

func TestFileProcessor_ResolveFilePath_WithExistingFile_ReturnsFile(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	indexPath := "index.html"

	// Create static directory
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	// Create test files
	err = os.WriteFile(filepath.Join(staticPath, "app.js"), []byte("js content"), 0644)
	require.NoError(t, err)

	requestPath := "/app.js"
	expectedPath := filepath.Join(staticPath, "app.js")
	expectedIsIndex := false
	expectError := false

	// Act
	path, isIndex, err := processor.ResolveFilePath(requestPath, staticPath, indexPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, path)
		assert.Equal(t, expectedIsIndex, isIndex)
	}
}

func TestFileProcessor_ResolveFilePath_WithExistingCSSFile_ReturnsFile(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	indexPath := "index.html"

	// Create static directory
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	// Create test files
	err = os.WriteFile(filepath.Join(staticPath, "style.css"), []byte("css content"), 0644)
	require.NoError(t, err)

	requestPath := "/style.css"
	expectedPath := filepath.Join(staticPath, "style.css")
	expectedIsIndex := false
	expectError := false

	// Act
	path, isIndex, err := processor.ResolveFilePath(requestPath, staticPath, indexPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, path)
		assert.Equal(t, expectedIsIndex, isIndex)
	}
}

func TestFileProcessor_ResolveFilePath_WithDirectoryWithIndex_ReturnsIndexFile(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	indexPath := "index.html"

	// Create static directory
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	// Create subdirectory
	subDir := filepath.Join(staticPath, "assets")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(subDir, "index.html"), []byte("subdir index"), 0644)
	require.NoError(t, err)

	requestPath := "/assets"
	expectedPath := filepath.Join(subDir, "index.html")
	expectedIsIndex := false
	expectError := false

	// Act
	path, isIndex, err := processor.ResolveFilePath(requestPath, staticPath, indexPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, path)
		assert.Equal(t, expectedIsIndex, isIndex)
	}
}

func TestFileProcessor_ResolveFilePath_WithDirectoryWithoutIndex_ReturnsMainIndexFile(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	indexPath := "index.html"

	// Create static directory
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	// Create test files
	err = os.WriteFile(filepath.Join(staticPath, "index.html"), []byte("index content"), 0644)
	require.NoError(t, err)

	requestPath := "/nonexistent"
	expectedPath := filepath.Join(staticPath, indexPath)
	expectedIsIndex := true
	expectError := false

	// Act
	path, isIndex, err := processor.ResolveFilePath(requestPath, staticPath, indexPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, path)
		assert.Equal(t, expectedIsIndex, isIndex)
	}
}

func TestFileProcessor_ResolveFilePath_WithPathTraversalAttempt_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	staticPath := filepath.Join(tempDir, "static")
	indexPath := "index.html"

	// Create static directory
	err := os.MkdirAll(staticPath, 0755)
	require.NoError(t, err)

	requestPath := "/../etc/passwd"
	expectedPath := ""
	expectedIsIndex := false
	expectError := true

	// Act
	path, isIndex, err := processor.ResolveFilePath(requestPath, staticPath, indexPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, expectedPath, path)
		assert.Equal(t, expectedIsIndex, isIndex)
	}
}

func TestFileProcessor_GetFileInfo(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary file
	tempFile := filepath.Join(t.TempDir(), "test.txt")
	content := "test content"
	err := os.WriteFile(tempFile, []byte(content), 0644)
	require.NoError(t, err)

	// Set modification time
	modTime := time.Now().Add(-time.Hour)
	err = os.Chtimes(tempFile, modTime, modTime)
	require.NoError(t, err)

	fileInfo, err := processor.GetFileInfo(tempFile)
	require.NoError(t, err)
	require.NotNil(t, fileInfo)

	assert.Equal(t, tempFile, fileInfo.Path)
	assert.Equal(t, int64(len(content)), fileInfo.Size)
	assert.WithinDuration(t, modTime, fileInfo.ModTime, time.Second)
	assert.Equal(t, "text/plain; charset=utf-8", fileInfo.ContentType)
	assert.NotEmpty(t, fileInfo.ETag)
	assert.False(t, fileInfo.Compressed)
}

func TestFileProcessor_GetFileInfo_NonExistent(t *testing.T) {
	processor := NewFileProcessor()

	fileInfo, err := processor.GetFileInfo("/nonexistent/file.txt")
	assert.Error(t, err)
	assert.Nil(t, fileInfo)
}

func TestFileProcessor_ValidateStaticPath_WithValidDirectory_ReturnsNoError(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	staticPath := tempDir
	expectError := false

	// Act
	err := processor.ValidateStaticPath(staticPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestFileProcessor_ValidateStaticPath_WithEmptyPath_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	staticPath := ""
	expectError := true

	// Act
	err := processor.ValidateStaticPath(staticPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestFileProcessor_ValidateStaticPath_WithNonexistentDirectory_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	staticPath := "/nonexistent/directory"
	expectError := true

	// Act
	err := processor.ValidateStaticPath(staticPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestFileProcessor_ValidateStaticPath_WithFileInsteadOfDirectory_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	staticPath := tempDir + "/file.txt"
	expectError := true

	// Create file for the test case
	err := os.WriteFile(staticPath, []byte("content"), 0644)
	require.NoError(t, err)

	// Act
	err = processor.ValidateStaticPath(staticPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestFileProcessor_ValidateIndexFile_WithValidIndexFile_ReturnsNoError(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	indexFile := filepath.Join(tempDir, "index.html")
	err := os.WriteFile(indexFile, []byte("index content"), 0644)
	require.NoError(t, err)

	staticPath := tempDir
	indexPath := "index.html"
	expectError := false

	// Act
	err = processor.ValidateIndexFile(staticPath, indexPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestFileProcessor_ValidateIndexFile_WithEmptyIndexPath_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	indexFile := filepath.Join(tempDir, "index.html")
	err := os.WriteFile(indexFile, []byte("index content"), 0644)
	require.NoError(t, err)

	staticPath := tempDir
	indexPath := ""
	expectError := true

	// Act
	err = processor.ValidateIndexFile(staticPath, indexPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestFileProcessor_ValidateIndexFile_WithNonexistentIndexFile_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	indexFile := filepath.Join(tempDir, "index.html")
	err := os.WriteFile(indexFile, []byte("index content"), 0644)
	require.NoError(t, err)

	staticPath := tempDir
	indexPath := "nonexistent.html"
	expectError := true

	// Act
	err = processor.ValidateIndexFile(staticPath, indexPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestFileProcessor_ValidateIndexFile_WithIndexPathAsDirectory_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	tempDir := t.TempDir()
	indexFile := filepath.Join(tempDir, "index.html")
	err := os.WriteFile(indexFile, []byte("index content"), 0644)
	require.NoError(t, err)

	staticPath := tempDir
	indexPath := "."
	expectError := true

	// Act
	err = processor.ValidateIndexFile(staticPath, indexPath)

	// Assert
	if expectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestFileProcessor_isPathSafe_WithSafePath_ReturnsTrue(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	staticPath := "/safe/static/path"
	requestedPath := "/safe/static/path/file.txt"
	expected := true

	// Act
	result := processor.isPathSafe(requestedPath, staticPath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_isPathSafe_WithPathTraversalAttempt_ReturnsFalse(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	staticPath := "/safe/static/path"
	requestedPath := "/safe/static/path/../../../etc/passwd"
	expected := false

	// Act
	result := processor.isPathSafe(requestedPath, staticPath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_isPathSafe_WithDifferentDirectory_ReturnsFalse(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	staticPath := "/safe/static/path"
	requestedPath := "/other/directory/file.txt"
	expected := false

	// Act
	result := processor.isPathSafe(requestedPath, staticPath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_isPathSafe_WithExactStaticPath_ReturnsTrue(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	staticPath := "/safe/static/path"
	requestedPath := "/safe/static/path"
	expected := true

	// Act
	result := processor.isPathSafe(requestedPath, staticPath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithHTMLFile_ReturnsHTMLContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "index.html"
	expected := "text/html; charset=utf-8"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithCSSFile_ReturnsCSSContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "style.css"
	expected := "text/css; charset=utf-8"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithJavaScriptFile_ReturnsJavaScriptContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "app.js"
	expected := "text/javascript; charset=utf-8"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithJSONFile_ReturnsJSONContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "data.json"
	expected := "application/json"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithXMLFile_ReturnsXMLContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "config.xml"

	// Act
	result := processor.getContentType(filePath)

	// Accept both possible MIME types for XML (varies by OS)
	assert.True(t, result == "application/xml" || result == "text/xml; charset=utf-8",
		"Expected XML content type, got: %s", result)
}

func TestFileProcessor_getContentType_WithSVGFile_ReturnsSVGContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "icon.svg"
	expected := "image/svg+xml"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithWOFFFont_ReturnsWOFFContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "font.woff"
	expected := "font/woff"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithWOFF2Font_ReturnsWOFF2ContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "font.woff2"
	expected := "font/woff2"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithTTFFont_ReturnsTTFContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "font.ttf"
	expected := "font/ttf"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithEOTFont_ReturnsEOTContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "font.eot"
	expected := "application/vnd.ms-fontobject"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithUnknownExtension_ReturnsOctetStreamContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "file.unknown"
	expected := "application/octet-stream"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithNoExtension_ReturnsOctetStreamContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := "file"
	expected := "application/octet-stream"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_getContentType_WithEmptyPath_ReturnsOctetStreamContentType(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	filePath := ""
	expected := "application/octet-stream"

	// Act
	result := processor.getContentType(filePath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_generateETag(t *testing.T) {
	processor := NewFileProcessor()

	// Create temporary file
	tempFile := filepath.Join(t.TempDir(), "test.txt")
	content := "test content"
	err := os.WriteFile(tempFile, []byte(content), 0644)
	require.NoError(t, err)

	fileInfo, err := os.Stat(tempFile)
	require.NoError(t, err)

	etag := processor.generateETag(fileInfo)

	// ETag should be in format "size-timestamp"
	assert.Regexp(t, `^"\d+-\d+"$`, etag)
	assert.Contains(t, etag, "12") // size of "test content"
}

func TestFileProcessor_ShouldServeIndex_WithRootPath_ReturnsTrue(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	requestPath := "/"
	expected := true

	// Act
	result := processor.ShouldServeIndex(requestPath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_ShouldServeIndex_WithEmptyPath_ReturnsTrue(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	requestPath := ""
	expected := true

	// Act
	result := processor.ShouldServeIndex(requestPath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_ShouldServeIndex_WithPathWithoutExtension_ReturnsTrue(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	requestPath := "/dashboard"
	expected := true

	// Act
	result := processor.ShouldServeIndex(requestPath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_ShouldServeIndex_WithPathWithJSExtension_ReturnsFalse(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	requestPath := "/app.js"
	expected := false

	// Act
	result := processor.ShouldServeIndex(requestPath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_ShouldServeIndex_WithPathWithCSSExtension_ReturnsFalse(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	requestPath := "/style.css"
	expected := false

	// Act
	result := processor.ShouldServeIndex(requestPath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_ShouldServeIndex_WithPathWithJSONExtension_ReturnsFalse(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	requestPath := "/data.json"
	expected := false

	// Act
	result := processor.ShouldServeIndex(requestPath)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_GetCompressionType_WithHTML_ReturnsGzip(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	contentType := "text/html"
	expected := "gzip"

	// Act
	result := processor.GetCompressionType(contentType)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_GetCompressionType_WithCSS_ReturnsGzip(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	contentType := "text/css"
	expected := "gzip"

	// Act
	result := processor.GetCompressionType(contentType)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_GetCompressionType_WithJavaScript_ReturnsGzip(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	contentType := "application/javascript"
	expected := "gzip"

	// Act
	result := processor.GetCompressionType(contentType)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_GetCompressionType_WithJSON_ReturnsGzip(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	contentType := "application/json"
	expected := "gzip"

	// Act
	result := processor.GetCompressionType(contentType)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_GetCompressionType_WithPlainText_ReturnsGzip(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	contentType := "text/plain"
	expected := "gzip"

	// Act
	result := processor.GetCompressionType(contentType)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_GetCompressionType_WithXML_ReturnsGzip(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	contentType := "application/xml"
	expected := "gzip"

	// Act
	result := processor.GetCompressionType(contentType)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_GetCompressionType_WithSVG_ReturnsGzip(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	contentType := "image/svg+xml"
	expected := "gzip"

	// Act
	result := processor.GetCompressionType(contentType)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_GetCompressionType_WithImage_ReturnsEmpty(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	contentType := "image/png"
	expected := ""

	// Act
	result := processor.GetCompressionType(contentType)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_GetCompressionType_WithVideo_ReturnsEmpty(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	contentType := "video/mp4"
	expected := ""

	// Act
	result := processor.GetCompressionType(contentType)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_GetCompressionType_WithFont_ReturnsEmpty(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	contentType := "font/woff2"
	expected := ""

	// Act
	result := processor.GetCompressionType(contentType)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_GetCompressionType_WithUnknown_ReturnsEmpty(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	contentType := "application/unknown"
	expected := ""

	// Act
	result := processor.GetCompressionType(contentType)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_NormalizePath_WithSimplePath_ReturnsSamePath(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	path := "/dashboard"
	expected := "/dashboard"

	// Act
	result := processor.NormalizePath(path)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_NormalizePath_WithPathWithQuery_ReturnsPathWithoutQuery(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	path := "/dashboard?param=value"
	expected := "/dashboard"

	// Act
	result := processor.NormalizePath(path)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_NormalizePath_WithPathWithFragment_ReturnsPathWithoutFragment(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	path := "/dashboard#section"
	expected := "/dashboard"

	// Act
	result := processor.NormalizePath(path)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_NormalizePath_WithPathWithBothQueryAndFragment_ReturnsPathWithoutBoth(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	path := "/dashboard?param=value#section"
	expected := "/dashboard"

	// Act
	result := processor.NormalizePath(path)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_NormalizePath_WithRelativePath_ReturnsAbsolutePath(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	path := "dashboard"
	expected := "/dashboard"

	// Act
	result := processor.NormalizePath(path)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_NormalizePath_WithPathWithDots_ReturnsCleanedPath(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	path := "/../dashboard"
	expected := "/dashboard"

	// Act
	result := processor.NormalizePath(path)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_NormalizePath_WithEmptyPath_ReturnsRootPath(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	path := ""
	expected := "/."

	// Act
	result := processor.NormalizePath(path)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_NormalizePath_WithRootPath_ReturnsRootPath(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	path := "/"
	expected := "/"

	// Act
	result := processor.NormalizePath(path)

	// Assert
	assert.Equal(t, expected, result)
}

func TestFileProcessor_NormalizePath_WithMultipleSlashes_ReturnsCleanedPath(t *testing.T) {
	// Arrange
	processor := NewFileProcessor()
	path := "//dashboard//"
	expected := "/dashboard"

	// Act
	result := processor.NormalizePath(path)

	// Assert
	assert.Equal(t, expected, result)
}
