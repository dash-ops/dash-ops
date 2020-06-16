import http from "../../helpers/http"

export function getInstances(config) {
  return http
    .get("/v1/ec2/instances", config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}

export function startInstance(instanceId) {
  return http.post(`/v1/ec2/instance/start/${instanceId}`)
}

export function stopInstance(instanceId) {
  return http.post(`/v1/ec2/instance/stop/${instanceId}`)
}
