package logic

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// QueryParser handles parsing of explorer queries
type QueryParser struct{}

// NewQueryParser creates a new query parser
func NewQueryParser() *QueryParser {
	return &QueryParser{}
}

// Parse parses an explorer query and returns the parsed structure
func (p *QueryParser) Parse(query string) (*wire.ParsedQuery, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	parsed := &wire.ParsedQuery{
		RawQuery: query,
		Filters:  make(map[string]interface{}),
	}

	// Determine data source from FROM clause
	queryUpper := strings.ToUpper(query)
	if strings.Contains(queryUpper, "FROM LOGS") || strings.Contains(queryUpper, "FROM LOG") {
		parsed.DataSource = "logs"
		// Parse WHERE conditions for SQL-style query
		p.parseWhereConditions(query, parsed)
	} else if strings.Contains(queryUpper, "FROM TRACES") || strings.Contains(queryUpper, "FROM TRACE") {
		parsed.DataSource = "traces"
		// Parse WHERE conditions for SQL-style query
		p.parseWhereConditions(query, parsed)
	} else if strings.Contains(queryUpper, "FROM METRICS") || strings.Contains(queryUpper, "FROM METRIC") {
		parsed.DataSource = "metrics"
		// Parse WHERE conditions for SQL-style query
		p.parseWhereConditions(query, parsed)
	} else {
		// Check if it's a direct LogQL query (starts with { or has {label="value"})
		if strings.HasPrefix(query, "{") || strings.Contains(query, "=\"") {
			parsed.DataSource = "logs"
			// Treat as raw LogQL query
			parsed.RawQuery = query
		} else {
			return nil, fmt.Errorf("invalid query: missing or unknown FROM clause (expected: FROM Logs, FROM Traces, FROM Metrics, or LogQL query)")
		}
	}

	return parsed, nil
}

// parseWhereConditions extracts WHERE clause conditions and converts to LogQL
func (p *QueryParser) parseWhereConditions(query string, parsed *wire.ParsedQuery) {
	// Simple regex-based parsing for WHERE conditions
	// Format: WHERE key = "value" [AND/OR key = "value"]
	wherePattern := regexp.MustCompile(`(?i)WHERE\s+(.+?)(?:\s+ORDER|\s+LIMIT|$)`)
	whereMatches := wherePattern.FindStringSubmatch(query)

	if len(whereMatches) < 2 {
		// No WHERE clause, use default
		return
	}

	whereClause := whereMatches[1]

	// Build LogQL query from WHERE conditions
	var logqlFilters []string

	// Parse individual conditions: key = "value" or key = 'value'
	conditionPattern := regexp.MustCompile(`(\w+)\s*=\s*["']([^"']+)["']`)
	conditions := conditionPattern.FindAllStringSubmatch(whereClause, -1)

	for _, match := range conditions {
		if len(match) >= 3 {
			key := strings.ToLower(match[1])
			value := match[2]
			parsed.Filters[key] = value
			// Build LogQL filter: key="value"
			logqlFilters = append(logqlFilters, fmt.Sprintf(`%s="%s"`, key, value))
		}
	}

	// Parse numeric conditions: key > value, key < value, etc.
	numericPattern := regexp.MustCompile(`(\w+)\s*([><]=?|=)\s*(\d+)`)
	numericConditions := numericPattern.FindAllStringSubmatch(whereClause, -1)

	for _, match := range numericConditions {
		if len(match) >= 4 {
			key := strings.ToLower(match[1])
			value := match[3]
			parsed.Filters[key] = value
			// Build LogQL filter: key=value (numeric)
			logqlFilters = append(logqlFilters, fmt.Sprintf(`%s=%s`, key, value))
		}
	}

	// Build complete LogQL query if we have filters
	if len(logqlFilters) > 0 {
		parsed.RawQuery = fmt.Sprintf("{%s}", strings.Join(logqlFilters, ","))
	}
}

// GetDataSource returns the data source from the query
func (p *QueryParser) GetDataSource(query string) (string, error) {
	parsed, err := p.Parse(query)
	if err != nil {
		return "", err
	}
	return parsed.DataSource, nil
}
