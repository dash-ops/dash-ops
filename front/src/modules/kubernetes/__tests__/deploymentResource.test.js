import http from '../../../helpers/http';
import {
  getDeployments,
  upDeployment,
  downDeployment,
} from '../deploymentResource';

vi.mock('../../../helpers/http');

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
  http.get.mockResolvedValue({
    data: mockResponse,
  });

  const context = 'prod';
  const namespace = 'default';
  const resp = await getDeployments({ context, namespace }, {});

  expect(http.get).toBeCalledWith(
    `/v1/k8s/${context}/deployments?namespace=${namespace}`,
    {}
  );
  expect(resp.data).toEqual(mockResponse);
});

it('should start pod when upDeployment called', async () => {
  const mockResponse = {};
  http.post.mockResolvedValue({
    data: mockResponse,
  });

  const context = 'prod';
  const name = 'my-microservice';
  const namespace = 'default';
  const resp = await upDeployment(context, name, namespace);

  expect(http.post).toBeCalledWith(
    `/v1/k8s/${context}/deployment/up/${namespace}/${name}`
  );
  expect(resp.data).toEqual(mockResponse);
});

it('should stop pod when downDeployment called', async () => {
  const mockResponse = {};
  http.post.mockResolvedValue({
    data: mockResponse,
  });

  const context = 'prod';
  const name = 'my-microservice';
  const namespace = 'default';
  const resp = await downDeployment(context, name, namespace);

  expect(http.post).toBeCalledWith(
    `/v1/k8s/${context}/deployment/down/${namespace}/${name}`
  );
  expect(resp.data).toEqual(mockResponse);
});
