package models

// SetupStatus represents the current setup status
type SetupStatus struct {
	SetupRequired bool
	PluginsCount  int
	HasAuth       bool
}

// SetupConfig represents the configuration for initial setup
type SetupConfig struct {
	Port    string
	Origin  string
	Headers []string
	Front   string
	Plugins SetupPluginsConfig
}

// SetupPluginsConfig contains plugin configurations
type SetupPluginsConfig struct {
	Auth           *AuthConfig           `yaml:"auth,omitempty" json:"auth,omitempty"`
	Kubernetes     []KubernetesConfig    `yaml:"kubernetes,omitempty" json:"kubernetes,omitempty"`
	AWS            []AWSConfig           `yaml:"aws,omitempty" json:"aws,omitempty"`
	ServiceCatalog *ServiceCatalogConfig `yaml:"service_catalog,omitempty" json:"service_catalog,omitempty"`
	Observability  *ObservabilityConfig  `yaml:"observability,omitempty" json:"observability,omitempty"`
	EnabledPlugins []string              `yaml:"plugins" json:"plugins"`
}

// AuthConfig represents authentication provider configuration
type AuthConfig struct {
	Provider      string   `yaml:"provider" json:"provider"`
	ClientID      string   `yaml:"clientId" json:"clientId"`
	ClientSecret  string   `yaml:"clientSecret" json:"clientSecret"`
	OrgPermission string   `yaml:"orgPermission,omitempty" json:"orgPermission,omitempty"`
	RedirectURL   string   `yaml:"redirectURL,omitempty" json:"redirectURL,omitempty"`
	Scopes        []string `yaml:"scopes,omitempty" json:"scopes,omitempty"`
}

// KubernetesConfig represents Kubernetes cluster configuration
type KubernetesConfig struct {
	Name           string `yaml:"name" json:"name"`
	Kubeconfig     string `yaml:"kubeconfig,omitempty" json:"kubeconfig,omitempty"`
	Context        string `yaml:"context,omitempty" json:"context,omitempty"`
	ConnectionType string `yaml:"connectionType,omitempty" json:"connectionType,omitempty"`
	Host           string `yaml:"host,omitempty" json:"host,omitempty"`
	Token          string `yaml:"token,omitempty" json:"token,omitempty"`
	Certificate    string `yaml:"certificate,omitempty" json:"certificate,omitempty"`
}

// AWSConfig represents AWS account configuration
type AWSConfig struct {
	Name            string `yaml:"name" json:"name"`
	Region          string `yaml:"region,omitempty" json:"region,omitempty"`
	AccessKeyID     string `yaml:"accessKeyId,omitempty" json:"accessKeyId,omitempty"`
	SecretAccessKey string `yaml:"secretAccessKey,omitempty" json:"secretAccessKey,omitempty"`
}

// ServiceCatalogConfig represents service catalog configuration
type ServiceCatalogConfig struct {
	Storage    StorageConfig     `yaml:"storage" json:"storage"`
	Versioning *VersioningConfig `yaml:"versioning,omitempty" json:"versioning,omitempty"`
}

// StorageConfig represents storage provider configuration
type StorageConfig struct {
	Provider   string             `yaml:"provider" json:"provider"`
	Filesystem *FilesystemStorage `yaml:"filesystem,omitempty" json:"filesystem,omitempty"`
	GitHub     *GitHubStorage     `yaml:"github,omitempty" json:"github,omitempty"`
	S3         *S3Storage         `yaml:"s3,omitempty" json:"s3,omitempty"`
}

// FilesystemStorage represents filesystem storage configuration
type FilesystemStorage struct {
	Directory string `yaml:"directory" json:"directory"`
}

// GitHubStorage represents GitHub storage configuration
type GitHubStorage struct {
	Repository string `yaml:"repository" json:"repository"`
	Branch     string `yaml:"branch,omitempty" json:"branch,omitempty"`
}

// S3Storage represents S3 storage configuration
type S3Storage struct {
	Bucket string `yaml:"bucket" json:"bucket"`
}

// VersioningConfig represents versioning configuration
type VersioningConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Provider string `yaml:"provider,omitempty" json:"provider,omitempty"`
}

// ObservabilityConfig represents observability configuration
type ObservabilityConfig struct {
	Enabled bool                  `yaml:"enabled" json:"enabled"`
	Logs    *ObservabilityLogs    `yaml:"logs,omitempty" json:"logs,omitempty"`
	Traces  *ObservabilityTraces  `yaml:"traces,omitempty" json:"traces,omitempty"`
	Metrics *ObservabilityMetrics `yaml:"metrics,omitempty" json:"metrics,omitempty"`
}

// ObservabilityLogs represents logs provider configuration
type ObservabilityLogs struct {
	Providers []ObservabilityProvider `yaml:"providers,omitempty" json:"providers,omitempty"`
}

// ObservabilityTraces represents traces provider configuration
type ObservabilityTraces struct {
	Providers []ObservabilityProvider `yaml:"providers,omitempty" json:"providers,omitempty"`
}

// ObservabilityMetrics represents metrics provider configuration
type ObservabilityMetrics struct {
	Providers []ObservabilityProvider `yaml:"providers,omitempty" json:"providers,omitempty"`
}

// ObservabilityProvider represents an observability provider
type ObservabilityProvider struct {
	Name      string            `yaml:"name" json:"name"`
	Type      string            `yaml:"type" json:"type"`
	URL       string            `yaml:"url" json:"url"`
	Timeout   string            `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Retention string            `yaml:"retention,omitempty" json:"retention,omitempty"`
	Enabled   bool              `yaml:"enabled" json:"enabled"`
	Auth      *ProviderAuth     `yaml:"auth,omitempty" json:"auth,omitempty"`
	Labels    map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
}

// ProviderAuth represents provider authentication
type ProviderAuth struct {
	Type     string `yaml:"type" json:"type"`
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
	Token    string `yaml:"token,omitempty" json:"token,omitempty"`
}
