import http from '../../helpers/http';
import { AxiosResponse, AxiosRequestConfig } from 'axios';
import { AWSTypes } from '@/types';
import { AccountFilter } from '@/types';

export function getInstances(
  filter: AccountFilter,
  config?: AxiosRequestConfig
): Promise<AxiosResponse<AWSTypes.Instance[]>> {
  return http
    .get(`/v1/aws/${filter.accountKey}/ec2/instances`, config)
    .then((resp) => {
      if (resp.data && resp.data.instances) {
        return { ...resp, data: resp.data.instances };
      }
      return { ...resp, data: [] };
    });
}

export function startInstance(
  accountKey: string,
  instanceId: string
): Promise<AxiosResponse<{ current_state: string }>> {
  return http.post(`/v1/aws/${accountKey}/ec2/instance/start/${instanceId}`);
}

export function stopInstance(
  accountKey: string,
  instanceId: string
): Promise<AxiosResponse<{ current_state: string }>> {
  return http.post(`/v1/aws/${accountKey}/ec2/instance/stop/${instanceId}`);
}
