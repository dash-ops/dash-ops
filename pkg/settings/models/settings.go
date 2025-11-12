package models

import (
	"fmt"
	"strings"
)

// SettingsConfig represents the current settings configuration
type SettingsConfig struct {
	Config  *DashConfig `json:"config"`
	Plugins []string    `json:"plugins"`
	CanEdit bool        `json:"can_edit"`
}

// Plugins represents a list of enabled plugins with helper methods.
type Plugins []string

// Has checks if a plugin is enabled (case-insensitive).
func (p Plugins) Has(pluginName string) bool {
	for _, plugin := range p {
		if strings.EqualFold(plugin, pluginName) {
			return true
		}
	}
	return false
}

// Add adds a plugin to the list if not already present.
func (p *Plugins) Add(pluginName string) {
	if pluginName == "" {
		return
	}

	if !p.Has(pluginName) {
		*p = append(*p, pluginName)
	}
}

// Remove removes a plugin from the list.
func (p *Plugins) Remove(pluginName string) {
	if pluginName == "" {
		return
	}

	for i, plugin := range *p {
		if strings.EqualFold(plugin, pluginName) {
			*p = append((*p)[:i], (*p)[i+1:]...)
			return
		}
	}
}

// List returns a copy of all enabled plugins.
func (p Plugins) List() []string {
	return append([]string(nil), p...)
}

// Count returns the number of enabled plugins.
func (p Plugins) Count() int {
	return len(p)
}

// DashConfig represents the full dash-ops configuration data.
type DashConfig struct {
	Port           string                `yaml:"port" json:"port"`
	Origin         string                `yaml:"origin" json:"origin"`
	Headers        []string              `yaml:"headers,omitempty" json:"headers,omitempty"`
	Front          string                `yaml:"front,omitempty" json:"front,omitempty"`
	Plugins        Plugins               `yaml:"plugins,omitempty" json:"plugins,omitempty"`
	Auth           []AuthConfig          `yaml:"auth,omitempty" json:"auth,omitempty"`
	Kubernetes     []KubernetesConfig    `yaml:"kubernetes,omitempty" json:"kubernetes,omitempty"`
	AWS            []AWSConfig           `yaml:"aws,omitempty" json:"aws,omitempty"`
	ServiceCatalog *ServiceCatalogConfig `yaml:"service_catalog,omitempty" json:"service_catalog,omitempty"`
	Observability  *ObservabilityConfig  `yaml:"observability,omitempty" json:"observability,omitempty"`
}

// Validate validates the DashConfig ensuring essential fields are present.
func (d *DashConfig) Validate() error {
	if d == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	if strings.TrimSpace(d.Port) == "" {
		return fmt.Errorf("port is required")
	}

	if strings.TrimSpace(d.Origin) == "" {
		return fmt.Errorf("origin is required")
	}

	return nil
}

// GetPort returns the configured port with default fallback.
func (d *DashConfig) GetPort() string {
	if d == nil || strings.TrimSpace(d.Port) == "" {
		return "8080"
	}
	return d.Port
}

// GetOrigin returns the configured origin with default fallback.
func (d *DashConfig) GetOrigin() string {
	if d == nil || strings.TrimSpace(d.Origin) == "" {
		return "http://localhost:5173"
	}
	return d.Origin
}

// GetHeaders returns configured headers or sensible defaults.
func (d *DashConfig) GetHeaders() []string {
	if d == nil || len(d.Headers) == 0 {
		return []string{"Content-Type", "Authorization"}
	}
	return d.Headers
}

// GetFront returns the frontend dist path with default fallback.
func (d *DashConfig) GetFront() string {
	if d == nil || strings.TrimSpace(d.Front) == "" {
		return "front/dist"
	}
	return d.Front
}

// IsPluginEnabled checks if a specific plugin is enabled.
func (d *DashConfig) IsPluginEnabled(pluginName string) bool {
	if d == nil {
		return false
	}
	return d.Plugins.Has(pluginName)
}

// Clone returns a deep copy of DashConfig.
func (d *DashConfig) Clone() *DashConfig {
	if d == nil {
		return nil
	}

	clone := *d
	clone.Headers = append([]string(nil), d.Headers...)
	clone.Plugins = Plugins(append([]string(nil), d.Plugins...))
	clone.Auth = append([]AuthConfig(nil), d.Auth...)
	clone.Kubernetes = append([]KubernetesConfig(nil), d.Kubernetes...)
	clone.AWS = append([]AWSConfig(nil), d.AWS...)

	if d.ServiceCatalog != nil {
		serviceCatalogClone := *d.ServiceCatalog
		clone.ServiceCatalog = &serviceCatalogClone
	}

	if d.Observability != nil {
		observabilityClone := *d.Observability
		clone.Observability = &observabilityClone
	}

	return &clone
}

// UpdateSettingsRequest represents a request to update settings
type UpdateSettingsRequest struct {
	Config         *UpdateGeneralConfig `json:"config,omitempty"`
	Plugins        *UpdatePlugins       `json:"plugins,omitempty"`
	EnabledPlugins []string             `json:"enabled_plugins,omitempty"`
}

// UpdateGeneralConfig represents updates to general configuration.
type UpdateGeneralConfig struct {
	Port    *string  `json:"port,omitempty"`
	Origin  *string  `json:"origin,omitempty"`
	Headers []string `json:"headers,omitempty"`
	Front   *string  `json:"front,omitempty"`
}

// UpdatePlugins represents updates to plugin configurations.
type UpdatePlugins struct {
	Auth           []UpdateAuthConfig       `json:"auth,omitempty"`
	Kubernetes     []UpdateKubernetesConfig `json:"kubernetes,omitempty"`
	AWS            []UpdateAWSConfig        `json:"aws,omitempty"`
	ServiceCatalog *ServiceCatalogConfig    `json:"service_catalog,omitempty"`
	Observability  *ObservabilityConfig     `json:"observability,omitempty"`
}

// UpdateAuthConfig represents an update to an auth provider configuration.
type UpdateAuthConfig struct {
	Provider          string   `json:"provider"`
	ClientID          *string  `json:"clientId,omitempty"`
	ClientSecret      *string  `json:"clientSecret,omitempty"`
	OrgPermission     *string  `json:"orgPermission,omitempty"`
	RedirectURL       *string  `json:"redirectURL,omitempty"`
	Scopes            []string `json:"scopes,omitempty"`
	ClearClientSecret bool     `json:"clearClientSecret,omitempty"`
}

// UpdateKubernetesConfig represents an update to a kubernetes cluster configuration.
type UpdateKubernetesConfig struct {
	Name           string  `json:"name"`
	ConnectionType *string `json:"connectionType,omitempty"`
	Kubeconfig     *string `json:"kubeconfig,omitempty"`
	Context        *string `json:"context,omitempty"`
	Host           *string `json:"host,omitempty"`
	Token          *string `json:"token,omitempty"`
	Certificate    *string `json:"certificate,omitempty"`
	ClearToken     bool    `json:"clearToken,omitempty"`
}

// UpdateAWSConfig represents an update to an aws account configuration.
type UpdateAWSConfig struct {
	Name                 string  `json:"name"`
	Region               *string `json:"region,omitempty"`
	AccessKeyID          *string `json:"accessKeyId,omitempty"`
	SecretAccessKey      *string `json:"secretAccessKey,omitempty"`
	ClearSecretAccessKey bool    `json:"clearSecretAccessKey,omitempty"`
}

// UpdateSettingsResponse represents the response after updating settings
type UpdateSettingsResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	RequiresRestart bool   `json:"requires_restart"`
}
