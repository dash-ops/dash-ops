package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/dash-ops/dash-ops/pkg/auth"
	"github.com/dash-ops/dash-ops/pkg/aws"
	"github.com/dash-ops/dash-ops/pkg/config"
	"github.com/dash-ops/dash-ops/pkg/github"
	"github.com/dash-ops/dash-ops/pkg/kubernetes"
	servicecatalog "github.com/dash-ops/dash-ops/pkg/service-catalog"
	"github.com/dash-ops/dash-ops/pkg/spa"
	spaModels "github.com/dash-ops/dash-ops/pkg/spa/models"
	"golang.org/x/oauth2"
)

func main() {
	// Initialize config module
	configModule, err := config.NewModule("")
	if err != nil {
		log.Fatalf("Failed to initialize config module: %v", err)
	}

	dashConfig := configModule.GetConfig()

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

	// Register config routes using hexagonal architecture
	configModule.RegisterRoutes(api)

	internal := api.PathPrefix("/v1").Subrouter()

	// Initialize plugins in dependency order
	var serviceCatalogInstance *servicecatalog.ServiceCatalog

	// Initialize plugins with dependency injection
	if dashConfig.Plugins.Has("OAuth2") {
		// Parse auth config using hexagonal architecture
		fileConfig := configModule.GetFileConfigBytes()
		authConfig, err := auth.ParseAuthConfigFromFileConfig(fileConfig)
		if err != nil {
			log.Fatalf("Failed to parse auth config: %v", err)
		}

		// Create OAuth2 config for GitHub module dependency
		oauthConfig := &oauth2.Config{
			ClientID:     authConfig.ClientID,
			ClientSecret: authConfig.ClientSecret,
			Scopes:       authConfig.Scopes,
			RedirectURL:  authConfig.RedirectURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:  authConfig.AuthURL,
				TokenURL: authConfig.TokenURL,
			},
		}

		// Initialize GitHub module (dependency)
		githubModule, err := github.NewModule(oauthConfig)
		if err != nil {
			log.Fatalf("Failed to create GitHub module: %v", err)
		}

		// Initialize auth module with GitHub dependency injection
		authModule, err := auth.NewModule(authConfig, githubModule)
		if err != nil {
			log.Fatalf("Failed to create auth module: %v", err)
		}

		// Register auth routes using hexagonal architecture
		authModule.RegisterRoutes(api, internal)
	}

	if dashConfig.Plugins.Has("ServiceCatalog") {
		// Service Catalog plugin - initialize first to provide context resolver
		fileConfig := configModule.GetFileConfigBytes()
		serviceCatalogInstance = servicecatalog.MakeServiceCatalogHandlers(internal, fileConfig)
	}

	if dashConfig.Plugins.Has("Kubernetes") {
		// Kubernetes plugin with service context integration
		fileConfig := configModule.GetFileConfigBytes()
		if serviceCatalogInstance != nil {
			// Use service context resolver for enhanced integration
			resolver := serviceCatalogInstance.GetKubernetesAdapter()
			kubernetes.MakeKubernetesHandlersWithResolver(internal, fileConfig, resolver)
		} else {
			// Fallback to basic kubernetes handlers
			kubernetes.MakeKubernetesHandlers(internal, fileConfig)
		}
	}

	if dashConfig.Plugins.Has("AWS") {
		// ToDo transform into isolated plugins
		fileConfig := configModule.GetFileConfigBytes()
		aws.MakeAWSInstanceHandlers(internal, fileConfig)
	}

	// Initialize SPA module using hexagonal architecture
	spaConfig := &spaModels.SPAConfig{
		StaticPath: dashConfig.Front,
		IndexPath:  "index.html",
	}
	spaModule, err := spa.NewModule(spaConfig)
	if err != nil {
		log.Fatalf("Failed to create SPA module: %v", err)
	}

	// Use SPA handler from hexagonal architecture
	router.PathPrefix("/").Handler(spaModule.Handler)

	fmt.Println("DashOps server running!!")
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", dashConfig.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
