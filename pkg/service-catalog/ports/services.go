package servicecatalog

import (
	"context"
	"time"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
)

// KubernetesService defines the interface for Kubernetes operations
type KubernetesService interface {
	// GetDeploymentHealth gets health information for a deployment
	GetDeploymentHealth(ctx context.Context, namespace, deploymentName, kubeContext string) (*scModels.DeploymentHealth, error)

	// GetEnvironmentHealth gets health information for all deployments in an environment
	GetEnvironmentHealth(ctx context.Context, service *scModels.Service, environment string) (*scModels.EnvironmentHealth, error)

	// GetServiceHealth gets aggregated health information for a service
	GetServiceHealth(ctx context.Context, service *scModels.Service) (*scModels.ServiceHealth, error)

	// ListNamespaces lists available namespaces in a context
	ListNamespaces(ctx context.Context, kubeContext string) ([]string, error)

	// ListDeployments lists deployments in a namespace
	ListDeployments(ctx context.Context, namespace, kubeContext string) ([]string, error)

	// ValidateContext validates if a Kubernetes context is accessible
	ValidateContext(ctx context.Context, kubeContext string) error
}

// GitHubService defines the interface for GitHub operations
type GitHubService interface {
	// GetTeamMembers gets members of a GitHub team
	GetTeamMembers(ctx context.Context, org, team string) ([]string, error)

	// ValidateTeamAccess validates if a user has access to a team
	ValidateTeamAccess(ctx context.Context, user, org, team string) (bool, error)

	// GetUserTeams gets all teams for a user in an organization
	GetUserTeams(ctx context.Context, user, org string) ([]string, error)

	// GetTeamInfo gets detailed information about a team
	GetTeamInfo(ctx context.Context, org, team string) (*TeamInfo, error)
}

// TeamInfo represents GitHub team information
type TeamInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
	URL         string   `json:"url"`
}

// StorageService defines the interface for storage operations
type StorageService interface {
	// Store stores a service definition
	Store(ctx context.Context, service *scModels.Service) error

	// Retrieve retrieves a service definition
	Retrieve(ctx context.Context, name string) (*scModels.Service, error)

	// Remove removes a service definition
	Remove(ctx context.Context, name string) error

	// ListAll lists all stored services
	ListAll(ctx context.Context) ([]scModels.Service, error)

	// Exists checks if a service exists in storage
	Exists(ctx context.Context, name string) (bool, error)

	// GetStorageInfo gets storage system information
	GetStorageInfo(ctx context.Context) (*StorageInfo, error)
}

// StorageInfo represents storage system information
type StorageInfo struct {
	Provider     string            `json:"provider"`
	Location     string            `json:"location"`
	ServiceCount int               `json:"service_count"`
	LastUpdated  time.Time         `json:"last_updated"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// ContextResolverService defines the interface for service context resolution
type ContextResolverService interface {
	// ResolveDeploymentService resolves which service owns a deployment
	ResolveDeploymentService(ctx context.Context, deploymentName, namespace, kubeContext string) (*scModels.ServiceContext, error)

	// ResolveServicesByNamespace resolves all services in a namespace
	ResolveServicesByNamespace(ctx context.Context, namespace, kubeContext string) ([]scModels.ServiceContext, error)

	// ResolveServiceDependencies resolves service dependencies
	ResolveServiceDependencies(ctx context.Context, serviceName string) ([]scModels.Service, error)
}

// NotificationService defines the interface for notification operations
type NotificationService interface {
	// NotifyServiceCreated notifies about service creation
	NotifyServiceCreated(ctx context.Context, service *scModels.Service, user *scModels.UserContext) error

	// NotifyServiceUpdated notifies about service updates
	NotifyServiceUpdated(ctx context.Context, service *scModels.Service, user *scModels.UserContext, changes []scModels.ServiceFieldChange) error

	// NotifyServiceDeleted notifies about service deletion
	NotifyServiceDeleted(ctx context.Context, serviceName string, user *scModels.UserContext) error

	// NotifyHealthAlert notifies about health issues
	NotifyHealthAlert(ctx context.Context, service *scModels.Service, health *scModels.ServiceHealth) error
}

// AuditService defines the interface for audit logging
type AuditService interface {
	// LogServiceOperation logs service operations
	LogServiceOperation(ctx context.Context, operation string, service *scModels.Service, user *scModels.UserContext) error

	// LogHealthCheck logs health check operations
	LogHealthCheck(ctx context.Context, serviceName string, health *scModels.ServiceHealth) error

	// LogSecurityEvent logs security-related events
	LogSecurityEvent(ctx context.Context, event string, user *scModels.UserContext, details map[string]interface{}) error
}
