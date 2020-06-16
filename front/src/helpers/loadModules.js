import { getPlugins } from "../modules/config/configResource"

export function loadModulesRouter() {
  return getPlugins().then(({ data }) => {
    const modulesRouters = data.map((plugin) => {
      if (plugin === "OAuth2") {
        // ToDo separate part of the authentication in the module settings
        return []
      }
      return import(`../modules/${plugin.toLowerCase()}`).then((module) => {
        return module.default.routers
      })
    })

    return Promise.all(modulesRouters).then((routers) => {
      let list = []
      routers.map((router) => {
        list = [...list, ...router]
        return router
      })
      return list
    })
  })
}
