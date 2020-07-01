import { getPlugins } from "../modules/config/configResource"

export function loadModulesConfig() {
  return getPlugins().then(({ data }) => {
    const modulesConfig = data.map((plugin) => {
      if (plugin === "OAuth2") {
        // ToDo separate part of the authentication in the module settings
        return []
      }
      return import(`../modules/${plugin.toLowerCase()}`).then((module) => {
        if (typeof module.default === "function") {
          // Module with dynamic route loading
          return module.default().then((config) => config)
        }
        return module.default
      })
    })

    return Promise.all(modulesConfig).then((configs) => {
      let menus = []
      let routers = []
      configs.map((config) => {
        if (config.menus) {
          menus = [...menus, ...config.menus]
        }
        if (config.routers) {
          routers = [...routers, ...config.routers]
        }
        return config
      })
      return { menus, routers }
    })
  })
}
