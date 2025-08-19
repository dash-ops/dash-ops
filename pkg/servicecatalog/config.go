package servicecatalog

import (
	"log"

	"gopkg.in/yaml.v2"
)

type dashYaml struct {
	ServiceCatalog []struct {
		Name        string `yaml:"name"`
		Storage     string `yaml:"storage"`     // file, database
		CatalogPath string `yaml:"catalogPath"` // for file storage
		Permission  struct {
			Read  []string `yaml:"read"`
			Write []string `yaml:"write"`
			Admin []string `yaml:"admin"`
		} `yaml:"permission"`
	} `yaml:"serviceCatalog"`
}

func loadConfig(file []byte) (dashYaml, error) {
	dc := dashYaml{}

	err := yaml.Unmarshal(file, &dc)
	if err != nil {
		log.Printf("Error parsing service catalog config: %v", err)
		return dc, err
	}

	return dc, nil
}
