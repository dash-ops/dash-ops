import http from "../../helpers/http"

export function getDeployments() {
  return http
    .get(`${process.env.REACT_APP_API_URL}/v1/k8s/deployments`)
    .then(resp => (resp.data ? resp : { ...resp, data: [] }))
}

export function upDeployment(name, namespace) {
  return http.post(`${process.env.REACT_APP_API_URL}/v1/k8s/deployment/up/${namespace}/${name}`)
}

export function downDeployment(name, namespace) {
  return http.post(`${process.env.REACT_APP_API_URL}/v1/k8s/deployment/down/${namespace}/${name}`)
}
