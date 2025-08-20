import http from '../../helpers/http';
import { AxiosResponse, AxiosRequestConfig } from 'axios';
import { KubernetesTypes } from '@/types';

export function getDeployments(
  { context, namespace }: KubernetesTypes.DeploymentFilter,
  config?: AxiosRequestConfig
): Promise<AxiosResponse<KubernetesTypes.Deployment[]>> {
  let url = `/v1/k8s/${context}/deployments`;

  const filterParams = new URLSearchParams({ namespace });
  url += filterParams.toString() ? `?${filterParams.toString()}` : '';

  return http
    .get(url, config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }));
}

export function upDeployment(
  context: string,
  name: string,
  namespace: string
): Promise<AxiosResponse<{ success: boolean }>> {
  return http.post(`/v1/k8s/${context}/deployment/up/${namespace}/${name}`);
}

export function downDeployment(
  context: string,
  name: string,
  namespace: string
): Promise<AxiosResponse<{ success: boolean }>> {
  return http.post(`/v1/k8s/${context}/deployment/down/${namespace}/${name}`);
}
