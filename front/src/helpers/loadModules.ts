import { toast } from 'sonner';
import { getPlugins } from '../modules/settings/resources/pluginsResource';
import {
  Menu,
  Router,
  AuthConfig,
  ModuleConfig,
  LoadedModulesConfig,
} from '@/types';

type ModuleOptions = Record<string, unknown> | undefined;

const pluginToFolderMap: Record<string, string> = {
  servicecatalog: 'service-catalog',
  auth: 'oauth2',
  kubernetes: 'kubernetes',
  aws: 'aws',
  observability: 'observability',
  settings: 'settings',
};

export async function loadModulesConfig(): Promise<LoadedModulesConfig> {
  const response = await getPlugins();
  const rawPlugins = Array.isArray(response.data) ? response.data : [];
  const validPlugins = rawPlugins.filter(
    (plugin) => typeof plugin === 'string' && plugin.trim() !== ''
  );

  const setupRequired = validPlugins.length === 0;
  const pluginsToLoad = setupRequired ? ['settings'] : [...validPlugins];

  if (
    !pluginsToLoad.some(
      (pluginName) => pluginName.toLowerCase() === 'settings'
    )
  ) {
    pluginsToLoad.push('settings');
  }

  const loadModule = async (
    pluginName: string,
    options?: ModuleOptions
  ): Promise<ModuleConfig> => {
      const moduleNameLower = pluginName.toLowerCase();
      const folderName = pluginToFolderMap[moduleNameLower] || moduleNameLower;

    const invokeLoader = async (
      loader: unknown,
      loaderOptions?: ModuleOptions
    ): Promise<ModuleConfig> => {
      if (typeof loader === 'function') {
        const result = (loader as (opts?: ModuleOptions) => ModuleConfig | Promise<ModuleConfig>)(loaderOptions);
        return await Promise.resolve(result);
      }
      return loader as ModuleConfig;
    };

    const tryImport = async (extension: string): Promise<ModuleConfig> => {
      const module = await import(`../modules/${folderName}/index.${extension}`);
      return invokeLoader(module.default, options);
    };

    try {
      return await tryImport('tsx');
    } catch (tsxError) {
      try {
        return await tryImport('ts');
      } catch (error) {
                  toast.error(
          `Failed to load plugin ${pluginName}: ${
            error instanceof Error ? error.message : String(error)
          }`
                  );
                  return {} as ModuleConfig;
      }
    }
  };

  const moduleConfigs = await Promise.all(
    pluginsToLoad.map((pluginName) =>
      loadModule(
        pluginName,
        pluginName.toLowerCase() === 'settings' && setupRequired
          ? { setupMode: true }
          : undefined
      )
    )
  );

      let auth: AuthConfig = { active: false };
      let menus: Menu[] = [];
      let routers: Router[] = [];

  moduleConfigs.forEach((config) => {
        if (config.auth) {
          auth = config.auth;
        }
        if (config.menus) {
          menus = [...menus, ...config.menus];
        }
        if (config.routers) {
          routers = [...routers, ...config.routers];
        }
      });

  return { auth, menus, routers, setupMode: setupRequired };
}
