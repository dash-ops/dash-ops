package servicecatalog

// CreateServiceRequest represents service creation request
type CreateServiceRequest struct {
	Name          string                `json:"name" validate:"required,min=1,max=100"`
	Description   string                `json:"description" validate:"required"`
	Tier          string                `json:"tier" validate:"required,oneof=TIER-1 TIER-2 TIER-3"`
	Team          TeamRequest           `json:"team" validate:"required"`
	Business      *BusinessRequest      `json:"business,omitempty"`
	Technology    *TechnologyRequest    `json:"technology,omitempty"`
	Kubernetes    *KubernetesRequest    `json:"kubernetes,omitempty"`
	Observability *ObservabilityRequest `json:"observability,omitempty"`
	Runbooks      []RunbookRequest      `json:"runbooks,omitempty"`
}

// UpdateServiceRequest represents service update request
type UpdateServiceRequest struct {
	Description   *string               `json:"description,omitempty"`
	Tier          *string               `json:"tier,omitempty" validate:"omitempty,oneof=TIER-1 TIER-2 TIER-3"`
	Team          *TeamRequest          `json:"team,omitempty"`
	Business      *BusinessRequest      `json:"business,omitempty"`
	Technology    *TechnologyRequest    `json:"technology,omitempty"`
	Kubernetes    *KubernetesRequest    `json:"kubernetes,omitempty"`
	Observability *ObservabilityRequest `json:"observability,omitempty"`
	Runbooks      *[]RunbookRequest     `json:"runbooks,omitempty"`
}

// TeamRequest represents team information in requests
type TeamRequest struct {
	GitHubTeam string `json:"github_team" validate:"required"`
}

// BusinessRequest represents business context in requests
type BusinessRequest struct {
	SLATarget    string   `json:"sla_target,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	Impact       string   `json:"impact,omitempty" validate:"omitempty,oneof=high medium low"`
}

// TechnologyRequest represents technology stack in requests
type TechnologyRequest struct {
	Language  string `json:"language,omitempty"`
	Framework string `json:"framework,omitempty"`
}

// KubernetesRequest represents Kubernetes configuration in requests
type KubernetesRequest struct {
	Environments []KubernetesEnvironmentRequest `json:"environments" validate:"required,min=1"`
}

// KubernetesEnvironmentRequest represents environment configuration in requests
type KubernetesEnvironmentRequest struct {
	Name      string                                `json:"name" validate:"required"`
	Context   string                                `json:"context" validate:"required"`
	Namespace string                                `json:"namespace" validate:"required"`
	Resources KubernetesEnvironmentResourcesRequest `json:"resources" validate:"required"`
}

// KubernetesEnvironmentResourcesRequest represents resources in requests
type KubernetesEnvironmentResourcesRequest struct {
	Deployments []KubernetesDeploymentRequest `json:"deployments" validate:"required,min=1"`
	Services    []string                      `json:"services,omitempty"`
	ConfigMaps  []string                      `json:"configmaps,omitempty"`
}

// KubernetesDeploymentRequest represents deployment configuration in requests
type KubernetesDeploymentRequest struct {
	Name      string                             `json:"name" validate:"required"`
	Replicas  int                                `json:"replicas" validate:"required,min=1"`
	Resources *KubernetesResourceRequestsRequest `json:"resources,omitempty"`
}

// KubernetesResourceRequestsRequest represents resource requests in requests
type KubernetesResourceRequestsRequest struct {
	Requests KubernetesResourceSpecRequest `json:"requests" validate:"required"`
	Limits   KubernetesResourceSpecRequest `json:"limits" validate:"required"`
}

// KubernetesResourceSpecRequest represents resource specifications in requests
type KubernetesResourceSpecRequest struct {
	CPU    string `json:"cpu" validate:"required"`
	Memory string `json:"memory" validate:"required"`
}

// ObservabilityRequest represents observability configuration in requests
type ObservabilityRequest struct {
	Metrics string `json:"metrics,omitempty"`
	Logs    string `json:"logs,omitempty"`
	Traces  string `json:"traces,omitempty"`
}

// RunbookRequest represents runbook information in requests
type RunbookRequest struct {
	Name string `json:"name" validate:"required"`
	URL  string `json:"url" validate:"required,url"`
}

// ServiceSearchRequest represents service search request
type ServiceSearchRequest struct {
	Query  string `json:"query,omitempty"`
	Team   string `json:"team,omitempty"`
	Tier   string `json:"tier,omitempty" validate:"omitempty,oneof=TIER-1 TIER-2 TIER-3"`
	Status string `json:"status,omitempty" validate:"omitempty,oneof=healthy degraded critical down unknown"`
	Limit  int    `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Offset int    `json:"offset,omitempty" validate:"omitempty,min=0"`
}

// BatchHealthRequest represents batch health check request
type BatchHealthRequest struct {
	ServiceNames []string `json:"service_names" validate:"required,min=1,max=20"`
}
