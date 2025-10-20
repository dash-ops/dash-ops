import type { LogEntry, LogsResponse } from '../types';

export const transformLogToDomain = (apiLog: any): LogEntry => {
  return {
    id: String(apiLog.id ?? apiLog._id ?? cryptoRandom()),
    timestamp: apiLog.timestamp ?? apiLog.ts ?? new Date().toISOString(),
    level: mapLevel(apiLog.level),
    service: apiLog.service ?? apiLog.source ?? 'unknown',
    message: apiLog.message ?? apiLog.msg ?? '',
    traceId: apiLog.traceId ?? apiLog.trace_id,
    spanId: apiLog.spanId ?? apiLog.span_id,
    source: apiLog.source,
    host: apiLog.host,
    metadata: apiLog.metadata ?? apiLog.meta ?? {},
    created_at: apiLog.created_at,
    updated_at: apiLog.updated_at,
  } as LogEntry;
};

export const transformLogsResponseToDomain = (apiResp: any): LogsResponse => {
  const items = Array.isArray(apiResp.items)
    ? apiResp.items.map(transformLogToDomain)
    : Array.isArray(apiResp.data)
      ? apiResp.data.map(transformLogToDomain)
      : [];

  return {
    items,
    total: apiResp.total ?? items.length,
    page: apiResp.page ?? 1,
    pageSize: apiResp.pageSize ?? items.length,
  };
};

const mapLevel = (level: any): 'info' | 'warn' | 'error' | 'debug' => {
  const l = String(level ?? '').toLowerCase();
  if (['error', 'err'].includes(l)) return 'error';
  if (['warn', 'warning'].includes(l)) return 'warn';
  if (['debug', 'dbg'].includes(l)) return 'debug';
  return 'info';
};

const cryptoRandom = (): string => Math.random().toString(36).slice(2);


