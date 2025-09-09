# DashOPS - Developer Experience Hub with AI (Beta)

![DashOps](https://github.com/dash-ops/dash-ops/workflows/DashOps/badge.svg)

> **⚠️ BETA VERSION - NOT RECOMMENDED FOR PRODUCTION USE**

> **The VS Code for infrastructure** - A unified, AI-powered hub that integrates all your developer tools in one intuitive interface. Focus on building features, not juggling between different platforms.

DashOPS is an **experimental integration platform** that connects your existing tools (Kubernetes, AWS, Grafana, ArgoCD) into a seamless developer experience, enhanced by contextual AI assistance and an extensible plugin system.

**🚧 This project is actively under development and should only be used for testing and evaluation purposes.**

## 🚀 Quick Start

### Option 1: Docker (Recommended)

```bash
# 1. Create configuration file
cp local.sample.yaml dash-ops.yaml

# 2. Set environment variables
export GITHUB_CLIENT_ID="your-github-client-id"
export GITHUB_CLIENT_SECRET="your-github-client-secret"
export AWS_ACCESS_KEY_ID="your-aws-access-key"
export AWS_SECRET_ACCESS_KEY="your-aws-secret-key"

# 3. Run with Docker
docker run --rm \
  -v $(pwd)/dash-ops.yaml:/dash-ops.yaml \
  -v ${HOME}/.kube/config:/.kube/config \
  -e GITHUB_CLIENT_ID \
  -e GITHUB_CLIENT_SECRET \
  -e AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY \
  -p 8080:8080 \
  -it dashops/dash-ops
```

### Option 2: Local Development

```bash
# 1. Backend (Go)
go run main.go

# 2. Frontend (React + TypeScript)
cd front
yarn
yarn dev

# Access: http://localhost:5173
```

## 🏗️ Architecture

### Backend Architecture

The DashOps backend follows a **Hexagonal Architecture** pattern with 8 consistent layers across all modules:

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

**Key Benefits**:

- **Consistent Structure**: Same pattern across all 8 modules
- **High Testability**: 80+ unit tests ensuring reliability
- **Extensibility**: Interface-based design for easy extension
- **Maintainability**: Clear separation of concerns

**📚 For Developers**: See [Backend Architecture Guide](./docs/backend-architecture.md) for detailed contribution guidelines.

### Frontend Architecture

DashOPS is built as an **Integration Hub** that connects your existing tools with AI-powered UX:

### **🎯 Integration Philosophy**

- **🔗 Connect, Don't Replace** - Integrate with tools you already use
- **🤖 AI-Enhanced** - Contextual assistance across all integrations
- **🎨 UX-First** - Intuitive interface that abstracts complexity
- **🧩 Extensible** - Plugin system for community contributions

### **Backend** (Go)

- **Integration Engine** - Smart aggregation of external APIs
- **Plugin System** - Extensible architecture for new tools
- **AI Context Layer** - Contextual data correlation across tools
- **Security Gateway** - OAuth2 integration with unified auth

### **Frontend** (React + TypeScript + AI)

- **Unified Dashboard** - All tools in one interface
- **AI Assistant** - Contextual help and automation
- **shadcn/ui Components** - Consistent, accessible design
- **Smart Caching** - Optimized performance across integrations

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────────────┐
│   Frontend      │    │  Integration     │    │    External Tools       │
│ React + TS + AI │◄──►│   Hub (Go)       │◄──►│ Grafana│ArgoCD│AWS│K8s │
│   Port 5173     │    │   Port 8080      │    │   Prometheus│Loki│etc   │
└─────────────────┘    └──────────────────┘    └─────────────────────────┘
```

## 🎯 Features

### **🆕 Latest Updates (v0.3.0-beta)**

**New Service Catalog Plugin:**

- ✅ **Service Registry** - Complete CRUD for service definitions with modern tabbed UI
- ✅ **Kubernetes Health Integration** - Real-time service health from K8s deployments
- ✅ **GitHub Teams** - Service ownership and permission management
- ✅ **Multi-Environment** - Service deployment across dev, staging, production
- ✅ **Git Versioning** - Automatic service definition versioning

**Enhanced Theme System:**

- ✅ **Dark/Light Mode** - Toggle between light and dark themes with header control
- ✅ **Color Themes** - 9 pre-built color palettes (Neutral, Red, Rose, Orange, Green, Blue, Yellow, Violet, Slate)
- ✅ **Persistence** - Theme preferences saved to localStorage
- ✅ **Responsive Logo** - Logo adapts to selected theme and mode

**Plugin UI Improvements:**

- ✅ **Kubernetes** - Single sidebar menu, optimized caching, modern pod/deployment interfaces
- ✅ **AWS** - Unified account selector, performance optimizations
- ✅ **All Plugins** - Consistent visual design, improved loading states

### **Core Features**

### **📋 Service Catalog**

- **Service Registry** - Complete service definitions with YAML-based storage
- **Team Ownership** - GitHub team-based service ownership and permissions
- **Health Monitoring** - Real-time health aggregation from Kubernetes deployments
- **Multi-Environment** - Service deployment tracking across environments
- **Search & Filter** - Advanced filtering by team, tier, technology, and health status

### **☁️ AWS Operations**

- **EC2 Management** - Start, stop, monitor instances with modern interface
- **Multi-Account** - Unified account selector and switching
- **Cost Optimization** - Instance lifecycle management

### **⚙️ Kubernetes Operations**

- **Multi-cluster Support** - Unified cluster context switching in single menu
- **Enhanced Workload Management** - Modern deployment and pod interfaces with restart/scale actions
- **Advanced Pod Logs** - Real-time log streaming with search, filter, and copy functionality
- **Node Monitoring** - Comprehensive resource usage with disk, conditions, and age information
- **Optimized Performance** - Intelligent API caching and shared namespace management

### **🔐 Authentication & Security**

- **GitHub OAuth2** - Enterprise SSO integration
- **Organization Permissions** - Team-based access control
- **Session Management** - Secure token handling
- **Audit Logging** - Track all operations

### **📊 Dashboard & Monitoring**

- **Unified Dashboard** - Cross-platform overview
- **Real-time Metrics** - Live system status
- **Resource Utilization** - Performance insights

### **🔮 Planned Features (Integration Roadmap)**

#### **Phase 2 - Observability Hub (Q2 2025)**

- **📊 Grafana Integration** - Embedded dashboards with service filtering
- **📈 Prometheus Integration** - Metrics aggregation with AI insights
- **🔍 Loki Integration** - Log search with service context
- **🤖 AI Assistant V1** - Troubleshooting automation

#### **Phase 3 - Pipeline Hub (Q3 2025)**

- **🔄 ArgoCD Integration** - GitOps workflow with service context
- **⚙️ GitHub Actions** - Status tracking and deployment history
- **🤖 AI Assistant V2** - Deployment intelligence and impact analysis

#### **Phase 4 - Multi-Cloud Hub (Q4 2025)**

- **☁️ GCP Integration** - Google Cloud resources and billing
- **🔷 Azure Integration** - Microsoft Azure management
- **💰 Cost Intelligence** - AI-powered optimization suggestions

#### **Phase 5 - Community Ecosystem (2026+)**

- **🧩 Plugin SDK** - Third-party integration framework
- **🌟 Plugin Marketplace** - Community registry
- **🤖 AI Assistant V3** - Cross-tool workflows with natural language

## 📖 Documentation

### **📚 General Documentation**

- **[Getting Started Guide](./docs/README.md)** - Detailed setup and configuration
- **[Plugin Documentation](./docs/plugins/README.md)** - Individual plugin guides

### **🔧 Development Documentation**

- **[Frontend Guide](./front/README.md)** - React/TypeScript development
- **[API Documentation](./docs/README.md#api)** - Backend Go development

### **🔌 Plugin Guides**

- **[OAuth2 Plugin](./docs/plugins/oauth2.md)** - Authentication setup
- **[Kubernetes Plugin](./docs/plugins/kubernetes.md)** - K8s configuration
- **[AWS Plugin](./docs/plugins/aws.md)** - AWS integration

## 🛠️ Configuration

### **Basic Configuration File** (`dash-ops.yaml`)

```yaml
# Server Configuration
port: 8080
origin: http://localhost:8080
headers:
  - 'Content-Type'
  - 'Authorization'
front: app

# Enable Plugins
plugins:
  - 'OAuth2'
  - 'Kubernetes'
  - 'AWS'

# GitHub OAuth2 Setup
oauth2:
  - provider: github
    clientId: ${GITHUB_CLIENT_ID}
    clientSecret: ${GITHUB_CLIENT_SECRET}
    authURL: 'https://github.com/login/oauth/authorize'
    tokenURL: 'https://github.com/login/oauth/access_token'
    redirectURL: 'http://localhost:8080/api/oauth/redirect'
    urlLoginSuccess: 'http://localhost:8080'
    orgPermission: 'your-github-org' # Replace with your org
    scopes: [user, repo, read:org]

# Kubernetes Configuration
kubernetes:
  - name: 'Development'
    kubeconfig: ${HOME}/.kube/config
    context: 'dev-cluster'
  - name: 'Production'
    kubeconfig: ${HOME}/.kube/config
    context: 'prod-cluster'

# AWS Configuration
aws:
  - name: 'Production Account'
    region: us-east-1
    accessKeyId: ${AWS_ACCESS_KEY_ID}
    secretAccessKey: ${AWS_SECRET_ACCESS_KEY}
    ec2Config:
      skipList:
        - 'EKSWorkerAutoScalingGroupSpot'
```

### **Environment Variables**

```bash
# Required for OAuth2
export GITHUB_CLIENT_ID="your-github-oauth-app-id"
export GITHUB_CLIENT_SECRET="your-github-oauth-app-secret"

# Required for AWS (if using AWS plugin)
export AWS_ACCESS_KEY_ID="your-aws-access-key"
export AWS_SECRET_ACCESS_KEY="your-aws-secret-key"

# Optional: Custom API URL for frontend
export VITE_API_URL="http://localhost:8080/api"
```

## 🐳 Deployment

> **⚠️ Warning**: This is a beta project under active development. **DO NOT USE IN PRODUCTION ENVIRONMENTS.**

### **Development Environment**

```bash
# Terminal 1: Backend
go run main.go

# Terminal 2: Frontend
cd front
yarn dev

# Access: http://localhost:5173
```

### **Testing/Staging with Docker**

```bash
# ⚠️ FOR TESTING ONLY - NOT PRODUCTION READY
docker run --rm \
  -v $(pwd)/dash-ops.yaml:/dash-ops.yaml \
  -v ${HOME}/.kube/config:/.kube/config \
  --env-file .env \
  -p 8080:8080 \
  dashops/dash-ops
```

### **Testing with Helm** (Development Clusters Only)

```bash
# ⚠️ FOR DEVELOPMENT/TESTING ONLY
helm repo add dash-ops-charts ./helm-charts

# Install with custom values
helm install dash-ops dash-ops-charts/dash-ops \
  --values ./your-values.yaml

# Access via ingress or port-forward
kubectl port-forward service/dash-ops 8080:8080
```

## 🤝 Contributing

We welcome contributions! Here's how to get started:

### **🔨 Development Setup**

```bash
# 1. Fork and clone the repository
git clone https://github.com/your-username/dash-ops.git
cd dash-ops

# 2. Backend setup (Go)
go mod download
go run main.go

# 3. Frontend setup (TypeScript/React)
cd front
yarn
yarn dev

# 4. Run tests
go test ./...           # Backend tests
cd front && yarn test   # Frontend tests
```

### **📋 Development Workflow**

```bash
# 1. Create feature branch
git checkout -b feature/amazing-new-feature

# 2. Make changes with quality checks
cd front
yarn type-check:watch  # Terminal 1: TypeScript validation
yarn dev               # Terminal 2: Development server

# 3. Ensure quality before commit
yarn quality           # Type check + lint + format
go test ./...          # Backend tests

# 4. Commit with semantic messages
git commit -m "feat: add amazing new feature"

# 5. Push and create PR
git push origin feature/amazing-new-feature
```

### **🎯 High-Priority Contribution Areas**

#### **🔥 Critical for Integration Hub**

- **🔗 Tool Integrations** - Grafana, Prometheus, Loki, ArgoCD connections
- **🤖 AI Context Layer** - Cross-tool data correlation and insights
- **🎨 UX Unification** - Consistent interface across all integrations
- **⚡ Performance** - Smart caching and aggregation optimizations

#### **✨ Integration Development**

- **📊 Observability Hub** - Grafana/Prometheus/Loki integration (Phase 2)
- **🔄 Pipeline Integration** - ArgoCD and GitHub Actions support (Phase 3)
- **☁️ Multi-Cloud** - GCP and Azure integrations (Phase 4)
- **🧩 Plugin System** - Community-extensible plugin framework

#### **🤖 AI & UX Enhancement**

- **AI Assistant** - Contextual help and troubleshooting automation
- **Developer UX** - Intuitive workflows that abstract tool complexity
- **📖 Documentation** - Integration guides and plugin development docs
- **🚀 Performance** - Cross-tool performance optimizations

### **💻 Code Standards**

#### **Backend (Go)**

- **Testing** - Unit tests required for new features
- **Documentation** - Godoc comments for public functions
- **Error Handling** - Proper error propagation
- **Code Review** - All changes require review

#### **Frontend (TypeScript/React)**

- **TypeScript Strict** - No `any` types allowed
- **Testing** - Component tests with Testing Library
- **Code Quality** - ESLint + Prettier enforced
- **Semantic Commits** - Conventional commit messages

## 🏷️ Plugin Development

### **Creating a New Plugin**

1. **Backend Plugin** (`pkg/your-plugin/`)

   ```go
   package yourplugin

   type Config struct {
       // Plugin configuration
   }

   func NewHandler(config Config) *Handler {
       // Implementation
   }
   ```

2. **Frontend Module** (`front/src/modules/your-plugin/`)

   ```typescript
   // index.tsx - Module configuration
   export default {
     menus: Menu[],
     routers: Router[],
   };
   ```

3. **Documentation** (`docs/plugins/your-plugin.md`)
   - Configuration options
   - API endpoints
   - Usage examples

See [Plugin Development Guide](./docs/plugins/README.md) for detailed instructions.

## 🛡️ Security

> **⚠️ Beta Security Notice**: Current security implementation is basic and **NOT suitable for production environments**.

### **Current Authentication (Beta)**

1. **GitHub OAuth2** - Basic SSO integration
2. **Organization Validation** - Simple team membership check
3. **Basic Permissions** - Limited role-based access
4. **Token Storage** - Browser localStorage (insecure for production)

### **Security Limitations (Beta)**

- **No encrypted storage** - Tokens stored in plain text
- **Limited audit logging** - Basic action tracking only
- **No session management** - Simple token-based auth
- **Missing RBAC** - Rudimentary permission system
- **No rate limiting** - API endpoints unprotected

### **Planned Security Enhancements**

- **Enterprise SSO** - SAML, OIDC, Active Directory
- **Encrypted Storage** - Secure credential management
- **Comprehensive RBAC** - Fine-grained permissions
- **Audit & Compliance** - Full action logging and reporting
- **API Security** - Rate limiting, WAF integration

## 📊 Monitoring & Observability

### **Health Endpoints**

- `GET /api/health` - Application health status
- `GET /api/version` - Version and build information

### **Metrics** (Planned)

- **Prometheus Integration** - System metrics
- **Application Metrics** - Custom business metrics
- **Tracing** - Request flow visualization

## 🔗 Useful Links

### **Project Resources**

- **[Homepage](https://dash-ops.github.io/)** - Project website
- **[GitHub Repository](https://github.com/dash-ops/dash-ops)** - Source code and issues
- **[Docker Hub](https://hub.docker.com/r/dashops/dash-ops)** - Beta container images
- **[Helm Charts](./helm-charts/)** - Development deployment charts

### **Community**

- **[Issues](https://github.com/dash-ops/dash-ops/issues)** - Bug reports and feature requests
- **[Discussions](https://github.com/dash-ops/dash-ops/discussions)** - Community discussions
- **[Contributing Guide](./CONTRIBUTING.md)** - Detailed contribution guidelines

## 🎊 Project Status

> **🚧 BETA VERSION** - Active development, breaking changes expected

| Component              | Status         | Maturity Level                         |
| ---------------------- | -------------- | -------------------------------------- |
| **Backend API**        | 🔄 Beta        | Go 1.21+ - Under development           |
| **Frontend**           | 🔄 Beta        | TypeScript + React 18 - Stable UI      |
| **AWS Plugin**         | 🔄 Alpha       | EC2 Operations - Basic features        |
| **Kubernetes Plugin**  | 🔄 Alpha       | Multi-cluster Support - Basic features |
| **OAuth2 Plugin**      | 🔄 Beta        | GitHub Integration - Functional        |
| **Docker Images**      | ✅ Available   | Multi-arch Support - Testing only      |
| **Helm Charts**        | 🔄 Alpha       | K8s Deployment - Development only      |
| **Documentation**      | 🔄 In Progress | Comprehensive Guides                   |
| **Service Catalog**    | 📋 Planned     | Q3 2025- Service lifecycle management  |
| **Observability**      | 📋 Planned     | Q4 2025 - Monitoring integrations      |
| **FinOps Integration** | 📋 Planned     | Q1 2026 - Cost management              |

### **Production Readiness**

❌ **NOT RECOMMENDED FOR PRODUCTION USE**

- Missing enterprise security features
- Limited error handling and recovery
- No data persistence layer
- Incomplete access control system
- Missing monitoring and alerting
- Breaking changes expected in updates

## 📄 License

This project is licensed under the [MIT License](./LICENSE) - see the license file for details.

---

**⚠️ Beta software - Use for testing and evaluation only!** 🧪

For detailed setup instructions, visit the [documentation directory](./docs/README.md).

For frontend development, see the [frontend guide](./front/README.md).
