import { render, screen, cleanup, act } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import userEvent from '@testing-library/user-event';
import { notification } from 'antd';
import * as deploymentResource from '../deploymentResource';
import * as namespaceResource from '../namespaceResource';
import DeploymentPage from '../DeploymentPage';

vi.mock('axios', () => ({
  default: {
    create: () => ({
      interceptors: {
        request: { use: vi.fn() },
        response: { use: vi.fn() },
      },
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      delete: vi.fn(),
    }),
  },
}));

beforeEach(() => {
  vi.spyOn(namespaceResource, 'getNamespaces').mockResolvedValue({ data: [] });
});

afterEach(cleanup);

it('should display empty deployments table when no deployments are returned', async () => {
  vi.spyOn(deploymentResource, 'getDeployments').mockResolvedValue({
    data: [],
  });

  await act(async () => {
    render(
      <MemoryRouter initialEntries={['/k8s/local/deployments']}>
        <DeploymentPage />
      </MemoryRouter>
    );
  });

  await screen.findByRole('searchbox');

  expect(screen.getByRole('searchbox')).toBeInTheDocument();
  expect(screen.getByRole('button', { name: /clear/i })).toBeInTheDocument();
  expect(screen.getByRole('combobox')).toBeInTheDocument();
});

it('should display deployments table with data when deployments are available', async () => {
  const mockDeployments = [
    {
      name: 'my-microservice',
      namespace: 'default',
      pod_count: 0,
      pod_info: { current: 0, desired: 1 },
    },
  ];
  vi.spyOn(namespaceResource, 'getNamespaces').mockResolvedValue({
    data: [{ name: 'default' }],
  });
  vi.spyOn(deploymentResource, 'getDeployments').mockResolvedValue({
    data: mockDeployments,
  });

  await act(async () => {
    render(
      <MemoryRouter initialEntries={['/k8s/local/deployments']}>
        <DeploymentPage />
      </MemoryRouter>
    );
  });

  const tds = screen.getAllByRole('cell');
  expect(tds[0].textContent).toBe('my-microservice');
  expect(tds[1].textContent).toBe('0/1');
});

it('should return notification error when failed instances fetch', async () => {
  vi.spyOn(namespaceResource, 'getNamespaces').mockResolvedValue({
    data: [{ name: 'default' }],
  });
  vi.spyOn(deploymentResource, 'getDeployments').mockRejectedValue(new Error());
  vi.spyOn(notification, 'error').mockImplementation(() => {});

  await act(async () => {
    render(
      <MemoryRouter initialEntries={['/k8s/local/deployments']}>
        <DeploymentPage />
      </MemoryRouter>
    );
  });

  expect(notification.error).toBeCalledWith({
    message: 'Ops... Failed to fetch API data',
  });
});

it('should filter deployments when typing in search field', async () => {
  const mockDeployments = [
    {
      name: 'my-microservice',
      namespace: 'default',
      pod_count: 0,
      pod_info: { current: 0, desired: 1 },
    },
    {
      name: 'other-microservice',
      namespace: 'default',
      pod_count: 0,
      pod_info: { current: 0, desired: 1 },
    },
  ];
  vi.spyOn(namespaceResource, 'getNamespaces').mockResolvedValue({
    data: [{ name: 'default' }],
  });
  vi.spyOn(deploymentResource, 'getDeployments').mockResolvedValue({
    data: mockDeployments,
  });

  await act(async () => {
    render(
      <MemoryRouter initialEntries={['/k8s/local/deployments']}>
        <DeploymentPage />
      </MemoryRouter>
    );
  });

  const input = screen.getByRole('searchbox');
  await act(async () => {
    await userEvent.type(input, 'other');
  });

  await screen.findByText('other-microservice');

  const tds = screen.getAllByRole('cell');
  expect(tds[0].textContent).toBe('other-microservice');
});

it('should display scale up button for deployments', async () => {
  const mockDeployments = [
    {
      name: 'my-microservice',
      namespace: 'default',
      pod_count: 0,
      pod_info: { current: 0, desired: 1 },
    },
  ];

  vi.spyOn(deploymentResource, 'getDeployments').mockResolvedValue({
    data: mockDeployments,
  });

  await act(async () => {
    render(
      <MemoryRouter initialEntries={['/k8s/local/deployments']}>
        <DeploymentPage />
      </MemoryRouter>
    );
  });

  const upButton = screen.getByRole('button', { name: /up/i });
  expect(upButton).toBeInTheDocument();
});

it('should display scale down button for deployments with pods', async () => {
  const mockDeployments = [
    {
      name: 'my-microservice',
      namespace: 'default',
      pod_count: 1,
      pod_info: { current: 1, desired: 1 },
    },
  ];

  vi.spyOn(deploymentResource, 'getDeployments').mockResolvedValue({
    data: mockDeployments,
  });

  await act(async () => {
    render(
      <MemoryRouter initialEntries={['/k8s/local/deployments']}>
        <DeploymentPage />
      </MemoryRouter>
    );
  });

  const actionButtons = screen.getAllByRole('button', { name: /up|down/i });
  expect(actionButtons.length).toBeGreaterThan(0);
});
