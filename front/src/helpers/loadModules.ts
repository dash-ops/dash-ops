import { toast } from 'sonner';
import { getPlugins } from '../modules/config/configResource';
import {
  Menu,
  Router,
  AuthConfig,
  ModuleConfig,
  LoadedModulesConfig,
} from '@/types';

// LoadedModulesConfig is now imported from @/types

export function loadModulesConfig(): Promise<LoadedModulesConfig> {
  return getPlugins().then(({ data }) => {
    // Filter valid plugin names
    const validPlugins = data.filter(
      (plugin) => plugin && plugin.trim() !== ''
    );

    const modulesConfig = validPlugins.map((pluginName) => {
      // Map plugin names to their actual folder names
      const pluginToFolderMap: Record<string, string> = {
        servicecatalog: 'service-catalog',
        oauth2: 'oauth2',
        kubernetes: 'kubernetes',
        aws: 'aws',
        config: 'config',
      };

      const moduleNameLower = pluginName.toLowerCase();
      const folderName = pluginToFolderMap[moduleNameLower] || moduleNameLower;

      // Try .tsx first, then .ts
      const tryImport = (extension: string): Promise<ModuleConfig> => {
        return import(`../modules/${folderName}/index.${extension}`).then(
          (module) => {
            if (typeof module.default === 'function') {
              // Module with dynamic route loading
              return module
                .default()
                .then((config: ModuleConfig) => config)
                .catch((e: Error) => {
                  toast.error(
                    `Failed to load plugin ${pluginName}: ${e.message}`
                  );
                  return {} as ModuleConfig;
                });
            }
            return module.default as ModuleConfig;
          }
        );
      };

      return tryImport('tsx').catch(() => tryImport('ts'));
    });

    return Promise.all(modulesConfig).then((configs) => {
      let auth: AuthConfig = { active: false };
      let menus: Menu[] = [];
      let routers: Router[] = [];

      configs.forEach((config) => {
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

      return { auth, menus, routers };
    });
  });
}
