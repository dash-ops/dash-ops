# API Reference

> **⚠️ Beta API** - Endpoints are experimental and may change without notice.

DashOPS provides a RESTful API built with Go, offering standardized endpoints for all cloud operations and authentication management.

## 🏗️ API Overview

### **Base Configuration**

- **Base URL**: `http://localhost:8080/api` (development)
- **Authentication**: Bearer token (Auth)
- **Content Type**: `application/json`
- **CORS**: Configurable origins

### **Standard Response Format**

All API endpoints return consistent JSON responses:

```json
{
  "data": [...],           # Response payload (varies by endpoint)
  "success": true,         # Operation success status
  "message": "Success",    # Optional descriptive message
  "errors": []             # Array of error details (if any)
}
```

### **Error Response Format**

```json
{
  "data": null,
  "success": false,
  "message": "Operation failed",
  "errors": [
    {
      "code": "INVALID_REQUEST",
      "message": "Missing required parameter: instanceId",
      "field": "instanceId"
    }
  ]
}
```

---

## 🔐 Authentication

### **Authentication Flow**

DashOPS uses Auth with GitHub for authentication:

```
1. GET /api/oauth/authorize → Redirect to GitHub
2. GitHub authorization → User approves access  
3. GET /api/oauth/redirect → GitHub callback with code
4. Backend exchanges code for token → Returns JWT to frontend
5. Subsequent requests → Include Bearer token in Authorization header
```

### **Authentication Endpoints**

#### **Initiate OAuth Flow**
```http
GET /api/oauth/authorize
```
Redirects to GitHub OAuth authorization URL.

#### **OAuth Callback** (Internal)
```http
GET /api/oauth/redirect?code={authorization_code}&state={state}
```
Handles GitHub OAuth callback and exchanges code for access token.

#### **Logout**
```http
POST /api/oauth/logout
Authorization: Bearer {token}
```
Invalidates the user session and clears authentication.

#### **User Information**
```http
GET /api/user
Authorization: Bearer {token}
```

**Response:**
```json
{
  "data": {
    "id": 12345,
    "name": "John Doe",
    "email": "john.doe@example.com",
    "avatar_url": "https://avatars.githubusercontent.com/u/12345",
    "organization": "dash-ops",
    "teams": ["developers", "sre"]
  }
}
```

---

## 🛠️ Core Endpoints

### **Health & Status**

#### **Application Health**
```http
GET /api/health
```

**Response:**
```json
{
  "data": {
    "status": "healthy",
    "version": "v0.1.0-beta",
    "uptime": "2h34m12s",
    "plugins": ["Auth", "AWS", "Kubernetes"]
  }
}
```

#### **Version Information**
```http
GET /api/version
```

**Response:**
```json
{
  "data": {
    "version": "v0.1.0-beta",
    "commit": "abc123def456",
    "buildDate": "2024-01-15T10:30:00Z",
    "goVersion": "go1.21.5"
  }
}
```

### **Configuration**

#### **Plugin Configuration**
```http
GET /api/config/plugins
Authorization: Bearer {token}
```

**Response:**
```json
{
  "data": [
    {
      "name": "Auth",
      "enabled": true,
      "version": "v0.1.0"
    },
    {
      "name": "AWS", 
      "enabled": true,
      "version": "v0.1.0"
    }
  ]
}
```

---

## ☁️ AWS Plugin API

### **Account Management**

#### **List AWS Accounts**
```http
GET /api/aws/accounts
Authorization: Bearer {token}
```

**Response:**
```json
{
  "data": [
    {
      "name": "Production Account",
      "region": "us-east-1",
      "accountId": "123456789012"
    },
    {
      "name": "Development Account", 
      "region": "us-west-2",
      "accountId": "123456789013"
    }
  ]
}
```

### **EC2 Operations**

#### **List EC2 Instances**
```http
GET /api/aws/instances?account={account_name}
Authorization: Bearer {token}
```

**Query Parameters:**
- `account` (required) - AWS account name from configuration
- `state` (optional) - Filter by instance state: `running`, `stopped`, `pending`

**Response:**
```json
{
  "data": [
    {
      "instanceId": "i-1234567890abcdef0",
      "name": "web-server-1",
      "state": "running",
      "instanceType": "t3.micro",
      "launchTime": "2024-01-15T10:30:00Z",
      "privateIpAddress": "10.0.1.100",
      "publicIpAddress": "54.123.45.67",
      "tags": {
        "Name": "web-server-1",
        "Environment": "development"
      }
    }
  ]
}
```

#### **Start EC2 Instance**
```http
POST /api/aws/instances/{instanceId}/start
Authorization: Bearer {token}
Content-Type: application/json

{
  "account": "Development Account"
}
```

**Response:**
```json
{
  "data": {
    "current_state": "pending",
    "previous_state": "stopped"
  },
  "success": true,
  "message": "Instance start initiated"
}
```

#### **Stop EC2 Instance**
```http
POST /api/aws/instances/{instanceId}/stop
Authorization: Bearer {token}
Content-Type: application/json

{
  "account": "Development Account"
}
```

**Response:**
```json
{
  "data": {
    "current_state": "stopping",
    "previous_state": "running"
  },
  "success": true,
  "message": "Instance stop initiated"
}
```

---

## ☸️ Kubernetes Plugin API

### **Cluster Management**

#### **List Clusters**
```http
GET /api/kubernetes/clusters
Authorization: Bearer {token}
```

**Response:**
```json
{
  "data": [
    {
      "name": "Development Cluster",
      "context": "dev",
      "status": "healthy",
      "version": "v1.28.2",
      "nodeCount": 3
    }
  ]
}
```

### **Workload Operations**

#### **List Deployments**
```http
GET /api/kubernetes/deployments?cluster={cluster_name}&namespace={namespace}
Authorization: Bearer {token}
```

**Query Parameters:**
- `cluster` (required) - Cluster name from configuration
- `namespace` (optional) - Kubernetes namespace, defaults to `default`

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
      "status": "running",
      "created": "2024-01-15T10:30:00Z"
    }
  ]
}
```

#### **Scale Deployment**
```http
POST /api/kubernetes/deployments/{deployment_name}/scale
Authorization: Bearer {token}
Content-Type: application/json

{
  "replicas": 5,
  "namespace": "default", 
  "cluster": "Development Cluster"
}
```

**Response:**
```json
{
  "data": {
    "success": true,
    "currentReplicas": 5,
    "previousReplicas": 3
  },
  "message": "Deployment scaled successfully"
}
```

### **Pod Operations**

#### **List Pods**
```http
GET /api/kubernetes/pods?cluster={cluster_name}&namespace={namespace}
Authorization: Bearer {token}
```

**Response:**
```json
{
  "data": [
    {
      "name": "api-server-7d4b8c8f9d-x2k9m",
      "namespace": "default",
      "status": "Running",
      "restarts": 0,
      "age": "2d3h",
      "node": "worker-node-1",
      "containers": [
        {
          "name": "api-server",
          "image": "api-server:v1.2.3",
          "ready": true
        }
      ]
    }
  ]
}
```

#### **Stream Pod Logs**
```http
GET /api/kubernetes/pods/{pod_name}/logs?cluster={cluster}&namespace={namespace}&follow=true
Authorization: Bearer {token}
```

**Query Parameters:**
- `cluster` (required) - Cluster name
- `namespace` (required) - Pod namespace
- `follow` (optional) - Stream live logs (`true`/`false`)
- `tail` (optional) - Number of recent lines
- `container` (optional) - Specific container name

**Response:** Server-Sent Events (SSE) stream
```
data: 2024-01-15T10:30:00Z INFO Starting application
data: 2024-01-15T10:30:01Z INFO Server listening on port 8080
data: 2024-01-15T10:30:02Z DEBUG Processing request
```

---

## 🔒 Security Headers

### **Required Authentication**

All protected endpoints require authentication:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### **CORS Configuration**

Configure allowed origins in `dash-ops.yaml`:

```yaml
origin: http://localhost:5173
headers:
  - 'Content-Type'
  - 'Authorization'
  - 'X-Requested-With'
```

### **Rate Limiting** (Planned)

Future API rate limiting:
- **Authentication endpoints**: 10 requests/minute
- **Read operations**: 100 requests/minute  
- **Write operations**: 20 requests/minute

---

## 📊 Error Handling

### **HTTP Status Codes**

| Code | Meaning | Usage |
|------|---------|-------|
| `200` | OK | Successful operation |
| `201` | Created | Resource created successfully |
| `400` | Bad Request | Invalid request parameters |
| `401` | Unauthorized | Authentication required or failed |
| `403` | Forbidden | Insufficient permissions |
| `404` | Not Found | Resource does not exist |
| `500` | Internal Server Error | Server-side error |

### **Error Response Examples**

#### **Authentication Error (401)**
```json
{
  "data": null,
  "success": false,
  "message": "Authentication required",
  "errors": [
    {
      "code": "UNAUTHORIZED",
      "message": "Valid Bearer token required"
    }
  ]
}
```

#### **Permission Error (403)**
```json
{
  "data": null,
  "success": false,
  "message": "Insufficient permissions",
  "errors": [
    {
      "code": "FORBIDDEN", 
      "message": "User does not have permission to stop EC2 instances",
      "resource": "ec2:instances:stop"
    }
  ]
}
```

#### **Validation Error (400)**
```json
{
  "data": null,
  "success": false,
  "message": "Validation failed",
  "errors": [
    {
      "code": "INVALID_PARAMETER",
      "message": "Instance ID is required",
      "field": "instanceId"
    }
  ]
}
```

---

## 🧪 Testing the API

### **Using cURL**

```bash
# 1. Get authentication token (through browser OAuth flow)
TOKEN="your-jwt-token-here"

# 2. Test health endpoint
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/health

# 3. List AWS instances
curl -H "Authorization: Bearer $TOKEN" \
     "http://localhost:8080/api/aws/instances?account=Development%20Account"

# 4. List Kubernetes deployments
curl -H "Authorization: Bearer $TOKEN" \
     "http://localhost:8080/api/kubernetes/deployments?cluster=Development%20Cluster"
```

### **Using Postman**

1. **Set up environment variables**:
   - `baseUrl`: `http://localhost:8080/api`
   - `token`: Your JWT token from browser localStorage

2. **Configure authentication**:
   - Authorization Type: `Bearer Token`
   - Token: `{{token}}`

3. **Test endpoints** using the documented requests above

---

## 🚨 Beta API Limitations

### **Current Restrictions**

❌ **Not Production Ready**

- **No API versioning** - Breaking changes possible
- **Limited error details** - Basic error information only
- **No rate limiting** - API calls not throttled
- **Missing pagination** - Large datasets not paginated
- **No API documentation** - No OpenAPI/Swagger specs
- **Basic validation** - Limited input validation

### **Security Limitations**

- **Simple authentication** - Basic JWT implementation
- **No API key management** - Only Auth tokens supported
- **Missing audit logs** - Limited request logging
- **No request signing** - No additional security layers
- **CORS restrictions** - Basic origin validation only

## 🛣️ API Roadmap

### **Short-term (Q1 2025)**
- **API versioning** - `/api/v1/` URL structure
- **Enhanced validation** - Request/response schema validation
- **Better error messages** - Detailed error context
- **Basic rate limiting** - Per-user request limits

### **Medium-term (Q2 2025)**
- **OpenAPI documentation** - Auto-generated API docs
- **Pagination support** - Handle large datasets efficiently
- **API key management** - Alternative authentication methods
- **Comprehensive audit logs** - Track all API operations

### **Long-term (Q3+ 2025)**
- **GraphQL support** - Flexible query interface
- **Webhook system** - Event-driven integrations
- **API analytics** - Usage metrics and monitoring
- **SDK generation** - Client libraries for popular languages

---

## 🤝 Contributing to the API

### **Backend Development**

```bash
# 1. Set up development environment
go mod download
air  # For hot reload during development

# 2. Add new endpoint
# Edit pkg/{plugin}/handler.go
func (h *Handler) HandleNewEndpoint(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

# 3. Test endpoint
go test ./pkg/{plugin}/...

# 4. Update documentation
# Add endpoint documentation to this file
```

### **API Testing**

```bash
# Run all backend tests
go test ./...

# Test specific plugin
go test ./pkg/aws/...

# Integration tests
go test -tags=integration ./...
```

## 📚 Additional Resources

- **[Go HTTP Package](https://pkg.go.dev/net/http)** - Go HTTP server documentation
- **[REST API Design](https://restfulapi.net/)** - REST API best practices
- **[OAuth2 Specification](https://tools.ietf.org/html/rfc6749)** - OAuth2 standard
- **[JWT Specification](https://tools.ietf.org/html/rfc7519)** - JSON Web Token standard

---

**⚠️ Beta Notice**: This API is experimental and intended for development use only. Endpoints and response formats may change without notice.
