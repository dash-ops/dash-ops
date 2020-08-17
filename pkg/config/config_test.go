package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGlobalConfig(t *testing.T) {
	fileConfig := []byte(`port: 8080
origin: http://localhost:3000
headers: 
  - "Content-Type"
  - "Authorization"
plugins:
  - "Oauth2"`)

	dashConfig := GetGlobalConfig(fileConfig)

	assert.Equal(t, "8080", dashConfig.Port)
	assert.Equal(t, "http://localhost:3000", dashConfig.Origin)
	assert.Equal(t, []string{"Content-Type", "Authorization"}, dashConfig.Headers)
}

func TestPluginHas(t *testing.T) {
	plugins := Plugins{"OAuth2", "Kubernetes"}

	assert.True(t, plugins.Has("Kubernetes"))
	assert.False(t, plugins.Has("AWS"))
}
