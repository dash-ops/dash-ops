package config

// DefaultConfig returns default configuration values
func DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"port":   "8080",
		"origin": "http://localhost:3000",
		"headers": []string{
			"Content-Type",
			"Authorization",
		},
		"plugins": []string{},
	}
}

// DefaultPort returns the default port
func DefaultPort() string {
	return "8080"
}

// DefaultOrigin returns the default origin
func DefaultOrigin() string {
	return "http://localhost:3000"
}

// DefaultHeaders returns the default CORS headers
func DefaultHeaders() []string {
	return []string{
		"Content-Type",
		"Authorization",
	}
}

// DefaultConfigFile returns the default config file path
func DefaultConfigFile() string {
	return "./dash-ops.yaml"
}

// SupportedConfigFormats returns supported configuration file formats
func SupportedConfigFormats() []string {
	return []string{"yaml", "yml"}
}

// RequiredFields returns the list of required configuration fields
func RequiredFields() []string {
	return []string{"port", "origin"}
}
