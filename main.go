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
	"github.com/dash-ops/dash-ops/pkg/kubernetes"
	servicecatalog "github.com/dash-ops/dash-ops/pkg/service-catalog"
	"github.com/dash-ops/dash-ops/pkg/spa"
	spaModels "github.com/dash-ops/dash-ops/pkg/spa/models"
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
	fileConfig := configModule.GetFileConfigBytes()

	internal := api.PathPrefix("/v1").Subrouter()

	// Phase 1: Initialize all modules
	modules := make(map[string]interface{})

	if dashConfig.Plugins.Has("Auth") {
		authModule, err := auth.NewModule(fileConfig)
		if err != nil {
			log.Fatalf("Failed to create auth module: %v", err)
		}
		modules["auth"] = authModule
		log.Println("Auth module initialized successfully")
	}

	if dashConfig.Plugins.Has("ServiceCatalog") {
		scModule, err := servicecatalog.NewModule(fileConfig)
		if err != nil {
			log.Fatalf("Failed to create service catalog module: %v", err)
		}
		modules["service-catalog"] = scModule
		log.Println("Service Catalog module initialized successfully")
	}

	if dashConfig.Plugins.Has("Kubernetes") {
		k8sModule, err := kubernetes.NewModule(fileConfig)
		if err != nil {
			log.Printf("Failed to create kubernetes module: %v", err)
		} else {
			modules["kubernetes"] = k8sModule
			log.Println("Kubernetes module initialized successfully")
		}
	}

	if dashConfig.Plugins.Has("AWS") {
		awsModule, err := aws.NewModule(fileConfig)
		if err != nil {
			log.Printf("Failed to create AWS module: %v", err)
		} else {
			modules["aws"] = awsModule
			log.Println("AWS module initialized successfully")
		}
	}

	// Phase 2: Register routes for all modules
	for name, module := range modules {
		if m, ok := module.(interface{ RegisterRoutes(*mux.Router) }); ok {
			m.RegisterRoutes(internal)
		} else if m, ok := module.(interface {
			RegisterRoutes(*mux.Router, *mux.Router)
		}); ok {
			m.RegisterRoutes(api, internal)
		}
		log.Printf("Routes registered for %s module", name)
	}

	// Phase 3: Load dependencies between modules
	for name, module := range modules {
		if m, ok := module.(interface {
			LoadDependencies(map[string]interface{}) error
		}); ok {
			if err := m.LoadDependencies(modules); err != nil {
				log.Printf("Warning: Failed to load dependencies for %s module: %v", name, err)
			} else {
				log.Printf("Dependencies loaded for %s module", name)
			}
		}
	}

	// Phase 4: Initialize SPA module (serves static files)
	spaConfig := &spaModels.SPAConfig{
		StaticPath: dashConfig.Front,
		IndexPath:  "index.html",
	}
	spaModule, err := spa.NewModule(spaConfig, api)
	if err != nil {
		log.Fatalf("Failed to create SPA module: %v", err)
	}
	log.Println("SPA module initialized successfully")

	// Register SPA routes with API middleware
	spaModule.RegisterRoutes(router)

	fmt.Println("DashOps server running!!")
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", dashConfig.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
