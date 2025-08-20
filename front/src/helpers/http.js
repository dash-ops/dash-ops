import axios from 'axios';
import { toast } from 'sonner';
import { getConfigBearerToken, cleanToken } from './oauth';

const http = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
});

http.interceptors.request.use(
  (config) => ({ ...config, ...getConfigBearerToken() }),
  (error) => Promise.reject(error)
);

http.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.message === 'Network Error') {
      return Promise.reject(error);
    }
    if (axios.isCancel(error)) {
      return Promise.reject(new Error('Request canceled'));
    }
    if (error.response && error.response.status === 401) {
      toast.error(`Unauthorized request: ${error.response.data.error}`);
      cleanToken();
      return Promise.reject(error);
    }
    return Promise.reject(error);
  }
);

export const cancelToken = axios.CancelToken;

export default http;
