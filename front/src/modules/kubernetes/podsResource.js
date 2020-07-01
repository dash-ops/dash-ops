import http from "../../helpers/http"

export function getPods({ context, namespace }, config) {
  let url = `/v1/k8s/${context}/pods`

  const filterParams = new URLSearchParams({ namespace })
  url += filterParams.toString() ? `?${filterParams.toString()}` : ""

  return http.get(url, config).then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}
