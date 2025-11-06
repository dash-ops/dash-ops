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

export const getTraceDetail = async (
  traceId: string,
  provider: string
): Promise<AxiosResponse<{ trace_id: string; spans: TraceSpan[] }>> => {
  if (!provider) {
    throw new Error('Provider is required to fetch trace details');
  }
  
  // Sanitize traceId: remove leading/trailing slashes to avoid double slashes in URL
  const sanitizedTraceId = traceId.replace(/^\/+|\/+$/g, '');
  
  if (!sanitizedTraceId) {
    throw new Error('Trace ID is required');
  }
  
  // URL encode the traceId to handle special characters in base64 (/, +, =)
  const encodedTraceId = encodeURIComponent(sanitizedTraceId);
  
  const params: Record<string, unknown> = {
    provider,
  };
  const response = await http.get(`${BASE_URL}/traces/${encodedTraceId}`, { params });
  return {
    ...response,
    data: response.data.data
  };
};

export const getTraceSpans = async (
  traceId: string,
  provider: string
): Promise<AxiosResponse<TraceSpan[]>> => {
  const detail = await getTraceDetail(traceId, provider);
  return {
    ...detail,
    data: detail.data.spans || []
  };
};


