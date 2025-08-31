package servicecatalog

import (
	"fmt"
	"log"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// ServiceCatalogConfig represents the complete service catalog configuration
type ServiceCatalogConfig struct {
	ServiceCatalog Config `yaml:"service_catalog"`
}

// MakeServiceCatalogHandlers creates and registers service catalog HTTP handlers
func MakeServiceCatalogHandlers(internal *mux.Router, fileConfig []byte) *ServiceCatalog {
	// Parse configuration
	config, err := parseServiceCatalogConfig(fileConfig)
	if err != nil {
		log.Printf("Failed to parse service catalog configuration: %v", err)
		return nil
	}

	// Create service catalog instance
	serviceCatalog, err := NewServiceCatalog(&config.ServiceCatalog)
	if err != nil {
		log.Printf("Failed to initialize service catalog: %v", err)
		return nil
	}

	// Create HTTP handler
	handler := NewHandler(serviceCatalog)

	// Register routes under /service-catalog prefix
	serviceCatalogRouter := internal.PathPrefix("/service-catalog").Subrouter()
	handler.RegisterRoutes(serviceCatalogRouter)

	log.Println("Service Catalog handlers registered successfully")
	return serviceCatalog
}

// parseServiceCatalogConfig parses the service catalog configuration from YAML
func parseServiceCatalogConfig(fileConfig []byte) (*ServiceCatalogConfig, error) {
	var config ServiceCatalogConfig

	if err := yaml.Unmarshal(fileConfig, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service catalog config: %w", err)
	}

	// Set default values if not provided
	if config.ServiceCatalog.Storage.Provider == "" {
		config.ServiceCatalog.Storage.Provider = "filesystem"
	}

	if config.ServiceCatalog.Storage.Filesystem.Directory == "" {
		config.ServiceCatalog.Storage.Filesystem.Directory = "./services"
	}

	return &config, nil
}

// ValidateServiceCatalogConfig validates the service catalog configuration
func ValidateServiceCatalogConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("service catalog configuration is required")
	}

	// Validate storage provider
	validProviders := map[string]bool{
		"filesystem": true,
		"github":     true,
		"s3":         true,
	}

	if !validProviders[config.Storage.Provider] {
		return fmt.Errorf("invalid storage provider '%s', must be one of: filesystem, github, s3", config.Storage.Provider)
	}

	// Validate provider-specific configuration
	switch config.Storage.Provider {
	case "filesystem":
		if config.Storage.Filesystem.Directory == "" {
			return fmt.Errorf("filesystem directory is required when using filesystem provider")
		}

	case "github":
		if config.Storage.GitHub.Repository == "" {
			return fmt.Errorf("github repository is required when using github provider")
		}
		if config.Storage.GitHub.Branch == "" {
			config.Storage.GitHub.Branch = "main" // Default branch
		}

	case "s3":
		if config.Storage.S3.Bucket == "" {
			return fmt.Errorf("s3 bucket is required when using s3 provider")
		}
	}

	return nil
}
