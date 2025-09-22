package servicecatalog

import (
	"context"
	"fmt"
	"time"

	scLogic "github.com/dash-ops/dash-ops/pkg/service-catalog/logic"
	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
)

// ServiceController handles service business logic orchestration
type ServiceController struct {
	serviceRepo    scPorts.ServiceRepository
	versioningRepo scPorts.VersioningRepository
	k8sService     scPorts.KubernetesService
	githubService  scPorts.GitHubService
	validator      *scLogic.ServiceValidator
	processor      *scLogic.ServiceProcessor
}

// NewServiceController creates a new service controller
func NewServiceController(
	serviceRepo scPorts.ServiceRepository,
	versioningRepo scPorts.VersioningRepository,
	k8sService scPorts.KubernetesService,
	githubService scPorts.GitHubService,
	validator *scLogic.ServiceValidator,
	processor *scLogic.ServiceProcessor,
) *ServiceController {
	return &ServiceController{
		serviceRepo:    serviceRepo,
		versioningRepo: versioningRepo,
		k8sService:     k8sService,
		githubService:  githubService,
		validator:      validator,
		processor:      processor,
	}
}

// CreateService creates a new service
func (sc *ServiceController) CreateService(ctx context.Context, service *scModels.Service, user *scModels.UserContext) (*scModels.Service, error) {
	// Validate for creation
	if err := sc.validator.ValidateForCreation(service); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if service already exists
	exists, err := sc.serviceRepo.Exists(ctx, service.Metadata.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check if service exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("service '%s' already exists", service.Metadata.Name)
	}

	// Prepare for creation
	preparedService, err := sc.processor.PrepareForCreation(service, user)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare service for creation: %w", err)
	}

	// Create in repository
	createdService, err := sc.serviceRepo.Create(ctx, preparedService)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	// Record change in versioning system
	if sc.versioningRepo != nil && sc.versioningRepo.IsEnabled() {
		if err := sc.versioningRepo.RecordChange(ctx, createdService, user, "create"); err != nil {
			// Log error but don't fail the operation
		}
	}

	return createdService, nil
}

// GetService retrieves a service by name
func (sc *ServiceController) GetService(ctx context.Context, name string) (*scModels.Service, error) {
	if name == "" {
		return nil, fmt.Errorf("service name is required")
	}

	service, err := sc.serviceRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	// Enrich with GitHub team information if available
	if sc.githubService != nil {
		if err := sc.enrichWithTeamInfo(ctx, service); err != nil {
			// Log error but don't fail the operation
		}
	}

	return service, nil
}

// UpdateService updates an existing service
func (sc *ServiceController) UpdateService(ctx context.Context, service *scModels.Service, user *scModels.UserContext) (*scModels.Service, error) {
	// Get existing service
	existingService, err := sc.serviceRepo.GetByName(ctx, service.Metadata.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing service: %w", err)
	}

	// Validate permissions
	if err := sc.validator.ValidateUserPermissions(existingService, user, "update"); err != nil {
		return nil, fmt.Errorf("permission denied: %w", err)
	}

	// Validate for update
	if err := sc.validator.ValidateForUpdate(service, existingService); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Prepare for update
	preparedService, err := sc.processor.PrepareForUpdate(service, existingService, user)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare service for update: %w", err)
	}

	// Update in repository
	updatedService, err := sc.serviceRepo.Update(ctx, preparedService)
	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	// Record change in versioning system
	if sc.versioningRepo != nil && sc.versioningRepo.IsEnabled() {
		if err := sc.versioningRepo.RecordChange(ctx, updatedService, user, "update"); err != nil {
			// Log error but don't fail the operation
		}
	}

	return updatedService, nil
}

// DeleteService deletes a service
func (sc *ServiceController) DeleteService(ctx context.Context, name string, user *scModels.UserContext) error {
	// Get existing service for permission check
	existingService, err := sc.serviceRepo.GetByName(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get service for deletion: %w", err)
	}

	// Validate permissions
	if err := sc.validator.ValidateUserPermissions(existingService, user, "delete"); err != nil {
		return fmt.Errorf("permission denied: %w", err)
	}

	// Delete from repository
	if err := sc.serviceRepo.Delete(ctx, name); err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	// Record deletion in versioning system
	if sc.versioningRepo != nil && sc.versioningRepo.IsEnabled() {
		if err := sc.versioningRepo.RecordDeletion(ctx, name, user); err != nil {
			// Log error but don't fail the operation
		}
	}

	return nil
}

// ListServices lists all services with optional filtering
func (sc *ServiceController) ListServices(ctx context.Context, filter *scModels.ServiceFilter) (*scModels.ServiceList, error) {
	services, err := sc.serviceRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	// Process the service list (filtering, pagination, etc.)
	serviceList := sc.processor.ProcessServiceList(services, filter)

	// Enrich with GitHub team information if available
	if sc.githubService != nil {
		for i := range serviceList.Services {
			if err := sc.enrichWithTeamInfo(ctx, &serviceList.Services[i]); err != nil {
				// Log error but continue with other services
			}
		}
	}

	return serviceList, nil
}

// ResolveDeploymentService resolves which service a deployment belongs to
func (sc *ServiceController) ResolveDeploymentService(ctx context.Context, deploymentName, namespace, kubeContext string) (*scModels.ServiceContext, error) {
	// Get all services
	serviceList, err := sc.ListServices(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	// Iterate through all services to find matching deployment
	for _, service := range serviceList.Services {
		if service.Spec.Kubernetes == nil {
			continue // Skip services without Kubernetes configuration
		}

		// Check each environment in the service
		for _, env := range service.Spec.Kubernetes.Environments {
			// Match context
			if env.Context != kubeContext {
				continue
			}

			// Match namespace
			if env.Namespace != namespace {
				continue
			}

			// Check if this environment has the deployment
			for _, deployment := range env.Resources.Deployments {
				if deployment.Name == deploymentName {
					// Found matching deployment - create service context
					return &scModels.ServiceContext{
						Service:     &service,
						Environment: env.Name,
						Namespace:   namespace,
						Context:     kubeContext,
						Found:       true,
					}, nil
				}
			}
		}
	}

	// No matching service found
	return nil, nil
}

// UpdateKubernetesService updates the kubernetes service dependency
func (sc *ServiceController) UpdateKubernetesService(k8sService scPorts.KubernetesService) {
	sc.k8sService = k8sService
}

// GetServiceHealth gets health information for a service
func (sc *ServiceController) GetServiceHealth(ctx context.Context, serviceName string) (*scModels.ServiceHealth, error) {
	// Get service definition
	service, err := sc.serviceRepo.GetByName(ctx, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	// Get health from Kubernetes if integration is available
	if sc.k8sService != nil {
		health, err := sc.k8sService.GetServiceHealth(ctx, service)
		if err != nil {
			// Return basic health info if K8s integration fails
			return &scModels.ServiceHealth{
				ServiceName:   serviceName,
				OverallStatus: scModels.StatusUnknown,
				Environments:  []scModels.EnvironmentHealth{},
				LastUpdated:   time.Now(),
			}, nil
		}
		return health, nil
	}

	// Return unknown status if no Kubernetes integration
	return &scModels.ServiceHealth{
		ServiceName:   serviceName,
		OverallStatus: scModels.StatusUnknown,
		Environments:  []scModels.EnvironmentHealth{},
		LastUpdated:   time.Now(),
	}, nil
}

// GetServiceHistory gets change history for a service
func (sc *ServiceController) GetServiceHistory(ctx context.Context, serviceName string) (*scModels.ServiceHistory, error) {
	if sc.versioningRepo == nil || !sc.versioningRepo.IsEnabled() {
		return nil, fmt.Errorf("versioning is not enabled")
	}

	changes, err := sc.versioningRepo.GetServiceHistory(ctx, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service history: %w", err)
	}

	return &scModels.ServiceHistory{
		ServiceName: serviceName,
		History:     changes,
	}, nil
}

// enrichWithTeamInfo enriches service with GitHub team information
func (sc *ServiceController) enrichWithTeamInfo(ctx context.Context, service *scModels.Service) error {
	if sc.githubService == nil {
		return nil
	}

	// TODO: Get organization from configuration
	org := "dash-ops" // This should come from configuration

	teamInfo, err := sc.githubService.GetTeamInfo(ctx, org, service.Spec.Team.GitHubTeam)
	if err != nil {
		return fmt.Errorf("failed to get team info: %w", err)
	}

	// Enrich service with team information
	service.Spec.Team.Members = teamInfo.Members
	service.Spec.Team.GitHubURL = teamInfo.URL

	return nil
}

// GetServiceRepository returns the service repository for external access
func (c *ServiceController) GetServiceRepository() scPorts.ServiceRepository {
	return c.serviceRepo
}
