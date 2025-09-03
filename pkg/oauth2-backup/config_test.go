package oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	fileOauthConfig := []byte(`oauth2:
  - provider: github
    clientId: 999
    clientSecret: 666
    authURL: "https://github.com/login/oauth/authorize"
    tokenURL: "https://github.com/login/oauth/access_token"
    urlLoginSuccess: "http://localhost:3000"
    scopes: 
      - user
      - repo
      - read:org`)

	dashConfig := loadConfig(fileOauthConfig)

	assert.Equal(t, "github", dashConfig.Oauth2[0].Provider)
	assert.Equal(t, "999", dashConfig.Oauth2[0].ClientID)
	assert.Equal(t, "666", dashConfig.Oauth2[0].ClientSecret)
	assert.Equal(t, "https://github.com/login/oauth/authorize", dashConfig.Oauth2[0].AuthURL)
	assert.Equal(t, "https://github.com/login/oauth/access_token", dashConfig.Oauth2[0].TokenURL)
	assert.Equal(t, "http://localhost:3000", dashConfig.Oauth2[0].URLLoginSuccess)
	assert.Equal(t, []string{"user", "repo", "read:org"}, dashConfig.Oauth2[0].Scopes)
}
