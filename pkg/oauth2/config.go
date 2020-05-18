package oauth2

import (
	"gopkg.in/yaml.v2"
	"github.com/apex/log"
)

type DashYaml struct {
	Oauth2 []struct {
		Provider        string   `yaml:"provider"`
		ClientID        string   `yaml:"clientId"`
		ClientSecret    string   `yaml:"clientSecret"`
		AuthURL         string   `yaml:"authURL"`
		TokenURL        string   `yaml:"tokenURL"`
		URLLoginSuccess string   `yaml:"urlLoginSuccess"`
		Scopes          []string `yaml:"scopes"`
	} `yaml:"oauth2"`
}

func loadConfig(file []byte) DashYaml {
	dc := DashYaml{}

	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		log.WithError(err).Fatal("parse yaml config")
	}

	return dc
}
