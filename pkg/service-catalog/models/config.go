package models

// ModuleConfig represents configuration for the service catalog module
type ModuleConfig struct {
	// Configuration data
	Directory string `yaml:"directory" json:"directory"`
}

// ServiceCatalogConfig represents the service catalog configuration structure
type ServiceCatalogConfig struct {
	ServiceCatalog struct {
		Storage struct {
			Provider   string `yaml:"provider"`
			Filesystem struct {
				Directory string `yaml:"directory"`
			} `yaml:"filesystem"`
		} `yaml:"storage"`
	} `yaml:"service_catalog"`
}

// ParsedConfig represents parsed configuration data
type ParsedConfig struct {
	Directory string
}
