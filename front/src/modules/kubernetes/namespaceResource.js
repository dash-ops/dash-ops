import http from "../../helpers/http"

export function getNamespaces(filter, config) {
  return http
    .get(`/v1/k8s/${filter.context}/namespaces`, config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}
