# User Authentication API - Demo for Service Catalog

Esta é uma API de demonstração criada para testar o **Service Catalog** do dash-ops com recursos reais do Kubernetes.

## 🎯 **Objetivo**

- Fornecer uma API real rodando no Kubernetes local
- Testar integração do Service Catalog com recursos K8s reais
- Demonstrar health checks, deployments e services funcionais
- Servir como exemplo para development e demos

## 🚀 **Quick Start**

### **Pré-requisitos:**

- Docker Desktop with Kubernetes enabled
- kubectl configured for docker-desktop context
- Go 1.21+ (para desenvolvimento)

### **Deploy Rápido:**

```bash
# Tornar o script executável e executar
chmod +x deploy.sh
./deploy.sh
```

### **Verificar Deployment:**

```bash
# Ver status dos recursos
kubectl -n auth get pods,svc,deploy

# Testar API
curl http://localhost:30080/health
curl http://localhost:30080/info
curl http://localhost:30080/api/status
```

## 📋 **API Endpoints**

| Endpoint           | Descrição                                            |
| ------------------ | ---------------------------------------------------- |
| `GET /`            | Welcome message e informações básicas                |
| `GET /health`      | **Liveness probe** - verifica se API está rodando    |
| `GET /ready`       | **Readiness probe** - verifica se API está pronta    |
| `GET /info`        | Informações detalhadas da API (uptime, versão, host) |
| `GET /api/status`  | Status do serviço com métricas simuladas             |
| `GET /api/version` | Versão e build info                                  |

## ☸️ **Recursos Kubernetes**

### **Deployments:**

- **auth-api**: 3 replicas (API principal)
- **auth-worker**: 2 replicas (worker background)

### **Services:**

- **auth-svc**: ClusterIP service (interno)
- **auth-nodeport**: NodePort (acesso externo via porta 30080)

### **ConfigMaps:**

- **auth-config**: Configurações da aplicação

### **Namespace:**

- **auth**: Namespace dedicado para o serviço

## 🔧 **Development**

### **Build Local:**

```bash
cd example-api
go mod tidy
go run main.go
```

### **Test Endpoints:**

```bash
# Health checks
curl http://localhost:8080/health
curl http://localhost:8080/ready

# API info
curl http://localhost:8080/info
curl http://localhost:8080/api/status
```

### **Build Docker Image:**

```bash
docker build -t user-authentication-api:latest .
```

## 📊 **Service Definition**

Esta API corresponde ao seguinte service definition no Service Catalog:

```yaml
metadata:
  name: user-authentication
  tier: TIER-1
spec:
  description: 'User authentication and authorization service'
  team:
    github_team: 'auth-squad'
  business:
    sla_target: '99.9%'
    dependencies: ['user-database', 'email-service']
    impact: 'high'
  kubernetes:
    environments:
      - name: 'local'
        context: 'docker-desktop'
        namespace: 'auth'
        resources:
          deployments:
            - name: 'auth-api'
              replicas: 3
              resources:
                requests: { cpu: '100m', memory: '128Mi' }
                limits: { cpu: '500m', memory: '256Mi' }
            - name: 'auth-worker'
              replicas: 2
              resources:
                requests: { cpu: '50m', memory: '64Mi' }
                limits: { cpu: '200m', memory: '128Mi' }
          services: ['auth-svc']
          configmaps: ['auth-config']
```

## 🔍 **Monitoring & Debug**

### **View Logs:**

```bash
# API logs
kubectl -n auth logs -l component=auth-api --tail=50 -f

# Worker logs
kubectl -n auth logs -l component=auth-worker --tail=50 -f

# All logs
kubectl -n auth logs --all-containers --tail=100 -f
```

### **Pod Status:**

```bash
# Detailed pod info
kubectl -n auth describe pods

# Resource usage
kubectl -n auth top pods
```

### **Port Forward (alternative access):**

```bash
# Forward API port
kubectl -n auth port-forward svc/auth-svc 8080:80
```

## 🧹 **Cleanup**

### **Remove Everything:**

```bash
kubectl delete namespace auth
```

### **Remove Docker Image:**

```bash
docker rmi user-authentication-api:latest
```

## 🎯 **Next Steps**

1. **Deploy this API** usando `./deploy.sh`
2. **Create corresponding Service Definition** no Service Catalog
3. **Test Service Catalog integration** com recursos K8s reais
4. **Verify health aggregation** functionality
5. **Demo the complete flow** Service Definition → K8s Resources

## 📝 **Notes**

- **ImagePullPolicy: Never** - Usa imagem local (não puxa do registry)
- **NodePort 30080** - Acesso externo para desenvolvimento
- **Health checks** - Liveness (15s delay) e Readiness (5s delay)
- **Resource limits** - CPU/Memory definidos para demo
- **Labels consistentes** - Para integração com Service Catalog
