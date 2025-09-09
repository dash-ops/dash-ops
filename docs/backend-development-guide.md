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

Every module follows this consistent 8-layer pattern:

```
pkg/{module}/
├── adapters/           # Data transformation & external integrations
│   ├── http/          # HTTP request/response adapters
│   ├── storage/       # Database/filesystem adapters
│   └── external/      # Third-party service adapters
├── controllers/       # Business logic orchestration
├── handlers/          # HTTP endpoints (entry points)
├── logic/             # Pure business logic (100% tested)
├── models/            # Domain entities with behavior
├── ports/             # Interfaces & contracts
├── wire/              # API DTOs (request/response)
└── module.go          # Module factory & initialization
```

### Layer Responsibilities

| Layer | Purpose | Testing Requirements |
|-------|---------|---------------------|
| **handlers** | HTTP endpoints, routing | Integration tests |
| **controllers** | Orchestration, workflow | Integration tests with mocks |
| **logic** | Pure business rules | 100% unit test coverage |
| **models** | Domain entities | Unit tests for methods |
| **adapters** | Data transformation | Unit tests |
| **ports** | Interfaces | No tests (interfaces) |
| **wire** | DTOs | No tests (data structures) |

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

#### New Logic File
```go
package {module}

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
