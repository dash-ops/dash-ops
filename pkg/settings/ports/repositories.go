package ports

import (
	"context"

	"github.com/dash-ops/dash-ops/pkg/settings/models"
)

// ConfigRepository handles configuration file operations
type ConfigRepository interface {
	// SaveConfig saves configuration to file
	SaveConfig(ctx context.Context, config *models.DashConfig, filePath string) error

	// LoadConfig loads configuration from file
	LoadConfig(ctx context.Context, filePath string) (*models.DashConfig, error)

	// ConfigExists checks if configuration file exists
	ConfigExists(ctx context.Context, filePath string) (bool, error)

	// GetDefaultConfigPath returns the default configuration file path
	GetDefaultConfigPath() string
}

// ConfigProvider exposes read-only access to the current configuration state.
type ConfigProvider interface {
	GetConfig() *models.DashConfig
	GetPlugins() models.Plugins
}

// ConfigCache allows read/write access to the configuration state.
type ConfigCache interface {
	ConfigProvider
	SetConfig(config *models.DashConfig)
}

// ConfigStatusRepository handles configuration status checks
type ConfigStatusRepository interface {
	// GetSetupStatus returns the current setup status
	GetSetupStatus(ctx context.Context) (*models.SetupStatus, error)

	// HasPluginsConfigured checks if any plugins are configured
	HasPluginsConfigured(ctx context.Context) (bool, error)

	// HasAuthConfigured checks if auth is configured
	HasAuthConfigured(ctx context.Context) (bool, error)
}
