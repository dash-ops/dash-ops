package kubernetes

import (
	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

type DashYaml struct {
	Kubernetes struct {
		Kubeconfig string `yaml:"kubeconfig"`
	} `yaml:"kubernetes"`
}

func loadConfig(file []byte) DashYaml {
	dc := DashYaml{}

	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		log.WithError(err).Fatal("parse yaml config")
	}

	return dc
}
