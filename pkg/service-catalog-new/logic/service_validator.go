package servicecatalog

import (
	"fmt"
	"regexp"
	"strings"

	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog-new/models"
)

// ServiceValidator provides service validation logic
type ServiceValidator struct{}

// NewServiceValidator creates a new service validator
func NewServiceValidator() *ServiceValidator {
	return &ServiceValidator{}
}

// ValidateForCreation validates a service for creation
func (sv *ServiceValidator) ValidateForCreation(service *scModels.Service) error {
	if service == nil {
		return fmt.Errorf("service cannot be nil")
	}

	// Validate basic fields
	if err := sv.validateBasicFields(service); err != nil {
		return fmt.Errorf("basic validation failed: %w", err)
	}

	// Validate service name format
	if err := sv.validateServiceName(service.Metadata.Name); err != nil {
		return fmt.Errorf("service name validation failed: %w", err)
	}

	// Validate tier
	if err := sv.validateTier(service.Metadata.Tier); err != nil {
		return fmt.Errorf("tier validation failed: %w", err)
	}

	// Validate team
	if err := sv.validateTeam(&service.Spec.Team); err != nil {
		return fmt.Errorf("team validation failed: %w", err)
	}

	// Validate Kubernetes configuration if present
	if service.Spec.Kubernetes != nil {
		if err := sv.validateKubernetes(service.Spec.Kubernetes); err != nil {
			return fmt.Errorf("kubernetes validation failed: %w", err)
		}
	}

	// Validate runbooks if present
	if err := sv.validateRunbooks(service.Spec.Runbooks); err != nil {
		return fmt.Errorf("runbooks validation failed: %w", err)
	}

	return nil
}

// ValidateForUpdate validates a service for update
func (sv *ServiceValidator) ValidateForUpdate(service *scModels.Service, existingService *scModels.Service) error {
	if service == nil {
		return fmt.Errorf("service cannot be nil")
	}

	if existingService == nil {
		return fmt.Errorf("existing service cannot be nil")
	}

	// Name cannot be changed
	if service.Metadata.Name != existingService.Metadata.Name {
		return fmt.Errorf("service name cannot be changed")
	}

	// Validate the updated service
	return sv.ValidateForCreation(service)
}

// validateBasicFields validates basic required fields
func (sv *ServiceValidator) validateBasicFields(service *scModels.Service) error {
	if service.Metadata.Name == "" {
		return fmt.Errorf("service name is required")
	}

	if service.Spec.Description == "" {
		return fmt.Errorf("service description is required")
	}

	if len(service.Spec.Description) > 500 {
		return fmt.Errorf("service description too long (max 500 characters)")
	}

	return nil
}

// validateServiceName validates service name format
func (sv *ServiceValidator) validateServiceName(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	// Check length
	if len(name) > 100 {
		return fmt.Errorf("service name too long (max 100 characters)")
	}

	if len(name) < 3 {
		return fmt.Errorf("service name too short (min 3 characters)")
	}

	// Check format (alphanumeric, hyphens, underscores)
	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_]*[a-zA-Z0-9]$`)
	if !validNameRegex.MatchString(name) {
		return fmt.Errorf("service name must contain only alphanumeric characters, hyphens, and underscores")
	}

	// Check for invalid characters in filename
	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		return fmt.Errorf("service name contains invalid characters")
	}

	return nil
}

// validateTier validates service tier
func (sv *ServiceValidator) validateTier(tier scModels.ServiceTier) error {
	validTiers := map[scModels.ServiceTier]bool{
		scModels.TierCritical:  true,
		scModels.TierImportant: true,
		scModels.TierStandard:  true,
	}

	if !validTiers[tier] {
		return fmt.Errorf("invalid tier '%s', must be TIER-1, TIER-2, or TIER-3", tier)
	}

	return nil
}

// validateTeam validates team configuration
func (sv *ServiceValidator) validateTeam(team *scModels.ServiceTeam) error {
	if team.GitHubTeam == "" {
		return fmt.Errorf("github team is required")
	}

	// Validate GitHub team name format
	validTeamRegex := regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
	if !validTeamRegex.MatchString(team.GitHubTeam) {
		return fmt.Errorf("invalid GitHub team name format")
	}

	return nil
}

// validateKubernetes validates Kubernetes configuration
func (sv *ServiceValidator) validateKubernetes(k8s *scModels.ServiceKubernetes) error {
	if len(k8s.Environments) == 0 {
		return fmt.Errorf("at least one environment is required")
	}

	envNames := make(map[string]bool)
	for i, env := range k8s.Environments {
		// Check for duplicate environment names
		if envNames[env.Name] {
			return fmt.Errorf("duplicate environment name '%s'", env.Name)
		}
		envNames[env.Name] = true

		if err := sv.validateEnvironment(&env, i); err != nil {
			return fmt.Errorf("environment[%d] validation failed: %w", i, err)
		}
	}

	return nil
}

// validateEnvironment validates a single Kubernetes environment
func (sv *ServiceValidator) validateEnvironment(env *scModels.KubernetesEnvironment, index int) error {
	if env.Name == "" {
		return fmt.Errorf("environment name is required")
	}

	if env.Context == "" {
		return fmt.Errorf("environment context is required")
	}

	if env.Namespace == "" {
		return fmt.Errorf("environment namespace is required")
	}

	// Validate deployments
	if len(env.Resources.Deployments) == 0 {
		return fmt.Errorf("at least one deployment is required")
	}

	deploymentNames := make(map[string]bool)
	for j, deployment := range env.Resources.Deployments {
		// Check for duplicate deployment names
		if deploymentNames[deployment.Name] {
			return fmt.Errorf("duplicate deployment name '%s' in environment '%s'", deployment.Name, env.Name)
		}
		deploymentNames[deployment.Name] = true

		if err := sv.validateDeployment(&deployment, j); err != nil {
			return fmt.Errorf("deployment[%d] validation failed: %w", j, err)
		}
	}

	return nil
}

// validateDeployment validates a single deployment
func (sv *ServiceValidator) validateDeployment(deployment *scModels.KubernetesDeployment, index int) error {
	if deployment.Name == "" {
		return fmt.Errorf("deployment name is required")
	}

	if deployment.Replicas <= 0 {
		return fmt.Errorf("deployment replicas must be greater than 0")
	}

	if deployment.Replicas > 100 {
		return fmt.Errorf("deployment replicas too high (max 100)")
	}

	// Validate resource specifications if present
	if deployment.Resources.Requests.CPU != "" || deployment.Resources.Requests.Memory != "" {
		if err := sv.validateResourceSpec(&deployment.Resources.Requests, "requests"); err != nil {
			return fmt.Errorf("resource requests validation failed: %w", err)
		}
	}

	if deployment.Resources.Limits.CPU != "" || deployment.Resources.Limits.Memory != "" {
		if err := sv.validateResourceSpec(&deployment.Resources.Limits, "limits"); err != nil {
			return fmt.Errorf("resource limits validation failed: %w", err)
		}
	}

	return nil
}

// validateResourceSpec validates Kubernetes resource specifications
func (sv *ServiceValidator) validateResourceSpec(spec *scModels.KubernetesResourceSpec, specType string) error {
	// Validate CPU format
	if spec.CPU != "" {
		if err := sv.validateCPUSpec(spec.CPU); err != nil {
			return fmt.Errorf("invalid CPU %s: %w", specType, err)
		}
	}

	// Validate memory format
	if spec.Memory != "" {
		if err := sv.validateMemorySpec(spec.Memory); err != nil {
			return fmt.Errorf("invalid memory %s: %w", specType, err)
		}
	}

	return nil
}

// validateCPUSpec validates CPU specification format
func (sv *ServiceValidator) validateCPUSpec(cpu string) error {
	// Accept formats like: 100m, 0.1, 1, 2.5
	cpuRegex := regexp.MustCompile(`^(\d+(\.\d+)?|\d+m)$`)
	if !cpuRegex.MatchString(cpu) {
		return fmt.Errorf("invalid CPU format '%s' (examples: 100m, 0.5, 1)", cpu)
	}
	return nil
}

// validateMemorySpec validates memory specification format
func (sv *ServiceValidator) validateMemorySpec(memory string) error {
	// Accept formats like: 128Mi, 1Gi, 512M, 1G
	memoryRegex := regexp.MustCompile(`^(\d+)(Mi|Gi|M|G|Ki|K|Ti|T)$`)
	if !memoryRegex.MatchString(memory) {
		return fmt.Errorf("invalid memory format '%s' (examples: 128Mi, 1Gi, 512M)", memory)
	}
	return nil
}

// validateRunbooks validates runbook configurations
func (sv *ServiceValidator) validateRunbooks(runbooks []scModels.ServiceRunbook) error {
	runbookNames := make(map[string]bool)

	for i, runbook := range runbooks {
		if runbook.Name == "" {
			return fmt.Errorf("runbook[%d].name is required", i)
		}

		if runbook.URL == "" {
			return fmt.Errorf("runbook[%d].url is required", i)
		}

		// Check for duplicate runbook names
		if runbookNames[runbook.Name] {
			return fmt.Errorf("duplicate runbook name '%s'", runbook.Name)
		}
		runbookNames[runbook.Name] = true

		// Basic URL format validation
		if !strings.HasPrefix(runbook.URL, "http://") && !strings.HasPrefix(runbook.URL, "https://") {
			return fmt.Errorf("runbook[%d].url must be a valid HTTP/HTTPS URL", i)
		}
	}

	return nil
}

// ValidateServiceExists validates that a service exists
func (sv *ServiceValidator) ValidateServiceExists(service *scModels.Service) error {
	if service == nil {
		return fmt.Errorf("service not found")
	}
	return nil
}

// ValidateUserPermissions validates user permissions for service operations
func (sv *ServiceValidator) ValidateUserPermissions(service *scModels.Service, user *scModels.UserContext, operation string) error {
	if user == nil {
		return fmt.Errorf("user context is required")
	}

	// For write operations, check team membership
	if operation == "create" || operation == "update" || operation == "delete" {
		if !service.CanBeModifiedBy(user.Teams) {
			return fmt.Errorf("user does not have permission to %s service '%s'", operation, service.Metadata.Name)
		}
	}

	return nil
}
