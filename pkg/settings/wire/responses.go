package wire

// SetupStatusResponse represents the setup status response
type SetupStatusResponse struct {
	SetupRequired bool `json:"setup_required"`
	PluginsCount  int  `json:"plugins_count"`
	HasAuth       bool `json:"has_auth"`
}

// SetupConfigureResponse represents the response after configuring setup
type SetupConfigureResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	ConfigPath string `json:"config_path,omitempty"`
}

// SettingsConfigResponse represents the settings configuration response
type SettingsConfigResponse struct {
	Config  *DashConfigData `json:"config"`
	Plugins []string        `json:"plugins"`
	CanEdit bool            `json:"can_edit"`
}

// DashConfigData represents dash config data in wire format
type DashConfigData struct {
	Port           string                  `json:"port"`
	Origin         string                  `json:"origin"`
	Headers        []string                `json:"headers,omitempty"`
	Front          string                  `json:"front,omitempty"`
	Plugins        []string                `json:"plugins,omitempty"`
	Auth           []AuthProviderData      `json:"auth,omitempty"`
	Kubernetes     []KubernetesClusterData `json:"kubernetes,omitempty"`
	AWS            []AWSAccountData        `json:"aws,omitempty"`
	ServiceCatalog *ServiceCatalogData     `json:"service_catalog,omitempty"`
	Observability  *ObservabilityData      `json:"observability,omitempty"`
}

// AuthProviderData represents sanitized auth provider information.
type AuthProviderData struct {
	Provider        string   `json:"provider"`
	ClientIDMasked  string   `json:"client_id_masked,omitempty"`
	OrgPermission   string   `json:"org_permission,omitempty"`
	HasClientSecret bool     `json:"has_client_secret"`
	RedirectURL     string   `json:"redirect_url,omitempty"`
	Scopes          []string `json:"scopes,omitempty"`
}

// KubernetesClusterData represents sanitized kubernetes cluster information.
type KubernetesClusterData struct {
	Name           string `json:"name"`
	ConnectionType string `json:"connection_type,omitempty"`
	Kubeconfig     string `json:"kubeconfig,omitempty"`
	Context        string `json:"context,omitempty"`
	Host           string `json:"host,omitempty"`
	HasToken       bool   `json:"has_token"`
}

// AWSAccountData represents sanitized aws account information.
type AWSAccountData struct {
	Name               string `json:"name"`
	Region             string `json:"region,omitempty"`
	AccessKeyIDMasked  string `json:"access_key_id_masked,omitempty"`
	HasSecretAccessKey bool   `json:"has_secret_access_key"`
}

// ServiceCatalogData represents service catalog wire data.
type ServiceCatalogData struct {
	Storage    ServiceCatalogStorageData     `json:"storage"`
	Versioning *ServiceCatalogVersioningData `json:"versioning,omitempty"`
}

// ServiceCatalogStorageData represents storage provider wire data.
type ServiceCatalogStorageData struct {
	Provider   string                 `json:"provider"`
	Directory  string                 `json:"directory,omitempty"`
	Repository string                 `json:"repository,omitempty"`
	Branch     string                 `json:"branch,omitempty"`
	Filesystem *FilesystemStorageData `json:"filesystem,omitempty"`
	GitHub     *GitHubStorageData     `json:"github,omitempty"`
	S3         *S3StorageData         `json:"s3,omitempty"`
}

// FilesystemStorageData represents filesystem storage details.
type FilesystemStorageData struct {
	Directory string `json:"directory"`
}

// GitHubStorageData represents GitHub storage details.
type GitHubStorageData struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch,omitempty"`
}

// S3StorageData represents S3 storage details.
type S3StorageData struct {
	Bucket string `json:"bucket"`
}

// ServiceCatalogVersioningData represents versioning wire data.
type ServiceCatalogVersioningData struct {
	Enabled  bool   `json:"enabled"`
	Provider string `json:"provider,omitempty"`
}

// ObservabilityData represents observability wire data.
type ObservabilityData struct {
	Enabled bool                        `json:"enabled"`
	Logs    []ObservabilityProviderData `json:"logs,omitempty"`
	Traces  []ObservabilityProviderData `json:"traces,omitempty"`
	Metrics []ObservabilityProviderData `json:"metrics,omitempty"`
}

// ObservabilityProviderData represents observability provider wire data.
type ObservabilityProviderData struct {
	Name      string                     `json:"name"`
	Type      string                     `json:"type"`
	URL       string                     `json:"url"`
	Timeout   string                     `json:"timeout,omitempty"`
	Retention string                     `json:"retention,omitempty"`
	Enabled   bool                       `json:"enabled"`
	Labels    map[string]string          `json:"labels,omitempty"`
	Auth      *ObservabilityProviderAuth `json:"auth,omitempty"`
}

// ObservabilityProviderAuth represents sanitized auth information.
type ObservabilityProviderAuth struct {
	Type        string `json:"type"`
	HasPassword bool   `json:"has_password"`
	HasToken    bool   `json:"has_token"`
	HasUsername bool   `json:"has_username"`
}

// UpdateSettingsResponse represents the response after updating settings
type UpdateSettingsResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	RequiresRestart bool   `json:"requires_restart"`
}
