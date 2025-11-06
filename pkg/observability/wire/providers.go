package wire

// ProvidersResponse represents the response for listing available providers
type ProvidersResponse struct {
	BaseResponse
	Data ProvidersResponseData `json:"data"`
}

// ProvidersResponseData represents the data portion of providers response
type ProvidersResponseData struct {
	LogsProviders    []string `json:"logs_providers"`
	TracesProviders  []string `json:"traces_providers"`
	MetricsProviders []string `json:"metrics_providers"`
}
