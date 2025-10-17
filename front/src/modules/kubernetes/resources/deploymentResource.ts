import http from '../../../helpers/http';
import { AxiosResponse, AxiosRequestConfig } from 'axios';
import { KubernetesTypes } from '@/types';

export function getDeployments(
  { context, namespace }: KubernetesTypes.DeploymentFilter,
  config?: AxiosRequestConfig
): Promise<AxiosResponse<KubernetesTypes.Deployment[]>> {
  let url = `/v1/k8s/clusters/${context}/deployments`;

  const filterParams = new URLSearchParams({ namespace });
  url += filterParams.toString() ? `?${filterParams.toString()}` : '';

  return http
    .get(url, config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }));
}

export function restartDeployment(
  context: string,
  name: string,
  namespace: string
): Promise<AxiosResponse<{ success: boolean }>> {
  return http.post(
    `/v1/k8s/clusters/${context}/deployment/restart/${namespace}/${name}`
  );
}

export function scaleDeployment(
  context: string,
  name: string,
  namespace: string,
  replicas: number
): Promise<AxiosResponse<{ success: boolean }>> {
  return http.post(
    `/v1/k8s/clusters/${context}/deployment/scale/${namespace}/${name}/${replicas}`
  );
}
