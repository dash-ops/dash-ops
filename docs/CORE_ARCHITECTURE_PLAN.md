# 🏗️ DashOps Core Architecture Plan

## Service Catalog as the Heart of DashOps

### 🎯 Vision Overview

Transform DashOps from a "Plugin-First Architecture" to a **"Service-Centric Architecture"** where:

- **Service Catalog = Core System** (not a plugin)
- **Plugins = Service Enhancers** (add capabilities to services)
- **Everything revolves around Services** (infrastructure, monitoring, CI/CD, etc.)

---

## 🔄 Current vs. Proposed Architecture

### Current Architecture Problems

```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   AWS Plugin    │  │  K8s Plugin     │  │ Service Plugin  │
│   (Isolated)    │  │   (Isolated)    │  │   (Isolated)    │
└─────────────────┘  └─────────────────┘  └─────────────────┘
      ↕️                      ↕️                      ↕️
┌─────────────────────────────────────────────────────────┐
│                 Plugin Manager                          │
└─────────────────────────────────────────────────────────┘
```

**Issues:**

- Plugins work in isolation
- No unified view of services
- Redundant information across plugins
- Hard to correlate infrastructure with services

### Proposed Architecture

```
                    ┌─────────────────────────────────────┐
                    │          Service Catalog            │
                    │             (CORE)                  │
                    │  ┌─────────────────────────────────┐ │
                    │  │        Service Registry         │ │
                    │  │     ┌─────────┬─────────────┐   │ │
                    │  │     │Service A│   Service B │   │ │
                    │  │     │  Context│    Context  │   │ │
                    │  │     └─────────┴─────────────┘   │ │
                    │  └─────────────────────────────────┘ │
                    └─────────────────────────────────────┘
                                      │
           ┌──────────────────────────┼──────────────────────────┐
           ▼                          ▼                          ▼
    ┌─────────────┐          ┌─────────────┐          ┌─────────────┐
    │AWS Plugin   │          │K8s Plugin   │          │CI/CD Plugin │
    │             │          │             │          │             │
    │• EC2 Info   │          │• Pod Status │          │• Pipeline   │
    │• RDS Data   │          │• Resources  │          │• Deploy     │
    │• CloudWatch │          │• Logs       │          │• History    │
    └─────────────┘          └─────────────┘          └─────────────┘
```

---

## 🏗️ Core Service Model (Enhanced)

### Service Definition (Extended)

```go
type Service struct {
    // Basic Info (already exists)
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    DisplayName string            `json:"displayName"`
    Description string            `json:"description"`
    Tier        string            `json:"tier"`
    Team        string            `json:"team"`
    Squad       string            `json:"squad"`
    Owner       string            `json:"owner"`
    Tags        []string          `json:"tags"`
    Status      string            `json:"status"`

    // NEW: Service Specifications
    Architecture ServiceArchitecture `json:"architecture"`
    Dependencies []ServiceDependency  `json:"dependencies"`
    Endpoints    []ServiceEndpoint    `json:"endpoints"`

    // NEW: Plugin Contexts (Dynamic)
    Contexts     map[string]interface{} `json:"contexts,omitempty"`

    // Lifecycle
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

type ServiceArchitecture struct {
    Type         string   `json:"type"`        // "api", "worker", "batch", "frontend"
    Runtime      string   `json:"runtime"`     // "docker", "lambda", "vm"
    Language     string   `json:"language"`    // "go", "python", "node"
    Framework    string   `json:"framework"`   // "gin", "express", "django"
    Database     []string `json:"database"`    // ["postgresql", "redis"]
}

type ServiceDependency struct {
    ServiceID   string `json:"serviceId"`
    Type        string `json:"type"`        // "sync", "async", "data"
    Required    bool   `json:"required"`
    Description string `json:"description"`
}

type ServiceEndpoint struct {
    Path        string `json:"path"`
    Method      string `json:"method"`
    Public      bool   `json:"public"`
    Description string `json:"description"`
}
```

---

## 🔌 Plugin System Redesign

### Plugin Interface (Service-Aware)

```go
type ServicePlugin interface {
    // Plugin Metadata
    Name() string
    Version() string
    Type() string // "infrastructure", "observability", "cicd", "security"

    // Service Integration
    CanEnhanceService(service *Service) bool
    GetServiceContext(serviceID string) (ServiceContext, error)
    GetServiceActions(serviceID string) ([]ServiceAction, error)
    GetServiceHealth(serviceID string) (HealthStatus, error)

    // Plugin Lifecycle
    Install(config PluginConfig) error
    Uninstall() error
    GetRoutes() []Route
    GetPermissions() []Permission
}

type ServiceContext struct {
    PluginName string                 `json:"pluginName"`
    Data       map[string]interface{} `json:"data"`
    LastUpdate time.Time              `json:"lastUpdate"`
    Status     string                 `json:"status"`
}

type ServiceAction struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Type        string            `json:"type"` // "button", "form", "link"
    Endpoint    string            `json:"endpoint"`
    Parameters  map[string]string `json:"parameters,omitempty"`
}
```

### Example Plugin Implementations

#### AWS Plugin (Infrastructure Context)

```go
type AWSPlugin struct {
    ec2Client *ec2.Client
    rdsClient *rds.Client
    cwClient  *cloudwatch.Client
}

func (p *AWSPlugin) GetServiceContext(serviceID string) (ServiceContext, error) {
    service := getServiceByID(serviceID)

    // Find AWS resources by service tags
    instances := p.findEC2InstancesByService(serviceID)
    databases := p.findRDSByService(serviceID)
    metrics := p.getCloudWatchMetrics(serviceID)

    return ServiceContext{
        PluginName: "aws",
        Data: map[string]interface{}{
            "instances": instances,
            "databases": databases,
            "metrics":   metrics,
            "costs":     p.getCostsByService(serviceID),
        },
        Status: "active",
        LastUpdate: time.Now(),
    }, nil
}

func (p *AWSPlugin) GetServiceActions(serviceID string) ([]ServiceAction, error) {
    return []ServiceAction{
        {
            ID:          "restart-instances",
            Name:        "Restart Instances",
            Description: "Restart all EC2 instances for this service",
            Type:        "button",
            Endpoint:    "/api/v1/plugins/aws/services/{serviceId}/restart",
        },
        {
            ID:          "scale-up",
            Name:        "Scale Up",
            Description: "Add more instances to handle load",
            Type:        "form",
            Endpoint:    "/api/v1/plugins/aws/services/{serviceId}/scale",
        },
    }, nil
}
```

#### Kubernetes Plugin (Orchestration Context)

```go
type KubernetesPlugin struct {
    client kubernetes.Interface
}

func (p *KubernetesPlugin) GetServiceContext(serviceID string) (ServiceContext, error) {
    // Find K8s resources by service labels
    deployments := p.findDeploymentsByService(serviceID)
    pods := p.findPodsByService(serviceID)
    services := p.findServicesByService(serviceID)

    return ServiceContext{
        PluginName: "kubernetes",
        Data: map[string]interface{}{
            "deployments": deployments,
            "pods":        pods,
            "services":    services,
            "ingresses":   p.findIngressesByService(serviceID),
            "health":      p.getHealthStatus(serviceID),
        },
        Status: "active",
        LastUpdate: time.Now(),
    }, nil
}
```

---

## 🎨 Frontend Service-Centric View

### Service Detail Page Architecture

```jsx
// ServiceDetailPage.jsx
function ServiceDetailPage({ serviceId }) {
  const service = useService(serviceId);
  const plugins = useServicePlugins(serviceId);

  return (
    <div className="service-detail">
      <ServiceHeader service={service} />

      <Tabs>
        <Tab label="Overview">
          <ServiceOverview service={service} />
          <ServiceDependencies dependencies={service.dependencies} />
        </Tab>

        <Tab label="Infrastructure">
          <PluginContext plugin="aws" serviceId={serviceId} />
          <PluginContext plugin="kubernetes" serviceId={serviceId} />
        </Tab>

        <Tab label="Observability">
          <PluginContext plugin="prometheus" serviceId={serviceId} />
          <PluginContext plugin="logs" serviceId={serviceId} />
        </Tab>

        <Tab label="CI/CD">
          <PluginContext plugin="github" serviceId={serviceId} />
          <PluginContext plugin="jenkins" serviceId={serviceId} />
        </Tab>

        <Tab label="Security">
          <PluginContext plugin="security" serviceId={serviceId} />
        </Tab>
      </Tabs>

      <ServiceActions service={service} plugins={plugins} />
    </div>
  );
}

// PluginContext.jsx - Dynamic plugin rendering
function PluginContext({ plugin, serviceId }) {
  const context = usePluginContext(plugin, serviceId);
  const actions = usePluginActions(plugin, serviceId);

  if (!context) return <div>Plugin not available</div>;

  return (
    <div className="plugin-context">
      <h3>{plugin.toUpperCase()} Context</h3>
      <PluginDataViewer data={context.data} />
      <PluginActions actions={actions} />
    </div>
  );
}
```

---

## 🚀 Implementation Phases

### Phase 1: Core Service Model Enhancement

- [ ] Expand Service model with Architecture, Dependencies, Endpoints
- [ ] Create Service Registry with enhanced metadata
- [ ] Build Service relationship mapping
- [ ] Implement Service search and filtering by plugins

### Phase 2: Plugin Interface Redesign

- [ ] Define new ServicePlugin interface
- [ ] Create Plugin Context system
- [ ] Build Plugin Action framework
- [ ] Implement Plugin Health checking

### Phase 3: Plugin Migration

- [ ] Refactor AWS Plugin to Service-aware
- [ ] Refactor Kubernetes Plugin to Service-aware
- [ ] Refactor OAuth2 Plugin (if service-related)
- [ ] Create new plugins (GitHub, Monitoring, etc.)

### Phase 4: Frontend Reconstruction

- [ ] Build Service-centric navigation
- [ ] Create dynamic Plugin Context components
- [ ] Implement Service Actions UI
- [ ] Build Service dependency visualization

### Phase 5: Advanced Features

- [ ] Service health aggregation from all plugins
- [ ] Cross-plugin correlation (K8s pod → AWS instance)
- [ ] Service impact analysis
- [ ] Automated service discovery

---

## 🎯 Benefits of New Architecture

### For Users

- **Single Source of Truth**: All service information in one place
- **Contextual Actions**: Actions relevant to each service
- **Holistic View**: See infrastructure, monitoring, CI/CD together
- **Impact Analysis**: Understand service dependencies

### For Developers

- **Plugin Simplicity**: Plugins focus on their domain
- **Service Integration**: Easy to correlate plugin data with services
- **Extensibility**: Add new plugins without core changes
- **Consistency**: Unified patterns across all plugins

### For Operations

- **Service-first Monitoring**: Monitor by business service, not infrastructure
- **Faster Troubleshooting**: All context in one view
- **Better Planning**: Understand service relationships
- **Cost Attribution**: See costs per service across all providers

---

## 📋 Migration Strategy

### Step 1: Backwards Compatibility

- Keep current plugin system working
- Run both architectures in parallel
- Migrate plugins one by one

### Step 2: Data Migration

- Export current service definitions
- Enhance with new architecture fields
- Import into new Service Registry

### Step 3: Plugin Migration

- Start with least complex plugin (AWS)
- Test thoroughly before next plugin
- Update frontend progressively

### Step 4: Cleanup

- Remove old plugin architecture
- Clean up deprecated endpoints
- Update documentation

---

This architecture transformation will make DashOps a true **Service Operations Platform** where everything revolves around your business services, not just infrastructure tools.
