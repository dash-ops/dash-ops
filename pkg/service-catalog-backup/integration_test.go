package servicecatalog

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

// TestServiceContextResolver_Integration tests the resolver with real service data
func TestServiceContextResolver_Integration(t *testing.T) {
	// Skip integration tests in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "service-catalog-test-*")
	if err != nil {
		t.Fatal("Failed to create temp directory:", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test service definition
	testService := Service{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: ServiceMetadata{
			Name: "user-authentication",
			Tier: "TIER-1",
		},
		Spec: ServiceSpec{
			Description: "Test authentication service",
			Team: ServiceTeam{
				GitHubTeam: "auth-squad",
			},
			Business: ServiceBusiness{
				SLATarget: "99.9%",
				Impact:    "high",
			},
			Kubernetes: &ServiceKubernetes{
				Environments: []KubernetesEnvironment{
					{
						Name:      "local",
						Context:   "docker-desktop",
						Namespace: "auth",
						Resources: KubernetesEnvironmentResources{
							Deployments: []KubernetesDeployment{
								{
									Name:     "auth-api",
									Replicas: 3,
									Resources: KubernetesResourceRequests{
										Requests: KubernetesResourceSpec{
											CPU:    "200m",
											Memory: "256Mi",
										},
										Limits: KubernetesResourceSpec{
											CPU:    "1000m",
											Memory: "512Mi",
										},
									},
								},
								{
									Name:     "auth-worker",
									Replicas: 2,
									Resources: KubernetesResourceRequests{
										Requests: KubernetesResourceSpec{
											CPU:    "100m",
											Memory: "128Mi",
										},
										Limits: KubernetesResourceSpec{
											CPU:    "500m",
											Memory: "256Mi",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Write service to file
	serviceFile := filepath.Join(tempDir, "user-authentication.yaml")
	serviceData, err := yaml.Marshal(testService)
	if err != nil {
		t.Fatal("Failed to marshal service:", err)
	}

	err = os.WriteFile(serviceFile, serviceData, 0644)
	if err != nil {
		t.Fatal("Failed to write service file:", err)
	}

	// Create service catalog config
	config := &Config{
		Storage: StorageConfig{
			Provider: "filesystem",
			Filesystem: FilesystemStorageConfig{
				Directory: tempDir,
			},
		},
	}

	// Initialize service catalog
	serviceCatalog, err := NewServiceCatalog(config)
	if err != nil {
		t.Fatal("Failed to create service catalog:", err)
	}

	// Test 1: Resolve existing deployment
	t.Run("resolve existing deployment", func(t *testing.T) {
		ctx, err := serviceCatalog.ResolveDeploymentService("auth-api", "auth", "docker-desktop")
		if err != nil {
			t.Fatal("Failed to resolve deployment service:", err)
		}

		if ctx == nil {
			t.Fatal("Expected service context but got nil")
		}

		// Validate resolved context
		if ctx.ServiceName != "user-authentication" {
			t.Errorf("Expected service name 'user-authentication', got '%s'", ctx.ServiceName)
		}

		if ctx.ServiceTier != "TIER-1" {
			t.Errorf("Expected service tier 'TIER-1', got '%s'", ctx.ServiceTier)
		}

		if ctx.Environment != "local" {
			t.Errorf("Expected environment 'local', got '%s'", ctx.Environment)
		}

		if ctx.Team != "auth-squad" {
			t.Errorf("Expected team 'auth-squad', got '%s'", ctx.Team)
		}
	})

	// Test 2: Resolve non-existing deployment
	t.Run("resolve non-existing deployment", func(t *testing.T) {
		ctx, err := serviceCatalog.ResolveDeploymentService("non-existent-api", "auth", "docker-desktop")
		if err != nil {
			t.Fatal("Failed to resolve deployment service:", err)
		}

		if ctx != nil {
			t.Fatal("Expected nil service context for non-existing deployment")
		}
	})

	// Test 3: Case insensitive match
	t.Run("case insensitive deployment match", func(t *testing.T) {
		ctx, err := serviceCatalog.ResolveDeploymentService("AUTH-API", "auth", "docker-desktop")
		if err != nil {
			t.Fatal("Failed to resolve deployment service:", err)
		}

		if ctx == nil {
			t.Fatal("Expected service context for case insensitive match but got nil")
		}

		if ctx.ServiceName != "user-authentication" {
			t.Errorf("Expected service name 'user-authentication', got '%s'", ctx.ServiceName)
		}
	})

	// Test 4: Get service deployments
	t.Run("get service deployments", func(t *testing.T) {
		deployments, err := serviceCatalog.GetContextResolver().GetServiceDeployments("user-authentication")
		if err != nil {
			t.Fatal("Failed to get service deployments:", err)
		}

		if len(deployments) != 2 {
			t.Fatalf("Expected 2 deployments, got %d", len(deployments))
		}

		// Check deployment names
		names := make(map[string]bool)
		for _, deployment := range deployments {
			names[deployment.DeploymentName] = true
		}

		if !names["auth-api"] {
			t.Error("Expected deployment 'auth-api' not found")
		}

		if !names["auth-worker"] {
			t.Error("Expected deployment 'auth-worker' not found")
		}
	})
}
