package kubernetes

import (
	"gopkg.in/yaml.v2"
)

type dashYaml struct {
	Kubernetes []config `yaml:"kubernetes"`
}

type config struct {
	Name       string     `yaml:"name"`
	Kubeconfig string     `yaml:"kubeconfig"`
	Context    string     `yaml:"context"`
	Permission permission `yaml:"permission"`
	Listen     string     `yaml:"-"`
}

type permission struct {
	Deployments deploymentsPermissions `yaml:"deployments" json:"deployments"`
}

type deploymentsPermissions struct {
	Namespaces []string `yaml:"namespaces" json:"namespaces"`
	Restart    []string `yaml:"restart" json:"restart"`
	Scale      []string `yaml:"scale" json:"scale"`
}

func loadConfig(file []byte) (dashYaml, error) {
	dc := dashYaml{}

	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		return dashYaml{}, err
	}

	return dc, nil
}
