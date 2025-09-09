package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	fileAwsConfig := []byte(`aws:
  - name: 'AWS Test'
    region: us-east-1
    accessKeyId: 1234
    secretAccessKey: 4321
    ec2Config:
      skipList:
        - "test"`)

	dashConfig, err := loadConfig(fileAwsConfig)

	assert.Nil(t, err)
	assert.Equal(t, "AWS Test", dashConfig.AWS[0].Name)
	assert.Equal(t, "us-east-1", dashConfig.AWS[0].Region)
	assert.Equal(t, "1234", dashConfig.AWS[0].AccessKeyID)
	assert.Equal(t, "4321", dashConfig.AWS[0].SecretAccessKey)
}

func TestLoadConfigWithYamlError(t *testing.T) {
	fileAwsConfig := []byte(`aws:
  - name 'AWS Test'`)

	_, err := loadConfig(fileAwsConfig)

	assert.Equal(t, "yaml: unmarshal errors:\n  line 2: cannot unmarshal !!str `name 'A...` into aws.config", err.Error())
}
