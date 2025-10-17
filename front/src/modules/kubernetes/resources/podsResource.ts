import http from '../../../helpers/http';
import { AxiosResponse, AxiosRequestConfig } from 'axios';
import { KubernetesTypes } from '@/types';

export function getPods(
  { context, namespace }: KubernetesTypes.PodFilter,
  config?: AxiosRequestConfig
): Promise<AxiosResponse<KubernetesTypes.Pod[]>> {
  let url = `/v1/k8s/clusters/${context}/pods`;

  const filterParams = new URLSearchParams({ namespace });
  url += filterParams.toString() ? `?${filterParams.toString()}` : '';

  return http
    .get(url, config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }));
}

export function getPodLogs(
  { context, name, namespace }: KubernetesTypes.PodLogsFilter,
  config?: AxiosRequestConfig
): Promise<AxiosResponse<KubernetesTypes.PodLogsResponse>> {
  let url = `/v1/k8s/clusters/${context}/namespaces/${namespace}/pods/${name}/logs`;

  return http.get(url, config);
}
