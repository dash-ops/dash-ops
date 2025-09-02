package servicecatalog

import "time"

// ServiceResponse represents service information response
type ServiceResponse struct {
	APIVersion string                  `json:"apiVersion"`
	Kind       string                  `json:"kind"`
	Metadata   ServiceMetadataResponse `json:"metadata"`
	Spec       ServiceSpecResponse     `json:"spec"`
}

// ServiceMetadataResponse represents service metadata in responses
type ServiceMetadataResponse struct {
	Name      string    `json:"name"`
	Tier      string    `json:"tier"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	Version   int       `json:"version"`
}

// ServiceSpecResponse represents service specification in responses
type ServiceSpecResponse struct {
	Description   string                 `json:"description"`
	Team          TeamResponse           `json:"team"`
	Business      BusinessResponse       `json:"business"`
	Technology    *TechnologyResponse    `json:"technology,omitempty"`
	Kubernetes    *KubernetesResponse    `json:"kubernetes,omitempty"`
	Observability *ObservabilityResponse `json:"observability,omitempty"`
	Runbooks      []RunbookResponse      `json:"runbooks,omitempty"`
}

// TeamResponse represents team information in responses
type TeamResponse struct {
	GitHubTeam string   `json:"github_team"`
	Members    []string `json:"members,omitempty"`
	GitHubURL  string   `json:"github_url,omitempty"`
}

// BusinessResponse represents business context in responses
type BusinessResponse struct {
	SLATarget    string   `json:"sla_target,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
	Impact       string   `json:"impact,omitempty"`
}

// TechnologyResponse represents technology stack in responses
type TechnologyResponse struct {
	Language  string `json:"language,omitempty"`
	Framework string `json:"framework,omitempty"`
}

// KubernetesResponse represents Kubernetes configuration in responses
type KubernetesResponse struct {
	Environments []KubernetesEnvironmentResponse `json:"environments"`
}

// KubernetesEnvironmentResponse represents environment configuration in responses
type KubernetesEnvironmentResponse struct {
	Name      string                                 `json:"name"`
	Context   string                                 `json:"context"`
	Namespace string                                 `json:"namespace"`
	Resources KubernetesEnvironmentResourcesResponse `json:"resources"`
}

// KubernetesEnvironmentResourcesResponse represents resources in responses
type KubernetesEnvironmentResourcesResponse struct {
	Deployments []KubernetesDeploymentResponse `json:"deployments"`
	Services    []string                       `json:"services,omitempty"`
	ConfigMaps  []string                       `json:"configmaps,omitempty"`
}

// KubernetesDeploymentResponse represents deployment configuration in responses
type KubernetesDeploymentResponse struct {
	Name      string                              `json:"name"`
	Replicas  int                                 `json:"replicas"`
	Resources *KubernetesResourceRequestsResponse `json:"resources,omitempty"`
}

// KubernetesResourceRequestsResponse represents resource requests in responses
type KubernetesResourceRequestsResponse struct {
	Requests KubernetesResourceSpecResponse `json:"requests"`
	Limits   KubernetesResourceSpecResponse `json:"limits"`
}

// KubernetesResourceSpecResponse represents resource specifications in responses
type KubernetesResourceSpecResponse struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// ObservabilityResponse represents observability configuration in responses
type ObservabilityResponse struct {
	Metrics string `json:"metrics,omitempty"`
	Logs    string `json:"logs,omitempty"`
	Traces  string `json:"traces,omitempty"`
}

// RunbookResponse represents runbook information in responses
type RunbookResponse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ServiceListResponse represents service list response
type ServiceListResponse struct {
	Services []ServiceResponse `json:"services"`
	Total    int               `json:"total"`
	Filters  interface{}       `json:"filters,omitempty"` // ServiceFilter from models
}

// ServiceHealthResponse represents service health response
type ServiceHealthResponse struct {
	ServiceName   string                      `json:"service_name"`
	OverallStatus string                      `json:"overall_status"`
	Environments  []EnvironmentHealthResponse `json:"environments"`
	LastUpdated   time.Time                   `json:"last_updated"`
}

// EnvironmentHealthResponse represents environment health in responses
type EnvironmentHealthResponse struct {
	Name        string                     `json:"name"`
	Context     string                     `json:"context"`
	Status      string                     `json:"status"`
	Deployments []DeploymentHealthResponse `json:"deployments"`
}

// DeploymentHealthResponse represents deployment health in responses
type DeploymentHealthResponse struct {
	Name            string    `json:"name"`
	ReadyReplicas   int       `json:"ready_replicas"`
	DesiredReplicas int       `json:"desired_replicas"`
	Status          string    `json:"status"`
	LastUpdated     time.Time `json:"last_updated"`
}

// ServiceHistoryResponse represents service history response
type ServiceHistoryResponse struct {
	ServiceName string                  `json:"service_name"`
	History     []ServiceChangeResponse `json:"history"`
}

// ServiceChangeResponse represents service change in responses
type ServiceChangeResponse struct {
	Commit    string                       `json:"commit"`
	Author    string                       `json:"author"`
	Email     string                       `json:"email"`
	Timestamp time.Time                    `json:"timestamp"`
	Message   string                       `json:"message"`
	Changes   []ServiceFieldChangeResponse `json:"changes,omitempty"`
}

// ServiceFieldChangeResponse represents field change in responses
type ServiceFieldChangeResponse struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"old_value"`
	NewValue interface{} `json:"new_value"`
}

// ServiceContextResponse represents service context in responses
type ServiceContextResponse struct {
	Service     *ServiceResponse `json:"service"`
	Environment string           `json:"environment"`
	Namespace   string           `json:"namespace"`
	Context     string           `json:"context"`
	Found       bool             `json:"found"`
}

// SystemStatusResponse represents system status response
type SystemStatusResponse struct {
	ServiceCount       int       `json:"service_count"`
	RepositoryStatus   string    `json:"repository_status"`
	StorageProvider    string    `json:"storage_provider"`
	VersioningEnabled  bool      `json:"versioning_enabled"`
	VersioningProvider string    `json:"versioning_provider"`
	LastUpdated        time.Time `json:"last_updated"`
}

// BatchHealthResponse represents batch health check response
type BatchHealthResponse struct {
	Services    []ServiceHealthResponse `json:"services"`
	Total       int                     `json:"total"`
	LastUpdated time.Time               `json:"last_updated"`
}
