package http

import (
	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog/models"
	scWire "github.com/dash-ops/dash-ops/pkg/service-catalog/wire"
)

// ServiceAdapter handles transformation between models and wire formats
type ServiceAdapter struct{}

// NewServiceAdapter creates a new service adapter
func NewServiceAdapter() *ServiceAdapter {
	return &ServiceAdapter{}
}

// RequestToModel converts CreateServiceRequest to Service model
func (sa *ServiceAdapter) RequestToModel(req scWire.CreateServiceRequest) (*scModels.Service, error) {
	service := &scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: req.Name,
			Tier: scModels.ServiceTier(req.Tier),
		},
		Spec: scModels.ServiceSpec{
			Description: req.Description,
			Team: scModels.ServiceTeam{
				GitHubTeam: req.Team.GitHubTeam,
			},
		},
	}

	// Convert business context
	if req.Business != nil {
		service.Spec.Business = scModels.ServiceBusiness{
			SLATarget:    req.Business.SLATarget,
			Dependencies: req.Business.Dependencies,
			Impact:       req.Business.Impact,
		}
	}

	// Convert technology stack
	if req.Technology != nil {
		service.Spec.Technology = scModels.ServiceTechnology{
			Language:  req.Technology.Language,
			Framework: req.Technology.Framework,
		}
	}

	// Convert Kubernetes configuration
	if req.Kubernetes != nil {
		service.Spec.Kubernetes = sa.convertKubernetesRequest(req.Kubernetes)
	}

	// Convert observability
	if req.Observability != nil {
		service.Spec.Observability = scModels.ServiceObservability{
			Metrics: req.Observability.Metrics,
			Logs:    req.Observability.Logs,
			Traces:  req.Observability.Traces,
		}
	}

	// Convert runbooks
	if len(req.Runbooks) > 0 {
		service.Spec.Runbooks = sa.convertRunbooksRequest(req.Runbooks)
	}

	return service, nil
}

// UpdateRequestToModel converts UpdateServiceRequest to Service model
func (sa *ServiceAdapter) UpdateRequestToModel(req scWire.UpdateServiceRequest, existingService *scModels.Service) (*scModels.Service, error) {
	// Start with existing service
	service := *existingService

	// Update fields that are provided
	if req.Description != nil {
		service.Spec.Description = *req.Description
	}

	if req.Tier != nil {
		service.Metadata.Tier = scModels.ServiceTier(*req.Tier)
	}

	if req.Team != nil {
		service.Spec.Team.GitHubTeam = req.Team.GitHubTeam
	}

	if req.Business != nil {
		service.Spec.Business = scModels.ServiceBusiness{
			SLATarget:    req.Business.SLATarget,
			Dependencies: req.Business.Dependencies,
			Impact:       req.Business.Impact,
		}
	}

	if req.Technology != nil {
		service.Spec.Technology = scModels.ServiceTechnology{
			Language:  req.Technology.Language,
			Framework: req.Technology.Framework,
		}
	}

	if req.Kubernetes != nil {
		service.Spec.Kubernetes = sa.convertKubernetesRequest(req.Kubernetes)
	}

	if req.Observability != nil {
		service.Spec.Observability = scModels.ServiceObservability{
			Metrics: req.Observability.Metrics,
			Logs:    req.Observability.Logs,
			Traces:  req.Observability.Traces,
		}
	}

	if req.Runbooks != nil {
		service.Spec.Runbooks = sa.convertRunbooksRequest(*req.Runbooks)
	}

	return &service, nil
}

// ModelToResponse converts Service model to ServiceResponse
func (sa *ServiceAdapter) ModelToResponse(service *scModels.Service) scWire.ServiceResponse {
	return scWire.ServiceResponse{
		APIVersion: service.APIVersion,
		Kind:       service.Kind,
		Metadata: scWire.ServiceMetadataResponse{
			Name:      service.Metadata.Name,
			Tier:      string(service.Metadata.Tier),
			CreatedAt: service.Metadata.CreatedAt,
			CreatedBy: service.Metadata.CreatedBy,
			UpdatedAt: service.Metadata.UpdatedAt,
			UpdatedBy: service.Metadata.UpdatedBy,
			Version:   service.Metadata.Version,
		},
		Spec: scWire.ServiceSpecResponse{
			Description: service.Spec.Description,
			Team: scWire.TeamResponse{
				GitHubTeam: service.Spec.Team.GitHubTeam,
				Members:    service.Spec.Team.Members,
				GitHubURL:  service.Spec.Team.GitHubURL,
			},
			Business: scWire.BusinessResponse{
				SLATarget:    service.Spec.Business.SLATarget,
				Dependencies: service.Spec.Business.Dependencies,
				Impact:       service.Spec.Business.Impact,
			},
			Technology:    sa.convertTechnologyResponse(&service.Spec.Technology),
			Kubernetes:    sa.convertKubernetesResponse(service.Spec.Kubernetes),
			Observability: sa.convertObservabilityResponse(&service.Spec.Observability),
			Runbooks:      sa.convertRunbooksResponse(service.Spec.Runbooks),
		},
	}
}

// ModelListToResponse converts ServiceList model to ServiceListResponse
func (sa *ServiceAdapter) ModelListToResponse(serviceList *scModels.ServiceList) scWire.ServiceListResponse {
	var services []scWire.ServiceResponse
	for _, service := range serviceList.Services {
		services = append(services, sa.ModelToResponse(&service))
	}

	return scWire.ServiceListResponse{
		Services: services,
		Total:    serviceList.Total,
		Filters:  serviceList.Filters,
	}
}

// HealthModelToResponse converts ServiceHealth model to ServiceHealthResponse
func (sa *ServiceAdapter) HealthModelToResponse(health *scModels.ServiceHealth) scWire.ServiceHealthResponse {
	var environments []scWire.EnvironmentHealthResponse
	for _, env := range health.Environments {
		var deployments []scWire.DeploymentHealthResponse
		for _, dep := range env.Deployments {
			deployments = append(deployments, scWire.DeploymentHealthResponse{
				Name:            dep.Name,
				ReadyReplicas:   dep.ReadyReplicas,
				DesiredReplicas: dep.DesiredReplicas,
				Status:          string(dep.Status),
				LastUpdated:     dep.LastUpdated,
			})
		}

		environments = append(environments, scWire.EnvironmentHealthResponse{
			Name:        env.Name,
			Context:     env.Context,
			Status:      string(env.Status),
			Deployments: deployments,
		})
	}

	return scWire.ServiceHealthResponse{
		ServiceName:   health.ServiceName,
		OverallStatus: string(health.OverallStatus),
		Environments:  environments,
		LastUpdated:   health.LastUpdated,
	}
}

// HistoryModelToResponse converts ServiceHistory model to ServiceHistoryResponse
func (sa *ServiceAdapter) HistoryModelToResponse(history *scModels.ServiceHistory) scWire.ServiceHistoryResponse {
	var changes []scWire.ServiceChangeResponse
	for _, change := range history.History {
		var fieldChanges []scWire.ServiceFieldChangeResponse
		for _, fieldChange := range change.Changes {
			fieldChanges = append(fieldChanges, scWire.ServiceFieldChangeResponse{
				Field:    fieldChange.Field,
				OldValue: fieldChange.OldValue,
				NewValue: fieldChange.NewValue,
			})
		}

		changes = append(changes, scWire.ServiceChangeResponse{
			Commit:    change.Commit,
			Author:    change.Author,
			Email:     change.Email,
			Timestamp: change.Timestamp,
			Message:   change.Message,
			Changes:   fieldChanges,
		})
	}

	return scWire.ServiceHistoryResponse{
		ServiceName: history.ServiceName,
		History:     changes,
	}
}

// convertKubernetesRequest converts KubernetesRequest to ServiceKubernetes
func (sa *ServiceAdapter) convertKubernetesRequest(req *scWire.KubernetesRequest) *scModels.ServiceKubernetes {
	if req == nil {
		return nil
	}

	var environments []scModels.KubernetesEnvironment
	for _, envReq := range req.Environments {
		var deployments []scModels.KubernetesDeployment
		for _, depReq := range envReq.Resources.Deployments {
			deployment := scModels.KubernetesDeployment{
				Name:     depReq.Name,
				Replicas: depReq.Replicas,
			}

			if depReq.Resources != nil {
				deployment.Resources = scModels.KubernetesResourceRequests{
					Requests: scModels.KubernetesResourceSpec{
						CPU:    depReq.Resources.Requests.CPU,
						Memory: depReq.Resources.Requests.Memory,
					},
					Limits: scModels.KubernetesResourceSpec{
						CPU:    depReq.Resources.Limits.CPU,
						Memory: depReq.Resources.Limits.Memory,
					},
				}
			}

			deployments = append(deployments, deployment)
		}

		environments = append(environments, scModels.KubernetesEnvironment{
			Name:      envReq.Name,
			Context:   envReq.Context,
			Namespace: envReq.Namespace,
			Resources: scModels.KubernetesEnvironmentResources{
				Deployments: deployments,
				Services:    envReq.Resources.Services,
				ConfigMaps:  envReq.Resources.ConfigMaps,
			},
		})
	}

	return &scModels.ServiceKubernetes{
		Environments: environments,
	}
}

// convertRunbooksRequest converts runbook requests to runbook models
func (sa *ServiceAdapter) convertRunbooksRequest(req []scWire.RunbookRequest) []scModels.ServiceRunbook {
	var runbooks []scModels.ServiceRunbook
	for _, rbReq := range req {
		runbooks = append(runbooks, scModels.ServiceRunbook{
			Name: rbReq.Name,
			URL:  rbReq.URL,
		})
	}
	return runbooks
}

// convertTechnologyResponse converts technology model to response
func (sa *ServiceAdapter) convertTechnologyResponse(tech *scModels.ServiceTechnology) *scWire.TechnologyResponse {
	if tech.Language == "" && tech.Framework == "" {
		return nil
	}
	return &scWire.TechnologyResponse{
		Language:  tech.Language,
		Framework: tech.Framework,
	}
}

// convertKubernetesResponse converts Kubernetes model to response
func (sa *ServiceAdapter) convertKubernetesResponse(k8s *scModels.ServiceKubernetes) *scWire.KubernetesResponse {
	if k8s == nil {
		return nil
	}

	var environments []scWire.KubernetesEnvironmentResponse
	for _, env := range k8s.Environments {
		var deployments []scWire.KubernetesDeploymentResponse
		for _, dep := range env.Resources.Deployments {
			deployment := scWire.KubernetesDeploymentResponse{
				Name:     dep.Name,
				Replicas: dep.Replicas,
			}

			if dep.Resources.Requests.CPU != "" || dep.Resources.Requests.Memory != "" {
				deployment.Resources = &scWire.KubernetesResourceRequestsResponse{
					Requests: scWire.KubernetesResourceSpecResponse{
						CPU:    dep.Resources.Requests.CPU,
						Memory: dep.Resources.Requests.Memory,
					},
					Limits: scWire.KubernetesResourceSpecResponse{
						CPU:    dep.Resources.Limits.CPU,
						Memory: dep.Resources.Limits.Memory,
					},
				}
			}

			deployments = append(deployments, deployment)
		}

		environments = append(environments, scWire.KubernetesEnvironmentResponse{
			Name:      env.Name,
			Context:   env.Context,
			Namespace: env.Namespace,
			Resources: scWire.KubernetesEnvironmentResourcesResponse{
				Deployments: deployments,
				Services:    env.Resources.Services,
				ConfigMaps:  env.Resources.ConfigMaps,
			},
		})
	}

	return &scWire.KubernetesResponse{
		Environments: environments,
	}
}

// convertObservabilityResponse converts observability model to response
func (sa *ServiceAdapter) convertObservabilityResponse(obs *scModels.ServiceObservability) *scWire.ObservabilityResponse {
	if obs.Metrics == "" && obs.Logs == "" && obs.Traces == "" {
		return nil
	}
	return &scWire.ObservabilityResponse{
		Metrics: obs.Metrics,
		Logs:    obs.Logs,
		Traces:  obs.Traces,
	}
}

// convertRunbooksResponse converts runbooks model to response
func (sa *ServiceAdapter) convertRunbooksResponse(runbooks []scModels.ServiceRunbook) []scWire.RunbookResponse {
	if len(runbooks) == 0 {
		return nil
	}

	var response []scWire.RunbookResponse
	for _, rb := range runbooks {
		response = append(response, scWire.RunbookResponse{
			Name: rb.Name,
			URL:  rb.URL,
		})
	}
	return response
}
