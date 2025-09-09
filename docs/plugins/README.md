# DashOPS Plugins

> **⚠️ Beta Plugins** - All plugins are experimental and under active development.

DashOPS uses a **modular plugin architecture** that allows you to enable only the cloud services you need. Each plugin provides both backend API integration and frontend user interfaces.

## 🔌 Available Plugins

| Plugin                                      | Status   | Features                                             | Production Ready    |
| ------------------------------------------- | -------- | ---------------------------------------------------- | ------------------- |
| **[Auth](./auth.md)**                     | 🔄 Beta  | GitHub authentication, org validation                | ❌ Testing only     |
| **[Service Catalog](./service-catalog.md)** | 🔄 Beta  | Service registry, health monitoring, K8s integration | ❌ Testing only     |
| **[AWS](./aws.md)**                         | 🔄 Alpha | EC2 management, multi-account, optimized UI          | ❌ Development only |
| **[Kubernetes](./kubernetes.md)**           | 🔄 Alpha | Multi-cluster, enhanced UI, pod logs, scaling        | ❌ Development only |

## 🚀 Quick Plugin Setup

### **1. Enable Plugins**

In your `dash-ops.yaml` configuration:

```yaml
plugins:
  - 'Auth' # Required - Authentication
  - 'ServiceCatalog' # Optional - Service registry
  - 'AWS' # Optional - AWS operations
  - 'Kubernetes' # Optional - K8s operations
```

### **2. Configure Each Plugin**

Each plugin requires specific configuration. See individual plugin documentation:

- **[Auth Setup](./auth.md#configuration)** - GitHub OAuth App configuration
- **[Service Catalog Setup](./service-catalog.md#configuration)** - Service registry and GitHub integration
- **[AWS Setup](./aws.md#configuration)** - AWS credentials and permissions
- **[Kubernetes Setup](./kubernetes.md#configuration)** - Cluster access and RBAC

### **3. Verify Plugin Status**

After starting DashOPS, check plugin status:

```bash
# Check if plugins loaded successfully
curl http://localhost:8080/api/health

# Verify plugin endpoints
curl http://localhost:8080/api/config/plugins
```

## 🏗️ Plugin Architecture

### **Backend Structure**

Each plugin follows a consistent Go package structure:

```
pkg/plugin-name/
├── client.go      # External service client (AWS SDK, K8s client-go)
├── config.go      # Configuration parsing and validation
├── handler.go     # HTTP endpoint handlers and routing
├── types.go       # Data structures and interfaces
└── *_test.go      # Unit tests for each component
```

### **Frontend Structure**

Each plugin includes a React module:

```
src/modules/plugin-name/
├── index.tsx           # Module configuration (menus, routes)
├── types.ts           # TypeScript interface definitions
├── PluginPage.tsx     # Main page component
├── pluginResource.ts  # API client functions
├── PluginActions.tsx  # Action components (buttons, forms)
└── __tests__/         # Component and integration tests
```

### **Plugin Registration**

Plugins are automatically discovered and loaded:

1. **Backend**: Plugin packages are imported in `main.go`
2. **Frontend**: Modules are dynamically loaded based on enabled plugins
3. **Configuration**: Plugin settings parsed from `dash-ops.yaml`

## 🔐 Permission System

### **Team-based Access Control**

All plugins use GitHub organization and team membership for permissions:

```yaml
plugin-name:
  - name: 'Resource Name'
    permission:
      operation-type:
        action: ['github-org*team-name']
```

### **Permission Formats**

| Format               | Description        | Example                            |
| -------------------- | ------------------ | ---------------------------------- |
| `org*team`           | Specific team only | `dash-ops*developers`              |
| `org*`               | All org members    | `dash-ops*`                        |
| `["team1", "team2"]` | Multiple teams     | `["dash-ops*sre", "dash-ops*ops"]` |

### **Common Permission Patterns**

```yaml
# Read-only access (no permission block)
kubernetes:
  - name: 'Production Cluster'
    # No permission = read-only

# Development operations
kubernetes:
  - name: 'Dev Cluster'
    permission:
      deployments:
        namespaces: ['dev', 'staging']
        start: ['dash-ops*developers']
        stop: ['dash-ops*developers']

# Production operations (restricted)
aws:
  - name: 'Production Account'
    permission:
      ec2:
        start: ['dash-ops*sre']  # Only SRE team
        stop: ['dash-ops*sre']
```

## 🛠️ Plugin Development

### **Creating a New Plugin**

#### **1. Backend Plugin**

```go
// pkg/newplugin/config.go
package newplugin

type Config struct {
    Name   string `yaml:"name"`
    ApiKey string `yaml:"apiKey"`
    URL    string `yaml:"url"`
}

// pkg/newplugin/client.go
type Client struct {
    config Config
    httpClient *http.Client
}

func NewClient(config Config) *Client {
    return &Client{
        config: config,
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }
}

// pkg/newplugin/handler.go
func NewHandler(config Config) *Handler {
    return &Handler{
        client: NewClient(config),
    }
}

func (h *Handler) HandleResources(w http.ResponseWriter, r *http.Request) {
    // API endpoint implementation
}
```

#### **2. Frontend Module**

```typescript
// src/modules/newplugin/index.tsx
import { Menu, Router } from '@/types';

export default {
  menus: [
    {
      label: 'New Plugin',
      key: 'newplugin',
      path: '/newplugin',
      icon: 'Settings',
    },
  ] as Menu[],
  routers: [
    {
      path: '/newplugin',
      component: 'NewPluginPage',
    },
  ] as Router[],
};

// src/modules/newplugin/types.ts
export interface NewPluginResource {
  id: string;
  name: string;
  status: 'active' | 'inactive';
}

// src/modules/newplugin/newPluginResource.ts
import http from '../../helpers/http';

export function getResources(): Promise<NewPluginResource[]> {
  return http.get('/newplugin/resources').then((res) => res.data);
}
```

#### **3. Register Plugin**

```yaml
# dash-ops.yaml
plugins:
  - 'NewPlugin'

newplugin:
  - name: 'My Service'
    apiKey: ${NEW_PLUGIN_API_KEY}
    url: 'https://api.service.com'
```

### **Plugin Development Guidelines**

#### **Backend Best Practices**

- **Error Handling** - Implement comprehensive error responses
- **Testing** - Unit tests for all handlers and clients
- **Configuration Validation** - Validate required fields at startup
- **Security** - Implement proper authentication and authorization
- **Logging** - Structured logging with context

#### **Frontend Best Practices**

- **TypeScript Strict** - No `any` types allowed
- **Component Testing** - Test all user interactions
- **Error Boundaries** - Graceful error handling
- **Loading States** - Proper loading indicators
- **Accessibility** - WCAG compliance with shadcn/ui

## 🚨 Beta Plugin Limitations

### **Current Issues Across All Plugins**

❌ **Not Production Ready**

- **Limited error recovery** - Basic error handling only
- **No data persistence** - Configuration-based setup only
- **Basic security model** - Simple team-based permissions
- **Missing monitoring** - No metrics or alerting
- **No plugin isolation** - Shared configuration and resources

### **Development Recommendations**

1. **Test in isolated environments** only
2. **Use minimal permissions** for plugin service accounts
3. **Monitor all operations** manually during testing
4. **Backup configurations** before testing new features
5. **Report issues** via GitHub Issues

## 🤝 Contributing to Plugins

### **High Priority Areas**

1. **🔒 Security** - Implement secure credential management
2. **🧪 Testing** - Comprehensive test coverage
3. **📊 Monitoring** - Health checks and metrics
4. **🚨 Error Handling** - Robust failure recovery
5. **📖 Documentation** - Usage examples and troubleshooting

### **Plugin Contribution Workflow**

```bash
# 1. Fork and clone
git clone https://github.com/your-username/dash-ops.git

# 2. Create plugin branch
git checkout -b feature/enhance-aws-plugin

# 3. Develop with testing
go test ./pkg/aws/...              # Backend tests
cd front && yarn test src/modules/aws/  # Frontend tests

# 4. Quality checks
yarn quality                       # Frontend quality
go test ./...                      # Backend tests
go vet ./...                       # Go static analysis

# 5. Submit PR with documentation updates
```

## 📚 Additional Resources

- **[Go Plugin Development](https://golang.org/doc/effective_go.html)** - Go best practices
- **[React Component Development](https://react.dev/learn)** - React patterns
- **[TypeScript Handbook](https://www.typescriptlang.org/docs/)** - TypeScript guide
- **[API Design Guidelines](https://github.com/microsoft/api-guidelines)** - REST API best practices

---

**⚠️ Beta Notice**: All plugins are experimental and intended for development use only. Do not use in production environments.
