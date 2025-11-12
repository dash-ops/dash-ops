package logic

import (
	"os"
	"path/filepath"
	"testing"

	settingsModels "github.com/dash-ops/dash-ops/pkg/settings/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigProcessor_ParseFromBytes_WithValidConfig_ReturnsConfig(t *testing.T) {
	processor := NewConfigProcessor()
	configData := `port: 8080
origin: http://localhost:5173
headers:
  - "Content-Type"
  - "Authorization"
plugins:
  - "Auth"
  - "Kubernetes"`

	config, err := processor.ParseFromBytes([]byte(configData))

	require.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "http://localhost:5173", config.Origin)
	assert.True(t, config.Plugins.Has("Auth"))
}

func TestConfigProcessor_ParseFromBytes_WithEnvironmentVariables_ReturnsConfigWithDefaults(t *testing.T) {
	processor := NewConfigProcessor()
	configData := `port: ${PORT:-8080}
origin: ${ORIGIN:-http://localhost:5173}
plugins:
  - "Auth"`

	config, err := processor.ParseFromBytes([]byte(configData))

	require.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, "8080", config.Port)                    // Default value
	assert.Equal(t, "http://localhost:5173", config.Origin) // Default value
	assert.True(t, config.Plugins.Has("Auth"))
}

func TestConfigProcessor_ParseFromBytes_WithEmptyConfig_ReturnsError(t *testing.T) {
	processor := NewConfigProcessor()

	config, err := processor.ParseFromBytes([]byte(""))

	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestConfigProcessor_ParseFromBytes_WithInvalidYAML_ReturnsError(t *testing.T) {
	processor := NewConfigProcessor()
	configData := `port: 8080
origin: [invalid yaml`

	config, err := processor.ParseFromBytes([]byte(configData))

	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestConfigProcessor_LoadFromFile_WithValidFile_ReturnsConfig(t *testing.T) {
	processor := NewConfigProcessor()
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test-config.yaml")

	configContent := `port: 9090
origin: http://test.local
headers:
  - "Content-Type"
plugins:
  - "TestPlugin"`

	err := os.WriteFile(configFile, []byte(configContent), 0o644)
	require.NoError(t, err)

	config, err := processor.LoadFromFile(configFile)

	require.NoError(t, err)
	require.NotNil(t, config)
	assert.Equal(t, "9090", config.Port)
	assert.Equal(t, "http://test.local", config.Origin)
	assert.Equal(t, []string{"Content-Type"}, config.Headers)
	assert.True(t, config.Plugins.Has("TestPlugin"))
}

func TestConfigProcessor_LoadFromFile_WithNonExistentFile_ReturnsError(t *testing.T) {
	processor := NewConfigProcessor()

	config, err := processor.LoadFromFile("/non/existent/file.yaml")

	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestConfigProcessor_GetConfigFilePath_DefaultPath(t *testing.T) {
	processor := NewConfigProcessor()

	path := processor.GetConfigFilePath()

	assert.Equal(t, "./dash-ops.yaml", path)
}

func TestConfigProcessor_GetConfigFilePath_WithEnvironmentVariable(t *testing.T) {
	processor := NewConfigProcessor()
	os.Setenv("DASH_CONFIG", "/custom/path/config.yaml")
	defer os.Unsetenv("DASH_CONFIG")

	path := processor.GetConfigFilePath()

	assert.Equal(t, "/custom/path/config.yaml", path)
}

func TestConfigProcessor_MergeConfigs_WithBaseAndOverride(t *testing.T) {
	processor := NewConfigProcessor()
	base := &settingsModels.DashConfig{
		Port:    "8080",
		Origin:  "http://localhost:5173",
		Headers: []string{"Content-Type"},
		Plugins: settingsModels.Plugins{"Auth"},
	}
	override := &settingsModels.DashConfig{
		Port:    "9090",                               // Override
		Front:   "/dist",                              // New value
		Plugins: settingsModels.Plugins{"Kubernetes"}, // Override
	}

	merged := processor.MergeConfigs(base, override)

	assert.Equal(t, "9090", merged.Port)                      // Overridden
	assert.Equal(t, "http://localhost:5173", merged.Origin)   // From base
	assert.Equal(t, "/dist", merged.Front)                    // From override
	assert.Equal(t, []string{"Content-Type"}, merged.Headers) // From base
	assert.True(t, merged.Plugins.Has("Kubernetes"))          // From override
	assert.False(t, merged.Plugins.Has("Auth"))               // Replaced
}
