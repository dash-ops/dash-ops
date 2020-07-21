package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	fileAwsConfig := []byte(`aws:
	- name: 'AWS Test'
		region: us-east-1
		accessKeyId: 666
		secretAccessKey: 999
		ec2Config:
			skipList:
				- "test"`)

	dashConfig := loadConfig(fileAwsConfig)

	assert.Equal(t, "AWS Test", dashConfig.AWS[0].Name)
	assert.Equal(t, "us-east-1", dashConfig.AWS[0].Region)
	assert.Equal(t, "666", dashConfig.AWS[0].AccessKeyID)
	assert.Equal(t, "999", dashConfig.AWS[0].SecretAccessKey)
}
