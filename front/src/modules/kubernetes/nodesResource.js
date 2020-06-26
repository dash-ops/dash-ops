import http from "../../helpers/http"

export function getNodes(config) {
  return http
    .get("/v1/k8s/nodes", config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}
