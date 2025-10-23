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
  return {
    id: apiSpan.id ?? apiSpan.spanId ?? '',
    traceId: apiSpan.traceId ?? apiSpan.trace_id ?? '',
    operationName: apiSpan.operationName ?? apiSpan.operation ?? 'unknown',
    service: apiSpan.service ?? 'unknown',
    startTime: apiSpan.startTime ?? apiSpan.start_time ?? 0,
    duration: apiSpan.duration ?? 0,
    status: mapStatus(apiSpan.status),
    tags: apiSpan.tags ?? {},
    parentId: apiSpan.parentId ?? apiSpan.parent_id,
    depth: apiSpan.depth,
    name: apiSpan.name ?? apiSpan.operationName ?? apiSpan.operation ?? 'unknown',
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


