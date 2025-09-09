package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugins_Has_WithExistingPluginExactCase_ReturnsTrue(t *testing.T) {
	// Arrange
	plugins := Plugins{"Auth", "Kubernetes", "AWS"}
	pluginName := "Auth"

	// Act
	result := plugins.Has(pluginName)

	// Assert
	assert.True(t, result)
}

func TestPlugins_Has_WithExistingPluginDifferentCase_ReturnsTrue(t *testing.T) {
	// Arrange
	plugins := Plugins{"Auth", "Kubernetes", "AWS"}
	pluginName := "auth"

	// Act
	result := plugins.Has(pluginName)

	// Assert
	assert.True(t, result)
}

func TestPlugins_Has_WithNonExistingPlugin_ReturnsFalse(t *testing.T) {
	// Arrange
	plugins := Plugins{"Auth", "Kubernetes", "AWS"}
	pluginName := "NonExistent"

	// Act
	result := plugins.Has(pluginName)

	// Assert
	assert.False(t, result)
}

func TestPlugins_Has_WithEmptyPluginName_ReturnsFalse(t *testing.T) {
	// Arrange
	plugins := Plugins{"Auth", "Kubernetes", "AWS"}
	pluginName := ""

	// Act
	result := plugins.Has(pluginName)

	// Assert
	assert.False(t, result)
}

func TestPlugins_Add_WithNewPlugin_AddsPlugin(t *testing.T) {
	// Arrange
	plugins := Plugins{"Auth"}

	// Act
	plugins.Add("Kubernetes")

	// Assert
	assert.True(t, plugins.Has("Kubernetes"))
	assert.Equal(t, 2, plugins.Count())
}

func TestPlugins_Add_WithDuplicatePlugin_DoesNotDuplicate(t *testing.T) {
	// Arrange
	plugins := Plugins{"Auth"}

	// Act
	plugins.Add("Auth")

	// Assert
	assert.Equal(t, 1, plugins.Count())
}

func TestPlugins_Add_WithDifferentCase_DoesNotDuplicate(t *testing.T) {
	// Arrange
	plugins := Plugins{"Auth"}

	// Act
	plugins.Add("auth")

	// Assert
	assert.Equal(t, 1, plugins.Count())
}

func TestPlugins_Remove_WithExistingPlugin_RemovesPlugin(t *testing.T) {
	// Arrange
	plugins := Plugins{"Auth", "Kubernetes", "AWS"}

	// Act
	plugins.Remove("Kubernetes")

	// Assert
	assert.False(t, plugins.Has("Kubernetes"))
	assert.Equal(t, 2, plugins.Count())
}

func TestPlugins_Remove_WithDifferentCase_RemovesPlugin(t *testing.T) {
	// Arrange
	plugins := Plugins{"Auth", "Kubernetes", "AWS"}

	// Act
	plugins.Remove("auth")

	// Assert
	assert.False(t, plugins.Has("Auth"))
	assert.Equal(t, 2, plugins.Count())
}

func TestPlugins_Remove_WithNonExistingPlugin_DoesNotError(t *testing.T) {
	// Arrange
	plugins := Plugins{"Auth", "Kubernetes", "AWS"}

	// Act
	plugins.Remove("NonExistent")

	// Assert
	assert.Equal(t, 3, plugins.Count())
}

func TestDashConfig_Validate_WithValidConfig_ReturnsNoError(t *testing.T) {
	// Arrange
	config := DashConfig{
		Port:   "8080",
		Origin: "http://localhost:3000",
	}

	// Act
	err := config.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestDashConfig_Validate_WithMissingPort_ReturnsError(t *testing.T) {
	// Arrange
	config := DashConfig{
		Origin: "http://localhost:3000",
	}

	// Act
	err := config.Validate()

	// Assert
	assert.Error(t, err)
}

func TestDashConfig_Validate_WithMissingOrigin_ReturnsError(t *testing.T) {
	// Arrange
	config := DashConfig{
		Port: "8080",
	}

	// Act
	err := config.Validate()

	// Assert
	assert.Error(t, err)
}

func TestDashConfig_Validate_WithEmptyConfig_ReturnsError(t *testing.T) {
	// Arrange
	config := DashConfig{}

	// Act
	err := config.Validate()

	// Assert
	assert.Error(t, err)
}

func TestDashConfig_GetPort_WithPortSet_ReturnsSetPort(t *testing.T) {
	// Arrange
	config := DashConfig{Port: "9090"}

	// Act
	result := config.GetPort()

	// Assert
	assert.Equal(t, "9090", result)
}

func TestDashConfig_GetPort_WithoutPortSet_ReturnsDefaultPort(t *testing.T) {
	// Arrange
	config := DashConfig{}

	// Act
	result := config.GetPort()

	// Assert
	assert.Equal(t, "8080", result) // Default
}

func TestDashConfig_GetOrigin_WithOriginSet_ReturnsSetOrigin(t *testing.T) {
	// Arrange
	config := DashConfig{Origin: "http://custom.local"}

	// Act
	result := config.GetOrigin()

	// Assert
	assert.Equal(t, "http://custom.local", result)
}

func TestDashConfig_GetOrigin_WithoutOriginSet_ReturnsDefaultOrigin(t *testing.T) {
	// Arrange
	config := DashConfig{}

	// Act
	result := config.GetOrigin()

	// Assert
	assert.Equal(t, "http://localhost:3000", result) // Default
}

func TestDashConfig_GetHeaders_WithHeadersSet_ReturnsSetHeaders(t *testing.T) {
	// Arrange
	config := DashConfig{Headers: []string{"Custom-Header"}}

	// Act
	result := config.GetHeaders()

	// Assert
	assert.Equal(t, []string{"Custom-Header"}, result)
}

func TestDashConfig_GetHeaders_WithoutHeadersSet_ReturnsDefaultHeaders(t *testing.T) {
	// Arrange
	config := DashConfig{}

	// Act
	result := config.GetHeaders()

	// Assert
	assert.Equal(t, []string{"Content-Type", "Authorization"}, result) // Default
}

func TestDashConfig_IsPluginEnabled_WithEnabledPluginExactCase_ReturnsTrue(t *testing.T) {
	// Arrange
	config := DashConfig{
		Plugins: Plugins{"Auth", "Kubernetes"},
	}

	// Act
	result := config.IsPluginEnabled("Auth")

	// Assert
	assert.True(t, result)
}

func TestDashConfig_IsPluginEnabled_WithEnabledPluginDifferentCase_ReturnsTrue(t *testing.T) {
	// Arrange
	config := DashConfig{
		Plugins: Plugins{"Auth", "Kubernetes"},
	}

	// Act
	result := config.IsPluginEnabled("auth") // Case insensitive

	// Assert
	assert.True(t, result)
}

func TestDashConfig_IsPluginEnabled_WithDisabledPlugin_ReturnsFalse(t *testing.T) {
	// Arrange
	config := DashConfig{
		Plugins: Plugins{"Auth", "Kubernetes"},
	}

	// Act
	result := config.IsPluginEnabled("AWS")

	// Assert
	assert.False(t, result)
}

func TestDashConfig_IsPluginEnabled_WithEmptyPluginName_ReturnsFalse(t *testing.T) {
	// Arrange
	config := DashConfig{
		Plugins: Plugins{"Auth", "Kubernetes"},
	}

	// Act
	result := config.IsPluginEnabled("")

	// Assert
	assert.False(t, result)
}
