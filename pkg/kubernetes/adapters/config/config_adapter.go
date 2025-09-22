package config

import (
	"fmt"

	k8sModels "github.com/dash-ops/dash-ops/pkg/kubernetes/models"
	"gopkg.in/yaml.v2"
)

// ConfigAdapter handles kubernetes configuration parsing
type ConfigAdapter struct{}

// NewConfigAdapter creates a new config adapter
func NewConfigAdapter() *ConfigAdapter {
	return &ConfigAdapter{}
}

// ParseKubernetesConfigFromFileConfig parses kubernetes config from file bytes
func (ca *ConfigAdapter) ParseKubernetesConfigFromFileConfig(fileConfig []byte) ([]k8sModels.KubernetesConfig, error) {
	var dashYaml struct {
		Kubernetes []k8sModels.KubernetesConfig `yaml:"kubernetes"`
	}

	err := yaml.Unmarshal(fileConfig, &dashYaml)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubernetes config: %w", err)
	}

	if len(dashYaml.Kubernetes) == 0 {
		return nil, fmt.Errorf("no kubernetes configuration found")
	}

	// Return all kubernetes configs
	return dashYaml.Kubernetes, nil
}

// ParseModuleConfig parses the complete module configuration
func (ca *ConfigAdapter) ParseModuleConfig(fileConfig []byte) (*k8sModels.ModuleConfig, error) {
	configs, err := ca.ParseKubernetesConfigFromFileConfig(fileConfig)
	if err != nil {
		return nil, err
	}

	return &k8sModels.ModuleConfig{
		Configs: configs,
	}, nil
}
