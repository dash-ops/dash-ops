package oauth2

import (
	"log"

	"gopkg.in/yaml.v2"
)

type dashYaml struct {
	Oauth2 []struct {
		Provider        string   `yaml:"provider"`
		ClientID        string   `yaml:"clientId"`
		ClientSecret    string   `yaml:"clientSecret"`
		AuthURL         string   `yaml:"authURL"`
		TokenURL        string   `yaml:"tokenURL"`
		RedirectURL     string   `yaml:"redirectURL"`
		URLLoginSuccess string   `yaml:"urlLoginSuccess"`
		OrgPermission   string   `yaml:"orgPermission"`
		Scopes          []string `yaml:"scopes"`
	} `yaml:"oauth2"`
}

func loadConfig(file []byte) dashYaml {
	dc := dashYaml{}

	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		log.Fatalln("parse yaml config", err)
	}

	return dc
}
