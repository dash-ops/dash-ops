package repositories

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dash-ops/dash-ops/pkg/settings/logic"
	"github.com/dash-ops/dash-ops/pkg/settings/models"
	"github.com/dash-ops/dash-ops/pkg/settings/ports"
)

// FileRepository implements ConfigRepository using local filesystem
type FileRepository struct {
	yamlProcessor *logic.YAMLProcessor
	defaultPath   string
}

// NewFileRepository creates a new file repository
func NewFileRepository(defaultPath string) ports.ConfigRepository {
	return &FileRepository{
		yamlProcessor: logic.NewYAMLProcessor(),
		defaultPath:   defaultPath,
	}
}

// SaveConfig saves configuration to file
func (fr *FileRepository) SaveConfig(ctx context.Context, config *models.DashConfig, filePath string) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if filePath == "" {
		filePath = fr.defaultPath
	}

	// Generate YAML
	yamlData, err := fr.yamlProcessor.GenerateYAML(config)
	if err != nil {
		return fmt.Errorf("failed to generate YAML: %w", err)
	}

	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(absPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadConfig loads configuration from file
func (fr *FileRepository) LoadConfig(ctx context.Context, filePath string) (*models.DashConfig, error) {
	if filePath == "" {
		filePath = fr.defaultPath
	}

	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Read file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	config, err := fr.yamlProcessor.ParseYAML(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return config, nil
}

// ConfigExists checks if configuration file exists
func (fr *FileRepository) ConfigExists(ctx context.Context, filePath string) (bool, error) {
	if filePath == "" {
		filePath = fr.defaultPath
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path: %w", err)
	}

	_, err = os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

// GetDefaultConfigPath returns the default configuration file path
func (fr *FileRepository) GetDefaultConfigPath() string {
	return fr.defaultPath
}
