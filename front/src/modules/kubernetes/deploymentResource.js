import http from "../../helpers/http"

export function getDeployments({ context, namespace }, config) {
  let url = `/v1/k8s/${context}/deployments`

  const filterParams = new URLSearchParams({ namespace })
  url += filterParams.toString() ? `?${filterParams.toString()}` : ""

  return http.get(url, config).then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}

export function upDeployment(context, name, namespace) {
  return http.post(`/v1/k8s/${context}/deployment/up/${namespace}/${name}`)
}

export function downDeployment(context, name, namespace) {
  return http.post(`/v1/k8s/${context}/deployment/down/${namespace}/${name}`)
}
