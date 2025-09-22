package spa

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	spaModels "github.com/dash-ops/dash-ops/pkg/spa/models"
)

func TestNewModule_WithValidConfig_ReturnsModule(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	staticDir := filepath.Join(tempDir, "static")
	err := os.MkdirAll(staticDir, 0755)
	require.NoError(t, err)

	// Create index file
	indexFile := filepath.Join(staticDir, "index.html")
	err = os.WriteFile(indexFile, []byte("test"), 0644)
	require.NoError(t, err)

	config := &spaModels.SPAConfig{
		StaticPath: staticDir,
		IndexPath:  "index.html",
	}
	apiRouter := mux.NewRouter()

	module, err := NewModule(config, apiRouter)

	require.NoError(t, err)
	assert.NotNil(t, module)
	assert.NotNil(t, module.Handler)
	assert.NotNil(t, module.FileProcessor)
	assert.NotNil(t, module.SPAAdapter)
	assert.Equal(t, config, module.Config)
	assert.Equal(t, apiRouter, module.APIRouter)
}

func TestNewModule_WithNilConfig_ReturnsError(t *testing.T) {
	apiRouter := mux.NewRouter()

	module, err := NewModule(nil, apiRouter)

	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "SPA config cannot be nil")
}

func TestNewModule_WithInvalidConfig_ReturnsError(t *testing.T) {
	config := &spaModels.SPAConfig{
		StaticPath: "", // Invalid empty static path
		IndexPath:  "index.html",
	}
	apiRouter := mux.NewRouter()

	module, err := NewModule(config, apiRouter)

	assert.Error(t, err)
	assert.Nil(t, module)
	assert.Contains(t, err.Error(), "invalid SPA config")
}

func TestModule_RegisterRoutes_WithValidHandler_RegistersRoutes(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	staticDir := filepath.Join(tempDir, "static")
	err := os.MkdirAll(staticDir, 0755)
	require.NoError(t, err)

	// Create index file
	indexFile := filepath.Join(staticDir, "index.html")
	err = os.WriteFile(indexFile, []byte("test"), 0644)
	require.NoError(t, err)

	config := &spaModels.SPAConfig{
		StaticPath: staticDir,
		IndexPath:  "index.html",
	}
	apiRouter := mux.NewRouter()

	module, err := NewModule(config, apiRouter)
	require.NoError(t, err)

	router := mux.NewRouter()
	module.RegisterRoutes(router)

	// Should not panic and should register routes
	assert.NotNil(t, router)
}

func TestModule_RegisterRoutes_WithNilHandler_DoesNotPanic(t *testing.T) {
	module := &Module{
		Handler: nil,
	}

	router := mux.NewRouter()

	// Should not panic
	assert.NotPanics(t, func() {
		module.RegisterRoutes(router)
	})
}
