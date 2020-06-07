package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	fileK8sConfig := []byte(`kubernetes:
  kubeconfig: /root/.kube/config`)

	dashConfig := loadConfig(fileK8sConfig)

	assert.Equal(t, "/root/.kube/config", dashConfig.Kubernetes.Kubeconfig)
}
