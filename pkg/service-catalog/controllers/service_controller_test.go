package servicecatalog

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	scLogic "github.com/dash-ops/dash-ops/pkg/service-catalog/logic"
	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scPorts "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
)

// MockServiceRepository is a mock implementation of ServiceRepository
type MockServiceRepository struct {
	CreateFunc     func(ctx context.Context, service *scModels.Service) (*scModels.Service, error)
	GetByNameFunc  func(ctx context.Context, name string) (*scModels.Service, error)
	UpdateFunc     func(ctx context.Context, service *scModels.Service) (*scModels.Service, error)
	DeleteFunc     func(ctx context.Context, name string) error
	ListFunc       func(ctx context.Context, filter *scModels.ServiceFilter) ([]scModels.Service, error)
	ExistsFunc     func(ctx context.Context, name string) (bool, error)
	ListByTeamFunc func(ctx context.Context, team string) ([]scModels.Service, error)
	ListByTierFunc func(ctx context.Context, tier scModels.ServiceTier) ([]scModels.Service, error)
	SearchFunc     func(ctx context.Context, query string, limit int) ([]scModels.Service, error)
}

func (m *MockServiceRepository) Create(ctx context.Context, service *scModels.Service) (*scModels.Service, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, service)
	}
	return service, nil
}

func (m *MockServiceRepository) GetByName(ctx context.Context, name string) (*scModels.Service, error) {
	if m.GetByNameFunc != nil {
		return m.GetByNameFunc(ctx, name)
	}
	return nil, nil
}

func (m *MockServiceRepository) List(ctx context.Context, filter *scModels.ServiceFilter) ([]scModels.Service, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, filter)
	}
	return []scModels.Service{}, nil
}

func (m *MockServiceRepository) Update(ctx context.Context, service *scModels.Service) (*scModels.Service, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, service)
	}
	return service, nil
}

func (m *MockServiceRepository) Delete(ctx context.Context, name string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, name)
	}
	return nil
}

func (m *MockServiceRepository) Exists(ctx context.Context, name string) (bool, error) {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, name)
	}
	return false, nil
}

func (m *MockServiceRepository) ListByTeam(ctx context.Context, team string) ([]scModels.Service, error) {
	if m.ListByTeamFunc != nil {
		return m.ListByTeamFunc(ctx, team)
	}
	return []scModels.Service{}, nil
}

func (m *MockServiceRepository) ListByTier(ctx context.Context, tier scModels.ServiceTier) ([]scModels.Service, error) {
	if m.ListByTierFunc != nil {
		return m.ListByTierFunc(ctx, tier)
	}
	return []scModels.Service{}, nil
}

func (m *MockServiceRepository) Search(ctx context.Context, query string, limit int) ([]scModels.Service, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, query, limit)
	}
	return []scModels.Service{}, nil
}

// MockVersioningRepository is a mock implementation of VersioningRepository
type MockVersioningRepository struct {
	IsEnabledFunc         func() bool
	RecordChangeFunc      func(ctx context.Context, service *scModels.Service, user *scModels.UserContext, action string) error
	RecordDeletionFunc    func(ctx context.Context, serviceName string, user *scModels.UserContext) error
	GetServiceHistoryFunc func(ctx context.Context, serviceName string) ([]scModels.ServiceChange, error)
	GetAllHistoryFunc     func(ctx context.Context) ([]scModels.ServiceChange, error)
	GetStatusFunc         func(ctx context.Context) (string, error)
}

func (m *MockVersioningRepository) IsEnabled() bool {
	if m.IsEnabledFunc != nil {
		return m.IsEnabledFunc()
	}
	return false
}

func (m *MockVersioningRepository) RecordChange(ctx context.Context, service *scModels.Service, user *scModels.UserContext, action string) error {
	if m.RecordChangeFunc != nil {
		return m.RecordChangeFunc(ctx, service, user, action)
	}
	return nil
}

func (m *MockVersioningRepository) RecordDeletion(ctx context.Context, serviceName string, user *scModels.UserContext) error {
	if m.RecordDeletionFunc != nil {
		return m.RecordDeletionFunc(ctx, serviceName, user)
	}
	return nil
}

func (m *MockVersioningRepository) GetServiceHistory(ctx context.Context, serviceName string) ([]scModels.ServiceChange, error) {
	if m.GetServiceHistoryFunc != nil {
		return m.GetServiceHistoryFunc(ctx, serviceName)
	}
	return []scModels.ServiceChange{}, nil
}

func (m *MockVersioningRepository) GetAllHistory(ctx context.Context) ([]scModels.ServiceChange, error) {
	if m.GetAllHistoryFunc != nil {
		return m.GetAllHistoryFunc(ctx)
	}
	return []scModels.ServiceChange{}, nil
}

func (m *MockVersioningRepository) GetStatus(ctx context.Context) (string, error) {
	if m.GetStatusFunc != nil {
		return m.GetStatusFunc(ctx)
	}
	return "enabled", nil
}

// MockKubernetesService is a mock implementation of KubernetesService
type MockKubernetesService struct {
	GetDeploymentHealthFunc  func(ctx context.Context, namespace, deploymentName, kubeContext string) (*scModels.DeploymentHealth, error)
	GetEnvironmentHealthFunc func(ctx context.Context, service *scModels.Service, environment string) (*scModels.EnvironmentHealth, error)
	GetServiceHealthFunc     func(ctx context.Context, service *scModels.Service) (*scModels.ServiceHealth, error)
	ListNamespacesFunc       func(ctx context.Context, kubeContext string) ([]string, error)
	ListDeploymentsFunc      func(ctx context.Context, namespace, kubeContext string) ([]string, error)
	ValidateContextFunc      func(ctx context.Context, kubeContext string) error
}

func (m *MockKubernetesService) GetDeploymentHealth(ctx context.Context, namespace, deploymentName, kubeContext string) (*scModels.DeploymentHealth, error) {
	if m.GetDeploymentHealthFunc != nil {
		return m.GetDeploymentHealthFunc(ctx, namespace, deploymentName, kubeContext)
	}
	return &scModels.DeploymentHealth{
		Name:            deploymentName,
		ReadyReplicas:   1,
		DesiredReplicas: 1,
		Status:          scModels.StatusHealthy,
	}, nil
}

func (m *MockKubernetesService) GetEnvironmentHealth(ctx context.Context, service *scModels.Service, environment string) (*scModels.EnvironmentHealth, error) {
	if m.GetEnvironmentHealthFunc != nil {
		return m.GetEnvironmentHealthFunc(ctx, service, environment)
	}
	return &scModels.EnvironmentHealth{
		Name:   environment,
		Status: scModels.StatusHealthy,
	}, nil
}

func (m *MockKubernetesService) GetServiceHealth(ctx context.Context, service *scModels.Service) (*scModels.ServiceHealth, error) {
	if m.GetServiceHealthFunc != nil {
		return m.GetServiceHealthFunc(ctx, service)
	}
	return &scModels.ServiceHealth{
		ServiceName:   service.Metadata.Name,
		OverallStatus: scModels.StatusHealthy,
	}, nil
}

func (m *MockKubernetesService) ListNamespaces(ctx context.Context, kubeContext string) ([]string, error) {
	if m.ListNamespacesFunc != nil {
		return m.ListNamespacesFunc(ctx, kubeContext)
	}
	return []string{"default"}, nil
}

func (m *MockKubernetesService) ListDeployments(ctx context.Context, namespace, kubeContext string) ([]string, error) {
	if m.ListDeploymentsFunc != nil {
		return m.ListDeploymentsFunc(ctx, namespace, kubeContext)
	}
	return []string{}, nil
}

func (m *MockKubernetesService) ValidateContext(ctx context.Context, kubeContext string) error {
	if m.ValidateContextFunc != nil {
		return m.ValidateContextFunc(ctx, kubeContext)
	}
	return nil
}

// MockGitHubService is a mock implementation of GitHubService
type MockGitHubService struct {
	GetTeamMembersFunc     func(ctx context.Context, org, team string) ([]string, error)
	ValidateTeamAccessFunc func(ctx context.Context, user, org, team string) (bool, error)
	GetUserTeamsFunc       func(ctx context.Context, user, org string) ([]string, error)
	GetTeamInfoFunc        func(ctx context.Context, org, team string) (*scPorts.TeamInfo, error)
}

func (m *MockGitHubService) GetTeamMembers(ctx context.Context, org, team string) ([]string, error) {
	if m.GetTeamMembersFunc != nil {
		return m.GetTeamMembersFunc(ctx, org, team)
	}
	return []string{}, nil
}

func (m *MockGitHubService) ValidateTeamAccess(ctx context.Context, user, org, team string) (bool, error) {
	if m.ValidateTeamAccessFunc != nil {
		return m.ValidateTeamAccessFunc(ctx, user, org, team)
	}
	return true, nil
}

func (m *MockGitHubService) GetUserTeams(ctx context.Context, user, org string) ([]string, error) {
	if m.GetUserTeamsFunc != nil {
		return m.GetUserTeamsFunc(ctx, user, org)
	}
	return []string{}, nil
}

func (m *MockGitHubService) GetTeamInfo(ctx context.Context, org, team string) (*scPorts.TeamInfo, error) {
	if m.GetTeamInfoFunc != nil {
		return m.GetTeamInfoFunc(ctx, org, team)
	}
	return &scPorts.TeamInfo{
		Name: team,
		Slug: team,
	}, nil
}

func TestNewServiceController_CreatesControllerWithDependencies(t *testing.T) {
	// Arrange
	mockServiceRepo := &MockServiceRepository{}
	mockVersioningRepo := &MockVersioningRepository{}
	mockK8sService := &MockKubernetesService{}
	mockGitHubService := &MockGitHubService{}
	validator := scLogic.NewServiceValidator()
	processor := scLogic.NewServiceProcessor()

	// Act
	controller := NewServiceController(
		mockServiceRepo,
		mockVersioningRepo,
		mockK8sService,
		mockGitHubService,
		validator,
		processor,
	)

	// Assert
	assert.NotNil(t, controller)
	assert.Equal(t, mockServiceRepo, controller.serviceRepo)
	assert.Equal(t, mockVersioningRepo, controller.versioningRepo)
	assert.Equal(t, mockK8sService, controller.k8sService)
	assert.Equal(t, mockGitHubService, controller.githubService)
	assert.Equal(t, validator, controller.validator)
	assert.Equal(t, processor, controller.processor)
}

func TestServiceController_CreateService_WithValidService_CreatesSuccessfully(t *testing.T) {
	// Arrange
	service := &scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: "test-service",
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Test Service",
			Team: scModels.ServiceTeam{
				GitHubTeam: "test-team",
			},
			Technology: scModels.ServiceTechnology{
				Language:  "go",
				Framework: "kubernetes",
			},
		},
	}

	user := &scModels.UserContext{
		Username: "john.doe",
		Email:    "john@example.com",
	}

	mockServiceRepo := &MockServiceRepository{
		ExistsFunc: func(ctx context.Context, name string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, service *scModels.Service) (*scModels.Service, error) {
			service.Metadata.CreatedAt = time.Now()
			service.Metadata.CreatedBy = user.Username
			return service, nil
		},
	}

	mockVersioningRepo := &MockVersioningRepository{
		IsEnabledFunc: func() bool {
			return true
		},
		RecordChangeFunc: func(ctx context.Context, service *scModels.Service, user *scModels.UserContext, operation string) error {
			return nil
		},
	}

	controller := NewServiceController(
		mockServiceRepo,
		mockVersioningRepo,
		nil,
		nil,
		scLogic.NewServiceValidator(),
		scLogic.NewServiceProcessor(),
	)

	// Act
	createdService, err := controller.CreateService(context.Background(), service, user)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdService)
	assert.Equal(t, "test-service", createdService.Metadata.Name)
	assert.Equal(t, user.Username, createdService.Metadata.CreatedBy)
	assert.NotZero(t, createdService.Metadata.CreatedAt)
}

func TestServiceController_CreateService_WhenServiceExists_ReturnsError(t *testing.T) {
	// Arrange
	service := &scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: "existing-service",
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Existing Service",
			Team: scModels.ServiceTeam{
				GitHubTeam: "test-team",
			},
		},
	}

	mockServiceRepo := &MockServiceRepository{
		ExistsFunc: func(ctx context.Context, name string) (bool, error) {
			return true, nil
		},
	}

	controller := NewServiceController(
		mockServiceRepo,
		nil,
		nil,
		nil,
		scLogic.NewServiceValidator(),
		scLogic.NewServiceProcessor(),
	)

	// Act
	createdService, err := controller.CreateService(context.Background(), service, &scModels.UserContext{})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdService)
	assert.Contains(t, err.Error(), "already exists")
}

func TestServiceController_CreateService_WithInvalidService_ReturnsValidationError(t *testing.T) {
	// Arrange
	service := &scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: "", // Invalid: empty name
		},
	}

	controller := NewServiceController(
		&MockServiceRepository{},
		nil,
		nil,
		nil,
		scLogic.NewServiceValidator(),
		scLogic.NewServiceProcessor(),
	)

	// Act
	createdService, err := controller.CreateService(context.Background(), service, &scModels.UserContext{})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdService)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestServiceController_CreateService_WhenRepositoryFails_ReturnsError(t *testing.T) {
	// Arrange
	service := &scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: "test-service",
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Test Service",
			Team: scModels.ServiceTeam{
				GitHubTeam: "test-team",
			},
		},
	}

	mockServiceRepo := &MockServiceRepository{
		ExistsFunc: func(ctx context.Context, name string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, service *scModels.Service) (*scModels.Service, error) {
			return nil, errors.New("database error")
		},
	}

	controller := NewServiceController(
		mockServiceRepo,
		nil,
		nil,
		nil,
		scLogic.NewServiceValidator(),
		scLogic.NewServiceProcessor(),
	)

	// Act
	createdService, err := controller.CreateService(context.Background(), service, &scModels.UserContext{})

	// Assert
	assert.Error(t, err)
	assert.Nil(t, createdService)
	assert.Contains(t, err.Error(), "failed to create service")
}

func TestServiceController_CreateService_WithVersioningDisabled_SkipsVersioning(t *testing.T) {
	// Arrange
	service := &scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: "test-service",
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Test Service",
			Team: scModels.ServiceTeam{
				GitHubTeam: "test-team",
			},
		},
	}

	mockServiceRepo := &MockServiceRepository{
		ExistsFunc: func(ctx context.Context, name string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, service *scModels.Service) (*scModels.Service, error) {
			return service, nil
		},
	}

	versioningCalled := false
	mockVersioningRepo := &MockVersioningRepository{
		IsEnabledFunc: func() bool {
			return false
		},
		RecordChangeFunc: func(ctx context.Context, service *scModels.Service, user *scModels.UserContext, operation string) error {
			versioningCalled = true
			return nil
		},
	}

	controller := NewServiceController(
		mockServiceRepo,
		mockVersioningRepo,
		nil,
		nil,
		scLogic.NewServiceValidator(),
		scLogic.NewServiceProcessor(),
	)

	// Act
	createdService, err := controller.CreateService(context.Background(), service, &scModels.UserContext{})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdService)
	assert.False(t, versioningCalled)
}

func TestServiceController_CreateService_WhenVersioningFails_StillCreatesService(t *testing.T) {
	// Arrange
	service := &scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: "test-service",
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Test Service",
			Team: scModels.ServiceTeam{
				GitHubTeam: "test-team",
			},
		},
	}

	mockServiceRepo := &MockServiceRepository{
		ExistsFunc: func(ctx context.Context, name string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, service *scModels.Service) (*scModels.Service, error) {
			return service, nil
		},
	}

	mockVersioningRepo := &MockVersioningRepository{
		IsEnabledFunc: func() bool {
			return true
		},
		RecordChangeFunc: func(ctx context.Context, service *scModels.Service, user *scModels.UserContext, operation string) error {
			return errors.New("versioning error")
		},
	}

	controller := NewServiceController(
		mockServiceRepo,
		mockVersioningRepo,
		nil,
		nil,
		scLogic.NewServiceValidator(),
		scLogic.NewServiceProcessor(),
	)

	// Act
	createdService, err := controller.CreateService(context.Background(), service, &scModels.UserContext{})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, createdService)
	assert.Equal(t, "test-service", createdService.Metadata.Name)
}
