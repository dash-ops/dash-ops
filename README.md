# DashOPS - Developer Experience Hub with AI (Beta)

![DashOps](https://github.com/dash-ops/dash-ops/workflows/DashOps/badge.svg)

> **⚠️ BETA VERSION - NOT RECOMMENDED FOR PRODUCTION USE**

> **The VS Code for infrastructure** - A unified, AI-powered hub that integrates all your developer tools in one intuitive interface. Focus on building features, not juggling between different platforms.

DashOPS is an **experimental integration platform** that connects your existing tools (Kubernetes, AWS, Prometheus) into a seamless developer experience, enhanced by contextual AI assistance and an extensible plugin system.
<img width="1226" height="715" alt="image" src="https://github.com/user-attachments/assets/82909cbb-8857-4c1d-b0d7-6afac3a50dc7" />

**🚧 This project is actively under development and should only be used for testing and evaluation purposes.**

## 🚀 Quick Start

### **New: Interactive Setup Wizard** ✨

No more YAML wrestling! Start DashOPS and configure everything through the UI:

```bash
# 1. Start the backend
go run main.go

# 2. Start the frontend (in another terminal)
cd front
yarn && yarn dev

# 3. Open browser at http://localhost:5173
# → First-time users are guided through an interactive setup wizard
```

**What the Setup Wizard does:**
- 🎯 **Select your plugins** - Pick only what you need (AWS, Kubernetes, Auth, etc.)
- 🔐 **Securely store credentials** - Secrets are masked in UI, never exposed via API
- 📝 **Generate dash-ops.yaml** - Configuration auto-saved to disk
- ♻️ **Update anytime** - Revisit Settings to adjust providers without revealing secrets

> **💡 Tip**: The wizard writes to `./dash-ops.yaml` (or `$DASH_CONFIG` if set). You can still edit the file manually or version it in Git!

### Option 2: Manual YAML Setup (For Infrastructure-as-Code fans)

If you prefer to version your config from day one:

```bash
# 1. Create configuration file
cp local.sample.yaml dash-ops.yaml

# 2. Edit dash-ops.yaml and set your credentials
# (See full configuration examples at https://dash-ops.github.io/)

# 3. Set environment variables (for sensitive values)
export GITHUB_CLIENT_ID="your-client-id"
export GITHUB_CLIENT_SECRET="your-client-secret"

# 4. Start services
go run main.go  # Backend on :8080
cd front && yarn dev  # Frontend on :5173
```

### Option 3: Docker (Quick Test)

```bash
docker run --rm \
  -v $(pwd)/dash-ops.yaml:/dash-ops.yaml \
  -v ${HOME}/.kube/config:/.kube/config \
  -e GITHUB_CLIENT_ID \
  -e GITHUB_CLIENT_SECRET \
  -p 8080:8080 \
  dashops/dash-ops
```

## 📖 Documentation

### **Complete Guides**
Visit **[dash-ops.github.io](https://dash-ops.github.io/)** for comprehensive documentation:

- **[Setup Wizard Guide](https://dash-ops.github.io/#setup-wizard)** - Step-by-step interactive setup
- **[Installation](https://dash-ops.github.io/#installation)** - Full installation guide
- **[Configuration Reference](https://dash-ops.github.io/#initial-setup)** - All config options explained
- **[Plugin Guides](https://dash-ops.github.io/#plugins-overview)** - AWS, Kubernetes, Auth, Observability
- **[API Reference](https://dash-ops.github.io/#api-intro)** - REST API documentation
- **[Contributing](https://dash-ops.github.io/#contributing)** - Contribution guidelines

### **Local Development Docs**
- **[Backend Architecture](./docs/backend-development-guide.md)** - Hexagonal architecture deep-dive
- **[Frontend Guide](./front/README.md)** - React/TypeScript development

## 🎯 Features

### **🆕 Latest Updates (v0.5.0-alpha)**

**Interactive Setup Module:**
- ✅ **Setup Wizard** - Guided first-run experience with 6 configuration tabs
- ✅ **Settings Page** - Update providers and credentials without exposing secrets
- ✅ **Secret Masking** - Never return sensitive values from API (shows `*****`)
- ✅ **Live Validation** - Real-time feedback on configuration errors
- ✅ **YAML Auto-Generation** - Writes to `dash-ops.yaml` or `$DASH_CONFIG` path

**Enhanced Architecture:**
- ✅ **Settings Module** - Replaces legacy config with hexagonal design
- ✅ **100% Test Coverage** - Comprehensive tests for critical module loading logic
- ✅ **Setup Mode Detection** - Frontend auto-detects first-run and shows wizard

### **Core Features**

### **📋 Service Catalog** ✅ Available
- Service registry with YAML-based storage
- GitHub team-based ownership and permissions
- Real-time health aggregation from Kubernetes
- Multi-environment deployment tracking
- Advanced search and filtering

### **📊 Observability Hub** ✅ Beta (v0.4.0)
- **Logs** - Loki integration with histogram visualization
- **Traces** - Tempo integration with trace timeline
- **Metrics** - Prometheus integration (backend complete)
- Service-aware monitoring with automatic filtering

### **☁️ AWS Operations** 🔄 Alpha
- EC2 instance management (start, stop, monitor)
- Multi-account support with unified selector
- Cost optimization through lifecycle management

### **⚙️ Kubernetes Operations** 🔄 Alpha
- Multi-cluster context switching
- Workload management (deployments, pods)
- Real-time log streaming with search
- Node monitoring with resource usage

### **🔐 Authentication** 🔄 Beta
- GitHub OAuth2 integration
- Organization-based access control
- Session management and audit logging

## 🏗️ Architecture

### Backend (Go + Hexagonal Architecture)

All modules follow a consistent 8-layer pattern:

```
pkg/{module}/
├── adapters/     # External integrations & data transformation
├── controllers/  # Business logic orchestration
├── handlers/     # HTTP endpoints (centralized in http.go)
├── logic/        # Pure business logic (100% tested)
├── models/       # Domain entities with behavior
├── ports/        # Interfaces & contracts
├── repositories/ # Data persistence
├── wire/         # API contracts (DTOs)
└── module.go     # Module factory & initialization
```

**Key Benefits:**
- ✅ Consistent structure across all 8 modules
- ✅ High testability with 80+ unit tests
- ✅ Interface-based design for easy extension
- ✅ Clear separation of concerns

### Frontend (React + TypeScript)

```
src/
├── modules/          # Feature modules (kubernetes, aws, settings, etc.)
│   └── {module}/
│       ├── components/    # UI components
│       ├── resources/     # API clients
│       ├── hooks/         # React hooks
│       ├── types.ts       # TypeScript types
│       └── index.tsx      # Module registration
├── components/       # Shared UI components (shadcn/ui)
├── helpers/          # Utilities (loadModules, http, oauth)
└── types/            # Global TypeScript types
```

**Features:**
- ✅ Plugin-based module loading
- ✅ TypeScript strict mode (no `any`)
- ✅ Comprehensive test coverage with Vitest
- ✅ Modern UI with shadcn/ui components

## 🤝 Contributing

We welcome contributions! Here's how to get started:

### Development Setup

```bash
# 1. Fork and clone
git clone https://github.com/your-username/dash-ops.git
cd dash-ops

# 2. Backend (Go 1.21+)
go mod download
go run main.go

# 3. Frontend (Node 18+)
cd front
yarn install
yarn dev

# 4. Run tests
go test ./...              # Backend
cd front && yarn test      # Frontend
```

### Development Workflow

```bash
# Create feature branch
git checkout -b feat/amazing-feature

# Frontend quality checks (Terminal 1)
cd front
yarn type-check:watch  # TypeScript validation
yarn test              # Run tests

# Development server (Terminal 2)
yarn dev               # Auto-reload on changes

# Before commit
yarn quality           # Type-check + lint + format
go test ./...          # Backend tests

# Commit with semantic messages (no emojis)
git commit -m "feat: add amazing new feature"
```

### High-Priority Areas

**🔥 Critical:**
- **Setup UX** - Improve wizard validation and error messages
- **Secret Management** - Encrypted storage for production readiness
- **Plugin System** - SDK for community-contributed integrations

**✨ Nice-to-Have:**
- **Grafana Integration** - Embedded dashboards
- **ArgoCD Plugin** - GitOps workflow management
- **Multi-Cloud** - GCP and Azure support
- **AI Assistant** - Contextual troubleshooting automation

### Code Standards

- **Backend**: Go conventions, godoc comments, proper error handling
- **Frontend**: TypeScript strict, no `any`, ESLint + Prettier
- **Tests**: Required for new features
- **Commits**: Semantic conventional commits (no emojis per team preference)

## 📊 Project Status

> **🚧 BETA** - Active development, breaking changes expected

| Component          | Status      | Maturity                                |
| ------------------ | ----------- | --------------------------------------- |
| **Backend API**    | 🔄 Beta     | Go 1.21+ - Stable core, evolving        |
| **Frontend**       | 🔄 Beta     | React 19 + TypeScript - Stable UI       |
| **Setup Module**   | ✅ Alpha    | Interactive wizard - v0.5.0             |
| **Settings**       | ✅ Alpha    | Provider management - v0.5.0            |
| **Service Catalog**| ✅ Beta     | Complete lifecycle management           |
| **Observability**  | ✅ Beta     | Logs & Traces - v0.4.0                  |
| **AWS Plugin**     | 🔄 Alpha    | EC2 operations - Basic features         |
| **Kubernetes**     | 🔄 Alpha    | Multi-cluster - Basic features          |
| **OAuth2**         | 🔄 Beta     | GitHub integration - Functional         |
| **Docker Images**  | ✅ Available| Multi-arch - Testing only               |
| **Helm Charts**    | 🔄 Alpha    | K8s deployment - Development only       |

### Production Readiness Checklist

❌ **NOT RECOMMENDED FOR PRODUCTION**

Missing features:
- [ ] Enterprise security (encrypted secrets, RBAC)
- [ ] Data persistence layer
- [ ] Comprehensive error recovery
- [ ] Monitoring and alerting
- [ ] Rate limiting and WAF
- [ ] Audit and compliance logging

## 🛡️ Security Notice

> **⚠️ Beta Security**: Current implementation is **NOT production-ready**

**Current Limitations:**
- Secrets stored in plain YAML files
- Basic OAuth2 without enterprise SSO
- No encrypted credential storage
- Limited audit logging
- No rate limiting on API endpoints

**Planned Enhancements:**
- Encrypted storage (HashiCorp Vault, AWS Secrets Manager)
- Enterprise SSO (SAML, OIDC)
- Comprehensive RBAC with fine-grained permissions
- Full audit trails and compliance reporting

## 🔗 Links

- **[Documentation](https://dash-ops.github.io/)** - Complete guides
- **[Issues](https://github.com/dash-ops/dash-ops/issues)** - Bug reports
- **[Discussions](https://github.com/dash-ops/dash-ops/discussions)** - Community
- **[Helm Charts](https://github.com/dash-ops/helm-charts)** - Kubernetes deployment
- **[Docker Hub](https://hub.docker.com/r/dashops/dash-ops)** - Container images

## 📄 License

MIT License - see [LICENSE](./LICENSE) for details.

---

**⚠️ Beta software - Use for testing and evaluation only!** 🧪

For detailed documentation, visit **[dash-ops.github.io](https://dash-ops.github.io/)**.
