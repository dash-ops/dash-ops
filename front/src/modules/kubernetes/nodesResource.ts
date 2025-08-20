import http from '../../helpers/http';
import { AxiosResponse, AxiosRequestConfig } from 'axios';
import { KubernetesTypes, ContextFilter } from '@/types';

export function getNodes(
  filter: ContextFilter,
  config?: AxiosRequestConfig
): Promise<AxiosResponse<KubernetesTypes.Node[]>> {
  return http
    .get(`/v1/k8s/${filter.context}/nodes`, config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }));
}
