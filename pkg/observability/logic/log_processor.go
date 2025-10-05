package logic

import (
	"github.com/dash-ops/dash-ops/pkg/observability/models"
)

// LogProcessor handles log processing business logic
type LogProcessor struct {
	// Add any configuration or dependencies here
}

// NewLogProcessor creates a new log processor
func NewLogProcessor() *LogProcessor {
	return &LogProcessor{}
}

// EnrichLogs enriches log entries with additional context
// This is pure business logic - no external dependencies
func (p *LogProcessor) EnrichLogs(logs []models.LogEntry) []models.LogEntry {
	enriched := make([]models.LogEntry, len(logs))

	for i, log := range logs {
		enriched[i] = log

		// Add enrichments based on business rules
		if enriched[i].Metadata == nil {
			enriched[i].Metadata = make(map[string]interface{})
		}

		// Enrich with severity classification
		enriched[i].Metadata["severity"] = p.classifySeverity(log.Level)

		// Enrich with error detection
		if p.isErrorLog(log) {
			enriched[i].Metadata["is_error"] = true
		}

		// Add correlation hints
		if log.TraceID != "" {
			enriched[i].Metadata["has_trace"] = true
		}
	}

	return enriched
}

// classifySeverity classifies log severity
func (p *LogProcessor) classifySeverity(level string) string {
	switch level {
	case "error":
		return "critical"
	case "warn":
		return "warning"
	case "info":
		return "informational"
	case "debug":
		return "debug"
	default:
		return "unknown"
	}
}

// isErrorLog checks if a log entry represents an error
func (p *LogProcessor) isErrorLog(log models.LogEntry) bool {
	return log.Level == "error" || log.Level == "fatal" || log.Level == "panic"
}

// FilterLogsByLevel filters logs by level
func (p *LogProcessor) FilterLogsByLevel(logs []models.LogEntry, level string) []models.LogEntry {
	if level == "" {
		return logs
	}

	filtered := make([]models.LogEntry, 0)
	for _, log := range logs {
		if log.Level == level {
			filtered = append(filtered, log)
		}
	}

	return filtered
}

// FilterLogsByService filters logs by service
func (p *LogProcessor) FilterLogsByService(logs []models.LogEntry, service string) []models.LogEntry {
	if service == "" {
		return logs
	}

	filtered := make([]models.LogEntry, 0)
	for _, log := range logs {
		if log.Service == service {
			filtered = append(filtered, log)
		}
	}

	return filtered
}
