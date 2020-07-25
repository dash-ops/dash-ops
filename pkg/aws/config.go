package aws

import (
	"log"

	"gopkg.in/yaml.v2"
)

type dashYaml struct {
	AWS []config `yaml:"aws"`
}

type config struct {
	Name            string        `yaml:"name"`
	Region          string        `yaml:"region"`
	AccessKeyID     string        `yaml:"accessKeyId"`
	SecretAccessKey string        `yaml:"secretAccessKey"`
	Permission      permission `yaml:"permission"`
	EC2Config       struct {
		SkipList []string `yaml:"skipList"`
	} `yaml:"ec2Config"`
}

type permission struct {
	EC2 struct {
		Start []string `yaml:"start" json:"start"`
		Stop  []string `yaml:"stop" json:"stop"`
	} `yaml:"ec2" json:"ec2"`
}

func loadConfig(file []byte) dashYaml {
	dc := dashYaml{}

	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		log.Fatalln("parse yaml config", err)
	}

	return dc
}
