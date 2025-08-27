package servicecatalog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// FilesystemProvider implements StorageProvider interface for local filesystem storage
type FilesystemProvider struct {
	directory string
}

// NewFilesystemProvider creates a new filesystem storage provider
func NewFilesystemProvider(directory string) (*FilesystemProvider, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create services directory: %w", err)
	}

	return &FilesystemProvider{
		directory: directory,
	}, nil
}

// CreateService creates a new service definition file
func (fp *FilesystemProvider) CreateService(service *Service) error {
	if err := fp.validateService(service); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	filePath := fp.getServiceFilePath(service.Metadata.Name)

	// Check if service already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("service '%s' already exists", service.Metadata.Name)
	}

	// Set default values
	if service.APIVersion == "" {
		service.APIVersion = "v1"
	}
	if service.Kind == "" {
		service.Kind = "Service"
	}

	// Set audit metadata
	now := time.Now().Format(time.RFC3339)
	service.Metadata.CreatedAt = now
	service.Metadata.UpdatedAt = now
	service.Metadata.Version = 1

	// Marshal to YAML
	data, err := yaml.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service to YAML: %w", err)
	}

	// Write to file
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	return nil
}

// GetService retrieves a service definition by name
func (fp *FilesystemProvider) GetService(name string) (*Service, error) {
	if err := fp.validateServiceName(name); err != nil {
		return nil, err
	}

	filePath := fp.getServiceFilePath(name)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("service '%s' not found", name)
	}

	// Read file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read service file: %w", err)
	}

	// Unmarshal from YAML
	var service Service
	if err := yaml.Unmarshal(data, &service); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service from YAML: %w", err)
	}

	return &service, nil
}

// UpdateService updates an existing service definition
func (fp *FilesystemProvider) UpdateService(service *Service) error {
	if err := fp.validateService(service); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	filePath := fp.getServiceFilePath(service.Metadata.Name)

	// Check if service exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("service '%s' not found", service.Metadata.Name)
	}

	// Get existing service to preserve creation info and increment version
	existingService, err := fp.GetService(service.Metadata.Name)
	if err != nil {
		return fmt.Errorf("failed to get existing service: %w", err)
	}

	// Preserve creation metadata and increment version
	service.Metadata.CreatedAt = existingService.Metadata.CreatedAt
	service.Metadata.CreatedBy = existingService.Metadata.CreatedBy
	service.Metadata.UpdatedAt = time.Now().Format(time.RFC3339)
	service.Metadata.Version = existingService.Metadata.Version + 1

	// Marshal to YAML
	data, err := yaml.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service to YAML: %w", err)
	}

	// Write to file
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	return nil
}

// DeleteService removes a service definition file
func (fp *FilesystemProvider) DeleteService(name string) error {
	if err := fp.validateServiceName(name); err != nil {
		return err
	}

	filePath := fp.getServiceFilePath(name)

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

// ListServices returns all service definitions
func (fp *FilesystemProvider) ListServices() ([]Service, error) {
	var services []Service

	// Read directory
	files, err := ioutil.ReadDir(fp.directory)
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
		service, err := fp.GetService(serviceName)
		if err != nil {
			// Skip invalid services but continue with others
			continue
		}

		services = append(services, *service)
	}

	return services, nil
}

// ServiceExists checks if a service definition file exists
func (fp *FilesystemProvider) ServiceExists(name string) bool {
	if err := fp.validateServiceName(name); err != nil {
		return false
	}

	filePath := fp.getServiceFilePath(name)
	_, err := os.Stat(filePath)
	return err == nil
}

// GetDirectory returns the storage directory path
func (fp *FilesystemProvider) GetDirectory() string {
	return fp.directory
}

// getServiceFilePath returns the full file path for a service
func (fp *FilesystemProvider) getServiceFilePath(serviceName string) string {
	return filepath.Join(fp.directory, serviceName+".yaml")
}

// validateServiceName validates service name format
func (fp *FilesystemProvider) validateServiceName(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	// Check for invalid characters in filename
	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		return fmt.Errorf("service name contains invalid characters")
	}

	// Check length
	if len(name) > 100 {
		return fmt.Errorf("service name too long (max 100 characters)")
	}

	return nil
}

// validateService validates a complete service definition
func (fp *FilesystemProvider) validateService(service *Service) error {
	if service == nil {
		return fmt.Errorf("service cannot be nil")
	}

	// Validate metadata
	if err := fp.validateServiceName(service.Metadata.Name); err != nil {
		return fmt.Errorf("invalid service name: %w", err)
	}

	// Validate tier
	validTiers := map[string]bool{
		"TIER-1": true,
		"TIER-2": true,
		"TIER-3": true,
	}
	if !validTiers[service.Metadata.Tier] {
		return fmt.Errorf("invalid tier '%s', must be TIER-1, TIER-2, or TIER-3", service.Metadata.Tier)
	}

	// Validate spec
	if service.Spec.Description == "" {
		return fmt.Errorf("service description is required")
	}

	// Validate team
	if service.Spec.Team.GitHubTeam == "" {
		return fmt.Errorf("github_team is required")
	}

	// Validate Kubernetes environments if present
	if service.Spec.Kubernetes != nil {
		for i, env := range service.Spec.Kubernetes.Environments {
			if env.Name == "" {
				return fmt.Errorf("environment[%d].name is required", i)
			}
			if env.Context == "" {
				return fmt.Errorf("environment[%d].context is required", i)
			}
			if env.Namespace == "" {
				return fmt.Errorf("environment[%d].namespace is required", i)
			}

			// Validate deployments
			for j, deploy := range env.Resources.Deployments {
				if deploy.Name == "" {
					return fmt.Errorf("environment[%d].deployments[%d].name is required", i, j)
				}
				if deploy.Replicas <= 0 {
					return fmt.Errorf("environment[%d].deployments[%d].replicas must be greater than 0", i, j)
				}

				// Validate resource specifications
				if deploy.Resources.Requests.CPU == "" || deploy.Resources.Requests.Memory == "" {
					return fmt.Errorf("environment[%d].deployments[%d] requires CPU and memory requests", i, j)
				}
				if deploy.Resources.Limits.CPU == "" || deploy.Resources.Limits.Memory == "" {
					return fmt.Errorf("environment[%d].deployments[%d] requires CPU and memory limits", i, j)
				}
			}
		}
	}

	return nil
}

// GetServicesDirectory returns the services directory path (for external use)
func (fp *FilesystemProvider) GetServicesDirectory() string {
	return fp.directory
}
