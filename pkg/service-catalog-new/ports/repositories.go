package servicecatalog

import (
	"context"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog-new/models"
)

// ServiceRepository defines the interface for service data access
type ServiceRepository interface {
	// Create creates a new service
	Create(ctx context.Context, service *scModels.Service) (*scModels.Service, error)

	// GetByName retrieves a service by name
	GetByName(ctx context.Context, name string) (*scModels.Service, error)

	// Update updates an existing service
	Update(ctx context.Context, service *scModels.Service) (*scModels.Service, error)

	// Delete deletes a service
	Delete(ctx context.Context, name string) error

	// List lists all services with optional filtering
	List(ctx context.Context, filter *scModels.ServiceFilter) ([]scModels.Service, error)

	// Exists checks if a service exists
	Exists(ctx context.Context, name string) (bool, error)

	// ListByTeam lists services owned by a specific team
	ListByTeam(ctx context.Context, team string) ([]scModels.Service, error)

	// ListByTier lists services of a specific tier
	ListByTier(ctx context.Context, tier scModels.ServiceTier) ([]scModels.Service, error)

	// Search searches services by text query
	Search(ctx context.Context, query string, limit int) ([]scModels.Service, error)
}

// VersioningRepository defines the interface for versioning data access
type VersioningRepository interface {
	// RecordChange records a service change
	RecordChange(ctx context.Context, service *scModels.Service, user *scModels.UserContext, action string) error

	// RecordDeletion records a service deletion
	RecordDeletion(ctx context.Context, serviceName string, user *scModels.UserContext) error

	// GetServiceHistory gets change history for a service
	GetServiceHistory(ctx context.Context, serviceName string) ([]scModels.ServiceChange, error)

	// GetAllHistory gets complete change history
	GetAllHistory(ctx context.Context) ([]scModels.ServiceChange, error)

	// IsEnabled returns whether versioning is enabled
	IsEnabled() bool

	// GetStatus returns versioning system status
	GetStatus(ctx context.Context) (string, error)
}

// HealthRepository defines the interface for health data access
type HealthRepository interface {
	// GetServiceHealth gets health information for a service
	GetServiceHealth(ctx context.Context, serviceName string) (*scModels.ServiceHealth, error)

	// GetBatchHealth gets health information for multiple services
	GetBatchHealth(ctx context.Context, serviceNames []string) ([]scModels.ServiceHealth, error)

	// UpdateServiceHealth updates health information for a service
	UpdateServiceHealth(ctx context.Context, health *scModels.ServiceHealth) error

	// GetEnvironmentHealth gets health for a specific environment
	GetEnvironmentHealth(ctx context.Context, serviceName, environment string) (*scModels.EnvironmentHealth, error)
}
