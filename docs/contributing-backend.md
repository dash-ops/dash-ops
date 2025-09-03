# Contributing to DashOps Backend

> **🏗️ Architecture**: **4 modules migrated** to pure Hexagonal Architecture with **dependency injection**. Zero legacy code in migrated modules.

## Quick Reference

### 📁 Module Structure (All modules follow this pattern)

```
pkg/{module}/
├── adapters/     # Data transformation & external integrations
├── controllers/  # Business logic orchestration
├── handlers/     # HTTP endpoints
├── logic/        # Pure business logic (100% tested)
├── models/       # Domain entities with behavior
├── ports/        # Interfaces & contracts
├── wire/         # API contracts (DTOs)
└── module.go     # Module factory
```

### 🔄 Data Flow

```
HTTP Request → Handler → Adapter → Controller → Logic → Repository
```

## 🚀 Adding New Features

### 1. Choose the Right Module

| Feature Type         | Module            | Examples                                 |
| -------------------- | ----------------- | ---------------------------------------- |
| Service management   | `service-catalog` | Add service types, new storage providers |
| Container operations | `kubernetes`      | New K8s resources, health checks         |
| Cloud infrastructure | `aws`             | S3, RDS, Lambda integration              |
| Authentication       | `auth`            | New OAuth providers, SAML, LDAP          |
| Configuration        | `config`          | New config sources, validation rules     |
| Git operations       | `github`          | Repository management, webhooks          |
| File serving         | `spa`             | Asset optimization, caching strategies   |
| Shared utilities     | `commons`         | HTTP helpers, validation utilities       |

### 2. Implementation Steps

#### Step 1: Define Domain Model

```go
// models/entities.go
type Resource struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Status      ResourceStatus `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}

// Add domain behavior
func (r *Resource) IsActive() bool {
    return r.Status == ResourceStatusActive
}

func (r *Resource) Validate() error {
    if r.Name == "" {
        return fmt.Errorf("name is required")
    }
    return nil
}
```

#### Step 2: Implement Business Logic

```go
// logic/resource_processor.go
type ResourceProcessor struct{}

func NewResourceProcessor() *ResourceProcessor {
    return &ResourceProcessor{}
}

func (rp *ResourceProcessor) PrepareForCreation(resource *Resource) (*Resource, error) {
    // Pure business logic here
    prepared := *resource
    prepared.ID = rp.generateID(resource.Name)
    prepared.CreatedAt = time.Now()
    prepared.Status = ResourceStatusActive
    return &prepared, nil
}

// logic/resource_processor_test.go - ALWAYS ADD TESTS!
func TestResourceProcessor_PrepareForCreation(t *testing.T) {
    processor := NewResourceProcessor()

    resource := &Resource{Name: "test"}
    result, err := processor.PrepareForCreation(resource)

    assert.NoError(t, err)
    assert.NotEmpty(t, result.ID)
    assert.Equal(t, ResourceStatusActive, result.Status)
}
```

#### Step 3: Define Interface

```go
// ports/repositories.go
type ResourceRepository interface {
    Create(ctx context.Context, resource *Resource) (*Resource, error)
    GetByID(ctx context.Context, id string) (*Resource, error)
    Update(ctx context.Context, resource *Resource) (*Resource, error)
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter *ResourceFilter) ([]Resource, error)
}
```

#### Step 4: Create Controller

```go
// controllers/resource_controller.go
type ResourceController struct {
    repo      ports.ResourceRepository
    processor *logic.ResourceProcessor
}

func (rc *ResourceController) CreateResource(ctx context.Context, resource *Resource) (*Resource, error) {
    // Validate
    if err := resource.Validate(); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Process
    prepared, err := rc.processor.PrepareForCreation(resource)
    if err != nil {
        return nil, fmt.Errorf("processing failed: %w", err)
    }

    // Store
    return rc.repo.Create(ctx, prepared)
}
```

#### Step 5: Add HTTP Handler

```go
// handlers/http.go
func (h *HTTPHandler) createResourceHandler(w http.ResponseWriter, r *http.Request) {
    var req wire.CreateResourceRequest
    if err := h.requestAdapter.ParseJSON(r, &req); err != nil {
        h.responseAdapter.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }

    resource, err := h.adapter.RequestToModel(req)
    if err != nil {
        h.responseAdapter.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }

    result, err := h.controller.CreateResource(r.Context(), resource)
    if err != nil {
        h.responseAdapter.WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }

    response := h.adapter.ModelToResponse(result)
    h.responseAdapter.WriteCreated(w, "/resources/"+result.ID, response)
}
```

## 🧪 Testing Best Practices

### Unit Tests (Required for Logic Layer)

```go
// ✅ DO: Test pure business logic
func TestResourceProcessor_ValidateResource(t *testing.T) {
    processor := NewResourceProcessor()

    tests := []struct {
        name        string
        resource    *Resource
        expectError bool
    }{
        {"valid resource", &Resource{Name: "test"}, false},
        {"missing name", &Resource{}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := processor.ValidateResource(tt.resource)
            if (err != nil) != tt.expectError {
                t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
            }
        })
    }
}

// ❌ DON'T: Test external dependencies directly
func TestController_CreateResource_BadExample(t *testing.T) {
    // Don't do this - tests real database
    db := sql.Open("postgres", "real-connection-string")
    controller := NewController(db)
    // ...
}

// ✅ DO: Use mocks for external dependencies
func TestController_CreateResource_GoodExample(t *testing.T) {
    mockRepo := &mocks.ResourceRepository{}
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(&Resource{ID: "123"}, nil)

    controller := NewController(mockRepo)
    // ...
}
```

### Integration Tests (Handlers)

```go
func TestHTTPHandler_CreateResource(t *testing.T) {
    // Setup test server
    mockRepo := &mocks.ResourceRepository{}
    module, _ := NewModule(&ModuleConfig{ResourceRepo: mockRepo})

    // Setup expectations
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(&Resource{
        ID: "test-id",
        Name: "test-resource",
    }, nil)

    // Create request
    requestBody := `{"name": "test-resource"}`
    req := httptest.NewRequest("POST", "/resources", strings.NewReader(requestBody))
    w := httptest.NewRecorder()

    // Execute
    module.GetHandler().ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)

    var response wire.ResourceResponse
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "test-resource", response.Name)
}
```

## 🔍 Code Examples by Module

### Service Catalog - Adding New Storage Provider

```go
// 1. Define interface (ports/repositories.go)
type ServiceRepository interface {
    Create(ctx context.Context, service *Service) (*Service, error)
    // ... other methods
}

// 2. Implement adapter (adapters/storage/s3_repository.go)
type S3Repository struct {
    bucket string
    client S3Client
}

func (s3 *S3Repository) Create(ctx context.Context, service *Service) (*Service, error) {
    // S3-specific implementation
}

// 3. Update module factory (module.go)
func newStorageProvider(config *Config) (ports.ServiceRepository, error) {
    switch config.Storage.Provider {
    case "filesystem":
        return storage.NewFilesystemRepository(config.Storage.Filesystem.Directory)
    case "s3":
        return storage.NewS3Repository(config.Storage.S3.Bucket)
    }
}
```

### Kubernetes - Adding New Resource Type

```go
// 1. Add to models (models/entities.go)
type StatefulSet struct {
    Name      string    `json:"name"`
    Namespace string    `json:"namespace"`
    Replicas  int32     `json:"replicas"`
    // ... other fields
}

func (ss *StatefulSet) IsHealthy() bool {
    // Domain logic
}

// 2. Add to logic (logic/health_calculator.go)
func (hc *HealthCalculator) CalculateStatefulSetHealth(ss *StatefulSet) HealthStatus {
    // Health calculation logic
}

// 3. Add to ports (ports/repositories.go)
type StatefulSetRepository interface {
    GetStatefulSet(ctx context.Context, namespace, name string) (*StatefulSet, error)
    ListStatefulSets(ctx context.Context, namespace string) ([]StatefulSet, error)
}

// 4. Implement adapter (adapters/external/k8s_client_adapter.go)
func (kca *KubernetesClientAdapter) GetStatefulSet(ctx context.Context, namespace, name string) (*StatefulSet, error) {
    // Kubernetes client implementation
}
```

### AWS - Adding New Service

```go
// 1. Add models (models/entities.go)
type S3Bucket struct {
    Name         string    `json:"name"`
    Region       string    `json:"region"`
    CreationDate time.Time `json:"creation_date"`
    // ... other fields
}

// 2. Add logic (logic/s3_processor.go)
type S3Processor struct{}

func (s3p *S3Processor) ValidateBucketName(name string) error {
    // S3 bucket naming validation
}

// 3. Add to controller (controllers/aws_controller.go)
func (ac *AWSController) ListS3Buckets(ctx context.Context, accountKey string) ([]S3Bucket, error) {
    // Implementation
}
```

## 🔄 **Dependency Inversion Best Practices**

### **✅ DO: Define Interfaces in Dependent Module**

```go
// ✅ CORRECT: auth/ports/services.go
package ports

type GitHubService interface {
    GetUser(ctx context.Context, token *oauth2.Token) (*github.User, error)
    GetUserTeams(ctx context.Context, token *oauth2.Token) ([]*github.Team, error)
}
```

### **✅ DO: Inject Dependencies in main.go**

```go
// ✅ CORRECT: main.go orchestrates all dependencies
func main() {
    // Initialize dependency first
    githubModule, _ := github.NewModule(oauthConfig)
    
    // Inject dependency into dependent module
    authModule, _ := auth.NewModule(authConfig, githubModule)
    
    // Register routes
    authModule.RegisterRoutes(api, internal)
}
```

### **✅ DO: Use Interfaces for Loose Coupling**

```go
// ✅ CORRECT: auth/controllers/auth_controller.go
type AuthController struct {
    githubService ports.GitHubService // Interface, not implementation
}

func (ac *AuthController) GetUserProfile(ctx context.Context, token *oauth2.Token) {
    // Use interface, works with any implementation
    return ac.githubService.GetUser(ctx, token)
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

## 🚨 Common Pitfalls

### ❌ What NOT to Do

1. **Don't mix layers**:

   ```go
   // BAD: Handler calling repository directly
   func (h *Handler) createResource(w http.ResponseWriter, r *http.Request) {
       resource := parseRequest(r)
       h.repository.Create(resource) // Skip controller!
   }
   ```

2. **Don't put business logic in handlers**:

   ```go
   // BAD: Business logic in handler
   func (h *Handler) createResource(w http.ResponseWriter, r *http.Request) {
       if resource.Name == "" { // Validation belongs in logic layer!
           http.Error(w, "Name required", 400)
           return
       }
   }
   ```

3. **Don't skip testing logic layer**:
   ```go
   // BAD: No tests for business logic
   func (p *Processor) ComplexBusinessLogic(data *Data) (*Result, error) {
       // Complex logic without tests - risky!
   }
   ```

### ✅ What TO Do

1. **Follow the data flow**:

   ```go
   // GOOD: Clear layer separation
   Handler → Adapter → Controller → Logic → Repository
   ```

2. **Test business logic thoroughly**:

   ```go
   // GOOD: Comprehensive logic tests
   func TestProcessor_ComplexBusinessLogic(t *testing.T) {
       // Test all edge cases and business rules
   }
   ```

3. **Use interfaces for external dependencies**:
   ```go
   // GOOD: Mockable dependencies
   type Controller struct {
       repo ports.Repository // Interface, not concrete type
   }
   ```

## 📋 Development Workflow

1. **🔍 Understand**: Read existing code in the module
2. **📝 Plan**: Design your changes following the layer pattern
3. **🧪 Test First**: Write tests for logic layer
4. **💻 Implement**: Code following established patterns
5. **🔄 Integrate**: Add handlers and wire everything together
6. **✅ Validate**: Run tests and ensure everything works

## 🎯 Quality Standards

- **Logic Layer**: 100% test coverage required
- **Models**: Domain methods should have tests
- **Handlers**: Integration tests for complex flows
- **Documentation**: Public APIs must be documented
- **Error Handling**: All errors must be wrapped with context
- **Performance**: No obvious performance anti-patterns

---

**Happy Coding! 🚀** The architecture is designed to make your life easier. When you follow the patterns, everything just works™️
