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

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	oauth2.MakeOauthHandlers(router, fileConfig)
	private := router.PathPrefix("/api").Subrouter()
	private.Use(oauth2.OAuthMiddleware)

	kubernetes.MakeKubernetesHandlers(private, fileConfig)
	aws.MakeAWSInstanceHandlers(private, fileConfig)

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
