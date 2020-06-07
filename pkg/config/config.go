package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

type DashYaml struct {
	Port    string   `yaml:"port"`
	Origin  string   `yaml:"origin"`
	Headers []string `yaml:"headers"`
	Plugins []string `yaml:"plugins"`
}

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

	return file
}

func GetGlobalConfig(file []byte) DashYaml {
	dc := DashYaml{}
	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		log.WithError(err).Fatal("parse yaml config")
	}

	return dc
}
