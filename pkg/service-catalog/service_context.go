package servicecatalog

import (
	"fmt"
	"strings"
)

// ServiceContext represents the context information linking a Kubernetes resource to a service
type ServiceContext struct {
	ServiceName string `json:"service_name"`
	ServiceTier string `json:"service_tier"`
	Environment string `json:"environment"`
	Context     string `json:"context"`
	Team        string `json:"team,omitempty"`
	Description string `json:"description,omitempty"`
}

// DeploymentContext represents context for a specific deployment
type DeploymentContext struct {
	DeploymentName string `json:"deployment_name"`
	Namespace      string `json:"namespace"`
	Context        string `json:"context"`
}

// ServiceContextResolver handles the resolution of Kubernetes resources to service context
type ServiceContextResolver struct {
	serviceCatalog *ServiceCatalog
}

// NewServiceContextResolver creates a new service context resolver
func NewServiceContextResolver(serviceCatalog *ServiceCatalog) *ServiceContextResolver {
	return &ServiceContextResolver{
		serviceCatalog: serviceCatalog,
	}
}

// ResolveDeploymentService resolves which service a deployment belongs to
// Parameters: deploymentName, namespace, context (k8s context)
// Returns: ServiceContext if found, nil if not found
func (scr *ServiceContextResolver) ResolveDeploymentService(deploymentName, namespace, context string) (*ServiceContext, error) {
	// Get all services
	serviceList, err := scr.serviceCatalog.ListServices()
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
			// Check if context matches (if provided)
			if context != "" && env.Context != context {
				continue
			}

			// Check if namespace matches (if provided)
			if namespace != "" && env.Namespace != namespace {
				continue
			}

			// Check each deployment in the environment
			for _, deployment := range env.Resources.Deployments {
				if matchesDeploymentName(deployment.Name, deploymentName) {
					// Found a match!
					return &ServiceContext{
						ServiceName: service.Metadata.Name,
						ServiceTier: service.Metadata.Tier,
						Environment: env.Name,
						Context:     env.Context,
						Team:        service.Spec.Team.GitHubTeam,
						Description: service.Spec.Description,
					}, nil
				}
			}
		}
	}

	// No matching service found
	return nil, nil
}

// ResolveMultipleDeployments resolves service context for multiple deployments at once
// This is useful for batch operations and improving performance
func (scr *ServiceContextResolver) ResolveMultipleDeployments(deployments []DeploymentContext) (map[string]*ServiceContext, error) {
	results := make(map[string]*ServiceContext)

	// Get all services once for efficiency
	serviceList, err := scr.serviceCatalog.ListServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	// Create a lookup key for each deployment
	for _, deployment := range deployments {
		key := createDeploymentKey(deployment.DeploymentName, deployment.Namespace, deployment.Context)

		// Find matching service for this deployment
		serviceContext := scr.findServiceForDeployment(serviceList.Services, deployment)
		if serviceContext != nil {
			results[key] = serviceContext
		}
	}

	return results, nil
}

// GetServiceDeployments returns all deployments for a given service
// This is useful for the reverse lookup: service -> deployments
func (scr *ServiceContextResolver) GetServiceDeployments(serviceName string) ([]DeploymentContext, error) {
	service, err := scr.serviceCatalog.GetService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s: %w", serviceName, err)
	}

	if service.Spec.Kubernetes == nil {
		return []DeploymentContext{}, nil
	}

	var deployments []DeploymentContext

	// Iterate through all environments
	for _, env := range service.Spec.Kubernetes.Environments {
		// Iterate through all deployments in the environment
		for _, deployment := range env.Resources.Deployments {
			deployments = append(deployments, DeploymentContext{
				DeploymentName: deployment.Name,
				Namespace:      env.Namespace,
				Context:        env.Context,
			})
		}
	}

	return deployments, nil
}

// Helper functions

// matchesDeploymentName checks if a deployment name matches
// This supports exact matches and could be extended for pattern matching
func matchesDeploymentName(configuredName, actualName string) bool {
	// Exact match
	if configuredName == actualName {
		return true
	}

	// Case-insensitive match
	if strings.EqualFold(configuredName, actualName) {
		return true
	}

	// TODO: Add support for pattern matching if needed
	// For example: auth-* could match auth-api, auth-worker, etc.

	return false
}

// findServiceForDeployment finds the service that contains a specific deployment
func (scr *ServiceContextResolver) findServiceForDeployment(services []Service, deployment DeploymentContext) *ServiceContext {
	for _, service := range services {
		if service.Spec.Kubernetes == nil {
			continue
		}

		for _, env := range service.Spec.Kubernetes.Environments {
			// Check context and namespace match
			if deployment.Context != "" && env.Context != deployment.Context {
				continue
			}
			if deployment.Namespace != "" && env.Namespace != deployment.Namespace {
				continue
			}

			// Check deployment match
			for _, serviceDeployment := range env.Resources.Deployments {
				if matchesDeploymentName(serviceDeployment.Name, deployment.DeploymentName) {
					return &ServiceContext{
						ServiceName: service.Metadata.Name,
						ServiceTier: service.Metadata.Tier,
						Environment: env.Name,
						Context:     env.Context,
						Team:        service.Spec.Team.GitHubTeam,
						Description: service.Spec.Description,
					}
				}
			}
		}
	}

	return nil
}

// createDeploymentKey creates a unique key for a deployment context
func createDeploymentKey(deploymentName, namespace, context string) string {
	return fmt.Sprintf("%s/%s/%s", context, namespace, deploymentName)
}

// ValidateServiceContext validates that a service context is complete
func ValidateServiceContext(ctx *ServiceContext) error {
	if ctx == nil {
		return fmt.Errorf("service context is nil")
	}

	if ctx.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}

	if ctx.ServiceTier == "" {
		return fmt.Errorf("service tier is required")
	}

	if ctx.Environment == "" {
		return fmt.Errorf("environment is required")
	}

	return nil
}
