package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugins_Has(t *testing.T) {
	plugins := Plugins{"OAuth2", "Kubernetes", "AWS"}

	tests := []struct {
		name       string
		pluginName string
		expected   bool
	}{
		{
			name:       "existing plugin exact case",
			pluginName: "OAuth2",
			expected:   true,
		},
		{
			name:       "existing plugin different case",
			pluginName: "oauth2",
			expected:   true,
		},
		{
			name:       "non-existing plugin",
			pluginName: "NonExistent",
			expected:   false,
		},
		{
			name:       "empty plugin name",
			pluginName: "",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := plugins.Has(tt.pluginName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPlugins_Add(t *testing.T) {
	plugins := Plugins{"OAuth2"}

	// Add new plugin
	plugins.Add("Kubernetes")
	assert.True(t, plugins.Has("Kubernetes"))
	assert.Equal(t, 2, plugins.Count())

	// Add duplicate plugin (should not duplicate)
	plugins.Add("OAuth2")
	assert.Equal(t, 2, plugins.Count())

	// Add with different case (should not duplicate)
	plugins.Add("oauth2")
	assert.Equal(t, 2, plugins.Count())
}

func TestPlugins_Remove(t *testing.T) {
	plugins := Plugins{"OAuth2", "Kubernetes", "AWS"}

	// Remove existing plugin
	plugins.Remove("Kubernetes")
	assert.False(t, plugins.Has("Kubernetes"))
	assert.Equal(t, 2, plugins.Count())

	// Remove with different case
	plugins.Remove("oauth2")
	assert.False(t, plugins.Has("OAuth2"))
	assert.Equal(t, 1, plugins.Count())

	// Remove non-existing plugin (should not error)
	plugins.Remove("NonExistent")
	assert.Equal(t, 1, plugins.Count())
}

func TestDashConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      DashConfig
		expectError bool
	}{
		{
			name: "valid config",
			config: DashConfig{
				Port:   "8080",
				Origin: "http://localhost:3000",
			},
			expectError: false,
		},
		{
			name: "missing port",
			config: DashConfig{
				Origin: "http://localhost:3000",
			},
			expectError: true,
		},
		{
			name: "missing origin",
			config: DashConfig{
				Port: "8080",
			},
			expectError: true,
		},
		{
			name:        "empty config",
			config:      DashConfig{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDashConfig_GetPort(t *testing.T) {
	tests := []struct {
		name     string
		config   DashConfig
		expected string
	}{
		{
			name:     "with port set",
			config:   DashConfig{Port: "9090"},
			expected: "9090",
		},
		{
			name:     "without port set",
			config:   DashConfig{},
			expected: "8080", // Default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetPort()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDashConfig_GetOrigin(t *testing.T) {
	tests := []struct {
		name     string
		config   DashConfig
		expected string
	}{
		{
			name:     "with origin set",
			config:   DashConfig{Origin: "http://custom.local"},
			expected: "http://custom.local",
		},
		{
			name:     "without origin set",
			config:   DashConfig{},
			expected: "http://localhost:3000", // Default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetOrigin()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDashConfig_GetHeaders(t *testing.T) {
	tests := []struct {
		name     string
		config   DashConfig
		expected []string
	}{
		{
			name:     "with headers set",
			config:   DashConfig{Headers: []string{"Custom-Header"}},
			expected: []string{"Custom-Header"},
		},
		{
			name:     "without headers set",
			config:   DashConfig{},
			expected: []string{"Content-Type", "Authorization"}, // Default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetHeaders()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDashConfig_IsPluginEnabled(t *testing.T) {
	config := DashConfig{
		Plugins: Plugins{"OAuth2", "Kubernetes"},
	}

	assert.True(t, config.IsPluginEnabled("OAuth2"))
	assert.True(t, config.IsPluginEnabled("oauth2")) // Case insensitive
	assert.False(t, config.IsPluginEnabled("AWS"))
	assert.False(t, config.IsPluginEnabled(""))
}
