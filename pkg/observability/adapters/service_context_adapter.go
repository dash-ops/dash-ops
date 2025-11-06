package adapters

import (
	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// ServiceContextAdapter handles transformation between wire DTOs and domain models for service context
type ServiceContextAdapter struct{}

// NewServiceContextAdapter creates a new service context adapter
func NewServiceContextAdapter() *ServiceContextAdapter {
	return &ServiceContextAdapter{}
}

// WireRequestToModel transforms wire.ServicesWithContextRequest to parameters for controller
func (a *ServiceContextAdapter) WireRequestToModel(req *wire.ServicesWithContextRequest) (search string, limit, offset int) {
	return req.Search, req.Limit, req.Offset
}

// ModelToWireResponse transforms domain models to wire.ServicesWithContextResponse
func (a *ServiceContextAdapter) ModelToWireResponse(services []models.ServiceWithContext, total int, hasMore bool) *wire.ServicesWithContextResponse {
	return &wire.ServicesWithContextResponse{
		BaseResponse: wire.BaseResponse{
			Success: true,
		},
		Data: wire.ServicesWithContextData{
			Services:   services,
			Total:      total,
			HasMore:    hasMore,
			NextOffset: 0, // Can be calculated if needed
		},
	}
}

// ErrorToWireResponse transforms an error to wire.ServicesWithContextResponse
func (a *ServiceContextAdapter) ErrorToWireResponse(err error) *wire.ServicesWithContextResponse {
	return &wire.ServicesWithContextResponse{
		BaseResponse: wire.BaseResponse{
			Success: false,
			Error:   err.Error(),
		},
	}
}
