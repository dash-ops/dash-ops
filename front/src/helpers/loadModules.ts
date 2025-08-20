import { toast } from 'sonner';
import { getPlugins } from '../modules/config/configResource';
import { Menu, Router, OAuth2Config, ModuleConfig } from '@/types';

interface LoadedModulesConfig {
  oAuth2: OAuth2Config;
  menus: Menu[];
  routers: Router[];
}

export function loadModulesConfig(): Promise<LoadedModulesConfig> {
  return getPlugins().then(({ data }) => {
    const modulesConfig = data.map((plugin) => {
      const pluginName = plugin.name?.toLowerCase() || 'unknown';

      // Try .tsx first, then .jsx, then .js
      const tryImport = (extension: string): Promise<ModuleConfig> => {
        return import(`../modules/${pluginName}/index.${extension}`).then(
          (module) => {
            if (typeof module.default === 'function') {
              // Module with dynamic route loading
              return module
                .default()
                .then((config: ModuleConfig) => config)
                .catch((e: Error) => {
                  toast.error(
                    `Failed to load plugin ${plugin.name || 'unknown'}: ${e.message}`
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
      let oAuth2: OAuth2Config = { active: false };
      let menus: Menu[] = [];
      let routers: Router[] = [];

      configs.forEach((config) => {
        if (config.oAuth2) {
          oAuth2 = config.oAuth2;
        }
        if (config.menus) {
          menus = [...menus, ...config.menus];
        }
        if (config.routers) {
          routers = [...routers, ...config.routers];
        }
      });

      return { oAuth2, menus, routers };
    });
  });
}
