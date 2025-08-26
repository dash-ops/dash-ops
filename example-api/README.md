# User Authentication Service - Demo API

Demo service for testing the Service Catalog integration with Kubernetes.

## Quick Deploy

Deploy to docker-desktop cluster:

```bash
kubectl apply -f k8s-manifests.yaml
```

## Verify Deployment

Check if pods are running:

```bash
kubectl -n auth get pods,deployments,services
```

## API Endpoints (via NodePort)

- **Health**: http://localhost:30080
- **API Status**: http://localhost:30080/api/status

## Clean Up

```bash
kubectl delete -f k8s-manifests.yaml
```

## Service Catalog Integration

This demo service is registered in the Service Catalog as:

- **Name**: `user-authentication`
- **Tier**: TIER-1
- **Team**: auth-squad
- **Environment**: local (docker-desktop context)
- **Deployments**: auth-api, auth-worker
