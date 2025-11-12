package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/dash-ops/dash-ops/pkg/auth"
	"github.com/dash-ops/dash-ops/pkg/aws"
	"github.com/dash-ops/dash-ops/pkg/kubernetes"
	"github.com/dash-ops/dash-ops/pkg/observability"
	servicecatalog "github.com/dash-ops/dash-ops/pkg/service-catalog"
	"github.com/dash-ops/dash-ops/pkg/settings"
	"github.com/dash-ops/dash-ops/pkg/spa"
	spaModels "github.com/dash-ops/dash-ops/pkg/spa/models"
)

func main() {
	settingsModule, err := settings.NewModule("")
	if err != nil {
		log.Fatalf("Failed to initialize settings module: %v", err)
	}

	dashConfig := settingsModule.GetConfig()
	if dashConfig == nil {
		log.Fatalf("Failed to load configuration: no configuration available")
	}

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

	settingsModule.RegisterRoutes(api)
	log.Println("Settings module initialized successfully")

	fileConfig := settingsModule.GetFileConfigBytes()

	internal := api.PathPrefix("/v1").Subrouter()

	// Phase 1: Initialize all modules
	modules := make(map[string]interface{})

	// Module factory registry - maps plugin names to their factory functions
	moduleFactories := map[string]func([]byte) (interface{}, error){
		"Auth":           func(config []byte) (interface{}, error) { return auth.NewModule(config) },
		"ServiceCatalog": func(config []byte) (interface{}, error) { return servicecatalog.NewModule(config) },
		"Kubernetes":     func(config []byte) (interface{}, error) { return kubernetes.NewModule(config) },
		"AWS":            func(config []byte) (interface{}, error) { return aws.NewModule(config) },
		"Observability":  func(config []byte) (interface{}, error) { return observability.NewModule(config) },
	}

	modulesLoaded := 0
	if fileConfig == nil {
		log.Println("No configuration file found; skipping plugin module initialization until setup is completed")
	} else {
		// Initialize modules dynamically based on active plugins
		for pluginName, factory := range moduleFactories {
			if dashConfig.Plugins.Has(pluginName) {
				module, err := factory(fileConfig)
				if err != nil {
					// All modules are optional - DashOps can run without any plugins
					log.Printf("Failed to create %s module: %v", pluginName, err)
					continue
				}

				// Use lowercase with dashes for consistency
				moduleKey := strings.ToLower(strings.ReplaceAll(pluginName, "Catalog", "-catalog"))
				modules[moduleKey] = module
				modulesLoaded++
				log.Printf("%s module initialized successfully", pluginName)
			}
		}
	}

	// Log summary of loaded modules
	if modulesLoaded == 0 {
		log.Println("No modules loaded - DashOps running in minimal mode")
	} else {
		log.Printf("Successfully loaded %d module(s)", modulesLoaded)
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
