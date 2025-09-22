package adapters

import (
	"context"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	scControllers "github.com/dash-ops/dash-ops/pkg/service-catalog/controllers"
)

// KubernetesServiceContextAdapter adapts service catalog controller to kubernetes ServiceContextResolver
type KubernetesServiceContextAdapter struct {
	controller *scControllers.ServiceController
}

// NewKubernetesServiceContextAdapter creates a new adapter
func NewKubernetesServiceContextAdapter(controller *scControllers.ServiceController) *KubernetesServiceContextAdapter {
	return &KubernetesServiceContextAdapter{
		controller: controller,
	}
}

// ResolveDeploymentService implements kubernetes.ServiceContextResolver interface
func (k *KubernetesServiceContextAdapter) ResolveDeploymentService(deploymentName, namespace, contextName string) (*k8sModels.ServiceContext, error) {
	// Use the controller to resolve the service context
	serviceContext, err := k.controller.ResolveDeploymentService(context.TODO(), deploymentName, namespace, contextName)
	if err != nil {
		// Return nil if no service found (not an error)
		return nil, nil
	}

	if serviceContext == nil {
		return nil, nil
	}

	// Convert from service-catalog ServiceContext to kubernetes ServiceContext
	return &k8sModels.ServiceContext{
		ServiceName: serviceContext.Service.Metadata.Name,
		ServiceTier: string(serviceContext.Service.Metadata.Tier),
		Environment: serviceContext.Environment,
		Context:     serviceContext.Context,
		Team:        serviceContext.Service.Spec.Team.GitHubTeam,
		Description: serviceContext.Service.Spec.Description,
		Found:       serviceContext.Found,
	}, nil
}
