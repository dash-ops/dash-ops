import { describe, it, expect, vi } from 'vitest';
import http from '../../../helpers/http';
import { getInstances, startInstance, stopInstance } from '../resources/instanceResource';

vi.mock('../../../helpers/http');

describe('instanceResource', () => {
  it('should return instances list', async () => {
    const mockResponse = [
      {
        instance_id: '666',
        name: 'app-ops',
        status: 'stopped',
      },
    ];
    vi.mocked(http.get).mockResolvedValue({
      data: { instances: mockResponse },
    } as any);

    const accID = 'prod';
    const config = {};
    const resp = await getInstances({ accountKey: accID }, config);

    expect(http.get).toHaveBeenCalledWith(`/v1/aws/${accID}/ec2/instances`, config);
    expect(resp.data).toEqual(mockResponse);
  });

  it('should return instance status when startInstance called', async () => {
    const mockResponse = { status: 'running' };
    vi.mocked(http.post).mockResolvedValue({
      data: mockResponse,
    } as any);

    const accID = 'prod';
    const instanceID = 666;
    const resp = await startInstance(accID, instanceID);

    expect(http.post).toHaveBeenCalledWith(
      `/v1/aws/${accID}/ec2/instance/start/${instanceID}`
    );
    expect(resp.data).toEqual(mockResponse);
  });

  it('should return instance status when stopInstance called', async () => {
    const mockResponse = { status: 'stopped' };
    vi.mocked(http.post).mockResolvedValue({
      data: mockResponse,
    } as any);

    const accID = 'prod';
    const instanceID = 666;
    const resp = await stopInstance(accID, instanceID);

    expect(http.post).toHaveBeenCalledWith(
      `/v1/aws/${accID}/ec2/instance/stop/${instanceID}`
    );
    expect(resp.data).toEqual(mockResponse);
  });
});
