package servicecatalog

import (
	"context"
	"fmt"

	scAdapters "github.com/dash-ops/dash-ops/pkg/service-catalog-new/adapters/storage"
	scModels "github.com/dash-ops/dash-ops/pkg/service-catalog-new/models"
)

// ExampleUsage demonstrates how to use the new architecture
func ExampleUsage() error {
	// 1. Initialize storage adapter
	fsRepo, err := scAdapters.NewFilesystemRepository("./test-services")
	if err != nil {
		return fmt.Errorf("failed to create filesystem repository: %w", err)
	}

	// 2. Create module with minimal configuration
	module, err := NewMinimalModule(fsRepo)
	if err != nil {
		return fmt.Errorf("failed to create module: %w", err)
	}

	// 3. Create a test service
	service := &scModels.Service{
		Metadata: scModels.ServiceMetadata{
			Name: "example-service",
			Tier: scModels.TierStandard,
		},
		Spec: scModels.ServiceSpec{
			Description: "Example service for testing new architecture",
			Team: scModels.ServiceTeam{
				GitHubTeam: "example-team",
			},
			Business: scModels.ServiceBusiness{
				Impact: "medium",
			},
		},
	}

	// 4. Create user context
	user := &scModels.UserContext{
		Username: "test-user",
		Name:     "Test User",
		Email:    "test@example.com",
		Teams:    []string{"example-team"},
	}

	// 5. Use the controller to create the service
	ctx := context.Background()
	createdService, err := module.GetController().CreateService(ctx, service, user)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	fmt.Printf("✅ Service created successfully: %s (Version: %d)\n",
		createdService.Metadata.Name, createdService.Metadata.Version)

	// 6. Retrieve the service
	retrievedService, err := module.GetController().GetService(ctx, "example-service")
	if err != nil {
		return fmt.Errorf("failed to get service: %w", err)
	}

	fmt.Printf("✅ Service retrieved: %s - %s\n",
		retrievedService.Metadata.Name, retrievedService.Spec.Description)

	// 7. List all services
	serviceList, err := module.GetController().ListServices(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to list services: %w", err)
	}

	fmt.Printf("✅ Total services: %d\n", serviceList.Total)

	return nil
}
