package auth

import (
	"fmt"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"

	authAdapters "github.com/dash-ops/dash-ops/pkg/auth/adapters/http"
	authControllers "github.com/dash-ops/dash-ops/pkg/auth/controllers"
	authHandlers "github.com/dash-ops/dash-ops/pkg/auth/handlers"
	authLogic "github.com/dash-ops/dash-ops/pkg/auth/logic"
	authModels "github.com/dash-ops/dash-ops/pkg/auth/models"
	commonsHttp "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	gh "github.com/dash-ops/dash-ops/pkg/github"
)

// Module represents the auth module - main entry point for the plugin
type Module struct {
	config     *authModels.AuthConfig
	controller *authControllers.AuthController
	handler    *authHandlers.HTTPHandler
}

// NewModule creates and initializes a new auth module (main factory)
func NewModule(config *authModels.AuthConfig) (*Module, error) {
	if config == nil {
		return nil, fmt.Errorf("auth config cannot be nil")
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid auth config: %w", err)
	}

	// Initialize logic components
	oauth2Processor := authLogic.NewOAuth2Processor()
	sessionManager := authLogic.NewSessionManager(24 * time.Hour)

	// Initialize GitHub client for provider integration
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       config.Scopes,
		RedirectURL:  config.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
	}
	githubClient, err := gh.NewClient(oauthConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// Initialize controller with dependencies
	controller := authControllers.NewAuthController(
		config,
		oauth2Processor,
		sessionManager,
		githubClient,
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

// GetConfig returns the current configuration (for compatibility)
func (m *Module) GetConfig() *authModels.AuthConfig {
	return m.config
}

// Legacy compatibility functions for existing main.go

// MakeOauthHandlers registers OAuth handlers - legacy compatibility
func MakeOauthHandlers(apiRouter, internalRouter *mux.Router, fileConfig []byte) {
	// Parse config from bytes (similar to original oauth2 module)
	config, err := parseAuthConfigFromBytes(fileConfig)
	if err != nil {
		panic(err) // Maintain same behavior as original
	}

	// Create module
	module, err := NewModule(config)
	if err != nil {
		panic(err) // Maintain same behavior as original
	}

	// Register routes
	module.RegisterRoutes(apiRouter, internalRouter)
}

// parseAuthConfigFromBytes parses auth config from YAML bytes (helper function)
func parseAuthConfigFromBytes(fileConfig []byte) (*authModels.AuthConfig, error) {
	// Parse YAML similar to the original oauth2.loadConfig
	type dashYaml struct {
		Oauth2 []struct {
			Provider        string   `yaml:"provider"`
			ClientID        string   `yaml:"clientId"`
			ClientSecret    string   `yaml:"clientSecret"`
			AuthURL         string   `yaml:"authURL"`
			TokenURL        string   `yaml:"tokenURL"`
			RedirectURL     string   `yaml:"redirectURL"`
			URLLoginSuccess string   `yaml:"urlLoginSuccess"`
			OrgPermission   string   `yaml:"orgPermission"`
			Scopes          []string `yaml:"scopes"`
		} `yaml:"oauth2"`
	}

	var config dashYaml
	if err := yaml.Unmarshal(fileConfig, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if len(config.Oauth2) == 0 {
		return nil, fmt.Errorf("no oauth2 configuration found")
	}

	oauth := config.Oauth2[0]
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
