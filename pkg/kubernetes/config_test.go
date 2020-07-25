package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	fileK8sConfig := []byte(`kubernetes:
  - name: Kubernetes Prod
    kubeconfig: /root/.kube/config
    context: k8s-prod`)

	dashConfig := loadConfig(fileK8sConfig)

	assert.Equal(t, "Kubernetes Prod", dashConfig.Kubernetes[0].Name)
	assert.Equal(t, "/root/.kube/config", dashConfig.Kubernetes[0].Kubeconfig)
	assert.Equal(t, "k8s-prod", dashConfig.Kubernetes[0].Context)
}
