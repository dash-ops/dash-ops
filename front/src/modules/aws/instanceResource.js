import http from "../../helpers/http"

export function getInstances(filter, config) {
  return http
    .get(`/v1/aws/${filter.accountKey}/ec2/instances`, config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}

export function startInstance(accountKey, instanceId) {
  return http.post(`/v1/aws/${accountKey}/ec2/instance/start/${instanceId}`)
}

export function stopInstance(accountKey, instanceId) {
  return http.post(`/v1/aws/${accountKey}/ec2/instance/stop/${instanceId}`)
}
