import http from '../../../helpers/http';
import type { AxiosResponse } from 'axios';
import type { LogsQueryFilters, LogsResponse } from '../types';

const BASE_URL = '/v1/observability';

export const getLogs = async (
  filters: LogsQueryFilters
): Promise<AxiosResponse<LogsResponse>> => {
  const params: Record<string, unknown> = {
    service: filters.service,
    level: filters.level === 'all' ? undefined : filters.level,
    search: filters.search,
    provider: filters.provider,
    from: filters.from,
    to: filters.to,
    limit: filters.limit,
    page: filters.page,
  };

  const response = await http.get(`${BASE_URL}/logs`, { params });
  return {
    ...response,
    data: response.data.data
  };
};


