package spa

import (
	"fmt"
	"time"

	"github.com/gorilla/mux"

	spaAdapters "github.com/dash-ops/dash-ops/pkg/spa/adapters/http"
	spaHandlers "github.com/dash-ops/dash-ops/pkg/spa/handlers"
	spaLogic "github.com/dash-ops/dash-ops/pkg/spa/logic"
	spaModels "github.com/dash-ops/dash-ops/pkg/spa/models"
)

// Module represents the SPA module with all its components
type Module struct {
	// Core components
	Handler *spaHandlers.HTTPHandler

	// Logic components
	FileProcessor *spaLogic.FileProcessor

	// Adapters
	SPAAdapter *spaAdapters.SPAAdapter

	// Configuration
	Config    *spaModels.SPAConfig
	APIRouter *mux.Router
}

// NewModule creates and initializes a new SPA module
func NewModule(config *spaModels.SPAConfig, apiRouter *mux.Router) (*Module, error) {
	if config == nil {
		return nil, fmt.Errorf("SPA config cannot be nil")
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid SPA config: %w", err)
	}

	// Initialize file processor
	fileProcessor := spaLogic.NewFileProcessor()

	// Validate static path and index file
	if err := fileProcessor.ValidateStaticPath(config.StaticPath); err != nil {
		return nil, fmt.Errorf("static path validation failed: %w", err)
	}

	if err := fileProcessor.ValidateIndexFile(config.StaticPath, config.IndexPath); err != nil {
		return nil, fmt.Errorf("index file validation failed: %w", err)
	}

	// Initialize SPA adapter
	spaAdapter := spaAdapters.NewSPAAdapter(
		config,
		fileProcessor,
		&spaModels.SPAStats{StartTime: time.Now()},
		apiRouter,
	)

	// Initialize handler
	handler := spaHandlers.NewHTTPHandler(spaAdapter)

	return &Module{
		Handler:       handler,
		FileProcessor: fileProcessor,
		SPAAdapter:    spaAdapter,
		Config:        config,
		APIRouter:     apiRouter,
	}, nil
}

// RegisterRoutes registers HTTP routes for the SPA module
func (m *Module) RegisterRoutes(router *mux.Router) {
	if m.Handler == nil {
		return
	}
	m.Handler.RegisterRoutes(router)
}
