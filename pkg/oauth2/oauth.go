package oauth2

import (
	"context"
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	mux_context "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type key string

const TokenKey key = "token"

func oauthHandler(oauthConfig *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := oauthConfig.AuthCodeURL(r.URL.Query().Get("redirect_url"))
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}

func oauthRedirectHandler(dc dashYaml, oauthConfig *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := oauthConfig.Exchange(context.Background(), r.URL.Query().Get("code"))
		if err != nil {
			commons.RespondError(w, http.StatusUnauthorized, "there was an issue getting your token, "+err.Error())
			return
		}

		if !token.Valid() {
			commons.RespondError(w, http.StatusUnauthorized, "retrieved invalid token: "+err.Error())
			return
		}

		http.Redirect(w, r, dc.Oauth2[0].URLLoginSuccess+r.URL.Query().Get("state")+"?access_token="+token.AccessToken, http.StatusPermanentRedirect)
	}
}

// OAuthMiddleware should valid authentication
func OAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const bearerSchema = "Bearer "
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			commons.RespondError(w, http.StatusUnauthorized, "retrieved invalid token")
			return
		}
		accessToken := authHeader[len(bearerSchema):]
		if accessToken == "" {
			commons.RespondError(w, http.StatusUnauthorized, "retrieved invalid token")
			return
		}
		token := &oauth2.Token{AccessToken: accessToken, TokenType: "Bearer"}
		if !token.Valid() {
			commons.RespondError(w, http.StatusUnauthorized, "retrieved invalid token")
			return
		}
		mux_context.Set(r, TokenKey, token)
		next.ServeHTTP(w, r)
	})
}

// MakeOauthHandlers Add outh endpoints
func MakeOauthHandlers(r *mux.Router, internal *mux.Router, fileConfig []byte) {
	dashConfig := loadConfig(fileConfig)

	oauthConfig := &oauth2.Config{
		ClientID:     dashConfig.Oauth2[0].ClientID,
		ClientSecret: dashConfig.Oauth2[0].ClientSecret,
		Scopes:       dashConfig.Oauth2[0].Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  dashConfig.Oauth2[0].AuthURL,
			TokenURL: dashConfig.Oauth2[0].TokenURL,
		},
	}

	r.HandleFunc("/oauth", oauthHandler(oauthConfig)).
		Name("oauth")
	r.HandleFunc("/oauth/redirect", oauthRedirectHandler(dashConfig, oauthConfig)).
		Name("oauthRedirect")

	if dashConfig.Oauth2[0].Provider == "github" {
		makeGithubHandlers(internal, oauthConfig)
	}
}
