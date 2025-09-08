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
	"github.com/dash-ops/dash-ops/pkg/config"
	"github.com/dash-ops/dash-ops/pkg/github"
	"github.com/dash-ops/dash-ops/pkg/kubernetes"
	k8sExternal "github.com/dash-ops/dash-ops/pkg/kubernetes/adapters/external"
	k8sPorts "github.com/dash-ops/dash-ops/pkg/kubernetes/ports"
	servicecatalog "github.com/dash-ops/dash-ops/pkg/service-catalog"
	"github.com/dash-ops/dash-ops/pkg/spa"
	spaModels "github.com/dash-ops/dash-ops/pkg/spa/models"
	"golang.org/x/oauth2"
)

// responseRecorder captures the status code for middleware
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

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

	// Initialize plugins with dependency injection
	var serviceCatalogModule *servicecatalog.ServiceCatalog
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
		// Initialize Service Catalog module using hexagonal architecture
		fileConfig := configModule.GetFileConfigBytes()

		// Parse service catalog config
		scConfig, err := servicecatalog.ParseServiceCatalogConfig(fileConfig)
		if err != nil {
			log.Fatalf("Failed to parse service catalog config: %v", err)
		}

		// Create filesystem repository (real implementation)
		serviceRepo, err := servicecatalog.NewFilesystemRepository(scConfig.Directory)
		if err != nil {
			log.Fatalf("Failed to create filesystem repository: %v", err)
		}

		// Create service catalog module with full dependencies
		moduleConfig := &servicecatalog.ModuleConfig{
			ServiceRepo: serviceRepo,
			// TODO: Add other dependencies (Kubernetes, GitHub, Versioning) when available
		}
		scModule, err := servicecatalog.NewModule(moduleConfig)
		if err != nil {
			log.Fatalf("Failed to create service catalog module: %v", err)
		}
		// Register routes using hexagonal architecture (prefix handled by module)
		scModule.RegisterRoutes(internal)
		// Keep reference for kubernetes integration
		serviceCatalogModule = scModule
	}

	if dashConfig.Plugins.Has("Kubernetes") {
		// Get file config bytes
		fileConfig := configModule.GetFileConfigBytes()

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

	// Create middleware to handle API routes vs SPA routes
	apiMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If it's an API route, let it be handled by the API router
			if strings.HasPrefix(r.URL.Path, "/api/") {
				// Create a response recorder to check if the route was handled
				recorder := &responseRecorder{ResponseWriter: w, statusCode: 0}
				api.ServeHTTP(recorder, r)

				// If no route was matched (status 0), return 404
				if recorder.statusCode == 0 {
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(map[string]string{"error": "API endpoint not found"})
					return
				}
				return
			}
			// Otherwise, serve the SPA
			next.ServeHTTP(w, r)
		})
	}

	// Use SPA handler with middleware for non-API routes
	router.PathPrefix("/").Handler(apiMiddleware(spaModule.Handler))

	fmt.Println("DashOps server running!!")
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", dashConfig.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
