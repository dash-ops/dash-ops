import http from "../../helpers/http"

export function getNamespaces(config) {
  return http
    .get("/v1/k8s/namespaces", config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}
