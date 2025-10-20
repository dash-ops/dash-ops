import http from '../../../helpers/http';
import { AxiosResponse } from 'axios';

export function getPlugins(): Promise<AxiosResponse<string[]>> {
  return http
    .get('/config/plugins')
    .then((resp: AxiosResponse<string[]>) =>
      resp.data ? resp : { ...resp, data: [] }
    );
}
