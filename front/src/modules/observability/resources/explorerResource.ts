import http from '../../../helpers/http';
import type { ExplorerQueryRequest, ExplorerQueryResponse } from '../types';

const BASE_URL = '/v1/observability';

export const executeQuery = async (
  request: ExplorerQueryRequest
): Promise<{ data: ExplorerQueryResponse }> => {
  const params: Record<string, unknown> = {
    query: request.query,
  };

  if (request.time_range) {
    params.time_range_from = request.time_range.from;
    params.time_range_to = request.time_range.to;
  }

  if (request.provider) {
    params.provider = request.provider;
  }

  const response = await http.get(`${BASE_URL}/explorer`, { params });
  return response.data;
};

