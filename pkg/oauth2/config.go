package oauth2

import (
	"os"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
	"github.com/apex/log"
	"golang.org/x/oauth2"
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

func loadConfig() (DashYaml, *oauth2.Config) {
	dc := DashYaml{}

	var dashYaml = "./dash-ops.yaml"
	if path := os.Getenv("DASH_CONFIG"); path != "" {
		dashYaml = path
	}

	filename, _ := filepath.Abs(dashYaml)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithError(err).Fatal("reading file config")
	}

	err = yaml.Unmarshal(file, &dc)
	if err != nil {
		log.WithError(err).Fatal("parse yaml config")
	}

	oauthConfig := &oauth2.Config{
		ClientID:     dc.Oauth2[0].ClientID,
		ClientSecret: dc.Oauth2[0].ClientSecret,
		Scopes:       dc.Oauth2[0].Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  dc.Oauth2[0].AuthURL,
			TokenURL: dc.Oauth2[0].TokenURL,
		},
	}

	return dc, oauthConfig
}
