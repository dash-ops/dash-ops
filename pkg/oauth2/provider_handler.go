package oauth2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	gh "github.com/dash-ops/dash-ops/pkg/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type providerClients struct {
	github gh.Client
	google ggl.Client
}

// func (pc providerClients) getProvider(p string) (interface{}, error){
// 	if (p == "github") {
// 		return pc.github
// 	}
// }

// type User struct {
// 	ID   *int64  `json:"id,omitempty"`
// 	Name *string `json:"name,omitempty"`
// }

// type Team struct {
// 	ID           *int64        `json:"id,omitempty"`
// 	Name         *string       `json:"name,omitempty"`
// 	Slug         *string       `json:"slug,omitempty"`
// 	Organization *Organization `json:"organization,omitempty"`
// }

// type Organization struct {
// 	Login *string `json:"login,omitempty"`
// 	ID    *int64  `json:"id,omitempty"`
// 	Name  *string `json:"name,omitempty"`
// }

// func NewProviderClient(p string, oauthConfig *oauth2.Config) (ProviderClients, error) {
// 	if p == "github" {
// 		provider, err := gh.NewClient(oauthConfig)
// 		if err != nil {
// 			return ProviderClients{}, fmt.Errorf("failed to load github oauth provider")
// 		}

// 		return ProviderClients{github: provider}, nil
// 	}
// 	return ProviderClients{}, fmt.Errorf("failed to load oauth provider")
// }

func meHandler(dashConfig dashYaml, configs oauth2Configs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value(commons.TokenKey).(*oauth2.Token)
		
		oauth2Provider := r.Context().Value(commons.Oauth2ProviderKey)
		if (oauth2Provider == "github") {
			provider, err := gh.NewClient(configs.github)
			if err != nil {
				commons.RespondError(w, http.StatusInternalServerError, err.Error())
			}
			
			user, err := provider.GetUserLogger(token)
			if err != nil {
				commons.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			commons.RespondJSON(w, http.StatusOK, user)
		}
		// if (oauth2Provider == "google") {
		// 	provider, err := ggl.NewClient(configs.google)
		// 	if err != nil {
		// 		commons.RespondError(w, http.StatusInternalServerError, err.Error())
		// 	}
			
		// 	user, err := provider.GetUserLogger(token)
		// 	if err != nil {
		// 		commons.RespondError(w, http.StatusInternalServerError, err.Error())
		// 		return
		// 	}
		// }
	}
}

func orgPermissionMiddleware(dashConfig dashYaml, configs oauth2Configs) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Context().Value(commons.TokenKey).(*oauth2.Token)

			auth2Provider := r.Context().Value(commons.Oauth2ProviderKey)
		if (oauth2Provider == "github") {
			provider, err := gh.NewClient(configs.github)
			if err != nil {
				commons.RespondError(w, http.StatusInternalServerError, err.Error())
			}
			
			userData := commons.UserData{Org: orgPermission}
			teams, err := providerClients.github.GetTeamsUserLogger(token)
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

func makeOauthProvideHandlers(r *mux.Router, dashConfig dashYaml, configs oauth2Configs) {
	r.HandleFunc("/me", meHandler(dashConfig, configs)).
		Methods("GET", "OPTIONS").
		Name("userLogger")

	if dashConfig.Oauth2[0].OrgPermission != "" {
		r.Use(orgPermissionMiddleware(dashConfig, configs))
	}
}
