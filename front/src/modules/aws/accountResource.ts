import http from '../../helpers/http';
import { AxiosResponse, AxiosRequestConfig } from 'axios';
import { AWSTypes } from '@/types';

export function getAccounts(
  config?: AxiosRequestConfig
): Promise<AxiosResponse<AWSTypes.Account[]>> {
  return http
    .get('/v1/aws/accounts', config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }));
}
