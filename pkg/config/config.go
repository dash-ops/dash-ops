package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

// DashYaml dash ops config
type DashYaml struct {
	Port    string   `yaml:"port"`
	Origin  string   `yaml:"origin"`
	Headers []string `yaml:"headers"`
	Front   string   `yaml:"front"`
	Plugins []string `yaml:"plugins"`
}

// GetFileGlobalConfig get dash ops file config
func GetFileGlobalConfig() []byte {
	var dashYaml = "./dash-ops.yaml"
	if path := os.Getenv("DASH_CONFIG"); path != "" {
		dashYaml = path
	}

	filename, _ := filepath.Abs(dashYaml)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithError(err).Fatal("reading file config")
	}

	return []byte(os.ExpandEnv(string(file)))
}

// GetGlobalConfig get dash ops global config
func GetGlobalConfig(file []byte) DashYaml {
	dc := DashYaml{}
	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		log.WithError(err).Fatal("parse yaml config")
	}

	return dc
}
