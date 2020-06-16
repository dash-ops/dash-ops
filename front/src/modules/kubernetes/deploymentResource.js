import http from "../../helpers/http"

export function getDeployments(filter, config) {
  const filterParams = new URLSearchParams(filter)

  let url = "/v1/k8s/deployments"
  url += filterParams.toString() ? `?${filterParams.toString()}` : ""

  return http.get(url, config).then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}

export function upDeployment(name, namespace) {
  return http.post(`/v1/k8s/deployment/up/${namespace}/${name}`)
}

export function downDeployment(name, namespace) {
  return http.post(`/v1/k8s/deployment/down/${namespace}/${name}`)
}
