# Service Catalog Plugin

> **🔄 Beta Plugin** - Service registry with Kubernetes integration. Feature-complete but still in testing phase.

The Service Catalog plugin provides a comprehensive service registry and management system, designed to give teams visibility and control over their services across multiple environments with deep Kubernetes integration.

## 🎯 Features

### **Current Capabilities (Beta)**

- **Service Registry** - Complete CRUD operations for service definitions
- **Kubernetes Integration** - Real-time health checks and deployment status
- **Multi-Environment Support** - Manage services across dev, staging, and production
- **Team-based Organization** - GitHub teams integration with ownership model
- **Service Tiers** - TIER-1 (Critical), TIER-2 (Important), TIER-3 (Standard) classification
- **Advanced Filtering** - Search, filter by tier, team, status, and more
- **Health Aggregation** - Contextual health status from Kubernetes deployments
- **Flexible Versioning** - Multiple versioning options: Git, Simple, or None

### **Recent Updates (v0.2.0-beta)**

**Modern UI Interface:**

- ✅ **Tabbed Forms**: Multi-step service creation with Basic Info, Kubernetes, Observability, and Review tabs
- ✅ **Real-time Health**: Live service health status with color-coded indicators
- ✅ **Smart Search**: Comprehensive filtering by name, team, technology, and status
- ✅ **Team Context**: "My Team" filtering based on GitHub organization membership
- ✅ **Responsive Design**: Modern card-based layout with shadcn/ui components

**Service Management:**

- ✅ **Service Tiers**: Three-tier classification system for business impact assessment
- ✅ **GitHub Integration**: Automatic team resolution and permission-based editing
- ✅ **Kubernetes Context**: Environment-specific configurations with namespace management
- ✅ **Observability Links**: Direct integration with metrics and logging systems
- ✅ **Dependency Tracking**: Visual dependency management between services

**Backend Enhancements:**

- ✅ **Flexible Storage**: Local YAML-based service definitions with multiple versioning options
- ✅ **Versioning Providers**: Git versioning, Simple versioning, or No versioning support
- ✅ **Health API**: Real-time health aggregation from Kubernetes deployments
- ✅ **Team Authorization**: GitHub team-based read/write permissions
- ✅ **Batch Operations**: Efficient health status retrieval for multiple services

### **Planned Features (Production Roadmap)**

- **GitHub Repository Storage** - Store service definitions in GitHub repositories
- **S3 Storage Provider** - AWS S3-based service definition storage
- **Advanced Drift Detection** - Configuration drift monitoring and sync workflows
- **Service Mesh Integration** - Istio/Linkerd service mesh awareness
- **Advanced Health Checks** - Custom health check definitions and alerting
- **Service Dependencies** - Automated dependency discovery and visualization
- **Compliance Reporting** - Service compliance tracking and reporting

## 🔧 Configuration

### **1. Basic Plugin Setup**

```yaml
# Enable Service Catalog plugin
plugins:
  - 'ServiceCatalog'

service-catalog:
  storage:
    provider: 'filesystem' # filesystem, github, s3
    filesystem:
      directory: './services' # Directory to store service definitions
  versioning:
    enabled: true # Enable/disable versioning (default: false)
    provider: 'simple' # git, simple, none (auto-detect if not specified)
```

### **2. GitHub Integration (Required)**

Service Catalog requires Auth plugin for GitHub team integration:

```yaml
plugins:
  - 'Auth' # Required for GitHub teams
  - 'ServiceCatalog'

auth:
  github:
    clientId: ${GITHUB_CLIENT_ID}
    clientSecret: ${GITHUB_CLIENT_SECRET}
    org: 'your-organization'
```

### **3. Kubernetes Integration (Recommended)**

For real-time health checks, configure Kubernetes plugin:

```yaml
plugins:
  - 'Auth'
  - 'Kubernetes'
  - 'ServiceCatalog'

kubernetes:
  - name: 'Development'
    context: 'dev-cluster'
    kubeconfig: ${HOME}/.kube/config
  - name: 'Production'
    context: 'prod-cluster'
    kubeconfig: ${HOME}/.kube/config
```

### **4. Advanced Configuration**

```yaml
service-catalog:
  storage:
    provider: 'filesystem'
    filesystem:
      directory: './services'

  versioning:
    enabled: true
    provider: 'git' # git, simple, none

  # Default service configuration
  defaults:
    tier: 'TIER-3'
    impact: 'medium'

  # Health check configuration
  health:
    timeout: '30s'
    cache_ttl: '5m'
    batch_size: 10
```

## 📋 Service Management

### **Service Definition Structure**

Services are defined using a structured YAML format:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: user-authentication
  tier: TIER-1
  created_at: '2024-01-15T10:30:00Z'
  created_by: 'john.doe'
  updated_at: '2024-01-15T14:20:00Z'
  updated_by: 'jane.smith'
  version: 3

spec:
  description: 'User authentication and authorization service'

  team:
    github_team: 'auth-squad'

  business:
    impact: 'high'
    sla_target: '99.9%'
    dependencies:
      - 'user-database'
      - 'oauth-provider'

  technology:
    language: 'Node.js'
    framework: 'Express'

  kubernetes:
    environments:
      - name: 'production'
        context: 'prod-cluster'
        namespace: 'auth'
        resources:
          deployments:
            - name: 'auth-api'
              replicas: 5
              resources:
                requests:
                  cpu: '100m'
                  memory: '128Mi'
                limits:
                  cpu: '500m'
                  memory: '512Mi'
      - name: 'staging'
        context: 'staging-cluster'
        namespace: 'auth-staging'
        resources:
          deployments:
            - name: 'auth-api'
              replicas: 2

  observability:
    metrics:
      url: 'https://grafana.company.com/d/auth-dashboard'
    logs:
      url: 'https://kibana.company.com/app/logs/auth-service'

  runbooks:
    - name: 'Incident Response'
      url: 'https://wiki.company.com/auth-service/incidents'
    - name: 'Deployment Guide'
      url: 'https://wiki.company.com/auth-service/deployment'
```

### **Service Tiers**

Services are classified into three tiers based on business impact:

| Tier       | Description | Characteristics                      | Example Services               |
| ---------- | ----------- | ------------------------------------ | ------------------------------ |
| **TIER-1** | Critical    | Revenue impact, customer-facing      | Payment API, Authentication    |
| **TIER-2** | Important   | Internal operations, indirect impact | User Management, Notifications |
| **TIER-3** | Standard    | Support functions, minimal impact    | Logging, Monitoring Tools      |

### **Team-based Ownership**

Services are owned by GitHub teams with automatic permission management:

- **Service Owners** - Can edit service definitions and configurations
- **Team Members** - Can view all team services and health status
- **Organization Members** - Can view public service information
- **External Users** - No access (requires authentication)

## 🔍 Service Discovery & Filtering

### **Search Capabilities**

The Service Catalog provides comprehensive search and filtering:

- **Text Search** - Service name, description, technology, team names
- **Team Filtering** - "My Team" services, specific team services, or all services
- **Tier Filtering** - Filter by service tier (TIER-1, TIER-2, TIER-3)
- **Status Filtering** - Filter by health status (healthy, degraded, critical, down)
- **Sort Options** - By name, tier, team, or last updated

### **Health Status Indicators**

Real-time health status with contextual meaning:

| Status       | Color     | Description                                   | Kubernetes Mapping                            |
| ------------ | --------- | --------------------------------------------- | --------------------------------------------- |
| **Healthy**  | 🟢 Green  | All replicas running as expected              | Desired replicas = Ready replicas             |
| **Drift**    | 🔵 Blue   | Service running but configuration out of sync | Replicas running < desired (but > 0)          |
| **Degraded** | 🟡 Yellow | Partial functionality, some issues            | Some replicas failing, errors in logs         |
| **Critical** | 🟠 Orange | Significant issues affecting functionality    | Most replicas failing, high error rates       |
| **Down**     | 🔴 Red    | Service completely unavailable                | No replicas running, deployment failed        |
| **Unknown**  | ⚪ Gray   | Health status cannot be determined            | No Kubernetes integration or data unavailable |

### **Statistics Dashboard**

Service overview with key metrics:

- **Total Services** - All registered services across teams
- **My Team Services** - Services owned by user's GitHub teams
- **Service Distribution** - Count by tier (TIER-1, TIER-2, TIER-3)
- **Critical Alerts** - Services requiring immediate attention
- **Editable Services** - Services user has permission to modify

## 📊 Health Monitoring

### **Kubernetes Integration**

The Service Catalog automatically aggregates health information from Kubernetes:

```yaml
# Service definition includes Kubernetes contexts
kubernetes:
  environments:
    - name: 'production'
      context: 'prod-cluster'
      namespace: 'api'
      resources:
        deployments:
          - name: 'user-api'
            replicas: 3
```

**Health Calculation Logic:**

1. **Deployment Status** - Compare desired vs ready replicas
2. **Pod Health** - Check pod restart counts and error rates
3. **Resource Utilization** - Monitor CPU/memory usage
4. **Age Analysis** - Consider deployment age for stability assessment

### **Multi-Environment Health**

Services can span multiple environments with aggregated health:

- **Production** - Weighted highest in overall health calculation
- **Staging** - Secondary priority for health assessment
- **Development** - Informational only, minimal impact on overall status

### **Health API Endpoints**

```bash
# Get health for single service
GET /api/v1/service-catalog/services/{service-name}/health

# Get batch health for multiple services
POST /api/v1/service-catalog/services/health/batch
{
  "service_names": ["auth-api", "user-api", "payment-api"]
}
```

**Response Format:**

```json
{
  "service_name": "user-authentication",
  "overall_status": "healthy",
  "environments": [
    {
      "name": "production",
      "status": "healthy",
      "deployments": [
        {
          "name": "auth-api",
          "status": "healthy",
          "replicas": {
            "desired": 5,
            "ready": 5,
            "available": 5
          },
          "last_updated": "2024-01-15T14:30:00Z"
        }
      ]
    }
  ],
  "last_check": "2024-01-15T15:00:00Z"
}
```

## 🔐 Security & Permissions

### **Authentication Requirements**

Service Catalog requires Auth authentication with GitHub:

```yaml
auth:
  github:
    clientId: ${GITHUB_CLIENT_ID}
    clientSecret: ${GITHUB_CLIENT_SECRET}
    org: 'your-organization' # Required for team membership
```

### **Permission Model**

**Read Permissions:**

- ✅ **Organization Members** - Can view all service information
- ✅ **Team Members** - Enhanced view for team-owned services
- ❌ **External Users** - No access without authentication

**Write Permissions:**

- ✅ **Service Owners** - GitHub team members can edit team services
- ❌ **Other Teams** - Cannot modify services owned by other teams
- ❌ **Read-only Users** - Cannot create or modify any services

### **Team Resolution**

GitHub team membership is resolved at runtime:

1. **User Authentication** - Auth GitHub token validation
2. **Organization Check** - Verify user is member of configured organization
3. **Team Membership** - Fetch user's team memberships via GitHub API
4. **Permission Mapping** - Match teams against service ownership

## 🗄️ Storage & Versioning

### **Storage Providers**

#### **Filesystem Storage (Current)**

Services are stored as YAML files with flexible versioning options:

```
./services/
├── user-authentication.yaml
├── payment-processor.yaml
├── notification-service.yaml
├── .git/                    # Git versioning (optional)
└── .history/               # Simple versioning (optional)
    ├── user-authentication.json
    └── payment-processor.json
```

#### **GitHub Repository Storage** (Roadmap):

```yaml
service-catalog:
  storage:
    provider: 'github'
    config:
      repository: 'your-org/service-definitions'
      branch: 'main'
      path: 'services/'
      auth_client: true # Use Auth GitHub client
```

#### **S3 Storage** (Roadmap):

```yaml
service-catalog:
  storage:
    provider: 's3'
    config:
      bucket: 'company-service-catalog'
      prefix: 'services/'
      aws_client: true # Use AWS plugin client
```

### **Versioning Providers**

The service catalog supports flexible versioning through multiple providers:

#### **1. Git Versioning (`git`)**

- **Best for**: Filesystem storage with full Git integration
- **Features**:
  - Full Git history with commits, authors, and timestamps
  - Automatic Git repository initialization
  - Proper commit messages with service metadata
- **Requirements**: Git must be installed and available in PATH
- **Storage compatibility**: Filesystem only

**Configuration:**

```yaml
service-catalog:
  storage:
    provider: 'filesystem'
    filesystem:
      directory: './services'
  versioning:
    enabled: true
    provider: 'git'
```

#### **2. Simple Versioning (`simple`)**

- **Best for**: Filesystem or S3 storage without Git dependency
- **Features**:
  - JSON-based history storage in `.history/` directory
  - Service change tracking with timestamps and user info
  - Limited to last 100 changes per service (configurable)
- **Requirements**: None (pure Go implementation)
- **Storage compatibility**: Filesystem (S3 support planned)

**Configuration:**

```yaml
service-catalog:
  storage:
    provider: 'filesystem'
    filesystem:
      directory: './services'
  versioning:
    enabled: true
    provider: 'simple'
```

#### **3. No Versioning (`none`)**

- **Best for**: Minimal setups or when versioning is not needed
- **Features**: Disables all versioning functionality
- **Requirements**: None
- **Storage compatibility**: All storage providers

**Configuration:**

```yaml
service-catalog:
  storage:
    provider: 'filesystem'
    filesystem:
      directory: './services'
  versioning:
    enabled: false
```

### **Auto-Detection (Default Behavior)**

When `versioning.provider` is not specified, the system automatically chooses:

- **Filesystem storage**: Uses `git` if Git is available, falls back to `simple`
- **Other storage**: Uses `simple`

### **Configuration Examples**

#### **Development Setup (No Versioning)**

```yaml
service-catalog:
  storage:
    provider: 'filesystem'
    filesystem:
      directory: './services'
  versioning:
    enabled: false
```

#### **Production Setup (Git Versioning)**

```yaml
service-catalog:
  storage:
    provider: 'filesystem'
    filesystem:
      directory: '/var/lib/dash-ops/services'
  versioning:
    enabled: true
    provider: 'git'
```

#### **Cloud Setup (Simple Versioning)**

```yaml
service-catalog:
  storage:
    provider: 's3'
    s3:
      bucket: 'company-service-definitions'
  versioning:
    enabled: true
    provider: 'simple' # Git not available in cloud environments
```

### **Migration from Git-Only System**

If you're upgrading from the previous Git-only system:

1. **No changes needed** for Git-based setups - the system will continue using Git
2. **To disable Git dependency**: Set `versioning.provider: 'simple'`
3. **To disable versioning**: Set `versioning.enabled: false`

### **Performance Considerations**

- **Git Versioning**: Slower for large histories, but provides full Git features
- **Simple Versioning**: Faster, limited history (100 entries per service by default)
- **No Versioning**: Fastest, no history tracking

## 📊 API Reference

### **Service CRUD Operations**

```bash
# List all services
GET /api/v1/service-catalog/services

# Get specific service
GET /api/v1/service-catalog/services/{service-name}

# Create new service
POST /api/v1/service-catalog/services
Content-Type: application/json

# Update existing service
PUT /api/v1/service-catalog/services/{service-name}
Content-Type: application/json

# Delete service
DELETE /api/v1/service-catalog/services/{service-name}
```

### **Health & Status Operations**

```bash
# Get service health
GET /api/v1/service-catalog/services/{service-name}/health

# Batch health check
POST /api/v1/service-catalog/services/health/batch
```

### **Team & Permission Operations**

```bash
# Get services by team
GET /api/v1/service-catalog/services?team={github-team}

# Get services by tier
GET /api/v1/service-catalog/services?tier={TIER-1|TIER-2|TIER-3}

# Get user permissions
GET /api/v1/service-catalog/user/permissions
```

## 🧪 Testing

### **Test Service Configuration**

```yaml
service-catalog:
  storage:
    provider: 'filesystem'
    filesystem:
      directory: './test-services'
  versioning:
    enabled: false # Disable for testing
```

### **Creating Test Services**

```bash
# Create a test service via API
curl -X POST http://localhost:8080/api/v1/service-catalog/services \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${GITHUB_TOKEN}" \
  -d '{
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
      "name": "test-api",
      "tier": "TIER-3"
    },
    "spec": {
      "description": "Test API for development",
      "team": {
        "github_team": "developers"
      },
      "business": {
        "impact": "low"
      }
    }
  }'
```

### **Health Check Testing**

Test health integration with a sample Kubernetes deployment:

```bash
# Deploy test service to Kubernetes
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-api
  namespace: default
spec:
  replicas: 2
  selector:
    matchLabels:
      app: test-api
  template:
    metadata:
      labels:
        app: test-api
    spec:
      containers:
      - name: api
        image: nginx:alpine
        ports:
        - containerPort: 80
EOF

# Check service health
curl http://localhost:8080/api/v1/service-catalog/services/test-api/health
```

## 🐛 Troubleshooting

### **Common Issues**

#### **Services Not Appearing**

- ✅ Check GitHub Auth authentication is working
- ✅ Verify user is member of configured GitHub organization
- ✅ Confirm service files exist in configured directory
- ✅ Check file permissions for service directory

#### **Health Status Shows Unknown**

- ✅ Verify Kubernetes plugin is enabled and configured
- ✅ Check Kubernetes context is accessible
- ✅ Confirm deployment names match service configuration
- ✅ Validate namespace permissions

#### **Permission Denied on Service Edit**

- ✅ Verify user is member of service's GitHub team
- ✅ Check GitHub team name matches exactly (case-sensitive)
- ✅ Confirm Auth token has org:read permissions
- ✅ Validate team membership via GitHub API

#### **Versioning Errors**

**Git Versioning Issues:**

- ✅ Ensure services directory is writable
- ✅ Check Git is installed and accessible
- ✅ Verify Git user configuration (name and email)
- ✅ Confirm no conflicting Git operations

**Simple Versioning Issues:**

- ✅ Check write permissions on services directory
- ✅ Verify `.history/` directory can be created
- ✅ Ensure sufficient disk space for history files

**General Versioning Issues:**

- ✅ Verify `versioning.enabled: true` in configuration
- ✅ Check `versioning.provider` is set correctly
- ✅ Confirm storage provider supports chosen versioning method

### **Debug Configuration**

```yaml
service-catalog:
  debug: true # Enable verbose logging
  storage:
    provider: 'filesystem'
    filesystem:
      directory: './services'
  versioning:
    enabled: false # Disable for troubleshooting
  health:
    timeout: '60s' # Increase timeout for debugging
```

### **Debug Commands**

```bash
# Check service catalog status
curl http://localhost:8080/api/health | jq .service_catalog

# Validate user permissions
curl -H "Authorization: Bearer ${TOKEN}" \
  http://localhost:8080/api/v1/service-catalog/user/permissions

# Check GitHub team membership
curl -H "Authorization: Bearer ${TOKEN}" \
  https://api.github.com/user/teams
```

## 🤝 Contributing

### **Priority Areas**

1. **🔒 Security** - Enhanced team permission models and audit trails
2. **📊 Monitoring** - Advanced health check definitions and alerting
3. **🧪 Testing** - Comprehensive integration test suite
4. **🔌 Storage** - GitHub and S3 storage provider implementations
5. **📖 Documentation** - Usage guides and best practice documentation

### **Development Setup**

```bash
# 1. Start backend with service catalog enabled
GITHUB_CLIENT_ID=your_id GITHUB_CLIENT_SECRET=your_secret go run main.go

# 2. Configure test services directory
mkdir -p ./test-services

# 3. Test API endpoints
curl http://localhost:8080/api/v1/service-catalog/services

# 4. Run frontend with service catalog
cd front && npm run dev
```

### **Adding New Features**

1. **Backend Changes** - Extend `pkg/service-catalog/` handlers and types
2. **Frontend Updates** - Enhance `front/src/modules/service-catalog/` components
3. **Documentation** - Update this documentation with new features
4. **Testing** - Add integration tests for new functionality

## 📚 Resources

- **[OpenAPI Service Specification](https://swagger.io/specification/)** - API documentation standards
- **[Kubernetes Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)** - CRD inspiration
- **[GitHub Teams API](https://docs.github.com/en/rest/teams)** - Team membership integration
- **[Git Integration Patterns](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects)** - Version control best practices

---

**🔄 Beta Notice**: The Service Catalog plugin is feature-complete but still in beta testing. While suitable for development and staging environments, production deployment should include proper backup and monitoring procedures.
