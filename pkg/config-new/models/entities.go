package config

import (
	"fmt"
	"strings"
)

// DashConfig represents the main DashOps configuration
type DashConfig struct {
	Port    string   `yaml:"port" json:"port"`
	Origin  string   `yaml:"origin" json:"origin"`
	Headers []string `yaml:"headers" json:"headers"`
	Front   string   `yaml:"front" json:"front"`
	Plugins Plugins  `yaml:"plugins" json:"plugins"`
}

// Plugins represents a list of enabled plugins
type Plugins []string

// Has checks if a plugin is enabled (case-insensitive)
func (p Plugins) Has(pluginName string) bool {
	for _, plugin := range p {
		if strings.EqualFold(plugin, pluginName) {
			return true
		}
	}
	return false
}

// Add adds a plugin to the list if not already present
func (p *Plugins) Add(pluginName string) {
	if !p.Has(pluginName) {
		*p = append(*p, pluginName)
	}
}

// Remove removes a plugin from the list
func (p *Plugins) Remove(pluginName string) {
	for i, plugin := range *p {
		if strings.EqualFold(plugin, pluginName) {
			*p = append((*p)[:i], (*p)[i+1:]...)
			return
		}
	}
}

// List returns all enabled plugins
func (p Plugins) List() []string {
	return []string(p)
}

// Count returns the number of enabled plugins
func (p Plugins) Count() int {
	return len(p)
}

// Validate validates the DashConfig
func (d *DashConfig) Validate() error {
	if d.Port == "" {
		return fmt.Errorf("port is required")
	}

	if d.Origin == "" {
		return fmt.Errorf("origin is required")
	}

	return nil
}

// GetPort returns the port with default fallback
func (d *DashConfig) GetPort() string {
	if d.Port == "" {
		return "8080"
	}
	return d.Port
}

// GetOrigin returns the origin with default fallback
func (d *DashConfig) GetOrigin() string {
	if d.Origin == "" {
		return "http://localhost:3000"
	}
	return d.Origin
}

// GetHeaders returns headers with default CORS headers if empty
func (d *DashConfig) GetHeaders() []string {
	if len(d.Headers) == 0 {
		return []string{"Content-Type", "Authorization"}
	}
	return d.Headers
}

// IsPluginEnabled checks if a specific plugin is enabled
func (d *DashConfig) IsPluginEnabled(pluginName string) bool {
	return d.Plugins.Has(pluginName)
}
