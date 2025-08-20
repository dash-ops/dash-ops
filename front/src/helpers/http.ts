import axios, {
  AxiosError,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from 'axios';
import { toast } from 'sonner';
import { getConfigBearerToken, cleanToken } from './oauth';

interface ImportMeta {
  env?: Record<string, string>;
}

const http = axios.create({
  baseURL: (import.meta as ImportMeta).env?.VITE_API_URL,
});

http.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    const bearerConfig = getConfigBearerToken();
    if (bearerConfig) {
      config.headers.Authorization = bearerConfig.headers.Authorization;
    }
    return config;
  },
  (error: AxiosError) => Promise.reject(error)
);

http.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error: unknown) => {
    if (error instanceof Error && error.message === 'Network Error') {
      return Promise.reject(error);
    }
    if (axios.isCancel(error)) {
      return Promise.reject(new Error('Request canceled'));
    }
    if (
      axios.isAxiosError(error) &&
      error.response &&
      error.response.status === 401
    ) {
      toast.error(
        `Unauthorized request: ${
          (error.response.data as { error?: string })?.error ||
          'Authentication failed'
        }`
      );
      cleanToken();
      return Promise.reject(error);
    }
    return Promise.reject(error);
  }
);

export const cancelToken = axios.CancelToken;

export default http;
