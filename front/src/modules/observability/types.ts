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
  query?: string; // LogQL query (takes precedence over other filters)
  provider?: string;
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
  provider?: string;
  durationMinMs?: number;
  durationMaxMs?: number;
  from?: string; // ISO
  to?: string;   // ISO
  limit?: number;
  page?: number;
}

export interface TracesResponse {
  traces: TraceInfo[];
  total: number;
  page?: number;
  pageSize?: number;
}

export interface ServiceContext {
  service_name: string;
  namespace?: string;
  cluster?: string;
  labels?: Record<string, string>;
  metadata?: Record<string, unknown>;
  health?: ServiceHealth;
}

export interface ServiceWithContext extends ServiceContext {
  log_count?: number;
  metric_count?: number;
  trace_count?: number;
  alert_count?: number;
}

export interface ServiceHealth {
  status: string;
  last_check: string;
  details?: Record<string, unknown>;
  metrics?: Record<string, number>;
  alerts?: string[];
}

export interface ServicesQueryFilters {
  search?: string;
  limit?: number;
  offset?: number;
}

export interface ServicesResponse {
  services: ServiceWithContext[];
  total: number;
  pagination: {
    total: number;
    limit: number;
    offset: number;
    has_more: boolean;
  };
}

// Explorer Query Types
export interface ExplorerQueryRequest {
  query: string;
  time_range?: {
    from: string; // ISO timestamp
    to: string;   // ISO timestamp
  };
  provider?: string;
}

export interface ExplorerQueryResponse {
  data_source: 'logs' | 'traces' | 'metrics';
  results: LogEntry[] | TraceSpan[] | any[];
  total: number;
  query: string;
  execution_time_ms: number;
}

