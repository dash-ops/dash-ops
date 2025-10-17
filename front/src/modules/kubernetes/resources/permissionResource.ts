import http from '../../../helpers/http';
import { AxiosResponse, AxiosRequestConfig } from 'axios';
import { KubernetesTypes, ContextFilter } from '@/types';

export function getPermissions(
  filter: ContextFilter,
  config?: AxiosRequestConfig
): Promise<AxiosResponse<KubernetesTypes.K8sPermission[]>> {
  return http
    .get(`/v1/k8s/clusters/${filter.context}/permissions`, config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }));
}
