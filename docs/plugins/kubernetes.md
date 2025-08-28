# Kubernetes Plugin

> **⚠️ Alpha Plugin** - Feature-rich cluster management. Improved UI and optimized performance but still not production-ready.

The Kubernetes plugin provides a comprehensive yet simplified interface for Kubernetes cluster management, focusing on essential operations for developers who need cluster visibility without deep Kubernetes expertise.

## 🎯 Features

### **Current Capabilities (Alpha)**

- **Multi-cluster Support** - Connect and manage multiple K8s clusters with unified context switching
- **Enhanced Workload Monitoring** - Modern UI for deployments, pods, and services with real-time status
- **Advanced Pod Logs** - Stream logs with search, filter, copy, and expand functionality
- **Smart Deployment Management** - Restart and scale operations with modal dialogs and validation
- **Comprehensive Resource Overview** - Enhanced node information with disk usage, conditions, age, and visual indicators
- **Optimized Performance** - Intelligent caching to prevent redundant API calls
- **Shared Namespace Management** - Cached namespace selector across all Kubernetes pages
- **Single Sidebar Integration** - Unified cluster context selector (no multiple menu items)

### **Recent Updates (v0.3.0-alpha)**

**Sidebar & Navigation Improvements:**

- ✅ **Single Menu Item**: Unified "Kubernetes" menu with cluster context selector dropdown
- ✅ **Context Switching**: Similar to `kubectl config use-context` but with modern UI
- ✅ **Optimized API Calls**: Intelligent caching prevents redundant cluster list requests
- ✅ **Performance**: Significant reduction in initial load and navigation times

**Enhanced Deployment Management:**

- ✅ **Modern Actions**: Replaced "Up/Down" with "Restart" and "Scale" using shadcn/ui dialogs
- ✅ **Input Validation**: Scale action includes replica count input with validation
- ✅ **Confirmation Dialogs**: Restart action has confirmation step for safety
- ✅ **Loading States**: Real-time feedback during operations
- ✅ **Permission Controls**: Separate `restart` and `scale` permission validation

**Advanced Pod Management:**

- ✅ **Visual Redesign**: Modern table layout matching deployment design patterns
- ✅ **Icon-based Actions**: Logs action converted to icon (eye icon) for better UX
- ✅ **Node Integration**: Node name displayed as table column instead of separate action
- ✅ **Container Information**: Simplified badges for ready/total containers (consistent styling)
- ✅ **Enhanced Data**: QoS class, controlled by, age, and restart count information

**Node Monitoring Enhancements:**

- ✅ **Comprehensive Metrics**: CPU, memory, and disk usage with visual progress bars
- ✅ **Health Conditions**: Real-time node condition monitoring (Ready, MemoryPressure, DiskPressure)
- ✅ **Age Calculation**: Backend-calculated node age for better performance
- ✅ **Visual Indicators**: Color-coded status bars and modern table layout

**Shared Infrastructure:**

- ✅ **Namespace Cache**: Shared caching mechanism across Pods and Deployments pages
- ✅ **"All" Option**: Namespace selector includes "All" option for viewing all resources
- ✅ **Performance**: Reduced redundant API calls for namespace and cluster data

### **Planned Features**

- **Advanced Workload Management** - StatefulSets, DaemonSets, Jobs
- **Resource Quotas** - Namespace limits and monitoring
- **ConfigMap/Secret Management** - Configuration editing interface
- **Helm Integration** - Chart deployment and management
- **Custom Resource Definitions** - CRD monitoring and management

## 🔧 Configuration

### **1. Cluster Access Setup**

#### **Method 1: External Cluster (kubeconfig)**

```yaml
# Enable Kubernetes plugin
plugins:
  - 'Kubernetes'

kubernetes:
  - name: 'Development Cluster'
    kubeconfig: ${HOME}/.kube/config
    context: 'dev-cluster-context'

  - name: 'Staging Cluster'
    kubeconfig: /path/to/staging/kubeconfig
    context: 'staging-context'
```

#### **Method 2: In-Cluster Configuration**

```yaml
kubernetes:
  - name: 'Current Cluster'
    kubeconfig: # Empty - uses in-cluster service account
```

> **📝 Note**: In-cluster configuration requires proper RBAC setup. Our [Helm charts](../../helm-charts/) include pre-configured ClusterRole permissions.

### **2. Permission Configuration**

Control which teams can perform operations:

```yaml
kubernetes:
  - name: 'Development Cluster'
    kubeconfig: ${HOME}/.kube/config
    context: 'dev'
    permission:
      deployments:
        namespaces: ['default', 'dev', 'staging'] # Allowed namespaces
        restart: ['dash-ops*developers', 'dash-ops*sre'] # Teams that can restart deployments
        scale: ['dash-ops*developers', 'dash-ops*sre'] # Teams that can scale deployments
```

### **3. Multi-Environment Setup**

```yaml
kubernetes:
  - name: 'Development'
    kubeconfig: ${HOME}/.kube/config
    context: 'dev'
    permission:
      deployments:
        namespaces: ['default', 'dev']
        restart: ['dash-ops*developers']
        scale: ['dash-ops*developers']

  - name: 'Staging'
    kubeconfig: ${HOME}/.kube/config
    context: 'staging'
    permission:
      deployments:
        namespaces: ['staging']
        restart: ['dash-ops*sre']
        scale: ['dash-ops*sre']

  - name: 'Production (Read-Only)'
    kubeconfig: ${HOME}/.kube/config
    context: 'prod'
    # No permission block = read-only access
```

## ☸️ Cluster Operations

### **Cluster Overview**

- **Node Status** - Available nodes and resource capacity
- **Namespace List** - All namespaces with resource usage
- **Cluster Info** - Kubernetes version, API server status
- **Resource Quotas** - Limits and current usage per namespace

### **Workload Management**

#### **Deployments**

- **List Deployments** - All deployments across allowed namespaces
- **Deployment Details** - Replica count, image, resource limits, age, conditions
- **Restart Operations** - Rolling restart of deployment pods
- **Scale Operations** - Increase/decrease replica count with validation
- **Rollout Status** - Current deployment state and history

#### **Pods**

- **Pod Listing** - All pods with status, resource usage, and container counts
- **Pod Details** - Container status, controlled by, QoS class, age, node placement
- **Log Streaming** - Real-time log viewing with search and copy functionality
- **Enhanced Information** - Restart counts, quality of service, resource allocation

#### **Services & Networking** (Planned)

- **Service Discovery** - List services and endpoints
- **Ingress Management** - Route configuration and status
- **Network Policies** - Security rule visualization

## 📊 Monitoring & Observability

### **Real-time Metrics**

- **Cluster Health** - Node and API server status
- **Resource Usage** - CPU, memory, storage per namespace
- **Pod Status** - Running, pending, failed pod counts
- **Event Monitoring** - Kubernetes events and warnings

### **Log Management**

- **Pod Log Streaming** - Live log tailing
- **Multi-container Support** - Select specific containers
- **Log Filtering** - Search and filter log content
- **Download Logs** - Export logs for analysis (planned)

## 🔐 Security & Permissions

### **RBAC Requirements**

#### **Minimum Read-Only Permissions**

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: dashops-reader
rules:
  - apiGroups: ['']
    resources: ['pods', 'nodes', 'namespaces', 'services']
    verbs: ['get', 'list', 'watch']
  - apiGroups: ['apps']
    resources: ['deployments', 'replicasets']
    verbs: ['get', 'list', 'watch']
  - apiGroups: ['']
    resources: ['pods/log']
    verbs: ['get']
```

#### **Deployment Management Permissions**

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: dashops-operator
rules:
  - apiGroups: ['apps']
    resources: ['deployments', 'deployments/scale']
    verbs: ['get', 'list', 'patch', 'update']
  - apiGroups: ['']
    resources: ['pods']
    verbs: ['get', 'list', 'delete'] # For restart functionality
```

### **Team-based Access Control**

```yaml
kubernetes:
  - name: 'Production Cluster'
    permission:
      deployments:
        namespaces: ['api', 'worker'] # Limit namespace access
        restart: ['dash-ops*sre'] # Only SRE can restart deployments
        scale: ['dash-ops*sre', 'dash-ops*ops'] # SRE and Ops can scale deployments
```

## 🚨 Alpha Limitations

### **Current Restrictions**

❌ **Not Production Ready**

- **Limited operations** - Only basic deployment scaling
- **No resource quotas** - Missing namespace limits
- **Basic error handling** - Limited failure recovery
- **No backup/restore** - No data protection features
- **Missing monitoring** - No alerting or metrics collection

### **Security Limitations**

- **Basic RBAC** - Simple team-based permissions only
- **No audit trail** - Limited operation logging
- **Credential exposure** - Kubeconfig in configuration files
- **No network policies** - Missing security controls
- **No admission control** - No policy enforcement

## 📊 API Endpoints

### **Cluster Operations**

```
GET /api/kubernetes/clusters
```

**Response:**

```json
{
  "data": [
    {
      "name": "Development Cluster",
      "context": "dev",
      "status": "healthy",
      "version": "v1.28.2"
    }
  ]
}
```

### **Workload Operations**

```
GET /api/kubernetes/deployments?cluster={name}&namespace={namespace}
```

**Response:**

```json
{
  "data": [
    {
      "name": "api-server",
      "namespace": "default",
      "replicas": 3,
      "readyReplicas": 3,
      "image": "api-server:v1.2.3",
      "status": "running"
    }
  ]
}
```

#### **Deployment Actions**

```
POST /api/kubernetes/{context}/deployment/restart/{namespace}/{name}
```

**Description:** Performs a rolling restart of the deployment

**Response:**

```json
{
  "message": "Deployment restart initiated"
}
```

```
POST /api/kubernetes/{context}/deployment/scale/{namespace}/{name}/{replicas}
```

**Description:** Scales the deployment to the specified number of replicas

**Response:**

```json
{
  "message": "Deployment scaled successfully"
}
```

### **Pod Operations**

```
GET /api/kubernetes/pods?cluster={name}&namespace={namespace}
GET /api/kubernetes/pods/{name}/logs?cluster={name}&namespace={namespace}
```

## 🧪 Testing Guidelines

### **Safe Testing Practices**

> **⚠️ Critical**: Only test on development or staging clusters, never production.

1. **Use dedicated test clusters** - Isolated K8s environments
2. **Limited permissions** - Restrict RBAC to test namespaces only
3. **Resource limits** - Set strict quotas on test namespaces
4. **Monitoring** - Watch cluster events during testing

### **Test Configuration**

```yaml
kubernetes:
  - name: 'Development Testing'
    kubeconfig: ${HOME}/.kube/dev-config
    context: 'dev-cluster'
    permission:
      deployments:
        namespaces: ['test', 'development'] # Safe namespaces only
        restart: ['dash-ops*developers'] # Allow restart for testing
        scale: ['dash-ops*developers'] # Allow scaling for testing
```

## 🐛 Troubleshooting

### **Common Issues**

#### **Connection Failed**

- ✅ Verify kubeconfig file exists and is valid
- ✅ Check cluster context is correct: `kubectl config current-context`
- ✅ Test cluster access: `kubectl cluster-info`

#### **Permission Denied**

- ✅ Verify RBAC permissions with: `kubectl auth can-i <verb> <resource>`
- ✅ Check service account permissions (in-cluster mode)
- ✅ Validate team membership for operations

#### **Pods/Deployments Not Visible**

- ✅ Check namespace permissions in configuration
- ✅ Verify user has access to specified namespaces
- ✅ Check if resources exist: `kubectl get deployments -A`

### **Debug Mode**

Enable verbose logging for troubleshooting:

```bash
# Backend debug logs
KUBERNETES_DEBUG=true go run main.go

# Check kubeconfig
kubectl config view --minify

# Test cluster connectivity
kubectl cluster-info dump
```

## 🤝 Contributing

### **Priority Areas**

1. **🔒 Security** - Enhanced RBAC and security controls
2. **📊 Monitoring** - Metrics collection and alerting
3. **🧪 Testing** - Kubernetes integration test suite
4. **🔌 Features** - Additional workload types (StatefulSets, Jobs)
5. **📖 Documentation** - Setup guides and best practices

### **Development Setup**

```bash
# 1. Set up test cluster (minikube/kind)
kind create cluster --name dashops-test

# 2. Configure kubeconfig
export KUBECONFIG=${HOME}/.kube/config

# 3. Run backend with debug
KUBERNETES_DEBUG=true go run main.go

# 4. Test API endpoints
curl http://localhost:8080/api/kubernetes/clusters
```

## 📚 Resources

- **[Kubernetes API Reference](https://kubernetes.io/docs/reference/kubernetes-api/)**
- **[kubectl Command Reference](https://kubernetes.io/docs/reference/kubectl/)**
- **[Kubernetes RBAC Documentation](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)**
- **[client-go Library](https://github.com/kubernetes/client-go)**

---

**⚠️ Alpha Notice**: This plugin is in early development. Use only for testing and evaluation in non-production environments.
