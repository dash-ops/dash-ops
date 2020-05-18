package config

import (
	"os"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
	"github.com/apex/log"
)

type DashYaml struct {
	Port    string   `yaml:"port"`
	Origin  string   `yaml:"origin"`
	Headers []string `yaml:"headers"`
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
