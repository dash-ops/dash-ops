package servicecatalog

import (
	"github.com/dash-ops/dash-ops/pkg/kubernetes"
)

// KubernetesServiceContextAdapter adapts ServiceContextResolver to work with kubernetes plugin
// This adapter converts between service-catalog types and kubernetes plugin types
type KubernetesServiceContextAdapter struct {
	resolver *ServiceContextResolver
}

// NewKubernetesServiceContextAdapter creates a new adapter
func NewKubernetesServiceContextAdapter(resolver *ServiceContextResolver) *KubernetesServiceContextAdapter {
	return &KubernetesServiceContextAdapter{
		resolver: resolver,
	}
}

// ResolveDeploymentService implements the kubernetes ServiceContextResolver interface
func (k *KubernetesServiceContextAdapter) ResolveDeploymentService(deploymentName, namespace, context string) (*kubernetes.ServiceContext, error) {
	// Use the underlying service catalog resolver
	serviceContext, err := k.resolver.ResolveDeploymentService(deploymentName, namespace, context)
	if err != nil {
		return nil, err
	}

	if serviceContext == nil {
		return nil, nil
	}

	// Convert from service-catalog ServiceContext to kubernetes ServiceContext
	return &kubernetes.ServiceContext{
		ServiceName: serviceContext.ServiceName,
		ServiceTier: serviceContext.ServiceTier,
		Environment: serviceContext.Environment,
		Context:     serviceContext.Context,
		Team:        serviceContext.Team,
		Description: serviceContext.Description,
	}, nil
}

// GetKubernetesAdapter returns an adapter that can be used by the kubernetes plugin
func (sc *ServiceCatalog) GetKubernetesAdapter() *KubernetesServiceContextAdapter {
	return NewKubernetesServiceContextAdapter(sc.contextResolver)
}
