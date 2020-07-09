package oauth2

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestMakeGithubHandlers(t *testing.T) {
	r := mux.NewRouter()
	makeGithubHandlers(r, dashYaml{}, &oauth2.Config{})

	path, err := r.GetRoute("userLogger").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/me", path)
}
