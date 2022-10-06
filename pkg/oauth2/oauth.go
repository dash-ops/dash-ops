package oauth2

import (
	"context"
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type oauth2Configs struct {
	github *oauth2.Config
	google *oauth2.Config
}

func (configs oauth2Configs) getProvider(provider string) *oauth2.Config {
	if provider == "github" {
		return configs.github
	}
	if provider == "google" {
		return configs.google
	}
	return nil
}

func oauthHandler(configs oauth2Configs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		provider := configs.getProvider(vars["provider"])
		url := provider.AuthCodeURL(r.URL.Query().Get("redirect_url"))
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}

func oauthRedirectHandler(dc dashYaml, configs oauth2Configs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		provider := configs.getProvider(vars["provider"])
		token, err := provider.Exchange(context.Background(), r.URL.Query().Get("code"))
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

func oAuthMiddleware(next http.Handler) http.Handler {
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
		oauth2Provider := r.URL.Query().Get("provider")
		if oauth2Provider == "" {
			commons.RespondError(w, http.StatusUnauthorized, "retrieved invalid oauth2 provider")
			return
		}
		tokenCtx := context.WithValue(r.Context(), commons.TokenKey, token)
		providerCtx := context.WithValue(r.Context(), commons.Oauth2ProviderKey, oauth2Provider)
		r = r.WithContext(tokenCtx).WithContext(providerCtx)
		next.ServeHTTP(w, r)
	})
}

// MakeOauthHandlers Add auth endpoints
func MakeOauthHandlers(r *mux.Router, internal *mux.Router, fileConfig []byte) {
	dashConfig := loadConfig(fileConfig)

	var configs oauth2Configs
	for i := 0; i < len(dashConfig.Oauth2); i++ {
		if dashConfig.Oauth2[i].Provider == "github" {
			config := &oauth2.Config{
				ClientID:     dashConfig.Oauth2[i].ClientID,
				ClientSecret: dashConfig.Oauth2[i].ClientSecret,
				Scopes:       dashConfig.Oauth2[i].Scopes,
				Endpoint: oauth2.Endpoint{
					AuthURL:  dashConfig.Oauth2[i].AuthURL,
					TokenURL: dashConfig.Oauth2[i].TokenURL,
				},
			}
			configs.github = config
		}

		if dashConfig.Oauth2[i].Provider == "google" {
			config := &oauth2.Config{
				ClientID:     dashConfig.Oauth2[i].ClientID,
				ClientSecret: dashConfig.Oauth2[i].ClientSecret,
				Scopes:       dashConfig.Oauth2[i].Scopes,
				Endpoint:     google.Endpoint,
			}
			configs.google = config
		}
	}

	r.HandleFunc("/oauth/{provider}", oauthHandler(configs)).
		Name("oauth")
	r.HandleFunc("/oauth/{provider}/redirect", oauthRedirectHandler(dashConfig, configs)).
		Name("oauthRedirect")
	internal.Use(oAuthMiddleware)

	makeOauthProvideHandlers(internal, dashConfig, configs)
}
