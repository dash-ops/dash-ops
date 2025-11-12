package http

import (
	"strings"

	settingsModels "github.com/dash-ops/dash-ops/pkg/settings/models"
	settingsWire "github.com/dash-ops/dash-ops/pkg/settings/wire"
)

// SettingsAdapter handles transformation between wire and models
type SettingsAdapter struct{}

// NewSettingsAdapter creates a new settings adapter
func NewSettingsAdapter() *SettingsAdapter {
	return &SettingsAdapter{}
}

// RequestToSetupConfig converts wire request to setup config model
func (sa *SettingsAdapter) RequestToSetupConfig(req *settingsWire.SetupConfigureRequest) *settingsModels.SetupConfig {
	if req == nil {
		return nil
	}

	setup := &settingsModels.SetupConfig{
		Port:    req.Config.Port,
		Origin:  req.Config.Origin,
		Headers: append([]string(nil), req.Config.Headers...),
		Front:   req.Config.Front,
		Plugins: settingsModels.SetupPluginsConfig{
			EnabledPlugins: append([]string(nil), req.EnabledPlugins...),
		},
	}

	if req.Plugins.Auth != nil {
		setup.Plugins.Auth = sa.toAuthConfig(*req.Plugins.Auth)
	}
	if len(req.Plugins.Kubernetes) > 0 {
		setup.Plugins.Kubernetes = sa.toKubernetesConfigs(req.Plugins.Kubernetes)
	}
	if len(req.Plugins.AWS) > 0 {
		setup.Plugins.AWS = sa.toAWSConfigs(req.Plugins.AWS)
	}
	if req.Plugins.ServiceCatalog != nil {
		setup.Plugins.ServiceCatalog = sa.toServiceCatalogConfig(*req.Plugins.ServiceCatalog)
	}
	if req.Plugins.Observability != nil {
		setup.Plugins.Observability = sa.toObservabilityConfig(*req.Plugins.Observability)
	}

	return setup
}

// SetupStatusToResponse converts setup status model to wire response
func (sa *SettingsAdapter) SetupStatusToResponse(status *settingsModels.SetupStatus) *settingsWire.SetupStatusResponse {
	if status == nil {
		return &settingsWire.SetupStatusResponse{
			SetupRequired: true,
			PluginsCount:  0,
			HasAuth:       false,
		}
	}

	return &settingsWire.SetupStatusResponse{
		SetupRequired: status.SetupRequired,
		PluginsCount:  status.PluginsCount,
		HasAuth:       status.HasAuth,
	}
}

// DashConfigToResponse converts DashConfig to wire response
func (sa *SettingsAdapter) DashConfigToResponse(config *settingsModels.DashConfig) *settingsWire.SetupConfigureResponse {
	if config == nil {
		return &settingsWire.SetupConfigureResponse{
			Success: false,
			Message: "configuration is nil",
		}
	}

	return &settingsWire.SetupConfigureResponse{
		Success:    true,
		Message:    "Configuration saved successfully",
		ConfigPath: "./dash-ops.yaml", // Default path
	}
}

// SettingsConfigToResponse converts SettingsConfig to wire response
func (sa *SettingsAdapter) SettingsConfigToResponse(config *settingsModels.SettingsConfig) *settingsWire.SettingsConfigResponse {
	if config == nil {
		return &settingsWire.SettingsConfigResponse{
			Config:  &settingsWire.DashConfigData{},
			Plugins: []string{},
			CanEdit: false,
		}
	}

	return &settingsWire.SettingsConfigResponse{
		Config:  sa.dashConfigToWire(config.Config),
		Plugins: config.Plugins,
		CanEdit: config.CanEdit,
	}
}

// dashConfigToWire converts DashConfig model to wire format
func (sa *SettingsAdapter) dashConfigToWire(config *settingsModels.DashConfig) *settingsWire.DashConfigData {
	if config == nil {
		return &settingsWire.DashConfigData{}
	}

	wireConfig := &settingsWire.DashConfigData{
		Port:    config.Port,
		Origin:  config.Origin,
		Headers: append([]string(nil), config.Headers...),
		Front:   config.GetFront(),
		Plugins: config.Plugins.List(),
	}

	if len(config.Auth) > 0 {
		wireConfig.Auth = make([]settingsWire.AuthProviderData, 0, len(config.Auth))
		for _, provider := range config.Auth {
			wireConfig.Auth = append(wireConfig.Auth, settingsWire.AuthProviderData{
				Provider:        provider.Provider,
				ClientIDMasked:  maskValue(provider.ClientID),
				OrgPermission:   provider.OrgPermission,
				RedirectURL:     provider.RedirectURL,
				Scopes:          append([]string(nil), provider.Scopes...),
				HasClientSecret: strings.TrimSpace(provider.ClientSecret) != "",
			})
		}
	}

	if len(config.Kubernetes) > 0 {
		wireConfig.Kubernetes = make([]settingsWire.KubernetesClusterData, 0, len(config.Kubernetes))
		for _, cluster := range config.Kubernetes {
			wireConfig.Kubernetes = append(wireConfig.Kubernetes, settingsWire.KubernetesClusterData{
				Name:           cluster.Name,
				ConnectionType: cluster.ConnectionType,
				Kubeconfig:     cluster.Kubeconfig,
				Context:        cluster.Context,
				Host:           cluster.Host,
				HasToken:       strings.TrimSpace(cluster.Token) != "",
			})
		}
	}

	if len(config.AWS) > 0 {
		wireConfig.AWS = make([]settingsWire.AWSAccountData, 0, len(config.AWS))
		for _, account := range config.AWS {
			wireConfig.AWS = append(wireConfig.AWS, settingsWire.AWSAccountData{
				Name:               account.Name,
				Region:             account.Region,
				AccessKeyIDMasked:  maskValue(account.AccessKeyID),
				HasSecretAccessKey: strings.TrimSpace(account.SecretAccessKey) != "",
			})
		}
	}

	if config.ServiceCatalog != nil {
		wireConfig.ServiceCatalog = sa.serviceCatalogToWire(config.ServiceCatalog)
	}

	if config.Observability != nil {
		wireConfig.Observability = sa.observabilityToWire(config.Observability)
	}

	return wireConfig
}

// RequestToUpdateSettingsRequest converts wire request to update settings model
func (sa *SettingsAdapter) RequestToUpdateSettingsRequest(req *settingsWire.UpdateSettingsRequest) *settingsModels.UpdateSettingsRequest {
	if req == nil {
		return nil
	}

	update := &settingsModels.UpdateSettingsRequest{
		EnabledPlugins: append([]string(nil), req.EnabledPlugins...),
	}

	if req.Config != nil {
		update.Config = &settingsModels.UpdateGeneralConfig{
			Port:    req.Config.Port,
			Origin:  req.Config.Origin,
			Headers: append([]string(nil), req.Config.Headers...),
			Front:   req.Config.Front,
		}
	}

	if req.Plugins != nil {
		update.Plugins = &settingsModels.UpdatePlugins{
			Auth:           sa.toUpdateAuthConfigs(req.Plugins.Auth),
			Kubernetes:     sa.toUpdateKubernetesConfigs(req.Plugins.Kubernetes),
			AWS:            sa.toUpdateAWSConfigs(req.Plugins.AWS),
			ServiceCatalog: sa.toUpdateServiceCatalogConfig(req.Plugins.ServiceCatalog),
			Observability:  sa.toUpdateObservabilityConfig(req.Plugins.Observability),
		}
	}

	return update
}

// UpdateSettingsResponseToWire converts update settings response model to wire
func (sa *SettingsAdapter) UpdateSettingsResponseToWire(resp *settingsModels.UpdateSettingsResponse) *settingsWire.UpdateSettingsResponse {
	if resp == nil {
		return &settingsWire.UpdateSettingsResponse{
			Success: false,
			Message: "response is nil",
		}
	}

	return &settingsWire.UpdateSettingsResponse{
		Success:         resp.Success,
		Message:         resp.Message,
		RequiresRestart: resp.RequiresRestart,
	}
}

func (sa *SettingsAdapter) toAuthConfig(req settingsWire.AuthProviderRequest) *settingsModels.AuthConfig {
	return &settingsModels.AuthConfig{
		Provider:      req.Provider,
		ClientID:      req.ClientID,
		ClientSecret:  req.ClientSecret,
		OrgPermission: req.OrgPermission,
		RedirectURL:   req.RedirectURL,
		Scopes:        append([]string(nil), req.Scopes...),
	}
}

func (sa *SettingsAdapter) toAWSConfigs(req []settingsWire.AWSAccountRequest) []settingsModels.AWSConfig {
	if len(req) == 0 {
		return nil
	}

	result := make([]settingsModels.AWSConfig, 0, len(req))
	for _, account := range req {
		result = append(result, settingsModels.AWSConfig{
			Name:            account.Name,
			Region:          account.Region,
			AccessKeyID:     account.AccessKeyID,
			SecretAccessKey: account.SecretAccessKey,
		})
	}
	return result
}

func (sa *SettingsAdapter) toKubernetesConfigs(req []settingsWire.KubernetesClusterRequest) []settingsModels.KubernetesConfig {
	if len(req) == 0 {
		return nil
	}

	result := make([]settingsModels.KubernetesConfig, 0, len(req))
	for _, cluster := range req {
		result = append(result, settingsModels.KubernetesConfig{
			Name:           cluster.Name,
			ConnectionType: cluster.ConnectionType,
			Kubeconfig:     cluster.Kubeconfig,
			Context:        cluster.Context,
			Host:           cluster.Host,
			Token:          cluster.Token,
		})
	}
	return result
}

func (sa *SettingsAdapter) toServiceCatalogConfig(req settingsWire.ServiceCatalogRequest) *settingsModels.ServiceCatalogConfig {
	storage := settingsModels.StorageConfig{
		Provider: req.Storage.Provider,
	}

	provider := strings.ToLower(req.Storage.Provider)

	if req.Storage.Filesystem != nil {
		storage.Filesystem = &settingsModels.FilesystemStorage{
			Directory: req.Storage.Filesystem.Directory,
		}
	} else if provider == "filesystem" && strings.TrimSpace(req.Storage.Directory) != "" {
		storage.Filesystem = &settingsModels.FilesystemStorage{
			Directory: req.Storage.Directory,
		}
	}

	if req.Storage.GitHub != nil {
		storage.GitHub = &settingsModels.GitHubStorage{
			Repository: req.Storage.GitHub.Repository,
			Branch:     req.Storage.GitHub.Branch,
		}
	} else if provider == "github" && strings.TrimSpace(req.Storage.Repository) != "" {
		storage.GitHub = &settingsModels.GitHubStorage{
			Repository: req.Storage.Repository,
			Branch:     req.Storage.Branch,
		}
	}

	if req.Storage.S3 != nil {
		storage.S3 = &settingsModels.S3Storage{
			Bucket: req.Storage.S3.Bucket,
		}
	}

	serviceCatalog := &settingsModels.ServiceCatalogConfig{
		Storage: storage,
	}

	if req.Versioning != nil {
		serviceCatalog.Versioning = &settingsModels.VersioningConfig{
			Enabled:  req.Versioning.Enabled,
			Provider: req.Versioning.Provider,
		}
	}

	return serviceCatalog
}

func (sa *SettingsAdapter) toObservabilityConfig(req settingsWire.ObservabilityRequest) *settingsModels.ObservabilityConfig {
	config := &settingsModels.ObservabilityConfig{
		Enabled: req.Enabled,
	}

	if len(req.Logs) > 0 {
		config.Logs = &settingsModels.ObservabilityLogs{
			Providers: sa.toObservabilityProviders(req.Logs),
		}
	}

	if len(req.Traces) > 0 {
		config.Traces = &settingsModels.ObservabilityTraces{
			Providers: sa.toObservabilityProviders(req.Traces),
		}
	}

	if len(req.Metrics) > 0 {
		config.Metrics = &settingsModels.ObservabilityMetrics{
			Providers: sa.toObservabilityProviders(req.Metrics),
		}
	}

	return config
}

func (sa *SettingsAdapter) toObservabilityProviders(req []settingsWire.ObservabilityProviderRequest) []settingsModels.ObservabilityProvider {
	result := make([]settingsModels.ObservabilityProvider, 0, len(req))
	for _, provider := range req {
		model := settingsModels.ObservabilityProvider{
			Name:      provider.Name,
			Type:      provider.Type,
			URL:       provider.URL,
			Timeout:   provider.Timeout,
			Retention: provider.Retention,
			Enabled:   provider.Enabled,
			Labels:    provider.Labels,
		}

		if provider.Auth != nil {
			model.Auth = &settingsModels.ProviderAuth{
				Type:     provider.Auth.Type,
				Username: valueOrEmpty(provider.Auth.Username),
				Password: valueOrEmpty(provider.Auth.Password),
				Token:    valueOrEmpty(provider.Auth.Token),
			}
		}

		result = append(result, model)
	}
	return result
}

func (sa *SettingsAdapter) serviceCatalogToWire(config *settingsModels.ServiceCatalogConfig) *settingsWire.ServiceCatalogData {
	if config == nil {
		return nil
	}

	wireData := &settingsWire.ServiceCatalogData{
		Storage: settingsWire.ServiceCatalogStorageData{
			Provider: config.Storage.Provider,
		},
	}

	switch strings.ToLower(config.Storage.Provider) {
	case "filesystem":
		if config.Storage.Filesystem != nil {
			wireData.Storage.Directory = config.Storage.Filesystem.Directory
			wireData.Storage.Filesystem = &settingsWire.FilesystemStorageData{
				Directory: config.Storage.Filesystem.Directory,
			}
		}
	case "github":
		if config.Storage.GitHub != nil {
			wireData.Storage.Repository = config.Storage.GitHub.Repository
			wireData.Storage.Branch = config.Storage.GitHub.Branch
			wireData.Storage.GitHub = &settingsWire.GitHubStorageData{
				Repository: config.Storage.GitHub.Repository,
				Branch:     config.Storage.GitHub.Branch,
			}
		}
	case "s3":
		if config.Storage.S3 != nil {
			wireData.Storage.S3 = &settingsWire.S3StorageData{
				Bucket: config.Storage.S3.Bucket,
			}
		}
	}

	if config.Versioning != nil {
		wireData.Versioning = &settingsWire.ServiceCatalogVersioningData{
			Enabled:  config.Versioning.Enabled,
			Provider: config.Versioning.Provider,
		}
	}

	return wireData
}

func (sa *SettingsAdapter) observabilityToWire(config *settingsModels.ObservabilityConfig) *settingsWire.ObservabilityData {
	if config == nil {
		return nil
	}

	wireData := &settingsWire.ObservabilityData{
		Enabled: config.Enabled,
	}

	if config.Logs != nil {
		wireData.Logs = sa.observabilityProvidersToWire(config.Logs.Providers)
	}
	if config.Traces != nil {
		wireData.Traces = sa.observabilityProvidersToWire(config.Traces.Providers)
	}
	if config.Metrics != nil {
		wireData.Metrics = sa.observabilityProvidersToWire(config.Metrics.Providers)
	}

	return wireData
}

func (sa *SettingsAdapter) observabilityProvidersToWire(providers []settingsModels.ObservabilityProvider) []settingsWire.ObservabilityProviderData {
	if len(providers) == 0 {
		return nil
	}

	result := make([]settingsWire.ObservabilityProviderData, 0, len(providers))
	for _, provider := range providers {
		wireProvider := settingsWire.ObservabilityProviderData{
			Name:      provider.Name,
			Type:      provider.Type,
			URL:       provider.URL,
			Timeout:   provider.Timeout,
			Retention: provider.Retention,
			Enabled:   provider.Enabled,
			Labels:    provider.Labels,
		}

		if provider.Auth != nil {
			wireProvider.Auth = &settingsWire.ObservabilityProviderAuth{
				Type:        provider.Auth.Type,
				HasUsername: strings.TrimSpace(provider.Auth.Username) != "",
				HasPassword: strings.TrimSpace(provider.Auth.Password) != "",
				HasToken:    strings.TrimSpace(provider.Auth.Token) != "",
			}
		}

		result = append(result, wireProvider)
	}

	return result
}

func (sa *SettingsAdapter) toUpdateAuthConfigs(req []settingsWire.UpdateAuthProviderRequest) []settingsModels.UpdateAuthConfig {
	if len(req) == 0 {
		return nil
	}

	result := make([]settingsModels.UpdateAuthConfig, 0, len(req))
	for _, provider := range req {
		result = append(result, settingsModels.UpdateAuthConfig{
			Provider:          provider.Provider,
			ClientID:          provider.ClientID,
			ClientSecret:      provider.ClientSecret,
			OrgPermission:     provider.OrgPermission,
			RedirectURL:       provider.RedirectURL,
			Scopes:            append([]string(nil), provider.Scopes...),
			ClearClientSecret: provider.ClearClientSecret,
		})
	}
	return result
}

func (sa *SettingsAdapter) toUpdateKubernetesConfigs(req []settingsWire.UpdateKubernetesClusterRequest) []settingsModels.UpdateKubernetesConfig {
	if len(req) == 0 {
		return nil
	}

	result := make([]settingsModels.UpdateKubernetesConfig, 0, len(req))
	for _, cluster := range req {
		result = append(result, settingsModels.UpdateKubernetesConfig{
			Name:           cluster.Name,
			ConnectionType: cluster.ConnectionType,
			Kubeconfig:     cluster.Kubeconfig,
			Context:        cluster.Context,
			Host:           cluster.Host,
			Token:          cluster.Token,
			Certificate:    cluster.Certificate,
			ClearToken:     cluster.ClearToken,
		})
	}
	return result
}

func (sa *SettingsAdapter) toUpdateAWSConfigs(req []settingsWire.UpdateAWSAccountRequest) []settingsModels.UpdateAWSConfig {
	if len(req) == 0 {
		return nil
	}

	result := make([]settingsModels.UpdateAWSConfig, 0, len(req))
	for _, account := range req {
		result = append(result, settingsModels.UpdateAWSConfig{
			Name:                 account.Name,
			Region:               account.Region,
			AccessKeyID:          account.AccessKeyID,
			SecretAccessKey:      account.SecretAccessKey,
			ClearSecretAccessKey: account.ClearSecretAccessKey,
		})
	}
	return result
}

func (sa *SettingsAdapter) toUpdateServiceCatalogConfig(req *settingsWire.UpdateServiceCatalogRequest) *settingsModels.ServiceCatalogConfig {
	if req == nil || req.Storage == nil {
		return nil
	}

	config := sa.toServiceCatalogConfig(settingsWire.ServiceCatalogRequest{
		Storage:    *req.Storage,
		Versioning: req.Versioning,
	})
	return config
}

func (sa *SettingsAdapter) toUpdateObservabilityConfig(req *settingsWire.UpdateObservabilityRequest) *settingsModels.ObservabilityConfig {
	if req == nil {
		return nil
	}

	config := &settingsModels.ObservabilityConfig{}
	if req.Enabled != nil {
		config.Enabled = *req.Enabled
	}

	if len(req.Logs) > 0 {
		config.Logs = &settingsModels.ObservabilityLogs{
			Providers: sa.toObservabilityProviders(req.Logs),
		}
	}

	if len(req.Traces) > 0 {
		config.Traces = &settingsModels.ObservabilityTraces{
			Providers: sa.toObservabilityProviders(req.Traces),
		}
	}

	if len(req.Metrics) > 0 {
		config.Metrics = &settingsModels.ObservabilityMetrics{
			Providers: sa.toObservabilityProviders(req.Metrics),
		}
	}

	return config
}

func maskValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	const maskPrefix = "****"
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}

	return maskPrefix + value[len(value)-4:]
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
