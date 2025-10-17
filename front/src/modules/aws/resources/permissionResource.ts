import http from '../../../helpers/http';
import { AxiosResponse, AxiosRequestConfig } from 'axios';
import { AWSTypes, AccountFilter } from '@/types';

export function getPermissions(
  filter: AccountFilter,
  config?: AxiosRequestConfig
): Promise<AxiosResponse<AWSTypes.AWSPermission[]>> {
  return http
    .get(`/v1/aws/${filter.account}/permissions`, config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }));
}
