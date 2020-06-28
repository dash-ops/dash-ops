import http from "../../helpers/http"

export function getPods(filter, config) {
  const filterParams = new URLSearchParams(filter)

  let url = "/v1/k8s/pods"
  url += filterParams.toString() ? `?${filterParams.toString()}` : ""

  return http.get(url, config).then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}
