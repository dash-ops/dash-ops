package config

import (
	"os"
	"path/filepath"
	"testing"

	configModels "github.com/dash-ops/dash-ops/pkg/config/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigProcessor_ParseFromBytes(t *testing.T) {
	processor := NewConfigProcessor()

	tests := []struct {
		name           string
		configData     string
		expectError    bool
		expectedValues map[string]interface{}
	}{
		{
			name: "valid config",
			configData: `port: 8080
origin: http://localhost:3000
headers: 
  - "Content-Type"
  - "Authorization"
plugins:
  - "OAuth2"
  - "Kubernetes"`,
			expectError: false,
			expectedValues: map[string]interface{}{
				"port":   "8080",
				"origin": "http://localhost:3000",
			},
		},
		{
			name: "config with environment variables",
			configData: `port: ${PORT:-8080}
origin: ${ORIGIN:-http://localhost:3000}
plugins:
  - "OAuth2"`,
			expectError: false,
			expectedValues: map[string]interface{}{
				"port":   "8080",                  // Default value
				"origin": "http://localhost:3000", // Default value
			},
		},
		{
			name:        "empty config",
			configData:  ``,
			expectError: true,
		},
		{
			name: "invalid YAML",
			configData: `port: 8080
origin: [invalid yaml`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := processor.ParseFromBytes([]byte(tt.configData))

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				require.NoError(t, err)
				require.NotNil(t, config)

				if expectedPort, ok := tt.expectedValues["port"]; ok {
					assert.Equal(t, expectedPort, config.Port)
				}
				if expectedOrigin, ok := tt.expectedValues["origin"]; ok {
					assert.Equal(t, expectedOrigin, config.Origin)
				}
			}
		})
	}
}

func TestConfigProcessor_LoadFromFile(t *testing.T) {
	processor := NewConfigProcessor()

	// Create a temporary config file
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

	// Test loading from file
	config, err := processor.LoadFromFile(configFile)
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "9090", config.Port)
	assert.Equal(t, "http://test.local", config.Origin)
	assert.Equal(t, []string{"Content-Type"}, config.Headers)
	assert.True(t, config.Plugins.Has("TestPlugin"))
}

func TestConfigProcessor_LoadFromFile_NonExistent(t *testing.T) {
	processor := NewConfigProcessor()

	config, err := processor.LoadFromFile("/non/existent/file.yaml")
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestConfigProcessor_GetConfigFilePath(t *testing.T) {
	processor := NewConfigProcessor()

	// Test default path
	path := processor.GetConfigFilePath()
	assert.Equal(t, "./dash-ops.yaml", path)

	// Test environment variable
	os.Setenv("DASH_CONFIG", "/custom/path/config.yaml")
	defer os.Unsetenv("DASH_CONFIG")

	path = processor.GetConfigFilePath()
	assert.Equal(t, "/custom/path/config.yaml", path)
}

func TestConfigProcessor_MergeConfigs(t *testing.T) {
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

	merged := processor.MergeConfigs(base, override)

	assert.Equal(t, "9090", merged.Port)                      // Overridden
	assert.Equal(t, "http://localhost:3000", merged.Origin)   // From base
	assert.Equal(t, "/dist", merged.Front)                    // From override
	assert.Equal(t, []string{"Content-Type"}, merged.Headers) // From base
	assert.True(t, merged.Plugins.Has("Kubernetes"))          // From override
	assert.False(t, merged.Plugins.Has("OAuth2"))             // Replaced
}
