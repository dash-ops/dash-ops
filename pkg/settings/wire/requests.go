package wire

// SetupConfigureRequest represents the request to configure initial setup
type SetupConfigureRequest struct {
	Config         SetupConfigRequest  `json:"config"`
	Plugins        SetupPluginsRequest `json:"plugins"`
	EnabledPlugins []string            `json:"enabled_plugins"`
}

// SetupConfigRequest represents the configuration section for setup.
type SetupConfigRequest struct {
	Port    string   `json:"port"`
	Origin  string   `json:"origin"`
	Headers []string `json:"headers"`
	Front   string   `json:"front"`
}

// SetupPluginsRequest represents plugin configuration for setup.
type SetupPluginsRequest struct {
	Auth           *AuthProviderRequest       `json:"auth,omitempty"`
	Kubernetes     []KubernetesClusterRequest `json:"kubernetes,omitempty"`
	AWS            []AWSAccountRequest        `json:"aws,omitempty"`
	ServiceCatalog *ServiceCatalogRequest     `json:"service_catalog,omitempty"`
	Observability  *ObservabilityRequest      `json:"observability,omitempty"`
}

// UpdateSettingsRequest represents the request to update settings
type UpdateSettingsRequest struct {
	Config         *UpdateConfigRequest  `json:"config,omitempty"`
	Plugins        *UpdatePluginsRequest `json:"plugins,omitempty"`
	EnabledPlugins []string              `json:"enabled_plugins,omitempty"`
}

// UpdateConfigRequest represents updates to general configuration.
type UpdateConfigRequest struct {
	Port    *string  `json:"port,omitempty"`
	Origin  *string  `json:"origin,omitempty"`
	Headers []string `json:"headers,omitempty"`
	Front   *string  `json:"front,omitempty"`
}

// UpdatePluginsRequest represents plugin-specific updates.
type UpdatePluginsRequest struct {
	Auth           []UpdateAuthProviderRequest      `json:"auth,omitempty"`
	Kubernetes     []UpdateKubernetesClusterRequest `json:"kubernetes,omitempty"`
	AWS            []UpdateAWSAccountRequest        `json:"aws,omitempty"`
	ServiceCatalog *UpdateServiceCatalogRequest     `json:"service_catalog,omitempty"`
	Observability  *UpdateObservabilityRequest      `json:"observability,omitempty"`
}

// AuthProviderRequest represents setup auth provider configuration.
type AuthProviderRequest struct {
	Provider      string   `json:"provider"`
	ClientID      string   `json:"clientId"`
	ClientSecret  string   `json:"clientSecret"`
	OrgPermission string   `json:"orgPermission,omitempty"`
	RedirectURL   string   `json:"redirectURL,omitempty"`
	Scopes        []string `json:"scopes,omitempty"`
}

// UpdateAuthProviderRequest represents updates to auth provider configuration.
type UpdateAuthProviderRequest struct {
	Provider          string   `json:"provider"`
	ClientID          *string  `json:"clientId,omitempty"`
	ClientSecret      *string  `json:"clientSecret,omitempty"`
	OrgPermission     *string  `json:"orgPermission,omitempty"`
	RedirectURL       *string  `json:"redirectURL,omitempty"`
	Scopes            []string `json:"scopes,omitempty"`
	ClearClientSecret bool     `json:"clearClientSecret,omitempty"`
}

// KubernetesClusterRequest represents setup kubernetes configuration.
type KubernetesClusterRequest struct {
	Name           string `json:"name"`
	ConnectionType string `json:"connectionType,omitempty"`
	Kubeconfig     string `json:"kubeconfig,omitempty"`
	Context        string `json:"context,omitempty"`
	Host           string `json:"host,omitempty"`
	Token          string `json:"token,omitempty"`
	Certificate    string `json:"certificate,omitempty"`
}

// UpdateKubernetesClusterRequest represents updates to kubernetes configuration.
type UpdateKubernetesClusterRequest struct {
	Name           string  `json:"name"`
	ConnectionType *string `json:"connectionType,omitempty"`
	Kubeconfig     *string `json:"kubeconfig,omitempty"`
	Context        *string `json:"context,omitempty"`
	Host           *string `json:"host,omitempty"`
	Token          *string `json:"token,omitempty"`
	Certificate    *string `json:"certificate,omitempty"`
	ClearToken     bool    `json:"clearToken,omitempty"`
}

// AWSAccountRequest represents setup aws configuration.
type AWSAccountRequest struct {
	Name            string `json:"name"`
	Region          string `json:"region,omitempty"`
	AccessKeyID     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

// UpdateAWSAccountRequest represents updates to aws configuration.
type UpdateAWSAccountRequest struct {
	Name                 string  `json:"name"`
	Region               *string `json:"region,omitempty"`
	AccessKeyID          *string `json:"accessKeyId,omitempty"`
	SecretAccessKey      *string `json:"secretAccessKey,omitempty"`
	ClearSecretAccessKey bool    `json:"clearSecretAccessKey,omitempty"`
}

// ServiceCatalogRequest represents setup service catalog configuration.
type ServiceCatalogRequest struct {
	Storage    ServiceCatalogStorageRequest     `json:"storage"`
	Versioning *ServiceCatalogVersioningRequest `json:"versioning,omitempty"`
}

// UpdateServiceCatalogRequest represents updates to service catalog configuration.
type UpdateServiceCatalogRequest struct {
	Storage    *ServiceCatalogStorageRequest    `json:"storage,omitempty"`
	Versioning *ServiceCatalogVersioningRequest `json:"versioning,omitempty"`
}

// ServiceCatalogStorageRequest represents storage provider configuration.
type ServiceCatalogStorageRequest struct {
	Provider   string                    `json:"provider"`
	Directory  string                    `json:"directory,omitempty"`
	Repository string                    `json:"repository,omitempty"`
	Branch     string                    `json:"branch,omitempty"`
	Filesystem *FilesystemStorageRequest `json:"filesystem,omitempty"`
	GitHub     *GitHubStorageRequest     `json:"github,omitempty"`
	S3         *S3StorageRequest         `json:"s3,omitempty"`
}

// ServiceCatalogVersioningRequest represents versioning configuration.
type ServiceCatalogVersioningRequest struct {
	Enabled  bool   `json:"enabled"`
	Provider string `json:"provider,omitempty"`
}

// FilesystemStorageRequest represents filesystem storage configuration.
type FilesystemStorageRequest struct {
	Directory string `json:"directory"`
}

// GitHubStorageRequest represents GitHub storage configuration.
type GitHubStorageRequest struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch,omitempty"`
}

// S3StorageRequest represents S3 storage configuration.
type S3StorageRequest struct {
	Bucket string `json:"bucket"`
}

// ObservabilityRequest represents setup observability configuration.
type ObservabilityRequest struct {
	Enabled bool                           `json:"enabled"`
	Logs    []ObservabilityProviderRequest `json:"logs,omitempty"`
	Traces  []ObservabilityProviderRequest `json:"traces,omitempty"`
	Metrics []ObservabilityProviderRequest `json:"metrics,omitempty"`
}

// UpdateObservabilityRequest represents updates to observability configuration.
type UpdateObservabilityRequest struct {
	Enabled *bool                          `json:"enabled,omitempty"`
	Logs    []ObservabilityProviderRequest `json:"logs,omitempty"`
	Traces  []ObservabilityProviderRequest `json:"traces,omitempty"`
	Metrics []ObservabilityProviderRequest `json:"metrics,omitempty"`
}

// ObservabilityProviderRequest represents observability provider configuration.
type ObservabilityProviderRequest struct {
	Name      string               `json:"name"`
	Type      string               `json:"type"`
	URL       string               `json:"url"`
	Timeout   string               `json:"timeout,omitempty"`
	Retention string               `json:"retention,omitempty"`
	Enabled   bool                 `json:"enabled"`
	Labels    map[string]string    `json:"labels,omitempty"`
	Auth      *ProviderAuthRequest `json:"auth,omitempty"`
}

// ProviderAuthRequest represents observability provider auth configuration.
type ProviderAuthRequest struct {
	Type          string  `json:"type"`
	Username      *string `json:"username,omitempty"`
	Password      *string `json:"password,omitempty"`
	Token         *string `json:"token,omitempty"`
	ClearPassword bool    `json:"clearPassword,omitempty"`
	ClearToken    bool    `json:"clearToken,omitempty"`
}
