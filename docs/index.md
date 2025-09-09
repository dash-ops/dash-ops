---
layout: default
title: Home
---

# DashOPS Documentation

> **⚠️ Beta Software** - Experimental cloud operations platform under active development.

Welcome to the **DashOPS Documentation**! This comprehensive guide covers everything you need to know about deploying, configuring, and contributing to DashOPS.

## 🚀 Quick Navigation

<div class="grid">
  <div class="card">
    <h3>🏁 Getting Started</h3>
    <p>Set up DashOPS locally in 5 minutes</p>
    <a href="./getting-started.html" class="btn">Start Here →</a>
  </div>
  
  <div class="card">
    <h3>🔌 Plugin Guides</h3>
    <p>Configure AWS, Kubernetes, and Auth plugins</p>
    <a href="./plugins/" class="btn">View Plugins →</a>
  </div>
  
  <div class="card">
    <h3>📖 API Reference</h3>
    <p>Complete API documentation and examples</p>
    <a href="./api-reference.html" class="btn">API Docs →</a>
  </div>
  
  <div class="card">
    <h3>🤝 Contributing</h3>
    <p>Help improve DashOPS development</p>
    <a href="https://github.com/dash-ops/dash-ops" class="btn">GitHub →</a>
  </div>
</div>

## 🎯 What is DashOPS?

DashOPS is a **cloud operations platform** designed to simplify infrastructure management for development teams. It provides:

- **🔗 Unified Interface** - Single dashboard for AWS, Kubernetes, and more
- **🔐 Secure Access** - Auth authentication with team-based permissions
- **🧩 Plugin Architecture** - Modular design for easy extensibility
- **👥 Developer-Focused** - Reduce cognitive load, increase productivity

## ⚠️ Beta Status

**DashOPS is experimental software:**

- ❌ **Not production ready** - Security and stability limitations
- 🧪 **Testing only** - Use in development environments
- 🔄 **Active development** - Features and APIs may change
- 🤝 **Community driven** - We welcome your feedback and contributions

## 🏁 Get Started in 3 Steps

### **1. Quick Setup**

```bash
git clone https://github.com/dash-ops/dash-ops.git
cd dash-ops && cp local.sample.yaml dash-ops.yaml
```

### **2. Configure GitHub OAuth**

```bash
# Set up GitHub OAuth App and environment variables
export GITHUB_CLIENT_ID="your-client-id"
export GITHUB_CLIENT_SECRET="your-client-secret"
```

### **3. Run DashOPS**

```bash
go run main.go &          # Backend
cd front && yarn dev      # Frontend
```

**→ Access at [http://localhost:5173](http://localhost:5173)**

## 📚 Documentation Structure

- **[📋 Main Guide](./README.html)** - Complete documentation overview
- **[🏁 Getting Started](./getting-started.html)** - Step-by-step setup guide
- **[📖 API Reference](./api-reference.html)** - Complete API documentation
- **[🔌 Plugin Guides](./plugins/)** - Individual plugin configuration
- **[💻 Frontend Guide](../front/README.md)** - React/TypeScript development

## 🔗 Useful Links

- **[GitHub Repository](https://github.com/dash-ops/dash-ops)** - Source code and issues
- **[Docker Images](https://hub.docker.com/r/dashops/dash-ops)** - Container images
- **[Helm Charts](../helm-charts/)** - Kubernetes deployment

---

**Ready to explore DashOPS?** Start with our [Getting Started Guide](./getting-started.html)!
