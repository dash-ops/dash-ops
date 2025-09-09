# DashOPS Documentation

> **⚠️ Beta Software** - This documentation covers an experimental platform under active development.

Welcome to **DashOPS**, a cloud operations platform designed to simplify infrastructure management for development teams. This documentation will help you understand, configure, and contribute to the project.

## 📚 Documentation Overview

### **🚀 Getting Started**

- **[Quick Start](#quick-start)** - Get DashOPS running in 5 minutes
- **[Installation](#installation)** - Detailed setup instructions
- **[Configuration](#configuration)** - Complete configuration guide
- **[First Steps](#first-steps)** - Your first DashOPS experience

### **🔌 Plugin Documentation**

- **[Auth Plugin](./plugins/oauth2.md)** - GitHub authentication setup
- **[Kubernetes Plugin](./plugins/kubernetes.md)** - Multi-cluster management
- **[AWS Plugin](./plugins/aws.md)** - EC2 and account operations

### **🛠️ Development**

- **[API Reference](#api-reference)** - Backend API endpoints
- **[Frontend Guide](../front/README.md)** - React/TypeScript development
- **[Plugin Development](#plugin-development)** - Creating new plugins

### **🤝 Community**

- **[Contributing](#contributing)** - How to contribute to the project
- **[Roadmap](#roadmap)** - Future features and timeline

---

## 🚀 Quick Start

Get DashOPS running locally in development mode:

### **Prerequisites**

- **Go** 1.21+ (for backend)
- **Node.js** 18+ and **Yarn** (for frontend)
- **Docker** (optional, for containerized setup)

### **5-Minute Setup**

```bash
# 1. Clone repository
git clone https://github.com/dash-ops/dash-ops.git
cd dash-ops

# 2. Create configuration
cp local.sample.yaml dash-ops.yaml

# 3. Set environment variables
export GITHUB_CLIENT_ID="your-github-oauth-app-id"
export GITHUB_CLIENT_SECRET="your-github-oauth-secret"

# 4. Run backend (Terminal 1)
go run main.go

# 5. Run frontend (Terminal 2)
cd front
yarn && yarn dev

# 6. Access: http://localhost:5173
```

### **What You'll See**

1. **Login Screen** - GitHub Auth authentication
2. **Dashboard** - Overview of connected services
3. **AWS Module** - EC2 instance management
4. **Kubernetes Module** - Cluster and workload monitoring

---

## 🛠️ Installation

### **Local Development Setup**

#### **Backend (Go)**

```bash
# Install dependencies
go mod download

# Run with hot reload (optional)
go install github.com/cosmtrek/air@latest
air

# Or run directly
go run main.go
```

#### **Frontend (React + TypeScript)**

```bash
cd front

# Install dependencies
yarn

# Development with hot reload
yarn dev

# Quality checks
yarn quality
```

### **Docker Setup**

```bash
# Build image
docker build -t dashops/dash-ops .

# Run container
docker run --rm \
  -v $(pwd)/dash-ops.yaml:/dash-ops.yaml \
  -v ${HOME}/.kube/config:/.kube/config \
  -p 8080:8080 \
  dashops/dash-ops
```

### **Kubernetes Setup** (Development Only)

```bash
# Using Helm charts
helm install dash-ops ./helm-charts/dash-ops \
  --values ./your-values.yaml

# Access via port-forward
kubectl port-forward service/dash-ops 8080:8080
```

---

## ⚙️ Configuration

### **Configuration File Structure**

DashOPS uses a single YAML configuration file (`dash-ops.yaml`) to manage all settings:

```yaml
# Basic server configuration
port: 8080
origin: http://localhost:8080
headers:
  - 'Content-Type'
  - 'Authorization'
front: app # Frontend build directory

# Enable plugins
plugins:
  - 'Auth' # Authentication (required)
  - 'Kubernetes' # K8s operations (optional)
  - 'AWS' # AWS operations (optional)

# Plugin-specific configurations
auth: [...] # See Auth plugin docs
kubernetes: [...] # See Kubernetes plugin docs
aws: [...] # See AWS plugin docs
```

### **Environment Variables**

| Variable                | Description                | Required                |
| ----------------------- | -------------------------- | ----------------------- |
| `GITHUB_CLIENT_ID`      | GitHub OAuth App Client ID | ✅                      |
| `GITHUB_CLIENT_SECRET`  | GitHub OAuth App Secret    | ✅                      |
| `AWS_ACCESS_KEY_ID`     | AWS Access Key             | ⚠️ AWS Plugin only      |
| `AWS_SECRET_ACCESS_KEY` | AWS Secret Key             | ⚠️ AWS Plugin only      |
| `VITE_API_URL`          | Frontend API URL           | 🔄 Development override |

### **Security Configuration**

> **⚠️ Beta Notice**: Current security features are limited and not suitable for production.

```yaml
oauth2:
  - provider: github
    orgPermission: 'your-github-org' # Required org membership
    scopes: [user, repo, read:org] # GitHub permissions
```

---

## 🏗️ API Reference

### **Health & Status Endpoints**

```
GET /api/health      # Application health check
GET /api/version     # Build version and info
```

### **Authentication Endpoints**

```
GET  /api/oauth/redirect     # Auth callback handler
POST /api/oauth/logout       # Session termination
GET  /api/user              # Current user information
```

### **Plugin Endpoints**

#### **AWS Plugin**

```
GET /api/aws/accounts          # List AWS accounts
GET /api/aws/instances         # List EC2 instances
POST /api/aws/instances/{id}/start  # Start instance
POST /api/aws/instances/{id}/stop   # Stop instance
```

#### **Kubernetes Plugin**

```
GET /api/kubernetes/clusters     # List configured clusters
GET /api/kubernetes/deployments # List deployments
GET /api/kubernetes/pods        # List pods
GET /api/kubernetes/pods/{id}/logs  # Stream pod logs
POST /api/kubernetes/deployments/{id}/scale  # Scale deployment
```

### **Response Format**

All API responses follow a consistent structure:

```json
{
  "data": [...],     # Response payload
  "success": true,   # Operation status
  "message": "...",  # Optional message
  "errors": [...]    # Error details (if any)
}
```

---

## 🔌 Plugin Development

### **Plugin Architecture**

DashOPS uses a modular plugin system with both backend and frontend components:

```
Plugin Structure:
├── Backend (Go)
│   └── pkg/plugin-name/
│       ├── client.go      # External API client
│       ├── config.go      # Configuration parsing
│       ├── handler.go     # HTTP endpoint handlers
│       └── types.go       # Data structures
└── Frontend (TypeScript)
    └── src/modules/plugin-name/
        ├── index.tsx      # Module configuration
        ├── types.ts       # TypeScript interfaces
        ├── PluginPage.tsx # Main page component
        └── pluginResource.ts  # API client functions
```

### **Creating a New Plugin**

1. **Backend Implementation**

   ```go
   // pkg/newplugin/config.go
   package newplugin

   type Config struct {
       Name   string `yaml:"name"`
       ApiKey string `yaml:"apiKey"`
   }

   // pkg/newplugin/handler.go
   func NewHandler(config Config) *Handler {
       return &Handler{config: config}
   }
   ```

2. **Frontend Module**

   ```typescript
   // src/modules/newplugin/index.tsx
   export default {
     menus: [{ label: 'New Plugin', path: '/newplugin', icon: 'Settings' }],
     routers: [{ path: '/newplugin', component: 'NewPluginPage' }],
   };
   ```

3. **Register Plugin**
   ```yaml
   # dash-ops.yaml
   plugins:
     - 'NewPlugin'
   ```

---

## 🤝 Contributing

### **Development Workflow**

1. **Fork the repository** on GitHub
2. **Create a feature branch**: `git checkout -b feature/awesome-feature`
3. **Make your changes** with proper testing
4. **Run quality checks**: `yarn quality` (frontend) + `go test ./...` (backend)
5. **Commit semantically**: `git commit -m "feat: add awesome feature"`
6. **Submit a pull request** with detailed description

### **Code Standards**

#### **Backend (Go)**

- Use `gofmt` for formatting
- Add unit tests for new features
- Follow Go naming conventions
- Add godoc comments for public functions

#### **Frontend (TypeScript)**

- Strict TypeScript mode (no `any` types)
- Use ESLint + Prettier formatting
- Component tests with Testing Library
- Follow React best practices

### **Documentation Standards**

- Write in **English** (following project rules)
- Use clear, concise language
- Include code examples
- Add screenshots for UI features
- Update relevant plugin docs

---

## 🗺️ Roadmap

### **🔥 Current Focus (Beta Stabilization)**

**Q4 2024**

- Security hardening and production-ready authentication
- Comprehensive error handling and recovery
- Enhanced test coverage (>80%)
- Performance optimization and monitoring

### **🛒 Service Catalog (Planned Q3 2025)**

A self-service platform for infrastructure provisioning:

- **Template Library** - Pre-configured service templates
- **Lifecycle Management** - Creation, updates, deletion workflows
- **Approval Workflows** - Governance and compliance controls
- **Cost Estimation** - Resource cost prediction
- **Dependency Management** - Service relationship mapping

### **📈 Observability Integration (Planned Q4 2025)**

Deep integration with monitoring and observability platforms:

- **Prometheus/Grafana** - Metrics collection and visualization
- **ELK/EFK Stack** - Centralized logging and search
- **Jaeger/Zipkin** - Distributed tracing integration
- **Custom Dashboards** - Business-specific monitoring
- **Alert Management** - Centralized alerting and escalation

### **💰 FinOps Integration (Planned Q1 2026)**

Cost optimization and financial operations:

- **Cost Analytics** - Real-time spending analysis
- **Budget Management** - Cost controls and alerts
- **Resource Optimization** - Right-sizing recommendations
- **Chargeback/Showback** - Team and project cost allocation
- **Waste Detection** - Unused resource identification

---

## 🛡️ Security Considerations

> **⚠️ Important**: Current security implementation is **not production-ready**.

### **Current Limitations**

- Basic GitHub Auth only
- Browser localStorage for token storage
- Limited role-based access control
- No API rate limiting
- Basic audit logging

### **Security Roadmap**

- **Enterprise SSO** - SAML, OIDC, Active Directory
- **Encrypted Storage** - Secure credential management
- **Advanced RBAC** - Fine-grained permission system
- **API Security** - Rate limiting, WAF, DDoS protection
- **Compliance** - SOC2, ISO27001 readiness

---

## 📞 Support & Community

### **Getting Help**

- **[GitHub Issues](https://github.com/dash-ops/dash-ops/issues)** - Bug reports and feature requests
- **[GitHub Discussions](https://github.com/dash-ops/dash-ops/discussions)** - Community questions and ideas
- **[Project Wiki](https://github.com/dash-ops/dash-ops/wiki)** - Community-maintained guides

### **Contributing Areas**

We welcome contributions in these key areas:

1. **🔒 Security** - Authentication, authorization, encryption
2. **🧪 Testing** - Unit tests, integration tests, e2e tests
3. **📖 Documentation** - Guides, tutorials, API docs
4. **🔌 Plugins** - New cloud provider integrations
5. **🎨 UI/UX** - Interface improvements and accessibility
6. **🚀 Performance** - Backend optimization and caching

---

## 📄 License

DashOPS is licensed under the [MIT License](../LICENSE).

---

**🚧 Remember**: This is **beta software** intended for **testing and evaluation only**.

For the latest updates and community discussions, visit our [GitHub repository](https://github.com/dash-ops/dash-ops).
