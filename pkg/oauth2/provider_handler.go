package oauth2

import (
	"fmt"
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	mux_context "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

func meHandler(githubClient GithubClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := mux_context.Get(r, TokenKey).(*oauth2.Token)
		user, err := githubClient.GetUserLogger(token)
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, user)
	}
}

func orgPermissionMiddleware(githubClient GithubClient, orgPermission string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := mux_context.Get(req, TokenKey).(*oauth2.Token)

			orgs, err := githubClient.GetOrgsUserLogger(token)
			if err != nil {
				commons.RespondError(w, http.StatusUnauthorized, "no organization found to validate your access permission, "+err.Error())
				return
			}

			for _, org := range orgs {
				if *org.Login == orgPermission {
					next.ServeHTTP(w, req)
					return
				}
			}

			commons.RespondError(w, http.StatusUnauthorized, "you need to be in organization "+orgPermission+" to have access")
			return
		})
	}
}

func makeGithubHandlers(r *mux.Router, dashConfig dashYaml, oauthConfig *oauth2.Config) {
	githubClient, err := NewGithubClient(oauthConfig)
	if err != nil {
		fmt.Println(err.Error())
	}

	r.HandleFunc("/me", meHandler(githubClient)).
		Methods("GET", "OPTIONS").
		Name("userLogger")

	if dashConfig.Oauth2[0].OrgPermission != "" {
		r.Use(orgPermissionMiddleware(githubClient, dashConfig.Oauth2[0].OrgPermission))
	}
}
