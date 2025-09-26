package http

import (
	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	k8sWire "github.com/dash-ops/dash-ops/pkg/kubernetes/wire"
)

// DeploymentListToResponse converts DeploymentList model to DeploymentListResponse
func DeploymentListToResponse(deploymentList *k8sModels.DeploymentList) k8sWire.DeploymentListResponse {
	var deployments []k8sWire.DeploymentResponse
	for _, deployment := range deploymentList.Deployments {
		deployments = append(deployments, *DeploymentToResponse(&deployment))
	}

	return k8sWire.DeploymentListResponse{
		Deployments: deployments,
		Total:       deploymentList.Total,
		Namespace:   deploymentList.Namespace,
		Filter:      deploymentList.Filter,
	}
}

// DeploymentToResponse converts Deployment model to DeploymentResponse
func DeploymentToResponse(deployment *k8sModels.Deployment) *k8sWire.DeploymentResponse {
	var conditions []k8sWire.DeploymentConditionResponse
	for _, condition := range deployment.Conditions {
		conditions = append(conditions, k8sWire.DeploymentConditionResponse{
			Type:           condition.Type,
			Status:         condition.Status,
			Reason:         condition.Reason,
			Message:        condition.Message,
			LastUpdateTime: condition.LastUpdateTime,
		})
	}

	var serviceContext *k8sWire.ServiceContextResponse
	if deployment.ServiceContext != nil {
		serviceContext = &k8sWire.ServiceContextResponse{
			ServiceName: deployment.ServiceContext.ServiceName,
			ServiceTier: deployment.ServiceContext.ServiceTier,
			Environment: deployment.ServiceContext.Environment,
			Context:     deployment.ServiceContext.Context,
			Team:        deployment.ServiceContext.Team,
			Description: deployment.ServiceContext.Description,
			Found:       deployment.ServiceContext.Found,
		}
	}

	return &k8sWire.DeploymentResponse{
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
		PodInfo: k8sWire.PodInfoResponse{
			Running: deployment.PodInfo.Running,
			Pending: deployment.PodInfo.Pending,
			Failed:  deployment.PodInfo.Failed,
			Total:   deployment.PodInfo.Total,
		},
		Replicas: k8sWire.DeploymentReplicasResponse{
			Desired:   deployment.Replicas.Desired,
			Current:   deployment.Replicas.Current,
			Ready:     deployment.Replicas.Ready,
			Available: deployment.Replicas.Available,
		},
		Age:                 deployment.Age,
		CreatedAt:           deployment.CreatedAt,
		Conditions:          conditions,
		ServiceContext:      serviceContext,
		AvailabilityPercent: deployment.GetAvailabilityPercentage(),
	}
}

// DeploymentsToResponse converts Deployment slice to DeploymentResponse slice
func DeploymentsToResponse(deployments []k8sModels.Deployment) []k8sWire.DeploymentResponse {
	var response []k8sWire.DeploymentResponse
	for _, deployment := range deployments {
		response = append(response, *DeploymentToResponse(&deployment))
	}
	return response
}
