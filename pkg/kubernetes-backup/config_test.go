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

	dashConfig, err := loadConfig(fileK8sConfig)

	assert.Nil(t, err)
	assert.Equal(t, "Kubernetes Prod", dashConfig.Kubernetes[0].Name)
	assert.Equal(t, "/root/.kube/config", dashConfig.Kubernetes[0].Kubeconfig)
	assert.Equal(t, "k8s-prod", dashConfig.Kubernetes[0].Context)
}

func TestLoadConfigWithYamlError(t *testing.T) {
	fileK8sConfig := []byte(`kubernetes:
  - name Kubernetes Prod`)

	_, err := loadConfig(fileK8sConfig)

	assert.Equal(t, "yaml: unmarshal errors:\n  line 2: cannot unmarshal !!str `name Ku...` into kubernetes.config", err.Error())
}
