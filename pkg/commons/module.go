package commons

import (
	httpAdapter "github.com/dash-ops/dash-ops/pkg/commons/adapters/http"
	commonsLogic "github.com/dash-ops/dash-ops/pkg/commons/logic"
)

// Module represents the commons module with all its components
type Module struct {
	// Adapters
	ResponseAdapter *httpAdapter.ResponseAdapter
	RequestAdapter  *httpAdapter.RequestAdapter

	// Logic
	PermissionChecker *commonsLogic.PermissionChecker
	StringProcessor   *commonsLogic.StringProcessor
}

// NewModule creates and initializes a new commons module
func NewModule() *Module {
	return &Module{
		// Initialize adapters
		ResponseAdapter: httpAdapter.NewResponseAdapter(),
		RequestAdapter:  httpAdapter.NewRequestAdapter(),

		// Initialize logic components
		PermissionChecker: commonsLogic.NewPermissionChecker(),
		StringProcessor:   commonsLogic.NewStringProcessor(),
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
func (m *Module) GetPermissionChecker() *commonsLogic.PermissionChecker {
	return m.PermissionChecker
}

// GetStringProcessor returns the string processor
func (m *Module) GetStringProcessor() *commonsLogic.StringProcessor {
	return m.StringProcessor
}
