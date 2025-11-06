import http from '../../../helpers/http';

const BASE_URL = '/v1/observability';

export interface ProvidersResponse {
  success: boolean;
  data: {
    logs_providers: string[];
    traces_providers: string[];
    metrics_providers: string[];
  };
}

export const getProviders = async (): Promise<ProvidersResponse> => {
  const response = await http.get<ProvidersResponse>(`${BASE_URL}/providers`);
  return response.data;
};
