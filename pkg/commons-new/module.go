package commons

import (
	httpAdapter "github.com/dash-ops/dash-ops/pkg/commons-new/adapters/http"
	"github.com/dash-ops/dash-ops/pkg/commons-new/logic"
)

// Module represents the commons module with all its components
type Module struct {
	// Adapters
	ResponseAdapter *httpAdapter.ResponseAdapter
	RequestAdapter  *httpAdapter.RequestAdapter

	// Logic
	PermissionChecker *logic.PermissionChecker
	StringProcessor   *logic.StringProcessor
}

// NewModule creates and initializes a new commons module
func NewModule() *Module {
	return &Module{
		// Initialize adapters
		ResponseAdapter: httpAdapter.NewResponseAdapter(),
		RequestAdapter:  httpAdapter.NewRequestAdapter(),

		// Initialize logic components
		PermissionChecker: logic.NewPermissionChecker(),
		StringProcessor:   logic.NewStringProcessor(),
	}
}

// GetResponseAdapter returns the response adapter
func (m *Module) GetResponseAdapter() *httpAdapter.ResponseAdapter {
	return m.ResponseAdapter
}

// GetRequestAdapter returns the request adapter
func (m *Module) GetRequestAdapter() *httpAdapter.RequestAdapter {
	return m.RequestAdapter
}

// GetPermissionChecker returns the permission checker
func (m *Module) GetPermissionChecker() *logic.PermissionChecker {
	return m.PermissionChecker
}

// GetStringProcessor returns the string processor
func (m *Module) GetStringProcessor() *logic.StringProcessor {
	return m.StringProcessor
}
