import type { ExplorerQueryRequest, ExplorerQueryResponse, LogEntry, TraceSpan } from '../types';

export const buildExplorerRequest = (
  query: string,
  timeRange?: { from: Date; to: Date },
  provider?: string
): ExplorerQueryRequest => {
  const request: ExplorerQueryRequest = {
    query,
  };

  if (timeRange) {
    request.time_range = {
      from: timeRange.from.toISOString(),
      to: timeRange.to.toISOString(),
    };
  }

  if (provider) {
    request.provider = provider;
  }

  return request;
};

export const transformExplorerResponseToDomain = (
  apiResponse: ExplorerQueryResponse
): {
  dataSource: 'logs' | 'traces' | 'metrics';
  results: LogEntry[] | TraceSpan[] | any[];
  total: number;
  query: string;
  executionTimeMs: number;
} => {
  return {
    dataSource: apiResponse.data_source,
    results: apiResponse.results,
    total: apiResponse.total,
    query: apiResponse.query,
    executionTimeMs: apiResponse.execution_time_ms,
  };
};

