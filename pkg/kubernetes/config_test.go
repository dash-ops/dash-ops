package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	fileOauthConfig := []byte(`kubernetes:
  kubeconfig: /root/.kube/config`)

	dashConfig := loadConfig(fileOauthConfig)

	assert.Equal(t, "/root/.kube/config", dashConfig.Kubernetes.Kubeconfig)
}
