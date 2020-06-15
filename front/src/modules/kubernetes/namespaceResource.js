import http from "../../helpers/http"

export function getNamespaces() {
  return http.get("/v1/k8s/namespaces").then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}
