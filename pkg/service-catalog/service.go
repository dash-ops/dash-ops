package servicecatalog

import (
	"fmt"
	"strings"
	"time"
)

// ServiceCatalog represents the main service catalog manager
type ServiceCatalog struct {
	storage       StorageProvider
	gitVersioning *GitVersioning
	config        *Config
}

// NewServiceCatalog creates a new service catalog instance
func NewServiceCatalog(config *Config) (*ServiceCatalog, error) {
	// Initialize storage provider
	storage, err := newStorageProvider(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage provider: %w", err)
	}

	// Initialize git versioning for filesystem provider
	var gitVersioning *GitVersioning
	if config.Storage.Provider == "filesystem" {
		// Get the directory path from filesystem provider
		if fsProvider, ok := storage.(*FilesystemProvider); ok {
			gitVersioning = NewGitVersioning(fsProvider.GetDirectory())

			// Validate git installation
			if err := ValidateGitInstallation(); err != nil {
				return nil, fmt.Errorf("git validation failed: %w", err)
			}

			// Initialize repository
			if err := gitVersioning.InitializeRepository(); err != nil {
				return nil, fmt.Errorf("failed to initialize git repository: %w", err)
			}

			// Setup git config
			if err := gitVersioning.SetupGitConfig(); err != nil {
				return nil, fmt.Errorf("failed to setup git config: %w", err)
			}
		}
	}

	return &ServiceCatalog{
		storage:       storage,
		gitVersioning: gitVersioning,
		config:        config,
	}, nil
}

// CreateService creates a new service definition
func (sc *ServiceCatalog) CreateService(service *Service, user *UserContext) error {
	// Set audit metadata
	if user != nil {
		service.Metadata.CreatedBy = user.Username
		service.Metadata.UpdatedBy = user.Username
	}

	// Create service in storage
	if err := sc.storage.CreateService(service); err != nil {
		return fmt.Errorf("failed to create service in storage: %w", err)
	}

	// Commit to git if versioning is enabled
	if sc.gitVersioning != nil {
		if err := sc.gitVersioning.CommitServiceChange(service, user, "create"); err != nil {
			// Log warning but don't fail the operation
			fmt.Printf("Warning: failed to commit service creation: %v\n", err)
		}
	}

	return nil
}

// GetService retrieves a service definition by name
func (sc *ServiceCatalog) GetService(name string) (*Service, error) {
	service, err := sc.storage.GetService(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get service from storage: %w", err)
	}

	// TODO: Resolve GitHub team information at runtime
	// This will be implemented when we add GitHub integration

	return service, nil
}

// UpdateService updates an existing service definition
func (sc *ServiceCatalog) UpdateService(service *Service, user *UserContext) error {
	// Set audit metadata
	if user != nil {
		service.Metadata.UpdatedBy = user.Username
	}

	// Update service in storage
	if err := sc.storage.UpdateService(service); err != nil {
		return fmt.Errorf("failed to update service in storage: %w", err)
	}

	// Commit to git if versioning is enabled
	if sc.gitVersioning != nil {
		if err := sc.gitVersioning.CommitServiceChange(service, user, "update"); err != nil {
			// Log warning but don't fail the operation
			fmt.Printf("Warning: failed to commit service update: %v\n", err)
		}
	}

	return nil
}

// DeleteService removes a service definition
func (sc *ServiceCatalog) DeleteService(name string, user *UserContext) error {
	// Check if service exists
	if !sc.storage.ServiceExists(name) {
		return fmt.Errorf("service '%s' not found", name)
	}

	// Delete service from storage
	if err := sc.storage.DeleteService(name); err != nil {
		return fmt.Errorf("failed to delete service from storage: %w", err)
	}

	// Commit to git if versioning is enabled
	if sc.gitVersioning != nil {
		if err := sc.gitVersioning.CommitServiceDeletion(name, user); err != nil {
			// Log warning but don't fail the operation
			fmt.Printf("Warning: failed to commit service deletion: %v\n", err)
		}
	}

	return nil
}

// ListServices returns all service definitions
func (sc *ServiceCatalog) ListServices() (*ServiceList, error) {
	services, err := sc.storage.ListServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list services from storage: %w", err)
	}

	// TODO: Filter services by user's GitHub teams
	// This will be implemented when we add GitHub integration

	return &ServiceList{
		Services: services,
		Total:    len(services),
	}, nil
}

// ListServicesByTeam returns services filtered by GitHub team
func (sc *ServiceCatalog) ListServicesByTeam(teamName string) (*ServiceList, error) {
	services, err := sc.storage.ListServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list services from storage: %w", err)
	}

	// Filter by team
	var filteredServices []Service
	for _, service := range services {
		if service.Spec.Team.GitHubTeam == teamName {
			filteredServices = append(filteredServices, service)
		}
	}

	return &ServiceList{
		Services: filteredServices,
		Total:    len(filteredServices),
	}, nil
}

// ListServicesByTier returns services filtered by tier
func (sc *ServiceCatalog) ListServicesByTier(tier string) (*ServiceList, error) {
	services, err := sc.storage.ListServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list services from storage: %w", err)
	}

	// Validate tier
	validTiers := map[string]bool{
		"TIER-1": true,
		"TIER-2": true,
		"TIER-3": true,
	}
	if !validTiers[tier] {
		return nil, fmt.Errorf("invalid tier '%s'", tier)
	}

	// Filter by tier
	var filteredServices []Service
	for _, service := range services {
		if service.Metadata.Tier == tier {
			filteredServices = append(filteredServices, service)
		}
	}

	return &ServiceList{
		Services: filteredServices,
		Total:    len(filteredServices),
	}, nil
}

// SearchServices searches services by name or description
func (sc *ServiceCatalog) SearchServices(query string) (*ServiceList, error) {
	services, err := sc.storage.ListServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list services from storage: %w", err)
	}

	// Search in name and description (case insensitive)
	query = strings.ToLower(query)
	var filteredServices []Service

	for _, service := range services {
		if strings.Contains(strings.ToLower(service.Metadata.Name), query) ||
			strings.Contains(strings.ToLower(service.Spec.Description), query) {
			filteredServices = append(filteredServices, service)
		}
	}

	return &ServiceList{
		Services: filteredServices,
		Total:    len(filteredServices),
	}, nil
}

// ServiceExists checks if a service exists
func (sc *ServiceCatalog) ServiceExists(name string) bool {
	return sc.storage.ServiceExists(name)
}

// GetServiceHealth returns aggregated health status for a service
func (sc *ServiceCatalog) GetServiceHealth(serviceName string) (*ServiceHealth, error) {
	service, err := sc.GetService(serviceName)
	if err != nil {
		return nil, err
	}

	// If service has no Kubernetes definition, return basic health
	if service.Spec.Kubernetes == nil {
		return &ServiceHealth{
			ServiceName:   serviceName,
			OverallStatus: "unknown",
			Environments:  []EnvironmentHealth{},
			LastUpdated:   time.Now(),
		}, nil
	}

	// TODO: Implement Kubernetes health aggregation
	// This will be implemented when we integrate with the Kubernetes plugin
	// For now, return a placeholder
	var environments []EnvironmentHealth

	for _, env := range service.Spec.Kubernetes.Environments {
		environments = append(environments, EnvironmentHealth{
			Name:        env.Name,
			Context:     env.Context,
			Status:      "unknown", // Will be populated from K8s plugin
			Deployments: []DeploymentHealth{},
		})
	}

	overallStatus := sc.calculateServiceStatus(environments, service.Metadata.Tier)

	return &ServiceHealth{
		ServiceName:   serviceName,
		OverallStatus: overallStatus,
		Environments:  environments,
		LastUpdated:   time.Now(),
	}, nil
}

// GetServiceHistory returns change history for a service
func (sc *ServiceCatalog) GetServiceHistory(serviceName string) (*ServiceHistory, error) {
	if sc.gitVersioning == nil {
		return nil, fmt.Errorf("git versioning is not available")
	}

	changes, err := sc.gitVersioning.GetServiceHistory(serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service history: %w", err)
	}

	return &ServiceHistory{
		ServiceName: serviceName,
		History:     changes,
	}, nil
}

// GetAllHistory returns complete change history
func (sc *ServiceCatalog) GetAllHistory() ([]ServiceChange, error) {
	if sc.gitVersioning == nil {
		return nil, fmt.Errorf("git versioning is not available")
	}

	changes, err := sc.gitVersioning.GetAllHistory()
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	return changes, nil
}

// ValidateService validates a service definition
func (sc *ServiceCatalog) ValidateService(service *Service) error {
	if service == nil {
		return fmt.Errorf("service cannot be nil")
	}

	// Use storage provider validation
	if fsProvider, ok := sc.storage.(*FilesystemProvider); ok {
		return fsProvider.validateService(service)
	}

	// Basic validation for other storage providers
	if service.Metadata.Name == "" {
		return fmt.Errorf("service name is required")
	}

	if service.Spec.Description == "" {
		return fmt.Errorf("service description is required")
	}

	if service.Spec.Team.GitHubTeam == "" {
		return fmt.Errorf("github_team is required")
	}

	return nil
}

// GetRepositoryStatus returns git repository status
func (sc *ServiceCatalog) GetRepositoryStatus() (string, error) {
	if sc.gitVersioning == nil {
		return "Git versioning is not available", nil
	}

	return sc.gitVersioning.GetRepositoryStatus()
}

// calculateServiceStatus determines overall service status based on environments and tier
func (sc *ServiceCatalog) calculateServiceStatus(environments []EnvironmentHealth, tier string) string {
	if len(environments) == 0 {
		return "unknown"
	}

	// Find production environment
	var prodHealth string
	var hasProduction bool

	for _, env := range environments {
		if env.Name == "production" {
			prodHealth = env.Status
			hasProduction = true
			break
		}
	}

	// If no production environment, use first environment
	if !hasProduction && len(environments) > 0 {
		prodHealth = environments[0].Status
	}

	// Apply tier-based logic
	switch tier {
	case "TIER-1":
		// Critical services: any production issues are critical
		if prodHealth == "down" {
			return "critical"
		}
		if prodHealth == "degraded" {
			return "critical"
		}
		return prodHealth

	case "TIER-2":
		// Important services: production issues are degraded
		if prodHealth == "down" {
			return "degraded"
		}
		return prodHealth

	case "TIER-3":
		// Standard services: only complete failure is concerning
		if prodHealth == "down" {
			return "degraded"
		}
		return "healthy"

	default:
		return prodHealth
	}
}

// newStorageProvider creates appropriate storage provider based on configuration
func newStorageProvider(config *Config) (StorageProvider, error) {
	switch config.Storage.Provider {
	case "filesystem":
		directory := config.Storage.Filesystem.Directory
		if directory == "" {
			directory = "./services" // Default directory
		}
		return NewFilesystemProvider(directory)

	case "github":
		// TODO: Implement GitHub storage provider
		return nil, fmt.Errorf("GitHub storage provider not yet implemented")

	case "s3":
		// TODO: Implement S3 storage provider
		return nil, fmt.Errorf("S3 storage provider not yet implemented")

	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", config.Storage.Provider)
	}
}
