# Backend Development Guide

> **Complete guide for developing and contributing to DashOps backend modules**

## 📚 Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Module Structure](#module-structure)
3. [Development Workflow](#development-workflow)
4. [Testing Standards](#testing-standards)
5. [Dependency Injection](#dependency-injection)
6. [Module-Specific Guides](#module-specific-guides)
7. [Best Practices](#best-practices)
8. [Common Patterns](#common-patterns)

## Architecture Overview

DashOps backend follows **Hexagonal Architecture** with **Dependency Injection** for all modules. This provides:

- ✅ **Consistency**: All modules follow the same pattern
- ✅ **Testability**: Pure business logic with 100% test coverage
- ✅ **Maintainability**: Clear separation of concerns
- ✅ **Extensibility**: Easy to add new features or modules
- ✅ **Modularity**: Domain-specific controllers and repositories for better organization

### 🎯 Core Principles

1. **Layer Separation**: Each layer has a specific responsibility
2. **Dependency Inversion**: Depend on interfaces, not implementations
3. **Pure Business Logic**: Logic layer contains no external dependencies
4. **Test-First Development**: Logic layer must be 100% tested

### 🔄 Data Flow

```
HTTP Request → Handler → Adapter → Controller → Logic → Repository
```

## Module Structure

Every module follows this consistent 9-layer pattern:

```
pkg/{module}/
├── adapters/           # Data transformation ONLY (wire ↔ models)
│   ├── http/          # HTTP request/response adapters
│   ├── storage/       # Database/filesystem adapters
│   └── external/      # Data transformation adapters
├── integrations/      # External communication & module integration
│   ├── external/      # External service integrations (APIs, services)
│   │   ├── github/    # GitHub API client
│   │   ├── aws/       # AWS API client
│   │   ├── kubernetes/# Kubernetes API client
│   │   ├── prometheus/# Prometheus API client
│   │   ├── loki/      # Loki API client
│   │   └── tempo/     # Tempo API client
│   ├── service_catalog/  # Service Catalog module integration
│   ├── kubernetes/       # Kubernetes module integration
│   └── auth/             # Auth module integration
├── controllers/       # Business logic orchestration
├── handlers/          # HTTP endpoints (entry points)
├── logic/             # Pure business logic (100% tested)
├── models/            # Domain entities with behavior
├── ports/             # Interfaces & contracts
├── wire/              # API DTOs (request/response)
└── module.go          # Module factory & initialization
```

## Modular Architecture (NEW)

### 🎯 Domain-Specific Organization

For better maintainability and separation of concerns, modules can be organized with domain-specific controllers and repositories:

```
pkg/{module}/
├── integrations/external/
│   └── {service}/
│       ├── {service}_client.go      # External service client
│       └── {service}_adapter.go     # Data transformation adapter
├── repositories/                     # Domain-specific repositories
│   ├── {domain}_repository.go       # Repository for specific domain
│   ├── {domain2}_repository.go      # Repository for another domain
│   └── ...
├── controllers/                      # Domain-specific controllers
│   ├── {domain}_controller.go       # Controller for specific domain
│   ├── {domain2}_controller.go      # Controller for another domain
│   └── ...
├── handlers/
│   └── http.go                       # HTTP handler with dependency injection
└── module.go                         # Module factory & initialization
```

### 🔧 Dependency Injection Pattern

```go
// handlers/http.go
func NewHTTPHandler(
    k8sClient *kubernetes.KubernetesClient,
    responseAdapter *commonsHttp.ResponseAdapter,
    requestAdapter *commonsHttp.RequestAdapter,
) *HTTPHandler {
    // Initialize repositories with client
    nodesRepo := repositories.NewNodesRepository(k8sClient)
    deploymentsRepo := repositories.NewDeploymentsRepository(k8sClient)
    podsRepo := repositories.NewPodsRepository(k8sClient)
    namespacesRepo := repositories.NewNamespacesRepository(k8sClient)
    
    // Initialize controllers with repositories
    nodesController := controllers.NewNodesController(nodesRepo)
    deploymentsController := controllers.NewDeploymentsController(deploymentsRepo)
    podsController := controllers.NewPodsController(podsRepo)
    namespacesController := controllers.NewNamespacesController(namespacesRepo)
    
    return &HTTPHandler{
        nodesController:       nodesController,
        deploymentsController: deploymentsController,
        podsController:        podsController,
        namespacesController:  namespacesController,
        responseAdapter:       responseAdapter,
        requestAdapter:        requestAdapter,
    }
}
```

### 📁 Example: Kubernetes Module Structure

```
pkg/kubernetes/
├── integrations/external/
│   └── kubernetes/
│       ├── kubernetes_client.go      # K8s API client
│       └── kubernetes_adapter.go     # Data transformation
├── repositories/
│   ├── nodes_repository.go           # Nodes-specific operations
│   ├── deployments_repository.go     # Deployments-specific operations
│   ├── pods_repository.go            # Pods-specific operations
│   └── namespaces_repository.go      # Namespaces-specific operations
├── controllers/
│   ├── nodes_controller.go           # Nodes business logic
│   ├── deployments_controller.go     # Deployments business logic
│   ├── pods_controller.go            # Pods business logic
│   └── namespaces_controller.go      # Namespaces business logic
├── handlers/
│   └── http.go                       # HTTP routing & dependency injection
└── module.go                         # Module initialization
```

### 🎯 Benefits of Modular Architecture

1. **Single Responsibility**: Each controller/repository has one clear purpose
2. **Easy Testing**: Components can be tested in isolation
3. **Maintainability**: Changes in one domain don't affect others
4. **Scalability**: Easy to add new domains or modify existing ones
5. **Clear Dependencies**: Explicit dependency injection makes relationships clear
6. **Code Organization**: Related functionality is grouped together

### Layer Responsibilities

| Layer | Purpose | Testing Requirements |
|-------|---------|---------------------|
| **handlers** | HTTP endpoints, routing | Integration tests |
| **controllers** | Orchestration, workflow | Integration tests with mocks |
| **logic** | Pure business rules | 100% unit test coverage |
| **models** | Domain entities | Unit tests for methods |
| **adapters** | Data transformation ONLY | Unit tests |
| **integrations** | External communication | Unit tests with mocks |
| **ports** | Interfaces | No tests (interfaces) |
| **wire** | DTOs | No tests (data structures) |

## Integrations vs Adapters

### 🎯 Clear Separation of Concerns

DashOps follows a **strict separation** between data transformation and external communication:

#### `adapters/` - Data Transformation ONLY
- **Purpose**: Transform data between different formats (wire ↔ models)
- **Responsibility**: Pure data conversion, no external communication
- **Examples**: 
  - Convert HTTP request to domain model
  - Convert domain model to HTTP response
  - Transform database row to entity
  - Map external API response to internal model

#### `integrations/` - External Communication
- **Purpose**: Handle communication with external services and modules
- **Responsibility**: API calls, network communication, service discovery
- **Examples**:
  - GitHub API calls
  - AWS service calls
  - Kubernetes API calls
  - Inter-module communication

### 📁 Integration Structure

#### External Integrations (`integrations/external/`)
For third-party services and APIs:

```
integrations/external/
├── github/
│   ├── github_client.go      # GitHub API client
│   └── github_adapter.go     # Data transformation for GitHub
├── aws/
│   ├── ec2_client.go         # EC2 API client
│   ├── s3_client.go          # S3 API client
│   └── aws_adapter.go        # Data transformation for AWS
└── kubernetes/
    ├── k8s_client.go         # Kubernetes API client
    └── k8s_adapter.go        # Data transformation for K8s
```

#### Internal Integrations (`integrations/{module}/`)
For communication between DashOps modules:

```
integrations/
├── service_catalog/
│   └── service_catalog_integration.go  # Service Catalog module client
├── kubernetes/
│   └── kubernetes_integration.go       # Kubernetes module client
└── auth/
    └── auth_integration.go             # Auth module client
```

### 🔄 Data Flow with Integrations

```
HTTP Request → Handler → Adapter → Controller → Logic → Repository
                    ↓
              Integration (External/Internal)
                    ↓
              External Service/Module
```

### ✅ Correct Patterns

#### External Integration Example
```go
// integrations/external/github/github_client.go
package github

import (
    "context"
    "github.com/google/go-github/v50/github"
    "golang.org/x/oauth2"
)

type GitHubClient struct {
    client *github.Client
}

func NewGitHubClient() *GitHubClient {
    return &GitHubClient{
        client: github.NewClient(nil),
    }
}

func (c *GitHubClient) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
    client := c.client.WithAuthToken(token.AccessToken)
    return client.Users.Get(ctx, "")
}

// integrations/external/github/github_adapter.go
package github

import (
    "context"
    "github.com/dash-ops/dash-ops/pkg/auth/ports"
    "github.com/google/go-github/v50/github"
    "golang.org/x/oauth2"
)

type GitHubAdapter struct {
    client *GitHubClient
}

func NewGitHubAdapter() ports.GitHubService {
    return &GitHubAdapter{
        client: NewGitHubClient(),
    }
}

// Pure data transformation - no external communication
func (a *GitHubAdapter) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
    return a.client.GetUser(ctx, token)
}
```

#### Internal Integration Example
```go
// integrations/service_catalog/service_catalog_integration.go
package servicecatalog

import (
    "context"
    "github.com/dash-ops/dash-ops/pkg/service-catalog/ports"
    "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
)

type ServiceCatalogIntegration struct {
    api ports.ExposedAPI
}

func NewServiceCatalogIntegration(api ports.ExposedAPI) *ServiceCatalogIntegration {
    return &ServiceCatalogIntegration{
        api: api,
    }
}

func (i *ServiceCatalogIntegration) ResolveServiceContext(ctx context.Context, 
    deployment, namespace, context string) (*models.ServiceContext, error) {
    
    return i.api.GetServiceContext(ctx, deployment, namespace, context)
}
```

### ❌ Anti-Patterns to Avoid

```go
// ❌ WRONG: Adapter doing external communication
type GitHubAdapter struct {
    client *github.Client
}

func (a *GitHubAdapter) GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error) {
    // This should be in integrations/external/github/
    client := a.client.WithAuthToken(token.AccessToken)
    return client.Users.Get(ctx, "")
}

// ❌ WRONG: Integration doing data transformation
type ServiceCatalogIntegration struct {
    // This should be in adapters/
    transformer *DataTransformer
}

func (i *ServiceCatalogIntegration) TransformService(service *Service) *ServiceResponse {
    // This should be in adapters/
    return i.transformer.Transform(service)
}
```

### 🎯 Benefits of This Separation

1. **Single Responsibility**: Each layer has one clear purpose
2. **Testability**: Easy to mock integrations and test adapters
3. **Reusability**: Integrations can be used by multiple adapters
4. **Maintainability**: Changes to external APIs don't affect data transformation
5. **Clarity**: Easy to understand what each component does
6. **Flexibility**: Easy to swap implementations (e.g., mock vs real API)

## Development Workflow

### 🚀 Adding a New Feature

#### 1. Define the Domain Model

```go
// models/entities.go
type Service struct {
    Metadata ServiceMetadata `yaml:"metadata" json:"metadata"`
    Spec     ServiceSpec     `yaml:"spec" json:"spec"`
}

// Add domain behavior
func (s *Service) Validate() error {
    if s.Metadata.Name == "" {
        return fmt.Errorf("service name is required")
    }
    if s.Metadata.Tier == "" {
        s.Metadata.Tier = TierStandard
    }
    return nil
}

func (s *Service) IsProduction() bool {
    return s.Metadata.Tier == TierCritical
}
```

#### 2. Implement Business Logic

```go
// logic/service_processor.go
type ServiceProcessor struct{}

func NewServiceProcessor() *ServiceProcessor {
    return &ServiceProcessor{}
}

func (sp *ServiceProcessor) PrepareForCreation(service *Service, user *UserContext) (*Service, error) {
    // Pure business logic - no external dependencies
    prepared := *service
    prepared.Metadata.CreatedAt = time.Now()
    prepared.Metadata.CreatedBy = user.Username
    prepared.Metadata.Version = 1
    
    // Apply business rules
    if err := sp.validateBusinessRules(&prepared); err != nil {
        return nil, fmt.Errorf("business rule validation failed: %w", err)
    }
    
    return &prepared, nil
}

func (sp *ServiceProcessor) validateBusinessRules(service *Service) error {
    // Complex business logic here
    if service.IsProduction() && service.Spec.Team.GitHubTeam == "" {
        return fmt.Errorf("production services must have a team assigned")
    }
    return nil
}
```

#### 3. Write Tests FIRST

```go
// logic/service_processor_test.go
package servicecatalog

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
)

func TestServiceProcessor_PrepareForCreation_WithValidService_PreparesCorrectly(t *testing.T) {
    // Arrange
    processor := NewServiceProcessor()
    service := &Service{
        Metadata: ServiceMetadata{
            Name: "test-service",
            Tier: TierStandard,
        },
        Spec: ServiceSpec{
            Team: ServiceTeam{
                GitHubTeam: "test-team",
            },
        },
    }
    user := &UserContext{
        Username: "john.doe",
        Email:    "john@example.com",
    }
    
    // Act
    result, err := processor.PrepareForCreation(service, user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "john.doe", result.Metadata.CreatedBy)
    assert.Equal(t, 1, result.Metadata.Version)
    assert.NotZero(t, result.Metadata.CreatedAt)
}

func TestServiceProcessor_PrepareForCreation_ProductionWithoutTeam_ReturnsError(t *testing.T) {
    // Arrange
    processor := NewServiceProcessor()
    service := &Service{
        Metadata: ServiceMetadata{
            Name: "prod-service",
            Tier: TierCritical, // Production tier
        },
        Spec: ServiceSpec{
            Team: ServiceTeam{
                GitHubTeam: "", // No team
            },
        },
    }
    user := &UserContext{Username: "john.doe"}
    
    // Act
    result, err := processor.PrepareForCreation(service, user)
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "production services must have a team assigned")
}
```

#### 4. Define Interfaces

```go
// ports/repositories.go
type ServiceRepository interface {
    Create(ctx context.Context, service *Service) (*Service, error)
    GetByName(ctx context.Context, name string) (*Service, error)
    Update(ctx context.Context, service *Service) (*Service, error)
    Delete(ctx context.Context, name string) error
    List(ctx context.Context, filter *ServiceFilter) ([]Service, error)
    Exists(ctx context.Context, name string) (bool, error)
}

// ports/services.go
type KubernetesService interface {
    GetServiceHealth(ctx context.Context, service *Service) (*ServiceHealth, error)
    GetDeploymentHealth(ctx context.Context, namespace, name, context string) (*DeploymentHealth, error)
}
```

#### 5. Create Controller

```go
// controllers/service_controller.go
type ServiceController struct {
    serviceRepo    ServiceRepository
    k8sService     KubernetesService
    processor      *ServiceProcessor
    validator      *ServiceValidator
}

func NewServiceController(
    serviceRepo ServiceRepository,
    k8sService KubernetesService,
    processor *ServiceProcessor,
    validator *ServiceValidator,
) *ServiceController {
    return &ServiceController{
        serviceRepo: serviceRepo,
        k8sService:  k8sService,
        processor:   processor,
        validator:   validator,
    }
}

func (sc *ServiceController) CreateService(ctx context.Context, service *Service, user *UserContext) (*Service, error) {
    // Validation
    if err := sc.validator.ValidateForCreation(service); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // Check existence
    if exists, err := sc.serviceRepo.Exists(ctx, service.Metadata.Name); err != nil {
        return nil, fmt.Errorf("failed to check existence: %w", err)
    } else if exists {
        return nil, fmt.Errorf("service '%s' already exists", service.Metadata.Name)
    }
    
    // Process
    prepared, err := sc.processor.PrepareForCreation(service, user)
    if err != nil {
        return nil, fmt.Errorf("preparation failed: %w", err)
    }
    
    // Store
    return sc.serviceRepo.Create(ctx, prepared)
}
```

#### 6. Add HTTP Handler

```go
// handlers/http.go
func (h *HTTPHandler) createServiceHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request
    var req wire.CreateServiceRequest
    if err := h.requestAdapter.ParseJSON(r, &req); err != nil {
        h.responseAdapter.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    // Get user context
    user := h.getUserContext(r)
    
    // Transform to domain model
    service, err := h.adapter.RequestToService(req)
    if err != nil {
        h.responseAdapter.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    // Call controller
    result, err := h.controller.CreateService(r.Context(), service, user)
    if err != nil {
        if strings.Contains(err.Error(), "already exists") {
            h.responseAdapter.WriteError(w, http.StatusConflict, err.Error())
        } else {
            h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }
    
    // Transform and respond
    response := h.adapter.ServiceToResponse(result)
    h.responseAdapter.WriteCreated(w, "/services/"+result.Metadata.Name, response)
}
```

## Testing Standards

### 🎯 Testing Philosophy

- **Test behavior, not implementation**
- **Use clear, descriptive test names**
- **Follow AAA pattern**: Arrange, Act, Assert
- **One assertion concept per test**
- **Use `testify/assert` for readability**

### Unit Testing (Logic Layer)

```go
package servicecatalog

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

// Test naming convention: Test{Struct}_{Method}_{Scenario}_{ExpectedResult}
func TestServiceValidator_ValidateForCreation_WithEmptyName_ReturnsError(t *testing.T) {
    // Arrange
    validator := NewServiceValidator()
    service := &Service{
        Metadata: ServiceMetadata{
            Name: "", // Empty name
        },
    }
    
    // Act
    err := validator.ValidateForCreation(service)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "name is required")
}

func TestServiceValidator_ValidateForCreation_WithValidService_ReturnsNoError(t *testing.T) {
    // Arrange
    validator := NewServiceValidator()
    service := &Service{
        Metadata: ServiceMetadata{
            Name: "valid-service",
            Tier: TierStandard,
        },
        Spec: ServiceSpec{
            Description: "A valid service",
        },
    }
    
    // Act
    err := validator.ValidateForCreation(service)
    
    // Assert
    assert.NoError(t, err)
}
```

### Integration Testing (Controllers)

```go
func TestServiceController_CreateService_WithValidInput_CreatesSuccessfully(t *testing.T) {
    // Arrange
    mockRepo := &MockServiceRepository{
        ExistsFunc: func(ctx context.Context, name string) (bool, error) {
            return false, nil
        },
        CreateFunc: func(ctx context.Context, service *Service) (*Service, error) {
            service.Metadata.CreatedAt = time.Now()
            return service, nil
        },
    }
    
    controller := NewServiceController(
        mockRepo,
        nil, // k8sService not needed for this test
        NewServiceProcessor(),
        NewServiceValidator(),
    )
    
    service := &Service{
        Metadata: ServiceMetadata{
            Name: "test-service",
            Tier: TierStandard,
        },
    }
    user := &UserContext{Username: "john.doe"}
    
    // Act
    result, err := controller.CreateService(context.Background(), service, user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "test-service", result.Metadata.Name)
    assert.Equal(t, "john.doe", result.Metadata.CreatedBy)
}
```

### HTTP Handler Testing

```go
func TestHTTPHandler_CreateService_WithValidRequest_Returns201(t *testing.T) {
    // Arrange
    mockController := &MockServiceController{
        CreateServiceFunc: func(ctx context.Context, service *Service, user *UserContext) (*Service, error) {
            service.Metadata.CreatedAt = time.Now()
            return service, nil
        },
    }
    
    handler := NewHTTPHandler(mockController, NewServiceAdapter(), responseAdapter, requestAdapter)
    
    requestBody := `{
        "name": "test-service",
        "tier": "TIER-3",
        "description": "Test service"
    }`
    
    req := httptest.NewRequest("POST", "/api/v1/services", strings.NewReader(requestBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    
    // Act
    handler.createServiceHandler(w, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    assert.Contains(t, w.Header().Get("Location"), "/services/test-service")
    
    var response wire.ServiceResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "test-service", response.Name)
}
```

### Mock Creation Guidelines

```go
// Create focused mocks for each interface
type MockServiceRepository struct {
    CreateFunc     func(ctx context.Context, service *Service) (*Service, error)
    GetByNameFunc  func(ctx context.Context, name string) (*Service, error)
    ExistsFunc     func(ctx context.Context, name string) (bool, error)
    // ... other methods
}

func (m *MockServiceRepository) Create(ctx context.Context, service *Service) (*Service, error) {
    if m.CreateFunc != nil {
        return m.CreateFunc(ctx, service)
    }
    return service, nil // Default behavior
}

func (m *MockServiceRepository) Exists(ctx context.Context, name string) (bool, error) {
    if m.ExistsFunc != nil {
        return m.ExistsFunc(ctx, name)
    }
    return false, nil // Default behavior
}
```

### Test Coverage Requirements

| Layer | Coverage Requirement | Testing Type |
|-------|---------------------|--------------|
| **logic/** | 100% | Unit tests |
| **models/** | 100% for methods | Unit tests |
| **controllers/** | 80%+ | Integration with mocks |
| **handlers/** | 70%+ | HTTP integration tests |
| **adapters/** | 80%+ | Unit tests |
| **ports/** | N/A | Interfaces only |
| **wire/** | N/A | DTOs only |

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./pkg/service-catalog/...

# Run tests with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestServiceProcessor_PrepareForCreation ./pkg/service-catalog/logic
```

## Dependency Injection

### Principles

1. **Depend on interfaces, not implementations**
2. **Define interfaces in the dependent module**
3. **Inject dependencies through constructors**
4. **Wire dependencies in main.go or module factory**
5. **Load cross-module dependencies after all modules are initialized**

### Module Initialization Phases

DashOps follows a **4-phase initialization pattern** to handle complex dependencies:

#### Phase 1: Module Initialization
All modules are initialized with `NewModule(fileConfig []byte)`:

```go
// main.go - Phase 1
modules := make(map[string]interface{})

if dashConfig.Plugins.Has("Auth") {
    authModule, err := auth.NewModule(fileConfig)
    if err != nil {
        log.Fatalf("Failed to create auth module: %v", err)
    }
    modules["auth"] = authModule
}

if dashConfig.Plugins.Has("ServiceCatalog") {
    scModule, err := servicecatalog.NewModule(fileConfig)
    if err != nil {
        log.Fatalf("Failed to create service catalog module: %v", err)
    }
    modules["service-catalog"] = scModule
}

if dashConfig.Plugins.Has("Kubernetes") {
    k8sModule, err := kubernetes.NewModule(fileConfig)
    if err != nil {
        log.Printf("Failed to create kubernetes module: %v", err)
    } else {
        modules["kubernetes"] = k8sModule
    }
}
```

#### Phase 2: Route Registration
All modules register their routes:

```go
// main.go - Phase 2
for name, module := range modules {
    if m, ok := module.(interface{ RegisterRoutes(*mux.Router) }); ok {
        m.RegisterRoutes(internal)
    } else if m, ok := module.(interface{ RegisterRoutes(*mux.Router, *mux.Router) }); ok {
        m.RegisterRoutes(api, internal)
    }
    log.Printf("Routes registered for %s module", name)
}
```

#### Phase 3: Cross-Module Dependencies
Load dependencies between modules:

```go
// main.go - Phase 3
for name, module := range modules {
    if m, ok := module.(interface{ LoadDependencies(map[string]interface{}) error }); ok {
        if err := m.LoadDependencies(modules); err != nil {
            log.Printf("Warning: Failed to load dependencies for %s module: %v", name, err)
        } else {
            log.Printf("Dependencies loaded for %s module", name)
        }
    }
}
```

#### Phase 4: SPA Initialization
Initialize SPA module to serve static files:

```go
// main.go - Phase 4
spaConfig := &spaModels.SPAConfig{
    StaticPath: dashConfig.Front,
    IndexPath:  "index.html",
}
spaModule, err := spa.NewModule(spaConfig, api)
if err != nil {
    log.Fatalf("Failed to create SPA module: %v", err)
}
spaModule.RegisterRoutes(router)
```

### Cross-Module Dependencies

Modules can have dependencies on other modules. These are loaded in Phase 3 using the `LoadDependencies` method:

#### Example: Service-Catalog ↔ Kubernetes Integration

**Service-Catalog Module**:
```go
// pkg/service-catalog/module.go
func (m *Module) LoadDependencies(modules map[string]interface{}) error {
    // Load kubernetes dependency if available
    if k8sModule, exists := modules["kubernetes"]; exists {
        if k8s, ok := k8sModule.(interface{ GetExposedAPI() k8sPorts.ExposedAPI }); ok {
            if api := k8s.GetExposedAPI(); api != nil {
                // Create internal integration
                k8sIntegration := integrations.NewKubernetesIntegration(api)
                m.controller.SetKubernetesIntegration(k8sIntegration)
            }
        }
    }
    return nil
}

// Expose API for other modules
func (m *Module) GetExposedAPI() ports.ExposedAPI {
    return &exposedAPI{
        controller: m.controller,
    }
}
```

**Kubernetes Module**:
```go
// pkg/kubernetes/module.go
func (m *Module) LoadDependencies(modules map[string]interface{}) error {
    // Load service-catalog dependency if available
    if scModule, exists := modules["service-catalog"]; exists {
        if sc, ok := scModule.(interface{ GetExposedAPI() scPorts.ExposedAPI }); ok {
            if api := sc.GetExposedAPI(); api != nil {
                // Create internal integration
                scIntegration := integrations.NewServiceCatalogIntegration(api)
                m.controller.SetServiceCatalogIntegration(scIntegration)
            }
        }
    }
    return nil
}

// Expose API for other modules
func (m *Module) GetExposedAPI() ports.ExposedAPI {
    return &exposedAPI{
        controller: m.controller,
    }
}
```

### External Service Dependencies

External services are injected through the module factory:

```go
// pkg/kubernetes/module.go
func NewModule(config []byte) (*Module, error) {
    // Create external integrations
    prometheusClient := prometheus.NewPrometheusClient(config.PrometheusEndpoint)
    lokiClient := loki.NewLokiClient(config.LokiEndpoint)
    tempoClient := tempo.NewTempoClient(config.TempoEndpoint)
    
    // Create adapters for external integrations
    prometheusAdapter := adapters.NewPrometheusAdapter(prometheusClient)
    lokiAdapter := adapters.NewLokiAdapter(lokiClient)
    tempoAdapter := adapters.NewTempoAdapter(tempoClient)
    
    // Inject into controller
    controller := controllers.NewKubernetesController(
        prometheusAdapter,
        lokiAdapter,
        tempoAdapter,
        // ... other dependencies
    )
    
    return &Module{
        controller: controller,
    }, nil
}
```

#### Benefits of This Pattern

1. **No Circular Dependencies**: All modules are initialized before dependencies are loaded
2. **Graceful Degradation**: Modules can run without their dependencies
3. **Plugin-Ready**: Easy to add/remove modules without breaking others
4. **Clear Separation**: Dependencies are explicit and manageable
5. **Testable**: Each phase can be tested independently

### ✅ Correct Patterns

#### Define Interface in Dependent Module

```go
// auth/ports/services.go
package ports

type GitHubService interface {
    GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error)
    GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error)
}
```

#### Use Interface in Controller

```go
// auth/controllers/auth_controller.go
type AuthController struct {
    githubService ports.GitHubService // Interface, not concrete type
}

func NewAuthController(githubService ports.GitHubService) *AuthController {
    return &AuthController{
        githubService: githubService,
    }
}
```

#### Wire Dependencies in main.go

```go
// main.go
func main() {
    // Create modules with dependencies
    githubModule, _ := github.NewModule(githubConfig)
    authModule, _ := auth.NewModule(authConfig, githubModule)
    
    // Service catalog depends on Kubernetes
    k8sModule, _ := kubernetes.NewModule(k8sConfig)
    scModule, _ := servicecatalog.NewModule(scConfig, k8sModule)
    
    // Register routes
    authModule.RegisterRoutes(router)
    scModule.RegisterRoutes(router)
}
```

### ❌ Anti-Patterns to Avoid

```go
// ❌ WRONG: Direct import of another module
import "github.com/dash-ops/dash-ops/pkg/github"

// ❌ WRONG: Concrete type instead of interface
type AuthController struct {
    githubController *github.GitHubController // Tight coupling!
}

// ❌ WRONG: Creating dependencies internally
func NewAuthController() *AuthController {
    githubClient := github.NewClient() // Should be injected!
    return &AuthController{
        githubClient: githubClient,
    }
}
```

## Architecture Standards

### 🎯 Controller Organization

DashOps follows a **domain-specific controller architecture** that ensures clean separation of concerns and maintainability.

#### Domain-Specific Controllers

Each module uses separate controllers for different domains, following the Single Responsibility Principle:

```go
// Controllers are organized by domain
type LogsController struct {
    logRepo      ports.LogRepository
    serviceRepo  ports.ServiceContextRepository
    logService   ports.LogService
    cacheService ports.CacheService
    processor    *logic.LogProcessor
}

type MetricsController struct {
    metricRepo    ports.MetricRepository
    serviceRepo   ports.ServiceContextRepository
    metricService ports.MetricService
    cacheService  ports.CacheService
    processor     *logic.MetricProcessor
}

type TracesController struct {
    traceRepo    ports.TraceRepository
    serviceRepo  ports.ServiceContextRepository
    traceService ports.TraceService
    cacheService ports.CacheService
    processor    *logic.TraceProcessor
}
```

**Benefits**:
- Each controller has a single, clear responsibility
- Easy to test and maintain
- Independent evolution of domains
- Better code organization and readability

### 📁 Model Organization

Models are organized by domain context for optimal maintainability and clarity:

```
models/
├── common.go           # BaseResponse, BaseMetadata
├── logs.go            # LogEntry, ProcessedLogEntry, LogsConfig
├── metrics.go         # MetricData, ProcessedMetric, DerivedMetric, MetricsConfig
├── traces.go          # TraceSpan, ProcessedTrace, TraceInfo, TracePerformance, TracesConfig
├── alerts.go          # Alert, ProcessedAlert, AlertRule, AlertEvaluation, AlertsConfig
├── dashboards.go      # Dashboard, Chart, DashboardTemplate, DashboardData, ChartData
├── service_context.go # ServiceContext, ServiceWithContext, ServiceHealth
├── config.go          # ObservabilityConfig, ServiceObservabilityConfig, CacheConfig, UIConfig, CacheStats
└── notifications.go   # NotificationChannel
```

**Benefits**:
- Easy to find related models
- Reduced file size and complexity
- Domain-focused development
- Minimal merge conflicts
- Clear domain boundaries

### 🔌 Wire Organization

Wire (API DTOs) are organized by domain context with clear request/response separation:

```
wire/
├── common.go           # BaseResponse, TimeSeriesData, ErrorResponse, PaginationInfo
├── logs.go            # LogsRequest, LogsResponse, LogStatsRequest, LogStatisticsResponse
├── metrics.go         # MetricsRequest, MetricsResponse, PrometheusQueryRequest, MetricStatsRequest
├── traces.go          # TracesRequest, TracesResponse, TraceDetailRequest, TraceAnalysisRequest
├── alerts.go          # AlertsRequest, AlertsResponse, CreateAlertRequest, AlertRulesRequest
├── dashboards.go      # DashboardsRequest, DashboardsResponse, CreateDashboardRequest, DashboardTemplatesRequest
├── service_context.go # ServiceContextRequest, ServiceContextResponse, ServicesWithContextRequest
├── config.go          # ConfigurationRequest, ConfigurationResponse, NotificationChannelsRequest
└── health.go          # HealthRequest, HealthResponse, CacheStatsRequest, CacheStatsResponse
```

**Naming Convention**:
- **Requests**: `{Domain}Request` (e.g., `LogsRequest`, `MetricsRequest`)
- **Responses**: `{Domain}Response` (e.g., `LogsResponse`, `MetricsResponse`)
- **Data Types**: `{Domain}Data` (e.g., `LogsData`, `MetricsData`)
- **Statistics**: `{Domain}Statistics` (e.g., `LogStatistics`, `MetricStatistics`)

**Benefits**:
- Clear separation of concerns
- Easy to find related wire types
- Consistent naming across domains
- Reduced file size and complexity
- Domain-focused development
- Minimal merge conflicts
- Shared common types in `common.go`

### 🏗️ Module Structure

Modules follow a simplified, focused structure that includes only essential components:

```go
// module.go - Clean and focused structure
type Module struct {
    // Core components
    LogsController    *controllers.LogsController
    MetricsController *controllers.MetricsController
    TracesController  *controllers.TracesController
    AlertsController  *controllers.AlertsController
    HealthController  *controllers.HealthController
    ConfigController  *controllers.ConfigController
    Handler           *handlers.HTTPHandler

    // Logic components
    LogProcessor       *logic.LogProcessor
    MetricProcessor    *logic.MetricProcessor
    TraceProcessor     *logic.TraceProcessor
    AlertProcessor     *logic.AlertProcessor
    DashboardProcessor *logic.DashboardProcessor

    // Adapters
    ResponseAdapter *commonsHttp.ResponseAdapter
    RequestAdapter  *commonsHttp.RequestAdapter

    // Repositories (interfaces - implementations injected)
    LogRepo       ports.LogRepository
    MetricRepo    ports.MetricRepository
    TraceRepo     ports.TraceRepository
    AlertRepo     ports.AlertRepository
    DashboardRepo ports.DashboardRepository
    ServiceRepo   ports.ServiceContextRepository

    // Services (interfaces - implementations injected)
    LogService           ports.LogService
    MetricService        ports.MetricService
    TraceService         ports.TraceService
    AlertService         ports.AlertService
    DashboardService     ports.DashboardService
    NotificationService  ports.NotificationService
    CacheService         ports.CacheService
    ConfigurationService ports.ConfigurationService
}

// NewModule creates and initializes a new module
func NewModule(config *ModuleConfig) (*Module, error) {
    // Validate required dependencies
    if config.LogRepo == nil {
        return nil, fmt.Errorf("log repository is required")
    }
    // ... other validations

    // Initialize logic components
    logProcessor := logic.NewLogProcessor()
    metricProcessor := logic.NewMetricProcessor()
    traceProcessor := logic.NewTraceProcessor()
    alertProcessor := logic.NewAlertProcessor()
    dashboardProcessor := logic.NewDashboardProcessor()

    // Initialize adapters
    responseAdapter := commonsHttp.NewResponseAdapter()
    requestAdapter := commonsHttp.NewRequestAdapter()

    // Initialize per-domain controllers
    logsController := controllers.NewLogsController(
        config.LogRepo,
        config.ServiceRepo,
        config.LogService,
        config.CacheService,
        logProcessor,
    )

    metricsController := controllers.NewMetricsController(
        config.MetricRepo,
        config.ServiceRepo,
        config.MetricService,
        config.CacheService,
        metricProcessor,
    )

    // ... other controllers

    // Initialize handler
    handler := handlers.NewHTTPHandler(
        logsController,
        metricsController,
        tracesController,
        alertsController,
        healthController,
        configController,
        responseAdapter,
        requestAdapter,
    )

    return &Module{
        LogsController:       logsController,
        MetricsController:    metricsController,
        TracesController:     tracesController,
        AlertsController:     alertsController,
        HealthController:     healthController,
        ConfigController:     configController,
        Handler:              handler,
        // ... other fields
    }, nil
}

// Essential getters
func (m *Module) GetHandler() *handlers.HTTPHandler {
    return m.Handler
}

func (m *Module) GetLogsController() *controllers.LogsController {
    return m.LogsController
}

// ... other getters

// RegisterRoutes registers HTTP routes for the module
func (m *Module) RegisterRoutes(router *mux.Router) {
    if m.Handler == nil {
        return
    }
    m.Handler.RegisterRoutes(router)
}
```

### 🎯 Architecture Principles

1. **Single Responsibility**: Each component has one clear purpose
2. **Dependency Injection**: Inject dependencies, don't create them
3. **Interface Segregation**: Use focused interfaces
4. **Clean Architecture**: Clear separation of concerns
5. **Testability**: Easy to test with mocks
6. **Maintainability**: Easy to understand and modify

## Module-Specific Guides

### 📦 Service Catalog Module

**Purpose**: Service registry with Kubernetes integration

```go
// Key entities
type Service struct {
    Metadata ServiceMetadata
    Spec     ServiceSpec
}

type ServiceMetadata struct {
    Name      string
    Tier      ServiceTier  // TIER-1, TIER-2, TIER-3
    CreatedAt time.Time
    CreatedBy string
    Version   int
}

// Business rules
- Services must have unique names
- Production services (TIER-1) require team assignment
- Versioning tracks all changes
- Health status aggregated from Kubernetes
```

### ☸️ Kubernetes Module

**Purpose**: Container orchestration and monitoring

```go
// Key operations
- Multi-cluster management
- Deployment health monitoring
- Resource metrics collection
- Service context resolution

// Integration with Service Catalog
type ServiceContextResolver interface {
    ResolveDeploymentService(ctx context.Context, 
        deployment, namespace, context string) (*ServiceContext, error)
}
```

### 🔐 Auth Module

**Purpose**: Multi-provider authentication

```go
// Supported providers
- GitHub Auth
- Google Auth (planned)
- SAML (planned)
- LDAP (planned)

// Key features
- Session management
- Team-based permissions
- Token validation
- User context propagation
```

### ☁️ AWS Module

**Purpose**: Cloud infrastructure management

```go
// Current features
- EC2 instance management
- Multi-account support
- Cost optimization
- Permission system

// Planned features
- S3 bucket management
- RDS database management
- Lambda function management
```

### 📊 Observability Module

**Purpose**: Unified observability with logs, metrics, traces, and alerts

```go
// Architecture: Domain-specific controllers
type LogsController struct {
    logRepo      ports.LogRepository
    serviceRepo  ports.ServiceContextRepository
    logService   ports.LogService
    cacheService ports.CacheService
    processor    *logic.LogProcessor
}

type MetricsController struct {
    metricRepo    ports.MetricRepository
    serviceRepo   ports.ServiceContextRepository
    metricService ports.MetricService
    cacheService  ports.CacheService
    processor     *logic.MetricProcessor
}

// Model organization by context
models/
├── common.go          # BaseResponse, BaseMetadata
├── logs.go            # LogEntry, ProcessedLogEntry, LogsConfig
├── metrics.go         # MetricData, ProcessedMetric, MetricsConfig
├── traces.go          # TraceSpan, ProcessedTrace, TraceInfo, TracesConfig
├── alerts.go          # Alert, ProcessedAlert, AlertRule, AlertsConfig
├── dashboards.go      # Dashboard, Chart, DashboardTemplate
├── service_context.go # ServiceContext, ServiceHealth
├── config.go          # ObservabilityConfig, CacheConfig
└── notifications.go   # NotificationChannel

// Wire organization by context
wire/
├── common.go          # BaseResponse, TimeSeriesData, ErrorResponse
├── logs.go            # LogsRequest, LogsResponse, LogStatsRequest
├── metrics.go         # MetricsRequest, MetricsResponse, PrometheusQueryRequest
├── traces.go          # TracesRequest, TracesResponse, TraceDetailRequest
├── alerts.go          # AlertsRequest, AlertsResponse, CreateAlertRequest
├── dashboards.go      # DashboardsRequest, DashboardsResponse, CreateDashboardRequest
├── service_context.go # ServiceContextRequest, ServiceContextResponse
├── config.go          # ConfigurationRequest, ConfigurationResponse
└── health.go          # HealthRequest, HealthResponse, CacheStatsRequest

// Key features
- Multi-source log aggregation (Loki integration)
- Metrics collection and analysis (Prometheus integration)
- Distributed tracing (Tempo integration)
- Alert management (Alertmanager integration)
- Dashboard creation and management
- Service context awareness
- Real-time streaming capabilities
- Caching for performance optimization
```

**Architecture**:
- Domain-specific controllers (6 focused controllers)
- Context-specific model files (9 organized files)
- Context-specific wire files (9 organized files)
- Simplified module structure (essential methods only)
- Clean Architecture principles
- Single Responsibility Principle
- Easy testing and maintenance

## Best Practices

### Code Organization

```go
// ✅ DO: Group related functionality
// logic/service_validator.go
type ServiceValidator struct{}

func (sv *ServiceValidator) ValidateForCreation(service *Service) error {}
func (sv *ServiceValidator) ValidateForUpdate(old, new *Service) error {}
func (sv *ServiceValidator) ValidateForDeletion(service *Service) error {}

// ❌ DON'T: Mix concerns
// logic/mixed.go
func ValidateService(service *Service) error {}
func ProcessPayment(payment *Payment) error {} // Different concern!
```

### Error Handling

```go
// ✅ DO: Wrap errors with context
func (sc *ServiceController) CreateService(ctx context.Context, service *Service) (*Service, error) {
    if err := sc.validator.Validate(service); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    created, err := sc.repo.Create(ctx, service)
    if err != nil {
        return nil, fmt.Errorf("failed to create service '%s': %w", service.Name, err)
    }
    
    return created, nil
}

// ❌ DON'T: Return errors without context
func (sc *ServiceController) CreateService(ctx context.Context, service *Service) (*Service, error) {
    if err := sc.validator.Validate(service); err != nil {
        return nil, err // No context!
    }
}
```

### Logging

```go
// ✅ DO: Log at appropriate levels with context
func (h *HTTPHandler) createServiceHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Creating service: method=%s path=%s", r.Method, r.URL.Path)
    
    service, err := h.controller.CreateService(r.Context(), req)
    if err != nil {
        log.Printf("ERROR: Failed to create service: %v", err)
        h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    log.Printf("Service created successfully: name=%s", service.Name)
}
```

### Performance

```go
// ✅ DO: Use context for cancellation
func (c *Controller) ProcessBatch(ctx context.Context, items []Item) error {
    for _, item := range items {
        select {
        case <-ctx.Done():
            return ctx.Err() // Respect cancellation
        default:
            if err := c.processItem(ctx, item); err != nil {
                return err
            }
        }
    }
    return nil
}

// ✅ DO: Batch operations when possible
func (r *Repository) GetMultiple(ctx context.Context, ids []string) ([]Service, error) {
    // Single query instead of N queries
    return r.db.Query("SELECT * FROM services WHERE id IN (?)", ids)
}
```

## Common Patterns

### Repository Pattern

```go
// ports/repositories.go
type ServiceRepository interface {
    Create(ctx context.Context, service *Service) (*Service, error)
    GetByName(ctx context.Context, name string) (*Service, error)
    Update(ctx context.Context, service *Service) (*Service, error)
    Delete(ctx context.Context, name string) error
}

// adapters/storage/filesystem_repository.go
type FilesystemRepository struct {
    directory string
}

func (fr *FilesystemRepository) Create(ctx context.Context, service *Service) (*Service, error) {
    // Implementation
}
```

### Factory Pattern

```go
// module.go
type Module struct {
    controller *ServiceController
    handler    *HTTPHandler
}

func NewModule(config *Config, deps Dependencies) (*Module, error) {
    // Create components
    processor := logic.NewServiceProcessor()
    validator := logic.NewServiceValidator()
    
    // Create repository based on config
    repo, err := createRepository(config)
    if err != nil {
        return nil, err
    }
    
    // Create controller
    controller := controllers.NewServiceController(
        repo,
        deps.KubernetesService,
        processor,
        validator,
    )
    
    // Create handler
    handler := handlers.NewHTTPHandler(controller, ...)
    
    return &Module{
        controller: controller,
        handler:    handler,
    }, nil
}
```

### Adapter Pattern

```go
// adapters/external/kubernetes_adapter.go
type KubernetesAdapter struct {
    client kubernetes.Interface
}

func (ka *KubernetesAdapter) GetDeploymentHealth(ctx context.Context, namespace, name string) (*DeploymentHealth, error) {
    // Transform Kubernetes API response to domain model
    deployment, err := ka.client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
    if err != nil {
        return nil, err
    }
    
    return &DeploymentHealth{
        Name:            deployment.Name,
        ReadyReplicas:   deployment.Status.ReadyReplicas,
        DesiredReplicas: *deployment.Spec.Replicas,
        Status:          ka.calculateStatus(deployment),
    }, nil
}
```

### Strategy Pattern

```go
// ports/versioning.go
type VersioningStrategy interface {
    RecordChange(ctx context.Context, service *Service, user *UserContext) error
    GetHistory(ctx context.Context, serviceName string) ([]Change, error)
}

// adapters/versioning/git_strategy.go
type GitStrategy struct {
    repo *git.Repository
}

func (gs *GitStrategy) RecordChange(ctx context.Context, service *Service, user *UserContext) error {
    // Git-based implementation
}

// adapters/versioning/simple_strategy.go
type SimpleStrategy struct {
    storage Storage
}

func (ss *SimpleStrategy) RecordChange(ctx context.Context, service *Service, user *UserContext) error {
    // Simple file-based implementation
}
```

## Troubleshooting

### Common Issues and Solutions

#### 1. Circular Dependencies

**Problem**: Module A depends on Module B, and Module B depends on Module A

**Solution**: Use interfaces and dependency inversion
```go
// Define interface in the dependent module
// service-catalog/ports/services.go
type KubernetesService interface {
    GetDeploymentHealth(ctx context.Context, namespace, name string) (*Health, error)
}

// Implement in the provider module
// kubernetes/adapters/service_catalog_adapter.go
type ServiceCatalogAdapter struct {
    // Implementation
}
```

#### 2. Test Failures After Refactoring

**Problem**: Tests fail after changing implementation

**Solution**: Tests should test behavior, not implementation
```go
// ✅ Good: Testing behavior
func TestServiceValidator_ValidateForCreation_RejectsEmptyName(t *testing.T) {
    validator := NewServiceValidator()
    service := &Service{Metadata: ServiceMetadata{Name: ""}}
    
    err := validator.ValidateForCreation(service)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "name is required")
}
```

#### 3. Complex Mock Setup

**Problem**: Tests require complex mock setup

**Solution**: Use builder pattern for test data
```go
// test/builders/service_builder.go
type ServiceBuilder struct {
    service Service
}

func NewServiceBuilder() *ServiceBuilder {
    return &ServiceBuilder{
        service: Service{
            Metadata: ServiceMetadata{
                Name: "default-service",
                Tier: TierStandard,
            },
        },
    }
}

func (sb *ServiceBuilder) WithName(name string) *ServiceBuilder {
    sb.service.Metadata.Name = name
    return sb
}

func (sb *ServiceBuilder) WithTier(tier ServiceTier) *ServiceBuilder {
    sb.service.Metadata.Tier = tier
    return sb
}

func (sb *ServiceBuilder) Build() *Service {
    return &sb.service
}

// Usage in tests
service := NewServiceBuilder().
    WithName("test-service").
    WithTier(TierCritical).
    Build()
```

## Quick Reference

### File Templates

#### Domain-Specific Controller Template
```go
package controllers

import (
    "context"
    "github.com/dash-ops/dash-ops/pkg/{module}/logic"
    "github.com/dash-ops/dash-ops/pkg/{module}/ports"
    "github.com/dash-ops/dash-ops/pkg/{module}/wire"
)

// {Domain}Controller handles {domain} operations
type {Domain}Controller struct {
    {domain}Repo    ports.{Domain}Repository
    serviceRepo     ports.ServiceContextRepository
    {domain}Service ports.{Domain}Service
    cacheService    ports.CacheService
    processor       *logic.{Domain}Processor
}

// New{Domain}Controller creates a new {domain} controller
func New{Domain}Controller(
    {domain}Repo ports.{Domain}Repository,
    serviceRepo ports.ServiceContextRepository,
    {domain}Service ports.{Domain}Service,
    cacheService ports.CacheService,
    processor *logic.{Domain}Processor,
) *{Domain}Controller {
    return &{Domain}Controller{
        {domain}Repo:    {domain}Repo,
        serviceRepo:     serviceRepo,
        {domain}Service: {domain}Service,
        cacheService:    cacheService,
        processor:       processor,
    }
}

// Get{Domain} retrieves {domain} data
func (c *{Domain}Controller) Get{Domain}(ctx context.Context, req *wire.{Domain}Request) (*wire.{Domain}Response, error) {
    // Implementation
    return nil, nil
}
```

#### Context-Specific Model File Template
```go
package models

import (
    "time"
)

// {Domain}Entity represents a {domain} entity
type {Domain}Entity struct {
    ID        string                 `json:"id"`
    Name      string                 `json:"name"`
    Timestamp time.Time              `json:"timestamp"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Processed{Domain}Entity represents a processed {domain} entity
type Processed{Domain}Entity struct {
    {Domain}Entity
    ProcessedAt time.Time              `json:"processed_at"`
    Enrichments map[string]interface{} `json:"enrichments,omitempty"`
}

// {Domain}Config represents {domain} configuration
type {Domain}Config struct {
    Enabled   bool     `json:"enabled"`
    Retention string   `json:"retention"`
    Settings  []string `json:"settings"`
}
```

#### New Logic File
```go
package logic

// {Name}Processor handles {description}
type {Name}Processor struct{}

// New{Name}Processor creates a new processor
func New{Name}Processor() *{Name}Processor {
    return &{Name}Processor{}
}

// Process{Operation} performs {description}
func (p *{Name}Processor) Process{Operation}(input *Model) (*Model, error) {
    // Implementation
    return input, nil
}
```

#### External Integration Template
```go
// integrations/external/{service}/{service}_client.go
package {service}

import (
    "context"
    // External service imports
)

// {Service}Client handles communication with {Service} API
type {Service}Client struct {
    client *external.Client
    config *Config
}

func New{Service}Client(config *Config) *{Service}Client {
    return &{Service}Client{
        client: external.NewClient(config),
        config: config,
    }
}

func (c *{Service}Client) Get{Resource}(ctx context.Context, id string) (*Resource, error) {
    // API call implementation
    return c.client.GetResource(ctx, id)
}

// integrations/external/{service}/{service}_adapter.go
package {service}

import (
    "context"
    "github.com/dash-ops/dash-ops/pkg/{module}/ports"
)

// {Service}Adapter transforms data between {Service} API and domain
type {Service}Adapter struct {
    client *{Service}Client
}

func New{Service}Adapter(client *{Service}Client) ports.{Service}Service {
    return &{Service}Adapter{
        client: client,
    }
}

func (a *{Service}Adapter) Get{Resource}(ctx context.Context, id string) (*models.Resource, error) {
    // Call client
    resource, err := a.client.Get{Resource}(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Transform data
    return a.transform{Resource}(resource), nil
}

func (a *{Service}Adapter) transform{Resource}(resource *external.Resource) *models.Resource {
    // Pure data transformation
    return &models.Resource{
        ID:   resource.ID,
        Name: resource.Name,
        // ... other fields
    }
}
```

#### Internal Integration Template
```go
// integrations/{module}/{module}_integration.go
package {module}

import (
    "context"
    "github.com/dash-ops/dash-ops/pkg/{module}/ports"
    "github.com/dash-ops/dash-ops/pkg/{module}/models"
)

// {Module}Integration handles communication with {module} module
type {Module}Integration struct {
    api ports.ExposedAPI
}

func New{Module}Integration(api ports.ExposedAPI) *{Module}Integration {
    return &{Module}Integration{
        api: api,
    }
}

func (i *{Module}Integration) Get{Resource}(ctx context.Context, id string) (*models.Resource, error) {
    return i.api.Get{Resource}(ctx, id)
}

func (i *{Module}Integration) Create{Resource}(ctx context.Context, resource *models.Resource) (*models.Resource, error) {
    return i.api.Create{Resource}(ctx, resource)
}
```

#### New Test File
```go
package {module}

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func Test{Name}Processor_{Method}_With{Scenario}_Returns{Result}(t *testing.T) {
    // Arrange
    processor := New{Name}Processor()
    input := &Model{/* test data */}
    
    // Act
    result, err := processor.{Method}(input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    // Additional assertions
}
```

### Command Cheatsheet

```bash
# Development
go run main.go -config ../workspace-dash-ops.yaml  # Run with workspace config
go build                                            # Build binary
go mod tidy                                        # Clean dependencies

# Testing
go test ./...                                      # Run all tests
go test -cover ./...                              # Run with coverage
go test -v ./pkg/{module}/...                    # Test specific module
go test -run {TestName} ./...                    # Run specific test

# Code Quality
go fmt ./...                                      # Format code
go vet ./...                                      # Run static analysis
golangci-lint run                                # Run linter (if installed)

# Coverage Reports
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out                 # Open HTML report
go tool cover -func=coverage.out                 # Terminal report
```

---

**📚 Additional Resources**

- [CLAUDE.md](../CLAUDE.md) - AI assistant context
- [ROADMAP.md](../ROADMAP.md) - Project vision and roadmap
- [Frontend Guide](./frontend-development-guide.md) - Frontend development

**Happy Coding! 🚀** Remember: Consistency is key. When in doubt, follow the existing patterns in the codebase.
