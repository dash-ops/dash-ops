import http from "../../helpers/http"

export function getNodes(filter, config) {
  return http
    .get(`/v1/k8s/${filter.context}/nodes`, config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}
