import http from '../../helpers/http';
import { AxiosResponse, AxiosRequestConfig } from 'axios';
import { KubernetesTypes } from '@/types';

export function getClusters(
  config?: AxiosRequestConfig
): Promise<AxiosResponse<KubernetesTypes.Cluster[]>> {
  return http
    .get('/v1/k8s/clusters', config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }));
}
