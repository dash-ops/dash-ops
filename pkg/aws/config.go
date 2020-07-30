package aws

import (
	"gopkg.in/yaml.v2"
)

type dashYaml struct {
	AWS []config `yaml:"aws"`
}

type config struct {
	Name            string     `yaml:"name"`
	Region          string     `yaml:"region"`
	AccessKeyID     string     `yaml:"accessKeyId"`
	SecretAccessKey string     `yaml:"secretAccessKey"`
	Permission      permission `yaml:"permission"`
	EC2Config       ec2Config  `yaml:"ec2Config"`
}

type permission struct {
	EC2 ec2Permissions `yaml:"ec2" json:"ec2"`
}

type ec2Permissions struct {
	Start []string `yaml:"start" json:"start"`
	Stop  []string `yaml:"stop" json:"stop"`
}

type ec2Config struct {
	SkipList []string `yaml:"skipList"`
}

func loadConfig(file []byte) (dashYaml, error) {
	dc := dashYaml{}

	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		return dashYaml{}, err
	}
	return dc, nil
}
