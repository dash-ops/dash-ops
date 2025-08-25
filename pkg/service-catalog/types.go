package servicecatalog

import (
	"time"
)

// Service represents a complete service definition
type Service struct {
	APIVersion string          `yaml:"apiVersion" json:"apiVersion"`
	Kind       string          `yaml:"kind" json:"kind"`
	Metadata   ServiceMetadata `yaml:"metadata" json:"metadata"`
	Spec       ServiceSpec     `yaml:"spec" json:"spec"`
}

// ServiceMetadata contains service identification and audit information
type ServiceMetadata struct {
	Name string `yaml:"name" json:"name"`
	Tier string `yaml:"tier" json:"tier"` // TIER-1, TIER-2, TIER-3

	// Audit fields (auto-populated)
	CreatedAt string `yaml:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedBy string `yaml:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedAt string `yaml:"updated_at,omitempty" json:"updated_at,omitempty"`
	UpdatedBy string `yaml:"updated_by,omitempty" json:"updated_by,omitempty"`
	Version   int    `yaml:"version,omitempty" json:"version,omitempty"`
}

// ServiceSpec contains the actual service definition
type ServiceSpec struct {
	Description string `yaml:"description" json:"description"`

	// GitHub team integration
	Team ServiceTeam `yaml:"team" json:"team"`

	// Business context
	Business ServiceBusiness `yaml:"business" json:"business"`

	// Technology stack
	Technology ServiceTechnology `yaml:"technology,omitempty" json:"technology,omitempty"`

	// Kubernetes integration
	Kubernetes *ServiceKubernetes `yaml:"kubernetes,omitempty" json:"kubernetes,omitempty"`

	// External links
	Observability ServiceObservability `yaml:"observability,omitempty" json:"observability,omitempty"`

	// Documentation
	Runbooks []ServiceRunbook `yaml:"runbooks,omitempty" json:"runbooks,omitempty"`
}

// ServiceTeam defines team ownership
type ServiceTeam struct {
	GitHubTeam string `yaml:"github_team" json:"github_team"`
	// Auto-resolved at runtime:
	// Members   []string `json:"members,omitempty"`
	// GitHubURL string   `json:"github_url,omitempty"`
}

// ServiceBusiness contains business context
type ServiceBusiness struct {
	SLATarget    string   `yaml:"sla_target,omitempty" json:"sla_target,omitempty"`
	Dependencies []string `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	Impact       string   `yaml:"impact,omitempty" json:"impact,omitempty"` // high, medium, low
}

// ServiceTechnology contains technology stack information
type ServiceTechnology struct {
	Language  string `yaml:"language,omitempty" json:"language,omitempty"`
	Framework string `yaml:"framework,omitempty" json:"framework,omitempty"`
}

// ServiceKubernetes contains Kubernetes integration configuration
type ServiceKubernetes struct {
	Environments []KubernetesEnvironment `yaml:"environments" json:"environments"`
}

// KubernetesEnvironment defines environment-specific configuration
type KubernetesEnvironment struct {
	Name      string                         `yaml:"name" json:"name"`       // staging, production
	Context   string                         `yaml:"context" json:"context"` // k8s context
	Namespace string                         `yaml:"namespace" json:"namespace"`
	Resources KubernetesEnvironmentResources `yaml:"resources" json:"resources"`
}

// KubernetesEnvironmentResources defines Kubernetes resources per environment
type KubernetesEnvironmentResources struct {
	Deployments []KubernetesDeployment `yaml:"deployments" json:"deployments"`
	Services    []string               `yaml:"services,omitempty" json:"services,omitempty"`
	ConfigMaps  []string               `yaml:"configmaps,omitempty" json:"configmaps,omitempty"`
}

// KubernetesDeployment defines deployment specifications
type KubernetesDeployment struct {
	Name      string                     `yaml:"name" json:"name"`
	Replicas  int                        `yaml:"replicas" json:"replicas"`
	Resources KubernetesResourceRequests `yaml:"resources" json:"resources"`
}

// KubernetesResourceRequests defines CPU and memory specifications
type KubernetesResourceRequests struct {
	Requests KubernetesResourceSpec `yaml:"requests" json:"requests"`
	Limits   KubernetesResourceSpec `yaml:"limits" json:"limits"`
}

// KubernetesResourceSpec defines specific resource values
type KubernetesResourceSpec struct {
	CPU    string `yaml:"cpu" json:"cpu"`
	Memory string `yaml:"memory" json:"memory"`
}

// ServiceObservability contains external monitoring links
type ServiceObservability struct {
	Metrics string `yaml:"metrics,omitempty" json:"metrics,omitempty"`
	Logs    string `yaml:"logs,omitempty" json:"logs,omitempty"`
	Traces  string `yaml:"traces,omitempty" json:"traces,omitempty"`
}

// ServiceRunbook contains documentation links
type ServiceRunbook struct {
	Name string `yaml:"name" json:"name"`
	URL  string `yaml:"url" json:"url"`
}

// ServiceList represents a list of services for API responses
type ServiceList struct {
	Services []Service `json:"services"`
	Total    int       `json:"total"`
}

// ServiceHealth represents aggregated service health status
type ServiceHealth struct {
	ServiceName   string              `json:"service_name"`
	OverallStatus string              `json:"overall_status"` // healthy, degraded, down, critical
	Environments  []EnvironmentHealth `json:"environments"`
	LastUpdated   time.Time           `json:"last_updated"`
}

// EnvironmentHealth represents health status for a specific environment
type EnvironmentHealth struct {
	Name        string             `json:"name"`    // staging, production
	Context     string             `json:"context"` // k8s context
	Status      string             `json:"status"`  // healthy, degraded, down
	Deployments []DeploymentHealth `json:"deployments"`
}

// DeploymentHealth represents health status for a specific deployment
type DeploymentHealth struct {
	Name            string    `json:"name"`
	ReadyReplicas   int       `json:"ready_replicas"`
	DesiredReplicas int       `json:"desired_replicas"`
	Status          string    `json:"status"`
	LastUpdated     time.Time `json:"last_updated"`
}

// ServiceHistory represents service change history
type ServiceHistory struct {
	ServiceName string          `json:"service_name"`
	History     []ServiceChange `json:"history"`
}

// ServiceChange represents a single change in service history
type ServiceChange struct {
	Commit    string               `json:"commit"`
	Author    string               `json:"author"`
	Email     string               `json:"email"`
	Timestamp time.Time            `json:"timestamp"`
	Message   string               `json:"message"`
	Changes   []ServiceFieldChange `json:"changes,omitempty"`
}

// ServiceFieldChange represents a specific field change
type ServiceFieldChange struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"old_value"`
	NewValue interface{} `json:"new_value"`
}

// StorageProvider defines the interface for service storage backends
type StorageProvider interface {
	CreateService(service *Service) error
	GetService(name string) (*Service, error)
	UpdateService(service *Service) error
	DeleteService(name string) error
	ListServices() ([]Service, error)
	ServiceExists(name string) bool
}

// Config represents service catalog configuration
type Config struct {
	Storage StorageConfig `yaml:"storage" json:"storage"`
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	Provider   string                  `yaml:"provider" json:"provider"` // filesystem, github, s3
	Filesystem FilesystemStorageConfig `yaml:"filesystem,omitempty" json:"filesystem,omitempty"`
	GitHub     GitHubStorageConfig     `yaml:"github,omitempty" json:"github,omitempty"`
	S3         S3StorageConfig         `yaml:"s3,omitempty" json:"s3,omitempty"`
}

// FilesystemStorageConfig represents local filesystem storage configuration
type FilesystemStorageConfig struct {
	Directory string `yaml:"directory" json:"directory"`
}

// GitHubStorageConfig represents GitHub repository storage configuration
type GitHubStorageConfig struct {
	Repository string `yaml:"repository" json:"repository"`
	Branch     string `yaml:"branch" json:"branch"`
}

// S3StorageConfig represents S3 bucket storage configuration
type S3StorageConfig struct {
	Bucket string `yaml:"bucket" json:"bucket"`
}

// UserContext represents user information from OAuth2
type UserContext struct {
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Teams    []string `json:"teams,omitempty"`
}
