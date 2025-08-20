import http from '../../helpers/http';
import { AxiosResponse } from 'axios';
import { ConfigTypes } from '@/types';

export function getPlugins(): Promise<AxiosResponse<ConfigTypes.Plugin[]>> {
  return http
    .get('/config/plugins')
    .then((resp: AxiosResponse<ConfigTypes.Plugin[]>) =>
      resp.data ? resp : { ...resp, data: [] }
    );
}
