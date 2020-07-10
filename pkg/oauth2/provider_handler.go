package oauth2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

func meHandler(githubClient GithubClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value(commons.TokenKey).(*oauth2.Token)
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
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Context().Value(commons.TokenKey).(*oauth2.Token)

			userData := commons.UserData{Org: orgPermission}
			teams, err := githubClient.GetTeamsUserLogger(token)
			if err != nil {
				commons.RespondError(w, http.StatusUnauthorized, "no organization found to validate your access permission, "+err.Error())
				return
			}

			for _, team := range teams {
				if *team.Organization.Login == orgPermission {
					userData.Groups = append(userData.Groups, fmt.Sprintf("%s%s%s", userData.Org, "*", *team.Slug))
				}
			}

			ctx := context.WithValue(r.Context(), commons.UserDataKey, userData)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
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
