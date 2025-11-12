package wire

// ConfigResponse represents the configuration API response.
type ConfigResponse struct {
	Port    string   `json:"port"`
	Origin  string   `json:"origin"`
	Headers []string `json:"headers"`
	Front   string   `json:"front,omitempty"`
	Plugins []string `json:"plugins"`
}

// PluginsResponse represents the plugins list API response.
type PluginsResponse struct {
	Plugins []string `json:"plugins"`
	Count   int      `json:"count"`
}

// PluginStatusResponse represents a single plugin status.
type PluginStatusResponse struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// SystemInfoResponse represents system information.
type SystemInfoResponse struct {
	Version     string   `json:"version,omitempty"`
	Environment string   `json:"environment,omitempty"`
	Uptime      string   `json:"uptime,omitempty"`
	Plugins     []string `json:"plugins"`
}
