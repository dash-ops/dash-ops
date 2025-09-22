package models

import (
	"time"
)

// Alert represents an alert instance
type Alert struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Status       string                 `json:"status"`
	Severity     string                 `json:"severity"`
	Service      string                 `json:"service"`
	Labels       map[string]string      `json:"labels"`
	Annotations  map[string]string      `json:"annotations"`
	StartsAt     time.Time              `json:"starts_at"`
	EndsAt       *time.Time             `json:"ends_at,omitempty"`
	GeneratorURL string                 `json:"generator_url,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	RuleID       string                 `json:"rule_id,omitempty"`
}

// ProcessedAlert represents a processed alert with additional context
type ProcessedAlert struct {
	Alert
	ProcessedAt  time.Time              `json:"processed_at"`
	Enrichments  map[string]interface{} `json:"enrichments,omitempty"`
	Correlations []string               `json:"correlations,omitempty"`
	Actions      []AlertAction          `json:"actions,omitempty"`
}

// AlertAction represents an action taken on an alert
type AlertAction struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	User        string                 `json:"user,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AlertRule represents an alert rule configuration
type AlertRule struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Query       string            `json:"query"`
	Threshold   float64           `json:"threshold"`
	Severity    string            `json:"severity"`
	Service     string            `json:"service,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Enabled     bool              `json:"enabled"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// AlertEvaluation represents the evaluation of an alert rule
type AlertEvaluation struct {
	RuleID      string    `json:"rule_id"`
	EvaluatedAt time.Time `json:"evaluated_at"`
	Result      bool      `json:"result"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	Message     string    `json:"message,omitempty"`
}

// AlertsConfig represents alerts configuration
type AlertsConfig struct {
	Enabled    bool     `json:"enabled"`
	Channels   []string `json:"channels"`
	Severities []string `json:"severities"`
	Cooldown   string   `json:"cooldown"`
	MaxAlerts  int      `json:"max_alerts"`
}
