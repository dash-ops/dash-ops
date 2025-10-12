package tempo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dash-ops/dash-ops/pkg/observability/models"
	"github.com/dash-ops/dash-ops/pkg/observability/wire"
)

// transformToTempoParams converts models.TraceQuery to wire.TempoQueryParams
func transformToTempoParams(query *models.TraceQuery) *wire.TempoQueryParams {
	params := &wire.TempoQueryParams{
		Start: query.StartTime.Unix(), // Tempo expects Unix seconds, not nanoseconds
		End:   query.EndTime.Unix(),   // Tempo expects Unix seconds, not nanoseconds
		Limit: query.Limit,
	}

	// Build TraceQL query
	params.Query = buildTraceQLQuery(query)

	return params
}

// buildTraceQLQuery builds a TraceQL query string from models.TraceQuery
func buildTraceQLQuery(query *models.TraceQuery) string {
	var conditions []string

	// Service filter
	if query.Service != "" {
		conditions = append(conditions, fmt.Sprintf(`resource.service.name="%s"`, query.Service))
	}

	// Operation/span name filter
	if query.Operation != "" {
		conditions = append(conditions, fmt.Sprintf(`name="%s"`, query.Operation))
	}

	// Tags/attributes filter
	for key, value := range query.Tags {
		conditions = append(conditions, fmt.Sprintf(`span.%s="%s"`, key, value))
	}

	// Duration filters
	if query.MinDuration != "" {
		if duration, err := time.ParseDuration(query.MinDuration); err == nil {
			conditions = append(conditions, fmt.Sprintf("duration>=%s", duration.String()))
		}
	}
	if query.MaxDuration != "" {
		if duration, err := time.ParseDuration(query.MaxDuration); err == nil {
			conditions = append(conditions, fmt.Sprintf("duration<=%s", duration.String()))
		}
	}

	// Combine conditions with AND
	if len(conditions) == 0 {
		return "{}" // Empty query returns all traces
	}

	return fmt.Sprintf("{%s}", strings.Join(conditions, " && "))
}

// transformSearchResponseToModels converts wire.TempoSearchResponse to []models.TraceSummary
func transformSearchResponseToModels(response *wire.TempoSearchResponse) []models.TraceSummary {
	summaries := make([]models.TraceSummary, 0, len(response.Traces))

	for _, trace := range response.Traces {
		summary := models.TraceSummary{
			TraceID:       trace.TraceID,
			RootService:   trace.RootServiceName,
			RootOperation: trace.RootTraceName,
			Duration:      time.Duration(trace.DurationMs) * time.Millisecond,
		}

		// Parse start time (Unix nanoseconds as string)
		if startNano, err := strconv.ParseInt(trace.StartTimeUnixNano, 10, 64); err == nil {
			summary.StartTime = time.Unix(0, startNano)
		}

		// Extract span count and services from SpanSet if available
		if trace.SpanSet != nil {
			summary.SpanCount = trace.SpanSet.Matched
			// Extract unique services from spans
			servicesMap := make(map[string]bool)
			for _, span := range trace.SpanSet.Spans {
				if serviceName, ok := span.Attributes["service.name"].(string); ok {
					servicesMap[serviceName] = true
				}
			}
			for service := range servicesMap {
				summary.Services = append(summary.Services, service)
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries
}

// transformTraceByIDResponseToModel converts wire.TempoTraceByIDResponse to models.Trace
func transformTraceByIDResponseToModel(response *wire.TempoTraceByIDResponse) *models.Trace {
	if len(response.Batches) == 0 {
		return nil
	}

	trace := &models.Trace{
		Spans:    make([]models.TraceSpan, 0),
		Services: make([]string, 0),
	}

	servicesMap := make(map[string]bool)
	var minStartTime time.Time
	var maxEndTime time.Time

	// Process all batches and spans
	for _, batch := range response.Batches {
		// Extract service name from resource attributes
		serviceName := extractServiceNameFromResource(&batch.Resource)

		// Process scope spans
		for _, scopeSpan := range batch.ScopeSpans {
			for _, tempoSpan := range scopeSpan.Spans {
				span := transformTempoSpanToModel(&tempoSpan, serviceName)
				trace.Spans = append(trace.Spans, span)

				// Track trace ID
				if trace.TraceID == "" {
					trace.TraceID = span.TraceID
				}

				// Track services
				if serviceName != "" {
					servicesMap[serviceName] = true
				}

				// Track time range
				if minStartTime.IsZero() || span.StartTime.Before(minStartTime) {
					minStartTime = span.StartTime
				}
				endTime := span.StartTime.Add(span.Duration)
				if maxEndTime.IsZero() || endTime.After(maxEndTime) {
					maxEndTime = endTime
				}
			}
		}

		// Process deprecated instrumentation library spans
		for _, ilSpan := range batch.InstrumentationLibrarySpans {
			for _, tempoSpan := range ilSpan.Spans {
				span := transformTempoSpanToModel(&tempoSpan, serviceName)
				trace.Spans = append(trace.Spans, span)

				// Track services
				if serviceName != "" {
					servicesMap[serviceName] = true
				}

				// Track time range
				if minStartTime.IsZero() || span.StartTime.Before(minStartTime) {
					minStartTime = span.StartTime
				}
				endTime := span.StartTime.Add(span.Duration)
				if maxEndTime.IsZero() || endTime.After(maxEndTime) {
					maxEndTime = endTime
				}
			}
		}
	}

	// Set trace-level properties
	trace.StartTime = minStartTime
	if !maxEndTime.IsZero() && !minStartTime.IsZero() {
		trace.Duration = maxEndTime.Sub(minStartTime)
	}
	for service := range servicesMap {
		trace.Services = append(trace.Services, service)
	}

	return trace
}

// extractServiceNameFromResource extracts the service name from resource attributes
func extractServiceNameFromResource(resource *wire.TempoResource) string {
	for _, attr := range resource.Attributes {
		if attr.Key == "service.name" {
			return extractStringValue(&attr.Value)
		}
	}
	return ""
}

// transformTempoSpanToModel converts wire.TempoTraceSpan to models.TraceSpan
func transformTempoSpanToModel(tempoSpan *wire.TempoTraceSpan, serviceName string) models.TraceSpan {
	span := models.TraceSpan{
		TraceID:       tempoSpan.TraceID,
		SpanID:        tempoSpan.SpanID,
		ParentSpanID:  tempoSpan.ParentSpanID,
		OperationName: tempoSpan.Name,
		ServiceName:   serviceName,
		Tags:          make(map[string]interface{}),
		Logs:          make([]models.SpanLog, 0),
		References:    make([]models.SpanReference, 0),
	}

	// Parse start time
	if startNano, err := strconv.ParseInt(tempoSpan.StartTimeUnixNano, 10, 64); err == nil {
		span.StartTime = time.Unix(0, startNano)
	}

	// Parse duration
	if endNano, err := strconv.ParseInt(tempoSpan.EndTimeUnixNano, 10, 64); err == nil {
		endTime := time.Unix(0, endNano)
		span.Duration = endTime.Sub(span.StartTime)
	}

	// Transform attributes to tags
	for _, attr := range tempoSpan.Attributes {
		span.Tags[attr.Key] = extractValue(&attr.Value)
	}

	// Transform events to logs
	for _, event := range tempoSpan.Events {
		log := models.SpanLog{
			Fields: make(map[string]interface{}),
		}
		if eventTime, err := strconv.ParseInt(event.TimeUnixNano, 10, 64); err == nil {
			log.Timestamp = time.Unix(0, eventTime)
		}
		log.Fields["event"] = event.Name
		for _, attr := range event.Attributes {
			log.Fields[attr.Key] = extractValue(&attr.Value)
		}
		span.Logs = append(span.Logs, log)
	}

	// Transform links to references
	for _, link := range tempoSpan.Links {
		ref := models.SpanReference{
			RefType: "FOLLOWS_FROM", // Tempo links are typically FOLLOWS_FROM
			TraceID: link.TraceID,
			SpanID:  link.SpanID,
		}
		span.References = append(span.References, ref)
	}

	// Transform status
	if tempoSpan.Status != nil {
		span.Status = models.SpanStatus{
			Code:    tempoSpan.Status.Code,
			Message: tempoSpan.Status.Message,
		}
	}

	return span
}

// extractValue extracts the actual value from TempoValue
func extractValue(value *wire.TempoValue) interface{} {
	if value.StringValue != "" {
		return value.StringValue
	}
	if value.IntValue != "" {
		if i, err := strconv.ParseInt(value.IntValue, 10, 64); err == nil {
			return i
		}
		return value.IntValue
	}
	if value.DoubleValue != 0 {
		return value.DoubleValue
	}
	if value.BoolValue {
		return value.BoolValue
	}
	if value.ArrayValue != nil {
		arr := make([]interface{}, len(value.ArrayValue.Values))
		for i, v := range value.ArrayValue.Values {
			arr[i] = extractValue(&v)
		}
		return arr
	}
	if value.KvlistValue != nil {
		kvMap := make(map[string]interface{})
		for _, kv := range value.KvlistValue.Values {
			kvMap[kv.Key] = extractValue(&kv.Value)
		}
		return kvMap
	}
	if value.BytesValue != "" {
		return value.BytesValue
	}
	return nil
}

// extractStringValue extracts a string value from TempoValue
func extractStringValue(value *wire.TempoValue) string {
	if v := extractValue(value); v != nil {
		if str, ok := v.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}
