package logic

import (
	"fmt"
	"strings"

	"github.com/dash-ops/dash-ops/pkg/settings/models"
)

// SettingsProcessor handles settings-related business logic
type SettingsProcessor struct {
	yamlProcessor *YAMLProcessor
}

// NewSettingsProcessor creates a new settings processor
func NewSettingsProcessor(yamlProcessor *YAMLProcessor) *SettingsProcessor {
	return &SettingsProcessor{
		yamlProcessor: yamlProcessor,
	}
}

// ProcessUpdateRequest processes an update settings request
func (sp *SettingsProcessor) ProcessUpdateRequest(
	currentConfig *models.DashConfig,
	request *models.UpdateSettingsRequest,
) (*models.DashConfig, error) {
	if currentConfig == nil {
		return nil, fmt.Errorf("current config cannot be nil")
	}

	if request == nil {
		return nil, fmt.Errorf("update request cannot be nil")
	}

	// Create updated config (copy current)
	updatedConfig := currentConfig.Clone()
	if updatedConfig == nil {
		return nil, fmt.Errorf("failed to clone current configuration")
	}

	// Update general configuration if provided
	if request.Config != nil {
		if request.Config.Port != nil {
			updatedConfig.Port = strings.TrimSpace(*request.Config.Port)
		}
		if request.Config.Origin != nil {
			updatedConfig.Origin = strings.TrimSpace(*request.Config.Origin)
		}
		if len(request.Config.Headers) > 0 {
			headers := make([]string, 0, len(request.Config.Headers))
			for _, header := range request.Config.Headers {
				header = strings.TrimSpace(header)
				if header != "" {
					headers = append(headers, header)
				}
			}
			if len(headers) > 0 {
				updatedConfig.Headers = headers
			}
		}
		if request.Config.Front != nil {
			updatedConfig.Front = strings.TrimSpace(*request.Config.Front)
		}
	}

	// Update enabled plugins if provided
	if len(request.EnabledPlugins) > 0 {
		updatedConfig.Plugins = models.Plugins(request.EnabledPlugins)
	}

	// Update plugin configurations if provided
	if request.Plugins != nil {
		if request.Plugins.Auth != nil {
			updatedConfig.Auth = sp.mergeAuthConfigs(currentConfig.Auth, request.Plugins.Auth)
		}
		if request.Plugins.Kubernetes != nil {
			updatedConfig.Kubernetes = sp.mergeKubernetesConfigs(currentConfig.Kubernetes, request.Plugins.Kubernetes)
		}
		if request.Plugins.AWS != nil {
			updatedConfig.AWS = sp.mergeAWSConfigs(currentConfig.AWS, request.Plugins.AWS)
		}
		if request.Plugins.ServiceCatalog != nil {
			updatedConfig.ServiceCatalog = request.Plugins.ServiceCatalog
		}
		if request.Plugins.Observability != nil {
			updatedConfig.Observability = request.Plugins.Observability
		}
	}

	// Validate updated configuration
	_, err := sp.yamlProcessor.GenerateYAML(updatedConfig)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return updatedConfig, nil
}

// RequiresRestart checks if the configuration changes require a restart
func (sp *SettingsProcessor) RequiresRestart(
	oldConfig *models.DashConfig,
	newConfig *models.DashConfig,
) bool {
	if oldConfig == nil || newConfig == nil {
		return false
	}

	// Check if plugins changed
	if !sp.pluginsEqual(oldConfig.Plugins, newConfig.Plugins) {
		return true
	}

	// Check if port changed
	if oldConfig.Port != newConfig.Port {
		return true
	}

	// Check if auth configuration changed
	if !sp.authEqual(oldConfig.Auth, newConfig.Auth) {
		return true
	}

	return false
}

// pluginsEqual compares two plugin lists
func (sp *SettingsProcessor) pluginsEqual(a, b models.Plugins) bool {
	if a.Count() != b.Count() {
		return false
	}

	// Create maps for comparison
	mapA := make(map[string]bool)
	for _, p := range a {
		mapA[p] = true
	}

	mapB := make(map[string]bool)
	for _, p := range b {
		mapB[p] = true
	}

	// Compare maps
	if len(mapA) != len(mapB) {
		return false
	}

	for k := range mapA {
		if !mapB[k] {
			return false
		}
	}

	return true
}

// authEqual compares two auth configurations
func (sp *SettingsProcessor) authEqual(a, b []models.AuthConfig) bool {
	if len(a) != len(b) {
		return false
	}

	// Simplified comparison - in production, you'd want deeper comparison
	for i := range a {
		if a[i].Provider != b[i].Provider {
			return false
		}
	}

	return true
}

func (sp *SettingsProcessor) mergeAuthConfigs(
	current []models.AuthConfig,
	updates []models.UpdateAuthConfig,
) []models.AuthConfig {
	if updates == nil {
		return current
	}
	if len(updates) == 0 {
		return []models.AuthConfig{}
	}

	currentByProvider := make(map[string]models.AuthConfig, len(current))
	for _, cfg := range current {
		currentByProvider[strings.ToLower(cfg.Provider)] = cfg
	}

	result := make([]models.AuthConfig, 0, len(updates))
	for _, update := range updates {
		existing := currentByProvider[strings.ToLower(update.Provider)]
		merged := existing
		merged.Provider = update.Provider

		if update.ClientID != nil {
			merged.ClientID = strings.TrimSpace(*update.ClientID)
		}
		if update.OrgPermission != nil {
			merged.OrgPermission = strings.TrimSpace(*update.OrgPermission)
		}
		if update.RedirectURL != nil {
			merged.RedirectURL = strings.TrimSpace(*update.RedirectURL)
		}
		if len(update.Scopes) > 0 {
			merged.Scopes = append([]string(nil), update.Scopes...)
		}

		if update.ClearClientSecret {
			merged.ClientSecret = ""
		} else if update.ClientSecret != nil {
			merged.ClientSecret = strings.TrimSpace(*update.ClientSecret)
		}

		result = append(result, merged)
	}

	return result
}

func (sp *SettingsProcessor) mergeKubernetesConfigs(
	current []models.KubernetesConfig,
	updates []models.UpdateKubernetesConfig,
) []models.KubernetesConfig {
	if updates == nil {
		return current
	}
	if len(updates) == 0 {
		return []models.KubernetesConfig{}
	}

	currentByName := make(map[string]models.KubernetesConfig, len(current))
	for _, cfg := range current {
		currentByName[strings.ToLower(cfg.Name)] = cfg
	}

	result := make([]models.KubernetesConfig, 0, len(updates))
	for _, update := range updates {
		existing := currentByName[strings.ToLower(update.Name)]
		merged := existing
		merged.Name = update.Name

		if update.ConnectionType != nil {
			merged.ConnectionType = strings.TrimSpace(*update.ConnectionType)
		}
		if update.Kubeconfig != nil {
			merged.Kubeconfig = *update.Kubeconfig
		}
		if update.Context != nil {
			merged.Context = strings.TrimSpace(*update.Context)
		}
		if update.Host != nil {
			merged.Host = strings.TrimSpace(*update.Host)
		}
		if update.Certificate != nil {
			merged.Certificate = *update.Certificate
		}

		if update.ClearToken {
			merged.Token = ""
		} else if update.Token != nil {
			merged.Token = *update.Token
		}

		result = append(result, merged)
	}

	return result
}

func (sp *SettingsProcessor) mergeAWSConfigs(
	current []models.AWSConfig,
	updates []models.UpdateAWSConfig,
) []models.AWSConfig {
	if updates == nil {
		return current
	}
	if len(updates) == 0 {
		return []models.AWSConfig{}
	}
	currentByName := make(map[string]models.AWSConfig, len(current))
	for _, cfg := range current {
		currentByName[strings.ToLower(cfg.Name)] = cfg
	}

	result := make([]models.AWSConfig, 0, len(updates))
	for _, update := range updates {
		existing := currentByName[strings.ToLower(update.Name)]
		merged := existing
		merged.Name = update.Name

		if update.Region != nil {
			merged.Region = strings.TrimSpace(*update.Region)
		}
		if update.AccessKeyID != nil {
			merged.AccessKeyID = strings.TrimSpace(*update.AccessKeyID)
		}

		if update.ClearSecretAccessKey {
			merged.SecretAccessKey = ""
		} else if update.SecretAccessKey != nil {
			merged.SecretAccessKey = strings.TrimSpace(*update.SecretAccessKey)
		}

		result = append(result, merged)
	}

	return result
}
