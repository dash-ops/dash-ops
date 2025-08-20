import http from '../../../helpers/http';
import { getInstances, startInstance, stopInstance } from '../instanceResource';

vi.mock('../../../helpers/http');

it('should return instances list', async () => {
  const mockResponse = [
    {
      instance_id: '666',
      name: 'app-ops',
      status: 'stopped',
    },
  ];
  http.get.mockResolvedValue({
    data: mockResponse,
  });

  const accID = 'prod';
  const config = {};
  const resp = await getInstances({ accountKey: accID }, config);

  expect(http.get).toBeCalledWith(`/v1/aws/${accID}/ec2/instances`, config);
  expect(resp.data).toEqual(mockResponse);
});

it('should return intance status when startInstance called', async () => {
  const mockResponse = { status: 'running' };
  http.post.mockResolvedValue({
    data: mockResponse,
  });

  const accID = 'prod';
  const instanceID = 666;
  const resp = await startInstance(accID, instanceID);

  expect(http.post).toBeCalledWith(
    `/v1/aws/${accID}/ec2/instance/start/${instanceID}`
  );
  expect(resp.data).toEqual(mockResponse);
});

it('should return intance status when stopInstance called', async () => {
  const mockResponse = { status: 'stopped' };
  http.post.mockResolvedValue({
    data: mockResponse,
  });

  const accID = 'prod';
  const instanceID = 666;
  const resp = await stopInstance(accID, instanceID);

  expect(http.post).toBeCalledWith(
    `/v1/aws/${accID}/ec2/instance/stop/${instanceID}`
  );
  expect(resp.data).toEqual(mockResponse);
});
