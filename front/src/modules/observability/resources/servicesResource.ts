import http from '../../../helpers/http';
import type { AxiosResponse } from 'axios';
import type { ServicesQueryFilters, ServicesResponse } from '../types';

const BASE_URL = '/v1/observability';

export const getServices = async (
  filters: ServicesQueryFilters
): Promise<AxiosResponse<ServicesResponse>> => {
  const params: Record<string, unknown> = {
    search: filters.search,
    limit: filters.limit,
    offset: filters.offset,
  };

  const response = await http.get(`${BASE_URL}/services`, { params });
  return {
    ...response,
    data: response.data.data,
  };
};

