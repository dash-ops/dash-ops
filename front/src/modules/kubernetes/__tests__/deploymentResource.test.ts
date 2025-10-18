import { describe, it, expect, vi } from 'vitest';
import http from '../../../helpers/http';
import {
  getDeployments,
  restartDeployment,
  scaleDeployment,
} from '../resources/deploymentResource';

vi.mock('../../../helpers/http');

describe('deploymentResource', () => {
  it('should return deployments list', async () => {
    const mockResponse = [
      {
        name: 'my-microservice',
        namespace: 'default',
        pod_count: 0,
      },
      {
        name: 'other-microservice',
        namespace: 'default',
        pod_count: 0,
      },
    ];
    vi.mocked(http.get).mockResolvedValue({
      data: mockResponse,
    } as any);

    const context = 'prod';
    const namespace = 'default';
    const resp = await getDeployments({ context, namespace }, {});

    expect(http.get).toHaveBeenCalledWith(
      `/v1/k8s/clusters/${context}/deployments?namespace=${namespace}`,
      {}
    );
    expect(resp.data).toEqual(mockResponse);
  });

  it('should restart deployment when restartDeployment called', async () => {
    const mockResponse = {};
    vi.mocked(http.post).mockResolvedValue({
      data: mockResponse,
    } as any);

    const context = 'prod';
    const name = 'my-microservice';
    const namespace = 'default';
    const resp = await restartDeployment(context, name, namespace);

    expect(http.post).toHaveBeenCalledWith(
      `/v1/k8s/clusters/${context}/deployment/restart/${namespace}/${name}`
    );
    expect(resp.data).toEqual(mockResponse);
  });

  it('should scale deployment when scaleDeployment called', async () => {
    const mockResponse = {};
    vi.mocked(http.post).mockResolvedValue({
      data: mockResponse,
    } as any);

    const context = 'prod';
    const name = 'my-microservice';
    const namespace = 'default';
    const replicas = 3;
    const resp = await scaleDeployment(context, name, namespace, replicas);

    expect(http.post).toHaveBeenCalledWith(
      `/v1/k8s/clusters/${context}/deployment/scale/${namespace}/${name}/${replicas}`
    );
    expect(resp.data).toEqual(mockResponse);
  });
});
