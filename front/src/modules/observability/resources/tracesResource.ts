import http from '../../../helpers/http';
import type { AxiosResponse } from 'axios';
import type { TracesQueryFilters, TracesResponse, TraceSpan } from '../types';

const BASE_URL = '/v1/observability';

export const getTraces = async (
  filters: TracesQueryFilters
): Promise<AxiosResponse<TracesResponse>> => {
  const params: Record<string, unknown> = {
    service: filters.service,
    status: filters.status === 'all' ? undefined : filters.status,
    search: filters.search,
    provider: filters.provider,
    durationMinMs: filters.durationMinMs,
    durationMaxMs: filters.durationMaxMs,
    from: filters.from,
    to: filters.to,
    limit: filters.limit,
    page: filters.page,
  };

  const response = await http.get(`${BASE_URL}/traces`, { params });
  return {
    ...response,
    data: response.data.data
  };
};

export const getTraceSpans = async (
  traceId: string
): Promise<AxiosResponse<TraceSpan[]>> => {
  return http.get(`${BASE_URL}/traces/${traceId}/spans`);
};


