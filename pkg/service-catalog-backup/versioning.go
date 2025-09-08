package servicecatalog

// VersioningProvider defines the interface for service versioning backends
type VersioningProvider interface {
	// Initialize sets up the versioning system
	Initialize() error

	// CommitServiceChange records a service change
	CommitServiceChange(service *Service, user *UserContext, action string) error

	// CommitServiceDeletion records a service deletion
	CommitServiceDeletion(serviceName string, user *UserContext) error

	// GetServiceHistory returns change history for a specific service
	GetServiceHistory(serviceName string) ([]ServiceChange, error)

	// GetAllHistory returns complete change history
	GetAllHistory() ([]ServiceChange, error)

	// IsEnabled returns whether versioning is enabled
	IsEnabled() bool

	// GetStatus returns current versioning system status
	GetStatus() (string, error)
}

// NoVersioning is a no-op versioning provider
type NoVersioning struct{}

// NewNoVersioning creates a disabled versioning provider
func NewNoVersioning() *NoVersioning {
	return &NoVersioning{}
}

func (nv *NoVersioning) Initialize() error {
	return nil
}

func (nv *NoVersioning) CommitServiceChange(service *Service, user *UserContext, action string) error {
	return nil
}

func (nv *NoVersioning) CommitServiceDeletion(serviceName string, user *UserContext) error {
	return nil
}

func (nv *NoVersioning) GetServiceHistory(serviceName string) ([]ServiceChange, error) {
	return []ServiceChange{}, nil
}

func (nv *NoVersioning) GetAllHistory() ([]ServiceChange, error) {
	return []ServiceChange{}, nil
}

func (nv *NoVersioning) IsEnabled() bool {
	return false
}

func (nv *NoVersioning) GetStatus() (string, error) {
	return "Versioning disabled", nil
}
