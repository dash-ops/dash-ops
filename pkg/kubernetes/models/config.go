package models

// KubernetesConfig represents kubernetes configuration
type KubernetesConfig struct {
	Name       string     `yaml:"name"`
	Kubeconfig string     `yaml:"kubeconfig"`
	Context    string     `yaml:"context"`
	Permission Permission `yaml:"permission"`
	Listen     string     `yaml:"-"`
}

// Permission represents kubernetes permissions
type Permission struct {
	Deployments DeploymentsPermissions `yaml:"deployments" json:"deployments"`
}

// DeploymentsPermissions represents deployment permissions
type DeploymentsPermissions struct {
	Namespaces []string `yaml:"namespaces" json:"namespaces"`
	Restart    []string `yaml:"restart" json:"restart"`
	Scale      []string `yaml:"scale" json:"scale"`
}

// ModuleConfig represents the kubernetes module configuration
type ModuleConfig struct {
	Configs []KubernetesConfig `yaml:"kubernetes_configs" json:"kubernetes_configs"`
}
