package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/dash-ops/dash-ops/pkg/aws"
	"github.com/dash-ops/dash-ops/pkg/config"
	"github.com/dash-ops/dash-ops/pkg/kubernetes"
	"github.com/dash-ops/dash-ops/pkg/oauth2"
	"github.com/dash-ops/dash-ops/pkg/spa"
)

func main() {
	fileConfig := config.GetFileGlobalConfig()
	dashConfig := config.GetGlobalConfig(fileConfig)

	router := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedHeaders(dashConfig.Headers),
		handlers.AllowedOrigins([]string{dashConfig.Origin}),
		handlers.AllowCredentials(),
	)
	router.Use(cors)

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	config.MakeConfigHandlers(api, dashConfig)

	internal := api.PathPrefix("/v1").Subrouter()
	if dashConfig.Plugins.Has("OAuth2") {
		// ToDo transform into isolated plugins
		oauth2.MakeOauthHandlers(api, internal, fileConfig)
	}
	if dashConfig.Plugins.Has("Kubernetes") {
		// ToDo transform into isolated plugins
		kubernetes.MakeKubernetesHandlers(internal, fileConfig)
	}
	if dashConfig.Plugins.Has("AWS") {
		// ToDo transform into isolated plugins
		aws.MakeAWSInstanceHandlers(internal, fileConfig)
	}

	spaHandler := spa.SpaHandler{StaticPath: dashConfig.Front, IndexPath: "index.html"}
	router.PathPrefix("/").Handler(spaHandler)

	fmt.Println("DashOps server running!!")
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", dashConfig.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
