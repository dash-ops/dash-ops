package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/dash-ops/dash-ops/pkg/config"
	"github.com/dash-ops/dash-ops/pkg/spa"
	oauth "github.com/dash-ops/dash-ops/pkg/oauth2"
)

func main() {
	dc := config.GetGlobalConfig()

	router := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedHeaders(dc.Config.Headers),
		handlers.AllowedOrigins([]string{dc.Config.Origin}),
		handlers.AllowCredentials(),
	)
	router.Use(cors)

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// OAuth API
	oauth.MakeOauthHandlers(router)
	private := router.PathPrefix("/api").Subrouter()
	private.Use(oauth.OAuthMiddleware)

	spaHandler := spa.SpaHandler{StaticPath: "front/build", IndexPath: "index.html"}
	router.PathPrefix("/").Handler(spaHandler)

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("localhost:%s", dc.Config.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
