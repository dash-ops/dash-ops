package oauth2

import (
	"context"
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	"github.com/google/go-github/github"
	mux_context "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

func meHandler(oauthConfig *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := mux_context.Get(r, TokenKey).(*oauth2.Token)
		client := github.NewClient(oauthConfig.Client(context.Background(), token))
		user, _, err := client.Users.Get(context.Background(), "")
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, *user)
	}
}

func makeGithubHandlers(r *mux.Router, oauthConfig *oauth2.Config) {
	r.HandleFunc("/v1/me", meHandler(oauthConfig)).
		Methods("GET", "OPTIONS").
		Name("userGithub")
}
