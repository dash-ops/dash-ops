# Backend Architecture Guide

> **🏗️ Current Status**: **4 modules migrated** to pure Hexagonal Architecture with **zero legacy code**. Remaining 4 modules pending migration.

## Quick Start for Developers

The DashOps backend is organized in 8 modules, each following the same architectural pattern. This consistency makes it easy to navigate, understand, and contribute to any module.

### 🗂️ Module Structure

```
pkg/{module-name}/
├── adapters/           # Data transformation & external integrations
│   ├── http/          # HTTP request/response adapters
│   ├── storage/       # Database/filesystem adapters
│   └── external/      # Third-party service adapters (AWS, K8s, GitHub)
├── config/            # Configuration management
├── controllers/       # Business logic orchestration
├── handlers/          # HTTP endpoints (entry points)
├── logic/             # Pure business logic (100% tested)
├── models/            # Domain entities with behavior
├── ports/             # Interfaces & contracts
├── wire/              # API DTOs (request/response)
└── module.go          # Module factory & initialization
```

## 🚀 Available Modules

### Core Modules

| Module      | Purpose                  | Status      | Key Features                              |
| ----------- | ------------------------ | ----------- | ----------------------------------------- |
| **commons** | Shared utilities         | ✅ Complete | HTTP adapters, permissions, string utils  |
| **config**  | Configuration management | ✅ Complete | Dynamic reload, validation, env vars      |
| **auth**    | Authentication           | ✅ Complete | Multi-provider OAuth2, session management |
| **github**  | Git integration          | ✅ Complete | Teams, repositories, permissions          |
| **spa**     | Static file serving      | ✅ Complete | SPA routing, security headers, compression |

### Platform Modules

| Module              | Purpose                 | Status          | Key Features                                     |
| ------------------- | ----------------------- | --------------- | ------------------------------------------------ |
| **kubernetes**      | Container orchestration | ✅ Complete     | Clusters, deployments, pods, health monitoring   |
| **aws**             | Cloud infrastructure    | ✅ Complete     | EC2 management, cost optimization, multi-account |
| **service-catalog** | Service registry        | ✅ Complete     | CRUD, versioning, K8s integration                |

## 🏗️ **Dependency Injection Best Practices**

### **✅ DO: Use Interfaces for Dependencies**

```go
// ✅ CORRECT: Define interface in dependent module
package auth

// ports/services.go
type GitHubService interface {
    GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error)
    GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error)
}

// controllers/auth_controller.go
type AuthController struct {
    githubService ports.GitHubService // Interface, not concrete type
}
```

### **✅ DO: Inject Dependencies in main.go**

```go
// ✅ CORRECT: main.go orchestrates all dependencies
func main() {
    // 1. Initialize dependency first
    githubModule, _ := github.NewModule(oauthConfig)
    
    // 2. Inject dependency into dependent module
    authModule, _ := auth.NewModule(authConfig, githubModule)
    
    // 3. Register routes
    authModule.RegisterRoutes(api, internal)
}
```

### **❌ DON'T: Import Business Modules Directly**

```go
// ❌ WRONG: Creates tight coupling
package auth

import "github.com/dash-ops/dash-ops/pkg/github" // BAD!

type AuthController struct {
    githubController *github.Controller // Tight coupling!
}
```

### **📋 DIP Checklist**

- [ ] **No direct imports** between business modules
- [ ] **Interfaces defined** in dependent module's `ports/`
- [ ] **Dependencies injected** in `main.go` or module factory
- [ ] **Easy to mock** for testing
- [ ] **No circular dependencies** between modules

## 📝 Development Guidelines

### Creating a New Feature

1. **Choose the Right Layer**:

   - **Logic**: Pure business functions (always test these!)
   - **Controllers**: Orchestration between components
   - **Adapters**: Data transformation or external calls
   - **Handlers**: HTTP endpoint logic

2. **Follow the Data Flow**:

   ```
   HTTP Request → Handler → Adapter → Controller → Logic → Repository
   ```

3. **Testing Strategy**:
   - **Unit Tests**: Logic and Models (aim for 100% coverage)
   - **Integration Tests**: Handlers (full HTTP cycle)
   - **Mocks**: Use interfaces for external dependencies

### Example: Adding a New Endpoint

```go
// 1. Define Wire DTO (wire/requests.go)
type CreateResourceRequest struct {
    Name string `json:"name" validate:"required"`
}

// 2. Add Handler method (handlers/http.go)
func (h *HTTPHandler) createResourceHandler(w http.ResponseWriter, r *http.Request) {
    var req wire.CreateResourceRequest
    if err := h.requestAdapter.ParseJSON(r, &req); err != nil {
        h.responseAdapter.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }

    // Transform to domain model
    resource, err := h.adapter.RequestToModel(req)
    if err != nil {
        h.responseAdapter.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }

    // Call controller
    result, err := h.controller.CreateResource(r.Context(), resource)
    if err != nil {
        h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }

    // Transform and respond
    response := h.adapter.ModelToResponse(result)
    h.responseAdapter.WriteCreated(w, "/resources/"+result.ID, response)
}

// 3. Add Controller method (controllers/resource_controller.go)
func (c *ResourceController) CreateResource(ctx context.Context, resource *models.Resource) (*models.Resource, error) {
    // Validate
    if err := c.validator.ValidateForCreation(resource); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Process
    processed, err := c.processor.PrepareForCreation(resource)
    if err != nil {
        return nil, fmt.Errorf("processing failed: %w", err)
    }

    // Store
    return c.repository.Create(ctx, processed)
}

// 4. Add Logic function (logic/resource_validator.go)
func (v *ResourceValidator) ValidateForCreation(resource *models.Resource) error {
    if resource.Name == "" {
        return errors.New("name is required")
    }
    return nil
}
```

## 🧪 Testing Patterns

### Unit Testing (Logic Layer)

```go
func TestResourceValidator_ValidateForCreation(t *testing.T) {
    validator := NewResourceValidator()

    tests := []struct {
        name        string
        resource    *models.Resource
        expectError bool
    }{
        {
            name: "valid resource",
            resource: &models.Resource{Name: "test"},
            expectError: false,
        },
        {
            name: "missing name",
            resource: &models.Resource{},
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validator.ValidateForCreation(tt.resource)
            if (err != nil) != tt.expectError {
                t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
            }
        })
    }
}
```

### Integration Testing (Handlers)

```go
func TestHTTPHandler_CreateResource(t *testing.T) {
    // Setup dependencies with mocks
    mockRepo := &mocks.ResourceRepository{}
    controller := NewResourceController(mockRepo, validator, processor)
    handler := NewHTTPHandler(controller, adapter, responseAdapter, requestAdapter)

    // Setup mock expectations
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(&models.Resource{
        ID: "test-id",
        Name: "test-resource",
    }, nil)

    // Create request
    requestBody := `{"name": "test-resource"}`
    req := httptest.NewRequest("POST", "/resources", strings.NewReader(requestBody))
    w := httptest.NewRecorder()

    // Execute
    handler.createResourceHandler(w, req)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    mockRepo.AssertExpectations(t)
}
```

## 🔧 Module-Specific Guides

### Service Catalog

- **Purpose**: Manage service definitions with Kubernetes integration
- **Key Features**: CRUD operations, versioning (Git/Simple/None), health monitoring
- **Testing**: 100% coverage in logic layer, integration tests for handlers
- **Extensions**: Easy to add new storage providers or versioning systems

### Kubernetes

- **Purpose**: Kubernetes cluster management and monitoring
- **Key Features**: Multi-cluster support, health calculations, resource management
- **Testing**: Comprehensive health and resource calculation tests
- **Extensions**: Ready for new K8s resources (StatefulSets, Jobs, etc.)

### AWS

- **Purpose**: AWS infrastructure management
- **Key Features**: Multi-account EC2, cost optimization, permission system
- **Testing**: Cost calculations and permission logic fully tested
- **Extensions**: Framework ready for S3, RDS, Lambda, etc.

### Auth

- **Purpose**: Authentication and authorization
- **Key Features**: Multi-provider OAuth2, session management, JWT support
- **Testing**: OAuth2 flows and session management tested
- **Extensions**: Ready for SAML, LDAP, custom providers

## 🤝 Contributing to Backend Modules

### Working with Existing Modules

**Before making changes**:

1. **Understand the module**: Read the models and logic to understand the domain
2. **Check existing tests**: Look at test files to understand expected behavior
3. **Identify the right layer**: Determine where your change belongs

**Common contribution patterns**:

#### Adding a New Feature

```go
// 1. Add to models if new entity needed (models/entities.go)
type NewEntity struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func (ne *NewEntity) Validate() error {
    if ne.Name == "" {
        return fmt.Errorf("name is required")
    }
    return nil
}

// 2. Add business logic (logic/new_processor.go)
func (np *NewProcessor) ProcessNewEntity(entity *NewEntity) (*NewEntity, error) {
    // Pure business logic here
    return entity, nil
}

// 3. Add tests (logic/new_processor_test.go) - REQUIRED!
func TestNewProcessor_ProcessNewEntity(t *testing.T) {
    // Test implementation
}

// 4. Add to controller (controllers/main_controller.go)
func (c *Controller) CreateNewEntity(ctx context.Context, entity *NewEntity) (*NewEntity, error) {
    // Orchestration logic
}

// 5. Add HTTP endpoint (handlers/http.go)
func (h *HTTPHandler) createNewEntityHandler(w http.ResponseWriter, r *http.Request) {
    // HTTP handling logic
}
```

### Adding a New Module

Follow the **8-layer pattern** used by all existing modules:

1. **Create Directory Structure**:

   ```bash
   mkdir -p pkg/new-module/{adapters/{http,storage,external},config,controllers,handlers,logic,models,ports,wire}
   ```

2. **Implementation Order**:
   - **Models**: Domain entities with behavior
   - **Logic**: Pure business functions (must be 100% tested)
   - **Ports**: Interfaces for external dependencies
   - **Adapters**: Data transformation and external integration
   - **Controllers**: Business operation orchestration
   - **Handlers**: HTTP endpoint implementation
   - **Wire**: API contracts (request/response DTOs)
   - **Module Factory**: Dependency injection and initialization

### Code Review Checklist

- [ ] **Layer Separation**: Each file belongs to correct layer
- [ ] **Testing**: Logic layer has unit tests
- [ ] **Interfaces**: External dependencies use interfaces
- [ ] **Error Handling**: Proper error wrapping with context
- [ ] **Documentation**: Public functions are documented
- [ ] **Naming**: Follows established conventions
- [ ] **Validation**: Input validation in appropriate layer

### Performance Guidelines

- **Logic Layer**: Keep functions pure and fast
- **Controllers**: Minimize external calls, use context for cancellation
- **Adapters**: Implement efficient data transformations
- **Handlers**: Validate input early, fail fast
- **Models**: Use value objects for immutable data

## 📚 Additional Resources

- **[Full Architecture Specification](../BACKEND_ARCHITECTURE.md)** - Complete technical details
- **[Testing Patterns](./testing-patterns.md)** - Testing best practices (TODO)
- **[Deployment Guide](./deployment.md)** - Production deployment (TODO)
- **[Monitoring Setup](./monitoring.md)** - Observability configuration (TODO)

---

**🎯 Remember**: Consistency is key! Follow the established patterns, and when in doubt, look at existing modules for reference. The architecture is designed to be predictable and easy to navigate.
