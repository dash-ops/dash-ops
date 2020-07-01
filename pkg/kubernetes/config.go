package kubernetes

import (
	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

type dashYaml struct {
	Kubernetes []kubernetesConfig `yaml:"kubernetes"`
}

type kubernetesConfig struct {
	Name       string `yaml:"name"`
	Kubeconfig string `yaml:"kubeconfig"`
	Context    string `yaml:"context"`
	Listen     string `yaml:"-"`
}

func loadConfig(file []byte) dashYaml {
	dc := dashYaml{}

	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		log.WithError(err).Fatal("parse yaml config")
	}

	return dc
}
