#!/bin/bash

# Deploy script for User Authentication API demo
# This script builds and deploys the example API to local Kubernetes (docker-desktop)

set -e

echo "🚀 Deploying User Authentication API to Kubernetes..."

# Check if docker-desktop context is available
if ! kubectl config get-contexts | grep -q docker-desktop; then
    echo "❌ docker-desktop Kubernetes context not found"
    echo "Please ensure Docker Desktop is running with Kubernetes enabled"
    exit 1
fi

# Switch to docker-desktop context
echo "📋 Switching to docker-desktop context..."
kubectl config use-context docker-desktop

# Build Docker image
echo "🐳 Building Docker image..."
cd "$(dirname "$0")"
docker build -t user-authentication-api:latest .

# Apply Kubernetes manifests
echo "☸️  Applying Kubernetes manifests..."
kubectl apply -f k8s-manifests.yaml

# Wait for deployments to be ready
echo "⏳ Waiting for deployments to be ready..."
kubectl -n auth rollout status deployment/auth-api --timeout=60s
kubectl -n auth rollout status deployment/auth-worker --timeout=60s

echo ""
echo "✅ Deployment completed successfully!"
echo ""
echo "📊 Check deployment status:"
echo "   kubectl -n auth get pods,svc,deploy"
echo ""
echo "🌐 API endpoints:"
echo "   Health:  http://localhost:30080/health"
echo "   Info:    http://localhost:30080/info"
echo "   Status:  http://localhost:30080/api/status"
echo ""
echo "🔍 View logs:"
echo "   kubectl -n auth logs -l component=auth-api --tail=50 -f"
echo ""
echo "🧹 To clean up:"
echo "   kubectl delete namespace auth"
echo ""
