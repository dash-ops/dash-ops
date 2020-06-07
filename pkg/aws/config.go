package aws

import (
	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

type dashYaml struct {
	AWS awsConfig `yaml:"aws"`
}

type awsConfig struct {
	Region          string `yaml:"region"`
	AccessKeyID     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	EC2Config       struct {
		Blacklist []string `yaml:"blacklist"`
	} `yaml:"ec2Config"`
}

func loadConfig(file []byte) dashYaml {
	dc := dashYaml{}

	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		log.WithError(err).Fatal("parse yaml config")
	}

	return dc
}
