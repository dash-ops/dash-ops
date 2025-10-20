import { BaseEntity } from '../../types/api';

export interface LogEntry extends BaseEntity {
  id: string;
  timestamp: string;
  level: 'info' | 'warn' | 'error' | 'debug';
  service: string;
  message: string;
  traceId?: string;
  spanId?: string;
  source?: string;
  host?: string;
  metadata?: Record<string, unknown>;
}

export interface LogsQueryFilters {
  service?: string;
  level?: LogEntry['level'] | 'all';
  search?: string;
  from?: string; // ISO
  to?: string;   // ISO
  limit?: number;
  page?: number;
}

export interface LogsResponse {
  items: LogEntry[];
  total: number;
  page: number;
  pageSize: number;
}

export interface TraceSpan extends BaseEntity {
  id: string;
  traceId: string;
  operationName: string;
  service: string;
  startTime: number;
  duration: number;
  status: 'ok' | 'error';
  tags: Record<string, unknown>;
  parentId?: string;
  depth?: number;
}

export interface TraceInfo extends BaseEntity {
  traceId: string;
  rootOperation: string;
  totalDuration: number;
  spanCount: number;
  serviceCount: number;
  status: 'ok' | 'error';
  timestamp: number;
  errors: number;
}

export interface TracesQueryFilters {
  service?: string;
  status?: 'ok' | 'error' | 'all';
  search?: string; // by traceId, operation, tags
  durationMinMs?: number;
  durationMaxMs?: number;
  from?: string; // ISO
  to?: string;   // ISO
  limit?: number;
  page?: number;
}

export interface TracesResponse {
  items: TraceInfo[];
  total: number;
  page: number;
  pageSize: number;
}

