import { notification } from 'antd';
import { getPlugins } from '../modules/config/configResource';

export function loadModulesConfig() {
  return getPlugins().then(({ data }) => {
    const modulesConfig = data.map((plugin) => {
      return import(`../modules/${plugin.toLowerCase()}/index.jsx`).then(
        (module) => {
          if (typeof module.default === 'function') {
            // Module with dynamic route loading
            return module
              .default()
              .then((config) => config)
              .catch((e) => {
                notification.error(
                  `Failed to load plugin ${plugin}: ${e.data.error}`
                );
                return {};
              });
          }
          return module.default;
        }
      );
    });

    return Promise.all(modulesConfig).then((configs) => {
      let oAuth2 = { active: false };
      let menus = [];
      let routers = [];
      configs.map((config) => {
        if (config.oAuth2) {
          oAuth2 = config.oAuth2;
        }
        if (config.menus) {
          menus = [...menus, ...config.menus];
        }
        if (config.routers) {
          routers = [...routers, ...config.routers];
        }
        return config;
      });
      return { oAuth2, menus, routers };
    });
  });
}
