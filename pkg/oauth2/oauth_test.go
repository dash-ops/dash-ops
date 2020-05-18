package oauth2

import (
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("DASH_CONFIG", "./test.yaml")
	os.Exit(m.Run())
}

func TestMakeOauthHandlers(t *testing.T) {
	r := mux.NewRouter()
	MakeOauthHandlers(r)

	path, err := r.GetRoute("oauth").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/api/oauth", path)

	path, err = r.GetRoute("oauthRedirect").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/oauth/redirect", path)
}
