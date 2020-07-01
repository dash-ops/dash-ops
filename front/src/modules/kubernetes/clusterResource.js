import http from "../../helpers/http"

export function getClusters(config) {
  return http
    .get("/v1/k8s/clusters", config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}
