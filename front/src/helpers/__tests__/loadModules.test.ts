import type { AxiosResponse } from 'axios';
import type { ReactElement } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
  toastErrorMock: vi.fn(),
  getPluginsMock: vi.fn<
    () => Promise<AxiosResponse<string[] | null | undefined>>
  >(),
  settingsLoader: vi.fn(),
  serviceCatalogLoader: vi.fn(),
  observabilityLoader: vi.fn(),
  authLoader: vi.fn(),
  kubernetesLoader: vi.fn(),
  awsLoader: vi.fn(),
}));

vi.mock('sonner', () => ({
  toast: {
    error: mocks.toastErrorMock,
  },
}));

vi.mock('../../modules/settings/resources/pluginsResource', () => ({
  getPlugins: mocks.getPluginsMock,
}));

vi.mock('../../modules/settings/index.tsx', () => ({
  default: mocks.settingsLoader,
}));

vi.mock('../../modules/service-catalog/index.tsx', () => ({
  default: mocks.serviceCatalogLoader,
}));

vi.mock('../../modules/observability/index.tsx', () => ({
  default: mocks.observabilityLoader,
}));

vi.mock('../../modules/oauth2/index.tsx', () => ({
  default: mocks.authLoader,
}));

vi.mock('../../modules/kubernetes/index.tsx', () => ({
  default: mocks.kubernetesLoader,
}));

vi.mock('../../modules/aws/index.tsx', () => ({
  default: mocks.awsLoader,
}));

import { loadModulesConfig } from '../loadModules';

const {
  toastErrorMock,
  getPluginsMock,
  settingsLoader,
  serviceCatalogLoader,
  observabilityLoader,
  authLoader,
  kubernetesLoader,
  awsLoader,
} = mocks;

function createResponse(data: unknown): AxiosResponse<string[] | null> {
  return {
    data: data as string[] | null,
    status: 200,
    statusText: 'OK',
    headers: {},
    config: {} as unknown as AxiosResponse['config'],
  };
}

describe('loadModulesConfig', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    getPluginsMock.mockReset();
    settingsLoader.mockReset();
    serviceCatalogLoader.mockReset();
    observabilityLoader.mockReset();
    authLoader.mockReset();
    kubernetesLoader.mockReset();
    awsLoader.mockReset();
  });

  describe('Setup Mode Detection', () => {
    it('activates setup mode when no plugins are returned', async () => {
      const icon = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(createResponse([]));

      settingsLoader.mockImplementationOnce((options?: Record<string, unknown>) => {
        expect(options).toEqual({ setupMode: true });
        return {
          auth: { active: false },
          menus: [{ label: 'Settings', icon, key: 'settings', link: '/settings' }],
          routers: [{ key: 'settings', path: '/settings', element: icon }],
        };
      });

      const result = await loadModulesConfig();

      expect(result.setupMode).toBe(true);
      expect(result.menus).toHaveLength(1);
      expect(result.routers).toHaveLength(1);
      expect(result.auth).toEqual({ active: false });
      expect(settingsLoader).toHaveBeenCalledTimes(1);
      expect(settingsLoader).toHaveBeenCalledWith({ setupMode: true });
      expect(serviceCatalogLoader).not.toHaveBeenCalled();
      expect(observabilityLoader).not.toHaveBeenCalled();
      expect(toastErrorMock).not.toHaveBeenCalled();
    });

    it('activates setup mode when plugins array is null', async () => {
      const icon = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(createResponse(null));

      settingsLoader.mockImplementationOnce((options?: Record<string, unknown>) => {
        expect(options).toEqual({ setupMode: true });
        return {
          auth: { active: false },
          menus: [{ label: 'Settings', icon, key: 'settings', link: '/settings' }],
          routers: [{ key: 'settings', path: '/settings', element: icon }],
        };
      });

      const result = await loadModulesConfig();

      expect(result.setupMode).toBe(true);
      expect(settingsLoader).toHaveBeenCalledWith({ setupMode: true });
    });

    it('filters out empty or whitespace-only plugin names', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['servicecatalog', '   ', '', 'settings'])
      );

      serviceCatalogLoader.mockImplementationOnce(() => ({
        menus: [
          { label: 'Service Catalog', icon: element, key: 'catalog', link: '/catalog' },
        ],
        routers: [{ key: 'catalog', path: '/catalog', element }],
      }));

      settingsLoader.mockImplementationOnce(() => ({
        auth: { active: true },
        menus: [{ label: 'Settings', icon: element, key: 'settings', link: '/settings' }],
        routers: [{ key: 'settings', path: '/settings', element }],
      }));

      const result = await loadModulesConfig();

      expect(result.setupMode).toBe(false);
      expect(serviceCatalogLoader).toHaveBeenCalledTimes(1);
      expect(settingsLoader).toHaveBeenCalledTimes(1);
      // Should only load 2 valid plugins (servicecatalog + settings)
      expect(result.menus).toHaveLength(2);
    });
  });

  describe('Plugin Loading', () => {
    it('loads declared plugins and automatically appends settings module', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['servicecatalog'])
      );

      serviceCatalogLoader.mockImplementationOnce(
        (options?: Record<string, unknown>) => {
          expect(options).toBeUndefined();
          return {
            menus: [
              { label: 'Service Catalog', icon: element, key: 'catalog', link: '/catalog' },
            ],
            routers: [
              { key: 'catalog', path: '/catalog', element },
            ],
          };
        }
      );

      settingsLoader.mockImplementationOnce((options?: Record<string, unknown>) => {
        expect(options).toBeUndefined();
        return {
          auth: { active: true },
          menus: [{ label: 'Settings', icon: element, key: 'settings', link: '/settings' }],
          routers: [{ key: 'settings', path: '/settings', element }],
        };
      });

      const result = await loadModulesConfig();

      expect(result.setupMode).toBe(false);
      expect(result.auth).toEqual({ active: true });
      expect(result.menus.map((menu) => menu.key)).toEqual(['catalog', 'settings']);
      expect(result.routers.map((router) => router.key)).toEqual(['catalog', 'settings']);
      expect(settingsLoader).toHaveBeenCalledTimes(1);
      expect(serviceCatalogLoader).toHaveBeenCalledTimes(1);
      expect(toastErrorMock).not.toHaveBeenCalled();
    });

    it('does not duplicate settings module if already declared', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['servicecatalog', 'settings'])
      );

      serviceCatalogLoader.mockImplementationOnce(() => ({
        menus: [
          { label: 'Service Catalog', icon: element, key: 'catalog', link: '/catalog' },
        ],
        routers: [{ key: 'catalog', path: '/catalog', element }],
      }));

      settingsLoader.mockImplementationOnce(() => ({
        auth: { active: true },
        menus: [{ label: 'Settings', icon: element, key: 'settings', link: '/settings' }],
        routers: [{ key: 'settings', path: '/settings', element }],
      }));

      const result = await loadModulesConfig();

      expect(result.setupMode).toBe(false);
      expect(settingsLoader).toHaveBeenCalledTimes(1);
      expect(result.menus).toHaveLength(2);
    });

    it('loads multiple plugins with different module types', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['servicecatalog', 'observability', 'kubernetes'])
      );

      serviceCatalogLoader.mockImplementationOnce(() => ({
        menus: [
          { label: 'Service Catalog', icon: element, key: 'catalog', link: '/catalog' },
        ],
        routers: [{ key: 'catalog', path: '/catalog', element }],
      }));

      observabilityLoader.mockImplementationOnce(() => ({
        menus: [
          { label: 'Observability', icon: element, key: 'observability', link: '/observability' },
        ],
        routers: [{ key: 'observability', path: '/observability', element }],
      }));

      kubernetesLoader.mockImplementationOnce(() => ({
        menus: [
          { label: 'Kubernetes', icon: element, key: 'kubernetes', link: '/kubernetes' },
        ],
        routers: [{ key: 'kubernetes', path: '/kubernetes', element }],
      }));

      settingsLoader.mockImplementationOnce(() => ({
        auth: { active: false },
        menus: [{ label: 'Settings', icon: element, key: 'settings', link: '/settings' }],
        routers: [{ key: 'settings', path: '/settings', element }],
      }));

      const result = await loadModulesConfig();

      expect(result.setupMode).toBe(false);
      expect(result.menus.map((menu) => menu.key)).toEqual([
        'catalog',
        'observability',
        'kubernetes',
        'settings',
      ]);
      expect(serviceCatalogLoader).toHaveBeenCalledTimes(1);
      expect(observabilityLoader).toHaveBeenCalledTimes(1);
      expect(kubernetesLoader).toHaveBeenCalledTimes(1);
      expect(settingsLoader).toHaveBeenCalledTimes(1);
    });

    it('uses pluginToFolderMap to resolve folder names', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['servicecatalog', 'aws'])
      );

      // 'servicecatalog' should map to 'service-catalog' folder
      serviceCatalogLoader.mockImplementationOnce(() => ({
        menus: [{ label: 'Service Catalog', icon: element, key: 'catalog', link: '/catalog' }],
        routers: [{ key: 'catalog', path: '/catalog', element }],
      }));

      // 'aws' should map to 'aws' folder
      awsLoader.mockImplementationOnce(() => ({
        menus: [{ label: 'AWS', icon: element, key: 'aws', link: '/aws' }],
        routers: [{ key: 'aws', path: '/aws', element }],
      }));

      settingsLoader.mockImplementationOnce(() => ({
        auth: { active: false },
        menus: [],
        routers: [],
      }));

      const result = await loadModulesConfig();

      expect(result.setupMode).toBe(false);
      expect(serviceCatalogLoader).toHaveBeenCalledTimes(1);
      expect(awsLoader).toHaveBeenCalledTimes(1);
      expect(result.menus.map((m) => m.key)).toEqual(['catalog', 'aws']);
    });
  });

  describe('Module Config Aggregation', () => {
    it('merges menus from multiple modules in loading order', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['servicecatalog', 'kubernetes'])
      );

      serviceCatalogLoader.mockImplementationOnce(() => ({
        menus: [
          { label: 'Catalog 1', icon: element, key: 'catalog1', link: '/catalog1' },
          { label: 'Catalog 2', icon: element, key: 'catalog2', link: '/catalog2' },
        ],
        routers: [],
      }));

      kubernetesLoader.mockImplementationOnce(() => ({
        menus: [
          { label: 'K8s 1', icon: element, key: 'k8s1', link: '/k8s1' },
        ],
        routers: [],
      }));

      settingsLoader.mockImplementationOnce(() => ({
        menus: [{ label: 'Settings', icon: element, key: 'settings', link: '/settings' }],
        routers: [],
      }));

      const result = await loadModulesConfig();

      expect(result.menus.map((m) => m.key)).toEqual([
        'catalog1',
        'catalog2',
        'k8s1',
        'settings',
      ]);
    });

    it('merges routers from multiple modules in loading order', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['servicecatalog', 'kubernetes'])
      );

      serviceCatalogLoader.mockImplementationOnce(() => ({
        menus: [],
        routers: [
          { key: 'catalog1', path: '/catalog1', element },
          { key: 'catalog2', path: '/catalog2', element },
        ],
      }));

      kubernetesLoader.mockImplementationOnce(() => ({
        menus: [],
        routers: [{ key: 'k8s1', path: '/k8s1', element }],
      }));

      settingsLoader.mockImplementationOnce(() => ({
        menus: [],
        routers: [{ key: 'settings', path: '/settings', element }],
      }));

      const result = await loadModulesConfig();

      expect(result.routers.map((r) => r.key)).toEqual([
        'catalog1',
        'catalog2',
        'k8s1',
        'settings',
      ]);
    });

    it('uses the last declared auth config when multiple modules provide it', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['servicecatalog', 'auth'])
      );

      serviceCatalogLoader.mockImplementationOnce(() => ({
        auth: { active: false },
        menus: [],
        routers: [],
      }));

      authLoader.mockImplementationOnce(() => ({
        auth: { active: true, LoginPage: () => element },
        menus: [],
        routers: [],
      }));

      settingsLoader.mockImplementationOnce(() => ({
        menus: [],
        routers: [],
      }));

      const result = await loadModulesConfig();

      expect(result.auth.active).toBe(true);
      expect(result.auth.LoginPage).toBeDefined();
    });

    it('handles modules with no menus or routers', async () => {
      getPluginsMock.mockResolvedValueOnce(
        createResponse(['servicecatalog'])
      );

      serviceCatalogLoader.mockImplementationOnce(() => ({}));

      settingsLoader.mockImplementationOnce(() => ({
        menus: [],
        routers: [],
      }));

      const result = await loadModulesConfig();

      expect(result.menus).toEqual([]);
      expect(result.routers).toEqual([]);
      expect(result.auth).toEqual({ active: false });
    });
  });

  describe('Error Handling', () => {
    it('emits a toast error when a plugin module fails to load', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(createResponse(['badPlugin']));

      settingsLoader.mockImplementationOnce(() => ({
        auth: { active: true },
        menus: [{ label: 'Settings', icon: element, key: 'settings', link: '/settings' }],
        routers: [{ key: 'settings', path: '/settings', element }],
      }));

      const result = await loadModulesConfig();

      expect(result.setupMode).toBe(false);
      expect(result.auth).toEqual({ active: true });
      expect(result.menus).toHaveLength(1);
      expect(settingsLoader).toHaveBeenCalledTimes(1);
      expect(toastErrorMock).toHaveBeenCalledTimes(1);
      expect(toastErrorMock.mock.calls[0][0]).toContain('Failed to load plugin badPlugin');
    });

    it('returns empty config for failed module and continues loading other modules', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['badPlugin', 'servicecatalog'])
      );

      serviceCatalogLoader.mockImplementationOnce(() => ({
        menus: [
          { label: 'Service Catalog', icon: element, key: 'catalog', link: '/catalog' },
        ],
        routers: [{ key: 'catalog', path: '/catalog', element }],
      }));

      settingsLoader.mockImplementationOnce(() => ({
        menus: [{ label: 'Settings', icon: element, key: 'settings', link: '/settings' }],
        routers: [{ key: 'settings', path: '/settings', element }],
      }));

      const result = await loadModulesConfig();

      // badPlugin should fail but servicecatalog and settings should load
      expect(result.menus.map((m) => m.key)).toEqual(['catalog', 'settings']);
      expect(toastErrorMock).toHaveBeenCalledTimes(1);
      expect(serviceCatalogLoader).toHaveBeenCalledTimes(1);
      expect(settingsLoader).toHaveBeenCalledTimes(1);
    });

    it('handles multiple failed plugins gracefully', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['badPlugin1', 'badPlugin2', 'servicecatalog'])
      );

      serviceCatalogLoader.mockImplementationOnce(() => ({
        menus: [
          { label: 'Service Catalog', icon: element, key: 'catalog', link: '/catalog' },
        ],
        routers: [{ key: 'catalog', path: '/catalog', element }],
      }));

      settingsLoader.mockImplementationOnce(() => ({
        menus: [],
        routers: [],
      }));

      const result = await loadModulesConfig();

      expect(result.menus).toHaveLength(1);
      expect(toastErrorMock).toHaveBeenCalledTimes(2);
      expect(toastErrorMock.mock.calls[0][0]).toContain('badPlugin1');
      expect(toastErrorMock.mock.calls[1][0]).toContain('badPlugin2');
    });
  });

  describe('Module Loader Invocation', () => {
    it('invokes function-based loaders with options', async () => {
      const element = {} as ReactElement;
      const loaderSpy = vi.fn(() => ({
        menus: [{ label: 'Test', icon: element, key: 'test', link: '/test' }],
        routers: [],
      }));

      getPluginsMock.mockResolvedValueOnce(createResponse([]));
      settingsLoader.mockImplementationOnce(loaderSpy);

      await loadModulesConfig();

      expect(loaderSpy).toHaveBeenCalledWith({ setupMode: true });
    });

    it('invokes async function-based loaders correctly', async () => {
      const element = {} as ReactElement;
      const asyncLoader = vi.fn(async () => {
        await new Promise((resolve) => setTimeout(resolve, 10));
        return {
          menus: [{ label: 'Test', icon: element, key: 'test', link: '/test' }],
          routers: [],
        };
      });

      getPluginsMock.mockResolvedValueOnce(createResponse(['servicecatalog']));
      serviceCatalogLoader.mockImplementationOnce(asyncLoader);
      settingsLoader.mockImplementationOnce(() => ({ menus: [], routers: [] }));

      const result = await loadModulesConfig();

      expect(asyncLoader).toHaveBeenCalledTimes(1);
      expect(result.menus).toHaveLength(1);
    });

    it('handles non-function module exports as static config', async () => {
      const element = {} as ReactElement;
      const staticConfig = {
        menus: [{ label: 'Static', icon: element, key: 'static', link: '/static' }],
        routers: [],
      };

      getPluginsMock.mockResolvedValueOnce(createResponse(['servicecatalog']));
      serviceCatalogLoader.mockReturnValueOnce(staticConfig);
      settingsLoader.mockImplementationOnce(() => ({ menus: [], routers: [] }));

      const result = await loadModulesConfig();

      expect(result.menus).toHaveLength(1);
      expect(result.menus[0].key).toBe('static');
    });
  });

  describe('Edge Cases', () => {
    it('handles plugin names with mixed case correctly', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['ServiceCatalog', 'KUBERNETES'])
      );

      serviceCatalogLoader.mockImplementationOnce(() => ({
        menus: [{ label: 'Catalog', icon: element, key: 'catalog', link: '/catalog' }],
        routers: [],
      }));

      kubernetesLoader.mockImplementationOnce(() => ({
        menus: [{ label: 'K8s', icon: element, key: 'k8s', link: '/k8s' }],
        routers: [],
      }));

      settingsLoader.mockImplementationOnce(() => ({
        menus: [],
        routers: [],
      }));

      const result = await loadModulesConfig();

      expect(result.menus).toHaveLength(2);
      expect(serviceCatalogLoader).toHaveBeenCalledTimes(1);
      expect(kubernetesLoader).toHaveBeenCalledTimes(1);
    });

    it('handles response.data as undefined gracefully', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce({
        data: undefined as unknown as string[],
        status: 200,
        statusText: 'OK',
        headers: {},
        config: {} as unknown as AxiosResponse['config'],
      });

      settingsLoader.mockImplementationOnce(() => ({
        menus: [{ label: 'Settings', icon: element, key: 'settings', link: '/settings' }],
        routers: [],
      }));

      const result = await loadModulesConfig();

      expect(result.setupMode).toBe(true);
      expect(settingsLoader).toHaveBeenCalledWith({ setupMode: true });
    });

    it('handles non-string values in plugins array', async () => {
      const element = {} as ReactElement;

      getPluginsMock.mockResolvedValueOnce(
        createResponse(['servicecatalog', 123, null, undefined, true] as unknown as string[])
      );

      serviceCatalogLoader.mockImplementationOnce(() => ({
        menus: [{ label: 'Catalog', icon: element, key: 'catalog', link: '/catalog' }],
        routers: [],
      }));

      settingsLoader.mockImplementationOnce(() => ({
        menus: [],
        routers: [],
      }));

      const result = await loadModulesConfig();

      // Should only load servicecatalog (valid string) + settings
      expect(result.menus).toHaveLength(1);
      expect(serviceCatalogLoader).toHaveBeenCalledTimes(1);
    });
  });
});
