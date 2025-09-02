package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog-new/models"
)

// FilesystemRepository implements ServiceRepository for filesystem storage
type FilesystemRepository struct {
	directory string
}

// NewFilesystemRepository creates a new filesystem repository
func NewFilesystemRepository(directory string) (*FilesystemRepository, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create services directory: %w", err)
	}

	return &FilesystemRepository{
		directory: directory,
	}, nil
}

// Create creates a new service
func (fr *FilesystemRepository) Create(ctx context.Context, service *scModels.Service) (*scModels.Service, error) {
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}

	filePath := fr.getServiceFilePath(service.Metadata.Name)

	// Check if service already exists
	if _, err := os.Stat(filePath); err == nil {
		return nil, fmt.Errorf("service '%s' already exists", service.Metadata.Name)
	}

	// Set defaults
	service.SetDefaults()

	// Marshal to YAML
	data, err := yaml.Marshal(service)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal service to YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write service file: %w", err)
	}

	return service, nil
}

// GetByName retrieves a service by name
func (fr *FilesystemRepository) GetByName(ctx context.Context, name string) (*scModels.Service, error) {
	if name == "" {
		return nil, fmt.Errorf("service name cannot be empty")
	}

	filePath := fr.getServiceFilePath(name)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("service '%s' not found", name)
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read service file: %w", err)
	}

	// Unmarshal from YAML
	var service scModels.Service
	if err := yaml.Unmarshal(data, &service); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service from YAML: %w", err)
	}

	return &service, nil
}

// Update updates an existing service
func (fr *FilesystemRepository) Update(ctx context.Context, service *scModels.Service) (*scModels.Service, error) {
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}

	filePath := fr.getServiceFilePath(service.Metadata.Name)

	// Check if service exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("service '%s' not found", service.Metadata.Name)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(service)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal service to YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write service file: %w", err)
	}

	return service, nil
}

// Delete deletes a service
func (fr *FilesystemRepository) Delete(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	filePath := fr.getServiceFilePath(name)

	// Check if service exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("service '%s' not found", name)
	}

	// Remove file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete service file: %w", err)
	}

	return nil
}

// List lists all services with optional filtering
func (fr *FilesystemRepository) List(ctx context.Context, filter *scModels.ServiceFilter) ([]scModels.Service, error) {
	var services []scModels.Service

	// Read directory
	files, err := os.ReadDir(fr.directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read services directory: %w", err)
	}

	// Process each YAML file
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		// Get service name from filename
		serviceName := strings.TrimSuffix(file.Name(), ".yaml")

		// Get service
		service, err := fr.GetByName(ctx, serviceName)
		if err != nil {
			// Skip invalid services but continue with others
			continue
		}

		services = append(services, *service)
	}

	return services, nil
}

// Exists checks if a service exists
func (fr *FilesystemRepository) Exists(ctx context.Context, name string) (bool, error) {
	if name == "" {
		return false, fmt.Errorf("service name cannot be empty")
	}

	filePath := fr.getServiceFilePath(name)
	_, err := os.Stat(filePath)
	return err == nil, nil
}

// ListByTeam lists services owned by a specific team
func (fr *FilesystemRepository) ListByTeam(ctx context.Context, team string) ([]scModels.Service, error) {
	allServices, err := fr.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	var teamServices []scModels.Service
	for _, service := range allServices {
		if strings.EqualFold(service.Spec.Team.GitHubTeam, team) {
			teamServices = append(teamServices, service)
		}
	}

	return teamServices, nil
}

// ListByTier lists services of a specific tier
func (fr *FilesystemRepository) ListByTier(ctx context.Context, tier scModels.ServiceTier) ([]scModels.Service, error) {
	allServices, err := fr.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	var tierServices []scModels.Service
	for _, service := range allServices {
		if service.Metadata.Tier == tier {
			tierServices = append(tierServices, service)
		}
	}

	return tierServices, nil
}

// Search searches services by text query
func (fr *FilesystemRepository) Search(ctx context.Context, query string, limit int) ([]scModels.Service, error) {
	allServices, err := fr.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var matchingServices []scModels.Service

	for _, service := range allServices {
		if fr.matchesQuery(service, query) {
			matchingServices = append(matchingServices, service)
			if limit > 0 && len(matchingServices) >= limit {
				break
			}
		}
	}

	return matchingServices, nil
}

// getServiceFilePath returns the full file path for a service
func (fr *FilesystemRepository) getServiceFilePath(serviceName string) string {
	return filepath.Join(fr.directory, serviceName+".yaml")
}

// matchesQuery checks if service matches search query
func (fr *FilesystemRepository) matchesQuery(service scModels.Service, query string) bool {
	searchFields := []string{
		strings.ToLower(service.Metadata.Name),
		strings.ToLower(service.Spec.Description),
		strings.ToLower(service.Spec.Team.GitHubTeam),
		strings.ToLower(service.Spec.Technology.Language),
		strings.ToLower(service.Spec.Technology.Framework),
	}

	for _, field := range searchFields {
		if strings.Contains(field, query) {
			return true
		}
	}

	return false
}
