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
	k8sExternal "github.com/dash-ops/dash-ops/pkg/kubernetes/adapters/external"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
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

	// Initialize plugins with dependency injection
	var serviceCatalogModule *servicecatalog.ServiceCatalog
	if dashConfig.Plugins.Has("Auth") {
		authModule, err := auth.NewModule(fileConfig)
		if err != nil {
			log.Fatalf("Failed to create auth module: %v", err)
		}

		// Register auth routes using hexagonal architecture
		authModule.RegisterRoutes(api, internal)
		log.Println("Auth module initialized successfully")
	}

	if dashConfig.Plugins.Has("ServiceCatalog") {
		scModule, err := servicecatalog.NewModule(fileConfig)
		if err != nil {
			log.Fatalf("Failed to create service catalog module: %v", err)
		}
		// Register routes using hexagonal architecture (prefix handled by module)
		scModule.RegisterRoutes(internal)
		log.Println("Service Catalog module initialized successfully")
		// Keep reference for kubernetes integration
		serviceCatalogModule = scModule
	}

	if dashConfig.Plugins.Has("Kubernetes") {
		// Parse kubernetes config
		k8sConfigs, err := kubernetes.ParseKubernetesConfigFromFileConfig(fileConfig)
		if err != nil {
			log.Printf("Failed to parse kubernetes config: %v", err)
			// Continue without kubernetes module
		} else {
			// Create kubernetes configurations
			var configs []k8sExternal.KubernetesConfig
			for _, k8sConfig := range k8sConfigs {
				configs = append(configs, k8sExternal.KubernetesConfig{
					Kubeconfig: k8sConfig.Kubeconfig,
					Context:    k8sConfig.Context,
				})
			}

			// Create repositories
			clusterRepo, err := k8sExternal.NewClusterRepository(configs)
			if err != nil {
				log.Printf("Failed to create cluster repository: %v", err)
			} else {
				nodeRepo := k8sExternal.NewNodeRepository(clusterRepo)
				namespaceRepo := k8sExternal.NewNamespaceRepository(clusterRepo)
				// Get service context resolver from service-catalog if available
				var serviceContextResolver k8sPorts.ServiceContextResolver
				if serviceCatalogModule != nil {
					serviceContextResolver = serviceCatalogModule.GetKubernetesAdapter()
				}

				deploymentRepo := k8sExternal.NewDeploymentRepository(clusterRepo, serviceContextResolver)
				podRepo := k8sExternal.NewPodRepository(clusterRepo)

				// Create kubernetes module with real repositories
				k8sModuleConfig := &kubernetes.ModuleConfig{
					ClusterRepo:    clusterRepo,
					NodeRepo:       nodeRepo,
					NamespaceRepo:  namespaceRepo,
					DeploymentRepo: deploymentRepo,
					PodRepo:        podRepo,
					// Services
					ClientService:  nil, // TODO: Implement client service
					MetricsService: nil, // TODO: Implement metrics service
					EventService:   nil, // TODO: Implement event service
					HealthService:  nil, // TODO: Implement health service
				}

				k8sModule, err := kubernetes.NewModule(k8sModuleConfig)
				if err != nil {
					log.Printf("Failed to create kubernetes module: %v", err)
				} else {
					// Register routes using hexagonal architecture (prefix handled by module)
					k8sModule.RegisterRoutes(internal)

					// Inject kubernetes service into service-catalog if available
					if serviceCatalogModule != nil {
						// Get kubernetes service adapter from kubernetes module
						k8sService := k8sModule.GetServiceCatalogAdapter()
						// Update service-catalog with kubernetes service
						serviceCatalogModule.UpdateKubernetesService(k8sService)
					}

					log.Println("Kubernetes module initialized successfully")
				}
			}
		}
	}

	if dashConfig.Plugins.Has("AWS") {
		// Create AWS module with minimal dependencies
		awsModule, err := aws.NewModule(fileConfig)
		if err != nil {
			log.Printf("Failed to create AWS module: %v", err)
		} else {
			// Register routes using hexagonal architecture
			awsModule.RegisterRoutes(internal)
			log.Println("AWS module initialized successfully")
		}
	}

	// Initialize SPA module using hexagonal architecture
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
