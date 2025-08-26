package servicecatalog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// KubernetesIntegration handles communication with Kubernetes plugin
type KubernetesIntegration struct {
	baseURL string // Base URL for internal API calls
}

// NewKubernetesIntegration creates a new Kubernetes integration instance
func NewKubernetesIntegration(baseURL string) *KubernetesIntegration {
	if baseURL == "" {
		baseURL = "http://localhost:8080" // Default to same server
	}

	return &KubernetesIntegration{
		baseURL: baseURL,
	}
}

// Use types from Kubernetes plugin (external types)
type K8sDeployment struct {
	Name       string         `json:"name"`
	Namespace  string         `json:"namespace"`
	PodInfo    K8sPodInfo     `json:"pod_info"`
	Replicas   K8sReplicas    `json:"replicas"`
	Age        string         `json:"age"`
	CreatedAt  time.Time      `json:"created_at"`
	Conditions []K8sCondition `json:"conditions"`
}

type K8sPodInfo struct {
	Current int32 `json:"current"`
	Desired int32 `json:"desired"`
}

type K8sReplicas struct {
	Ready     int32 `json:"ready"`
	Available int32 `json:"available"`
	Updated   int32 `json:"updated"`
	Desired   int32 `json:"desired"`
}

type K8sCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// GetDeployments retrieves deployments from Kubernetes plugin
func (ki *KubernetesIntegration) GetDeployments(context, namespace string, authHeader string) ([]K8sDeployment, error) {
	url := fmt.Sprintf("%s/api/v1/k8s/%s/deployments", ki.baseURL, context)
	if namespace != "" {
		url += fmt.Sprintf("?namespace=%s", namespace)
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Forward authorization header
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	// Make request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Kubernetes API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Kubernetes API returned status %d", resp.StatusCode)
	}

	// Parse response
	var deployments []K8sDeployment
	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return nil, fmt.Errorf("failed to parse Kubernetes response: %w", err)
	}

	return deployments, nil
}

// AggregateServiceHealth aggregates health from Kubernetes deployments
func (ki *KubernetesIntegration) AggregateServiceHealth(service *Service, authHeader string) (*ServiceHealth, error) {
	if service.Spec.Kubernetes == nil {
		// No Kubernetes definition - return unknown status
		return &ServiceHealth{
			ServiceName:   service.Metadata.Name,
			OverallStatus: "unknown",
			Environments:  []EnvironmentHealth{},
			LastUpdated:   time.Now(),
		}, nil
	}

	var environments []EnvironmentHealth

	for _, env := range service.Spec.Kubernetes.Environments {
		// Get deployments for this environment
		k8sDeployments, err := ki.GetDeployments(env.Context, env.Namespace, authHeader)
		if err != nil {
			// If we can't get K8s data, mark as unknown
			environments = append(environments, EnvironmentHealth{
				Name:        env.Name,
				Context:     env.Context,
				Status:      "unknown",
				Deployments: []DeploymentHealth{},
			})
			continue
		}

		// Filter deployments that match our service definition
		var matchingDeployments []DeploymentHealth
		for _, expectedDeploy := range env.Resources.Deployments {
			// Find matching deployment in K8s
			var k8sDeploy *K8sDeployment
			for i := range k8sDeployments {
				if k8sDeployments[i].Name == expectedDeploy.Name {
					k8sDeploy = &k8sDeployments[i]
					break
				}
			}

			if k8sDeploy == nil {
				// Deployment not found in cluster
				matchingDeployments = append(matchingDeployments, DeploymentHealth{
					Name:            expectedDeploy.Name,
					ReadyReplicas:   0,
					DesiredReplicas: int(expectedDeploy.Replicas),
					Status:          "NotFound",
					LastUpdated:     time.Now(),
				})
			} else {
				// Calculate deployment status
				status := ki.calculateDeploymentStatus(k8sDeploy, expectedDeploy.Replicas)

				matchingDeployments = append(matchingDeployments, DeploymentHealth{
					Name:            k8sDeploy.Name,
					ReadyReplicas:   int(k8sDeploy.Replicas.Ready),
					DesiredReplicas: int(k8sDeploy.Replicas.Desired),
					Status:          status,
					LastUpdated:     k8sDeploy.CreatedAt,
				})
			}
		}

		// Calculate environment status
		envStatus := ki.calculateEnvironmentStatus(matchingDeployments)

		environments = append(environments, EnvironmentHealth{
			Name:        env.Name,
			Context:     env.Context,
			Status:      envStatus,
			Deployments: matchingDeployments,
		})
	}

	// Calculate overall service status
	overallStatus := ki.calculateServiceStatus(environments, service.Metadata.Tier)

	return &ServiceHealth{
		ServiceName:   service.Metadata.Name,
		OverallStatus: overallStatus,
		Environments:  environments,
		LastUpdated:   time.Now(),
	}, nil
}

// calculateDeploymentStatus determines status based on K8s deployment data
func (ki *KubernetesIntegration) calculateDeploymentStatus(k8sDeploy *K8sDeployment, expectedReplicas int) string {
	// Check if deployment is available
	available := false
	progressing := false

	for _, condition := range k8sDeploy.Conditions {
		if condition.Type == "Available" && condition.Status == "True" {
			available = true
		}
		if condition.Type == "Progressing" && condition.Status == "True" {
			progressing = true
		}
	}

	// Check replica health
	ready := k8sDeploy.Replicas.Ready
	desired := k8sDeploy.Replicas.Desired

	if !available || ready == 0 {
		return "down"
	}

	if ready < desired || int(ready) < expectedReplicas {
		return "degraded"
	}

	if available && progressing && ready == desired {
		return "healthy"
	}

	return "unknown"
}

// calculateEnvironmentStatus determines environment status from deployments
func (ki *KubernetesIntegration) calculateEnvironmentStatus(deployments []DeploymentHealth) string {
	if len(deployments) == 0 {
		return "unknown"
	}

	healthyCount := 0
	downCount := 0
	degradedCount := 0

	for _, deploy := range deployments {
		switch deploy.Status {
		case "healthy":
			healthyCount++
		case "down", "NotFound":
			downCount++
		case "degraded":
			degradedCount++
		}
	}

	// If any deployment is down, environment is down
	if downCount > 0 {
		return "down"
	}

	// If any deployment is degraded, environment is degraded
	if degradedCount > 0 {
		return "degraded"
	}

	// If all deployments are healthy
	if healthyCount == len(deployments) {
		return "healthy"
	}

	return "unknown"
}

// calculateServiceStatus determines overall service status from environments and tier
func (ki *KubernetesIntegration) calculateServiceStatus(environments []EnvironmentHealth, tier string) string {
	if len(environments) == 0 {
		return "unknown"
	}

	// Find production environment status
	var prodStatus string
	for _, env := range environments {
		// Check for production environment (common names)
		if env.Name == "production" || env.Name == "prod" {
			prodStatus = env.Status
			break
		}
	}

	// If no production, use first environment
	if prodStatus == "" && len(environments) > 0 {
		prodStatus = environments[0].Status
	}

	// Apply tier-based logic (same as before)
	switch tier {
	case "TIER-1":
		// Critical services: any issues are critical
		if prodStatus == "down" || prodStatus == "degraded" {
			return "critical"
		}
		return prodStatus

	case "TIER-2":
		// Important services: production down is degraded
		if prodStatus == "down" {
			return "degraded"
		}
		return prodStatus

	case "TIER-3":
		// Standard services: only complete failure matters
		if prodStatus == "down" {
			return "degraded"
		}
		return "healthy"

	default:
		return prodStatus
	}
}
