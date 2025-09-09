# Getting Started with DashOPS

> **⚠️ Beta Software** - This guide covers experimental software not recommended for production use.

Welcome to **DashOPS**! This comprehensive guide will help you get DashOPS running locally and understand its core concepts.

## 🎯 What is DashOPS?

DashOPS is a **cloud operations platform** designed to:

- **Reduce cognitive load** for developers working with cloud infrastructure
- **Simplify common operations** across AWS, Kubernetes, and other platforms
- **Provide unified interface** for multi-cloud environments
- **Enable self-service** infrastructure operations with proper governance

### **Who Should Use DashOPS?**

- **Development Teams** - Need occasional infrastructure access without deep expertise
- **SRE Teams** - Want to provide controlled self-service capabilities
- **DevOps Engineers** - Looking for unified tooling across cloud platforms
- **Organizations** - Seeking to reduce infrastructure management overhead

---

## 🚀 Prerequisites

### **System Requirements**

| Component   | Version | Purpose                             |
| ----------- | ------- | ----------------------------------- |
| **Go**      | 1.21+   | Backend API development             |
| **Node.js** | 18.0+   | Frontend development                |
| **Yarn**    | 1.22+   | Package management (preferred)      |
| **Docker**  | 20.0+   | Containerized deployment (optional) |

### **Cloud Access Requirements**

#### **For GitHub OAuth (Required)**

- GitHub organization with OAuth App creation permissions
- Administrative access to create OAuth applications

#### **For AWS Plugin (Optional)**

- AWS account with programmatic access
- IAM permissions for EC2 operations
- AWS CLI configured (recommended)

#### **For Kubernetes Plugin (Optional)**

- Kubernetes cluster access
- Valid kubeconfig file
- Cluster admin permissions (for RBAC setup)

---

## 📦 Installation

### **Option 1: Local Development (Recommended)**

This method is best for development, testing, and contribution.

#### **Step 1: Clone Repository**

```bash
git clone https://github.com/dash-ops/dash-ops.git
cd dash-ops
```

#### **Step 2: Backend Setup**

```bash
# Download Go dependencies
go mod download

# Verify Go installation
go version  # Should be 1.21+

# Optional: Install Air for hot reload
go install github.com/cosmtrek/air@latest
```

#### **Step 3: Frontend Setup**

```bash
cd front

# Install dependencies
yarn

# Verify installation
yarn --version  # Should be 1.22+
node --version  # Should be 18.0+
```

#### **Step 4: Configuration**

```bash
# Create configuration file
cp local.sample.yaml dash-ops.yaml

# Edit configuration with your settings
nano dash-ops.yaml  # or your preferred editor
```

### **Option 2: Docker Setup**

This method is good for quick testing without local Go/Node setup.

```bash
# 1. Clone and configure
git clone https://github.com/dash-ops/dash-ops.git
cd dash-ops
cp local.sample.yaml dash-ops.yaml

# 2. Edit configuration
nano dash-ops.yaml

# 3. Run with Docker
docker run --rm \
  -v $(pwd)/dash-ops.yaml:/dash-ops.yaml \
  -v ${HOME}/.kube/config:/.kube/config \
  -p 8080:8080 \
  dashops/dash-ops
```

---

## ⚙️ Configuration

### **Basic Configuration**

Create your `dash-ops.yaml` configuration file:

```yaml
# Server settings
port: 8080
origin: http://localhost:8080
headers:
  - 'Content-Type'
  - 'Authorization'
front: app

# Enable plugins (start with Auth only)
plugins:
  - 'Auth'

# GitHub Auth setup (required)
auth:
  - provider: github
    clientId: ${GITHUB_CLIENT_ID}
    clientSecret: ${GITHUB_CLIENT_SECRET}
    authURL: 'https://github.com/login/oauth/authorize'
    tokenURL: 'https://github.com/login/oauth/access_token'
    redirectURL: 'http://localhost:8080/api/oauth/redirect'
    urlLoginSuccess: 'http://localhost:5173'
    orgPermission: 'your-github-org' # Replace with your organization
    scopes: [user, repo, read:org]
```

### **Environment Variables**

Create a `.env` file or export variables:

```bash
# Required for Auth
export GITHUB_CLIENT_ID="your-github-oauth-client-id"
export GITHUB_CLIENT_SECRET="your-github-oauth-client-secret"

# Optional for AWS plugin
export AWS_ACCESS_KEY_ID="your-aws-access-key"
export AWS_SECRET_ACCESS_KEY="your-aws-secret-key"
```

### **GitHub OAuth App Setup**

1. **Go to GitHub Organization Settings**

   - Navigate to `Settings` → `Developer settings` → `OAuth Apps`

2. **Create New OAuth App**

   - **Application name**: `DashOPS`
   - **Homepage URL**: `http://localhost:5173` (local) or your domain
   - **Authorization callback URL**: `http://localhost:8080/api/oauth/redirect`

3. **Copy Credentials**
   - Save the `Client ID` and `Client Secret`
   - Add them to your environment variables

---

## 🎮 First Run

### **Starting the Application**

#### **Local Development**

```bash
# Terminal 1: Backend
go run main.go
# ✅ API server starts on http://localhost:8080

# Terminal 2: Frontend
cd front
yarn dev
# ✅ Frontend starts on http://localhost:5173
```

#### **With Docker**

```bash
docker run --rm \
  -v $(pwd)/dash-ops.yaml:/dash-ops.yaml \
  --env-file .env \
  -p 8080:8080 \
  dashops/dash-ops
# ✅ Full application on http://localhost:8080
```

### **Accessing DashOPS**

1. **Open browser** to `http://localhost:5173` (local) or `http://localhost:8080` (Docker)
2. **Click "Login"** - Redirects to GitHub OAuth
3. **Authorize Application** - Grant permissions to DashOPS
4. **Access Dashboard** - You'll see the main DashOPS interface

### **Expected Interface**

After successful login, you should see:

- **📊 Dashboard** - Main overview page
- **👤 User Profile** - Your GitHub profile and permissions
- **🔌 Enabled Plugins** - Available modules based on your configuration

---

## 🔌 Adding More Plugins

Once you have the basic setup working, you can enable additional plugins:

### **Adding AWS Plugin**

```yaml
# 1. Add to plugins list
plugins:
  - 'Auth'
  - 'AWS' # Add this line

# 2. Add AWS configuration
aws:
  - name: 'Development Account'
    region: us-east-1
    accessKeyId: ${AWS_ACCESS_KEY_ID}
    secretAccessKey: ${AWS_SECRET_ACCESS_KEY}
```

### **Adding Kubernetes Plugin**

```yaml
# 1. Add to plugins list
plugins:
  - 'Auth'
  - 'Kubernetes' # Add this line

# 2. Add Kubernetes configuration
kubernetes:
  - name: 'Development Cluster'
    kubeconfig: ${HOME}/.kube/config
    context: 'dev-cluster-context'
```

### **Restart and Verify**

After adding plugins:

```bash
# Restart backend
# Ctrl+C to stop, then:
go run main.go

# Frontend automatically reloads
# Check for new menu items in the sidebar
```

---

## 🐛 Troubleshooting

### **Common Startup Issues**

#### **Backend Won't Start**

```bash
# Check Go version
go version  # Must be 1.21+

# Verify configuration file
go run main.go
# Look for YAML parsing errors
```

#### **Frontend Build Errors**

```bash
cd front

# Check Node.js version
node --version  # Must be 18.0+

# Clear cache and reinstall
rm -rf node_modules yarn.lock
yarn
```

#### **Auth Login Fails**

- ✅ Check GitHub OAuth App configuration
- ✅ Verify callback URL matches exactly
- ✅ Ensure organization membership
- ✅ Check browser console for errors

#### **Plugin Not Loading**

- ✅ Verify plugin name in `plugins` list
- ✅ Check plugin-specific configuration
- ✅ Review backend logs for errors
- ✅ Validate environment variables

### **Debug Mode**

Enable detailed logging for troubleshooting:

```bash
# Backend debug logs
LOG_LEVEL=debug go run main.go

# Frontend debug mode
# In browser console:
localStorage.setItem('dash-ops:debug', 'true')
```

### **Health Checks**

Verify system health:

```bash
# Backend health
curl http://localhost:8080/api/health

# Plugin status
curl http://localhost:8080/api/config/plugins

# Frontend accessibility
curl http://localhost:5173
```

---

## 🎓 Next Steps

Once you have DashOPS running:

### **1. Explore the Interface**

- Navigate through the sidebar menu
- Try the Dashboard overview
- Explore enabled plugin features

### **2. Configure Additional Plugins**

- Set up AWS plugin for EC2 management
- Configure Kubernetes plugin for cluster monitoring
- Customize permissions for your team

### **3. Development Setup** (If Contributing)

- Set up development tools (Air for Go hot reload)
- Configure VS Code with recommended extensions
- Run test suites to ensure everything works

### **4. Learn More**

- Read individual [plugin documentation](./plugins/)
- Explore the [API reference](./README.md#api-reference)
- Check out [plugin development guide](./plugins/README.md#plugin-development)

---

## 📞 Getting Help

### **Community Support**

- **[GitHub Issues](https://github.com/dash-ops/dash-ops/issues)** - Bug reports and feature requests
- **[GitHub Discussions](https://github.com/dash-ops/dash-ops/discussions)** - Questions and community help
- **[Documentation](./README.md)** - Complete reference documentation

### **Before Asking for Help**

1. Check this getting started guide
2. Review relevant plugin documentation
3. Search existing GitHub issues
4. Enable debug logging and check logs
5. Try with a minimal configuration

### **When Reporting Issues**

Include this information in your issue report:

- **DashOPS version** - Git commit hash or release version
- **Operating system** - OS and version
- **Configuration** - Your `dash-ops.yaml` (remove sensitive data)
- **Error logs** - Backend and frontend error messages
- **Steps to reproduce** - Clear reproduction steps

---

**🚧 Remember**: DashOPS is **beta software**. Use only for testing and development. We appreciate your feedback and contributions!

For more detailed information, see the [main documentation](./README.md).
