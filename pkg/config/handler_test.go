package config

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestPermissionsHandler(t *testing.T) {
	mockPlugins := Plugins{"AWS", "Kubernetes"}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/config/plugins", nil)

	handler := configPluginsHandler(mockPlugins)
	handler.ServeHTTP(rr, req)

	var resultPlugins Plugins
	json.NewDecoder(rr.Body).Decode(&resultPlugins)

	assert.Equal(t, mockPlugins, resultPlugins, "return plugins")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestMakeAWSInstanceHandlers(t *testing.T) {
	dashConfig := DashYaml{
		Plugins: Plugins{"AWS", "Kubernetes"},
	}

	r := mux.NewRouter()
	MakeConfigHandlers(r, dashConfig)

	path, err := r.GetRoute("configPlugins").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/config/plugins", path)
}
