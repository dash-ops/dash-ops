package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	fileAwsConfig := []byte(`aws:
  region: us-east-1
  accessKeyId: 666
  secretAccessKey: 999
  ec2Config:
    skiplist:
      - "test"`)

	dashConfig := loadConfig(fileAwsConfig)

	assert.Equal(t, "us-east-1", dashConfig.AWS.Region)
	assert.Equal(t, "666", dashConfig.AWS.AccessKeyID)
	assert.Equal(t, "999", dashConfig.AWS.SecretAccessKey)
}
