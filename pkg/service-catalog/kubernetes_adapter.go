package servicecatalog

// KubernetesServiceContextAdapter adapts ServiceContextResolver to work with kubernetes plugin
// This adapter converts between service-catalog types and kubernetes plugin types
type KubernetesServiceContextAdapter struct {
	resolver *ServiceContextResolver
}

// KubernetesServiceContext represents service context in kubernetes plugin format
// This mirrors the ServiceContext struct in kubernetes/deployment.go to avoid import cycles
type KubernetesServiceContext struct {
	ServiceName string `json:"service_name,omitempty"`
	ServiceTier string `json:"service_tier,omitempty"`
	Environment string `json:"environment,omitempty"`
	Context     string `json:"context,omitempty"`
	Team        string `json:"team,omitempty"`
	Description string `json:"description,omitempty"`
}

// NewKubernetesServiceContextAdapter creates a new adapter
func NewKubernetesServiceContextAdapter(resolver *ServiceContextResolver) *KubernetesServiceContextAdapter {
	return &KubernetesServiceContextAdapter{
		resolver: resolver,
	}
}

// ResolveDeploymentService implements the kubernetes ServiceContextResolver interface
func (k *KubernetesServiceContextAdapter) ResolveDeploymentService(deploymentName, namespace, context string) (*KubernetesServiceContext, error) {
	// Use the underlying service catalog resolver
	serviceContext, err := k.resolver.ResolveDeploymentService(deploymentName, namespace, context)
	if err != nil {
		return nil, err
	}

	if serviceContext == nil {
		return nil, nil
	}

	// Convert from service-catalog ServiceContext to kubernetes ServiceContext
	return &KubernetesServiceContext{
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
