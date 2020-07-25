package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// DashYaml dash ops config
type DashYaml struct {
	Port    string   `yaml:"port"`
	Origin  string   `yaml:"origin"`
	Headers []string `yaml:"headers"`
	Front   string   `yaml:"front"`
	Plugins Plugins  `yaml:"plugins"`
}

// Plugins ...
type Plugins []string

// Has ...
func (list Plugins) Has(a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
		log.Fatalln("reading file config", err)
	}

	return []byte(os.ExpandEnv(string(file)))
}

// GetGlobalConfig get dash ops global config
func GetGlobalConfig(file []byte) DashYaml {
	dc := DashYaml{}
	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		log.Fatalln("parse yaml config", err)
	}

	return dc
}
