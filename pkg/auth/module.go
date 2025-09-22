package auth

import (
	"fmt"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"

	authAdapters "github.com/dash-ops/dash-ops/pkg/auth/adapters/http"
	authControllers "github.com/dash-ops/dash-ops/pkg/auth/controllers"
	authHandlers "github.com/dash-ops/dash-ops/pkg/auth/handlers"
	authLogic "github.com/dash-ops/dash-ops/pkg/auth/logic"
	authModels "github.com/dash-ops/dash-ops/pkg/auth/models"
	authPorts "github.com/dash-ops/dash-ops/pkg/auth/ports"
	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
)

// Module represents the auth module - main entry point for the plugin
type Module struct {
	config     *authModels.AuthConfig
	controller *authControllers.AuthController
	handler    *authHandlers.HTTPHandler
}

// NewModule creates and initializes a new auth module (main factory)
func NewModule(config *authModels.AuthConfig, githubService authPorts.GitHubService) (*Module, error) {
	if config == nil {
		return nil, fmt.Errorf("auth config cannot be nil")
	}

	if githubService == nil {
		return nil, fmt.Errorf("github service is required")
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid auth config: %w", err)
	}

	// Initialize logic components
	oauth2Processor := authLogic.NewOAuth2Processor()
	sessionManager := authLogic.NewSessionManager(24 * time.Hour)

	// Initialize controller with dependencies (using injected GitHub service)
	controller := authControllers.NewAuthController(
		config,
		oauth2Processor,
		sessionManager,
		githubService, // Injected dependency
	)

	// Initialize adapters
	authAdapter := authAdapters.NewAuthAdapter()
	responseAdapter := commonsHttp.NewResponseAdapter()
	requestAdapter := commonsHttp.NewRequestAdapter()

	// Initialize HTTP handler
	handler := authHandlers.NewHTTPHandler(
		controller,
		authAdapter,
		responseAdapter,
		requestAdapter,
	)

	return &Module{
		config:     config,
		controller: controller,
		handler:    handler,
	}, nil
}

// RegisterRoutes registers HTTP routes for the auth module
func (m *Module) RegisterRoutes(apiRouter, internalRouter *mux.Router) {
	// Delegate to handler (following hexagonal architecture)
	m.handler.RegisterRoutes(apiRouter, internalRouter)

	// Add organization permission middleware if configured
	if m.config.OrgPermission != "" {
		internalRouter.Use(m.handler.OrgPermissionMiddleware)
	}
}

// ParseAuthConfigFromFileConfig parses auth config from YAML bytes (exported for main.go)
func ParseAuthConfigFromFileConfig(fileConfig []byte) (*authModels.AuthConfig, error) {
	// Parse YAML similar to the original oauth2.loadConfig
	type dashYaml struct {
		Auth []struct {
			Provider        string   `yaml:"provider"`
			ClientID        string   `yaml:"clientId"`
			ClientSecret    string   `yaml:"clientSecret"`
			AuthURL         string   `yaml:"authURL"`
			TokenURL        string   `yaml:"tokenURL"`
			RedirectURL     string   `yaml:"redirectURL"`
			URLLoginSuccess string   `yaml:"urlLoginSuccess"`
			OrgPermission   string   `yaml:"orgPermission"`
			Scopes          []string `yaml:"scopes"`
		} `yaml:"auth"`
	}

	var config dashYaml
	if err := yaml.Unmarshal(fileConfig, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if len(config.Auth) == 0 {
		return nil, fmt.Errorf("no auth configuration found")
	}

	oauth := config.Auth[0]
	return &authModels.AuthConfig{
		Provider:        authModels.ProviderGitHub,
		Method:          authModels.MethodOAuth2,
		Enabled:         true,
		ClientID:        oauth.ClientID,
		ClientSecret:    oauth.ClientSecret,
		AuthURL:         oauth.AuthURL,
		TokenURL:        oauth.TokenURL,
		RedirectURL:     oauth.RedirectURL,
		URLLoginSuccess: oauth.URLLoginSuccess,
		OrgPermission:   oauth.OrgPermission,
		Scopes:          oauth.Scopes,
	}, nil
}
