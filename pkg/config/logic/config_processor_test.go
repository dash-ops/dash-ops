package config

import (
	"os"
	"path/filepath"
	"testing"

	configModels "github.com/dash-ops/dash-ops/pkg/config/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigProcessor_ParseFromBytes_WithValidConfig_ReturnsConfig(t *testing.T) {
	// Arrange
	processor := NewConfigProcessor()
	configData := `port: 8080
origin: http://localhost:3000
headers: 
  - "Content-Type"
  - "Authorization"
plugins:
  - "OAuth2"
  - "Kubernetes"`

	// Act
	config, err := processor.ParseFromBytes([]byte(configData))

	// Assert
	require.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "http://localhost:3000", config.Origin)
}

func TestConfigProcessor_ParseFromBytes_WithEnvironmentVariables_ReturnsConfigWithDefaults(t *testing.T) {
	// Arrange
	processor := NewConfigProcessor()
	configData := `port: ${PORT:-8080}
origin: ${ORIGIN:-http://localhost:3000}
plugins:
  - "OAuth2"`

	// Act
	config, err := processor.ParseFromBytes([]byte(configData))

	// Assert
	require.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, "8080", config.Port)                    // Default value
	assert.Equal(t, "http://localhost:3000", config.Origin) // Default value
}

func TestConfigProcessor_ParseFromBytes_WithEmptyConfig_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewConfigProcessor()
	configData := ""

	// Act
	config, err := processor.ParseFromBytes([]byte(configData))

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestConfigProcessor_ParseFromBytes_WithInvalidYAML_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewConfigProcessor()
	configData := `port: 8080
origin: [invalid yaml`

	// Act
	config, err := processor.ParseFromBytes([]byte(configData))

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestConfigProcessor_LoadFromFile_WithValidFile_ReturnsConfig(t *testing.T) {
	// Arrange
	processor := NewConfigProcessor()
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test-config.yaml")

	configContent := `port: 9090
origin: http://test.local
headers: 
  - "Content-Type"
plugins:
  - "TestPlugin"`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// Act
	config, err := processor.LoadFromFile(configFile)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, "9090", config.Port)
	assert.Equal(t, "http://test.local", config.Origin)
	assert.Equal(t, []string{"Content-Type"}, config.Headers)
	assert.True(t, config.Plugins.Has("TestPlugin"))
}

func TestConfigProcessor_LoadFromFile_WithNonExistentFile_ReturnsError(t *testing.T) {
	// Arrange
	processor := NewConfigProcessor()

	// Act
	config, err := processor.LoadFromFile("/non/existent/file.yaml")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestConfigProcessor_GetConfigFilePath_WithDefaultPath_ReturnsDefaultPath(t *testing.T) {
	// Arrange
	processor := NewConfigProcessor()

	// Act
	path := processor.GetConfigFilePath()

	// Assert
	assert.Equal(t, "./dash-ops.yaml", path)
}

func TestConfigProcessor_GetConfigFilePath_WithEnvironmentVariable_ReturnsCustomPath(t *testing.T) {
	// Arrange
	processor := NewConfigProcessor()
	os.Setenv("DASH_CONFIG", "/custom/path/config.yaml")
	defer os.Unsetenv("DASH_CONFIG")

	// Act
	path := processor.GetConfigFilePath()

	// Assert
	assert.Equal(t, "/custom/path/config.yaml", path)
}

func TestConfigProcessor_MergeConfigs_WithBaseAndOverride_ReturnsMergedConfig(t *testing.T) {
	// Arrange
	processor := NewConfigProcessor()
	base := &configModels.DashConfig{
		Port:    "8080",
		Origin:  "http://localhost:3000",
		Headers: []string{"Content-Type"},
		Plugins: configModels.Plugins{"OAuth2"},
	}
	override := &configModels.DashConfig{
		Port:    "9090",                             // Override
		Front:   "/dist",                            // New value
		Plugins: configModels.Plugins{"Kubernetes"}, // Override
	}

	// Act
	merged := processor.MergeConfigs(base, override)

	// Assert
	assert.Equal(t, "9090", merged.Port)                      // Overridden
	assert.Equal(t, "http://localhost:3000", merged.Origin)   // From base
	assert.Equal(t, "/dist", merged.Front)                    // From override
	assert.Equal(t, []string{"Content-Type"}, merged.Headers) // From base
	assert.True(t, merged.Plugins.Has("Kubernetes"))          // From override
	assert.False(t, merged.Plugins.Has("OAuth2"))             // Replaced
}
