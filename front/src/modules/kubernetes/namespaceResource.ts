import http from '../../helpers/http';
import { AxiosResponse, AxiosRequestConfig } from 'axios';
import { KubernetesTypes, ContextFilter } from '@/types';

export function getNamespaces(
  filter: ContextFilter,
  config?: AxiosRequestConfig
): Promise<AxiosResponse<KubernetesTypes.Namespace[]>> {
  return http
    .get(`/v1/k8s/clusters/${filter.context}/namespaces`, config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }));
}
