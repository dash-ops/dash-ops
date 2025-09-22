package models

import (
	"fmt"
	"strings"
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
	Name string      `yaml:"name" json:"name"`
	Tier ServiceTier `yaml:"tier" json:"tier"`

	// Audit fields (auto-populated)
	CreatedAt time.Time `yaml:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedBy string    `yaml:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedAt time.Time `yaml:"updated_at,omitempty" json:"updated_at,omitempty"`
	UpdatedBy string    `yaml:"updated_by,omitempty" json:"updated_by,omitempty"`
	Version   int       `yaml:"version,omitempty" json:"version,omitempty"`
}

// ServiceSpec contains the actual service definition
type ServiceSpec struct {
	Description string `yaml:"description" json:"description"`

	// Team ownership
	Team ServiceTeam `yaml:"team" json:"team"`

	// Business context
	Business ServiceBusiness `yaml:"business" json:"business"`

	// Technology stack
	Technology ServiceTechnology `yaml:"technology,omitempty" json:"technology,omitempty"`

	// Kubernetes integration
	Kubernetes *ServiceKubernetes `yaml:"kubernetes,omitempty" json:"kubernetes,omitempty"`

	// Observability
	Observability ServiceObservability `yaml:"observability,omitempty" json:"observability,omitempty"`

	// Documentation
	Runbooks []ServiceRunbook `yaml:"runbooks,omitempty" json:"runbooks,omitempty"`
}

// ServiceTier represents service business tier
type ServiceTier string

const (
	TierCritical  ServiceTier = "TIER-1"
	TierImportant ServiceTier = "TIER-2"
	TierStandard  ServiceTier = "TIER-3"
)

// ServiceStatus represents service operational status
type ServiceStatus string

const (
	StatusHealthy  ServiceStatus = "healthy"
	StatusDegraded ServiceStatus = "degraded"
	StatusCritical ServiceStatus = "critical"
	StatusDown     ServiceStatus = "down"
	StatusUnknown  ServiceStatus = "unknown"
)

// ServiceTeam defines team ownership
type ServiceTeam struct {
	GitHubTeam string   `yaml:"github_team" json:"github_team"`
	Members    []string `json:"members,omitempty"`
	GitHubURL  string   `json:"github_url,omitempty"`
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

// KubernetesEnvironment defines environment-specific configuration
type KubernetesEnvironment struct {
	Name      string                         `yaml:"name" json:"name"`
	Context   string                         `yaml:"context" json:"context"`
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

// Domain methods for Service

// IsHighPriority checks if service is high priority (TIER-1 or TIER-2)
func (s *Service) IsHighPriority() bool {
	return s.Metadata.Tier == TierCritical || s.Metadata.Tier == TierImportant
}

// CanBeModifiedBy checks if service can be modified by user teams
func (s *Service) CanBeModifiedBy(userTeams []string) bool {
	for _, team := range userTeams {
		if strings.EqualFold(team, s.Spec.Team.GitHubTeam) {
			return true
		}
	}
	return false
}

// GetEnvironmentByName returns Kubernetes environment by name
func (s *Service) GetEnvironmentByName(envName string) (*KubernetesEnvironment, error) {
	if s.Spec.Kubernetes == nil {
		return nil, fmt.Errorf("service has no Kubernetes configuration")
	}

	for _, env := range s.Spec.Kubernetes.Environments {
		if strings.EqualFold(env.Name, envName) {
			return &env, nil
		}
	}

	return nil, fmt.Errorf("environment '%s' not found", envName)
}

// GetDeploymentByName returns deployment by name in a specific environment
func (s *Service) GetDeploymentByName(envName, deploymentName string) (*KubernetesDeployment, error) {
	env, err := s.GetEnvironmentByName(envName)
	if err != nil {
		return nil, err
	}

	for _, deployment := range env.Resources.Deployments {
		if strings.EqualFold(deployment.Name, deploymentName) {
			return &deployment, nil
		}
	}

	return nil, fmt.Errorf("deployment '%s' not found in environment '%s'", deploymentName, envName)
}

// HasDependency checks if service depends on another service
func (s *Service) HasDependency(serviceName string) bool {
	for _, dep := range s.Spec.Business.Dependencies {
		if strings.EqualFold(dep, serviceName) {
			return true
		}
	}
	return false
}

// Validate validates the service definition
func (s *Service) Validate() error {
	if s.Metadata.Name == "" {
		return fmt.Errorf("service name is required")
	}

	if s.Spec.Description == "" {
		return fmt.Errorf("service description is required")
	}

	if s.Spec.Team.GitHubTeam == "" {
		return fmt.Errorf("github team is required")
	}

	// Validate tier
	if !s.isValidTier() {
		return fmt.Errorf("invalid tier '%s', must be TIER-1, TIER-2, or TIER-3", s.Metadata.Tier)
	}

	// Validate Kubernetes configuration if present
	if s.Spec.Kubernetes != nil {
		if err := s.validateKubernetes(); err != nil {
			return fmt.Errorf("kubernetes validation failed: %w", err)
		}
	}

	return nil
}

// isValidTier checks if the service tier is valid
func (s *Service) isValidTier() bool {
	return s.Metadata.Tier == TierCritical ||
		s.Metadata.Tier == TierImportant ||
		s.Metadata.Tier == TierStandard
}

// validateKubernetes validates Kubernetes configuration
func (s *Service) validateKubernetes() error {
	if len(s.Spec.Kubernetes.Environments) == 0 {
		return fmt.Errorf("at least one environment is required")
	}

	for i, env := range s.Spec.Kubernetes.Environments {
		if env.Name == "" {
			return fmt.Errorf("environment[%d].name is required", i)
		}
		if env.Context == "" {
			return fmt.Errorf("environment[%d].context is required", i)
		}
		if env.Namespace == "" {
			return fmt.Errorf("environment[%d].namespace is required", i)
		}

		// Validate deployments
		for j, deploy := range env.Resources.Deployments {
			if deploy.Name == "" {
				return fmt.Errorf("environment[%d].deployments[%d].name is required", i, j)
			}
			if deploy.Replicas <= 0 {
				return fmt.Errorf("environment[%d].deployments[%d].replicas must be greater than 0", i, j)
			}
		}
	}

	return nil
}

// SetDefaults sets default values for the service
func (s *Service) SetDefaults() {
	if s.APIVersion == "" {
		s.APIVersion = "v1"
	}
	if s.Kind == "" {
		s.Kind = "Service"
	}
	if s.Metadata.Version == 0 {
		s.Metadata.Version = 1
	}
	if s.Metadata.CreatedAt.IsZero() {
		s.Metadata.CreatedAt = time.Now()
	}
	s.Metadata.UpdatedAt = time.Now()
}
