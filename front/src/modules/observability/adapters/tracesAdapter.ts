import type { TraceInfo, TraceSpan } from '../types';

export const transformTraceInfoToDomain = (apiTrace: any): TraceInfo => {
  const spans: any[] = apiTrace.spans || [];
  const services = new Set(spans.map((s) => s.service));
  const hasErrors = spans.some((s) => mapStatus(s.status) === 'error');

  return {
    traceId: apiTrace.traceId ?? apiTrace.trace_id ?? '',
    rootOperation: apiTrace.rootOperation ?? apiTrace.root_operation ?? apiTrace.operation ?? spans[0]?.operationName ?? 'unknown',
    totalDuration: apiTrace.totalDuration ?? apiTrace.duration ?? 0,
    spanCount: apiTrace.spanCount ?? spans.length ?? 0,
    serviceCount: apiTrace.serviceCount ?? services.size ?? 0,
    status: mapStatus(apiTrace.status ?? (hasErrors ? 'error' : 'ok')),
    timestamp: apiTrace.timestamp ?? Date.now(),
    errors: apiTrace.errors ?? apiTrace.error_count ?? spans.filter((s) => mapStatus(s.status) === 'error').length ?? 0,
    name: apiTrace.name ?? apiTrace.rootOperation ?? apiTrace.root_operation ?? 'unknown',
    created_at: apiTrace.created_at,
    updated_at: apiTrace.updated_at,
  } as TraceInfo;
};

export const transformTraceSpanToDomain = (apiSpan: any): TraceSpan => {
  // Handle span_id -> id mapping
  const spanId = apiSpan.id ?? apiSpan.spanId ?? apiSpan.span_id ?? '';
  
  // Handle parent_span_id -> parentId mapping
  const parentId = apiSpan.parentId ?? apiSpan.parent_id ?? apiSpan.parent_span_id;
  
  // Handle service_name -> service mapping
  const service = apiSpan.service ?? apiSpan.serviceName ?? apiSpan.service_name ?? 'unknown';
  
  // Handle operation_name -> operationName mapping
  const operationName = apiSpan.operationName ?? apiSpan.operation_name ?? apiSpan.operation ?? 'unknown';
  
  // Handle start_time: can be string ISO8601 or number (milliseconds or nanoseconds)
  let startTime: number = 0;
  if (apiSpan.startTime) {
    startTime = typeof apiSpan.startTime === 'string' 
      ? new Date(apiSpan.startTime).getTime()
      : apiSpan.startTime;
  } else if (apiSpan.start_time) {
    startTime = typeof apiSpan.start_time === 'string'
      ? new Date(apiSpan.start_time).getTime()
      : apiSpan.start_time;
  }
  
  // Handle duration: backend sends in nanoseconds, convert to milliseconds
  let duration: number = 0;
  if (apiSpan.duration) {
    // If duration is very large (> 1e9), it's likely in nanoseconds, convert to ms
    if (apiSpan.duration > 1e9) {
      duration = apiSpan.duration / 1e6; // nanoseconds to milliseconds
    } else if (apiSpan.duration > 1e6) {
      duration = apiSpan.duration / 1000; // microseconds to milliseconds
    } else {
      duration = apiSpan.duration; // already in milliseconds
    }
  }
  
  // Handle status: can be object with code/message or string
  let status: 'ok' | 'error' = 'ok';
  if (apiSpan.status) {
    if (typeof apiSpan.status === 'object') {
      // Backend sends SpanStatus with code: 0=OK, 1=ERROR
      status = apiSpan.status.code === 1 ? 'error' : 'ok';
    } else {
      status = mapStatus(apiSpan.status);
    }
  }

  return {
    id: spanId,
    traceId: apiSpan.traceId ?? apiSpan.trace_id ?? '',
    operationName,
    service,
    startTime,
    duration,
    status,
    tags: apiSpan.tags ?? {},
    parentId,
    depth: apiSpan.depth,
    name: operationName,
    created_at: apiSpan.created_at,
    updated_at: apiSpan.updated_at,
  } as TraceSpan;
};

export const transformTraceSpansToDomain = (apiSpans: any[]): TraceSpan[] => {
  return (apiSpans || []).map(transformTraceSpanToDomain);
};

const mapStatus = (status: any): 'ok' | 'error' => {
  const s = String(status ?? 'ok').toLowerCase();
  return s === 'error' ? 'error' : 'ok';
};


