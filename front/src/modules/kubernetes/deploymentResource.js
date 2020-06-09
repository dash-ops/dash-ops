import http from "../../helpers/http"

export function getDeployments() {
  return http
    .get('/v1/k8s/deployments')
    .then(resp => (resp.data ? resp : { ...resp, data: [] }))
}

export function upDeployment(name, namespace) {
  return http.post(`/v1/k8s/deployment/up/${namespace}/${name}`)
}

export function downDeployment(name, namespace) {
  return http.post(`/v1/k8s/deployment/down/${namespace}/${name}`)
}
