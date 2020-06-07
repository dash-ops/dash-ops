import http from "../../helpers/http"

export function getInstances() {
  return http
    .get(`${process.env.REACT_APP_API_URL}/v1/ec2/instances`)
    .then(resp => (resp.data ? resp : { ...resp, data: [] }))
}

export function startInstance(instanceId) {
  return http.post(`${process.env.REACT_APP_API_URL}/v1/ec2/instance/start/${instanceId}`)
}

export function stopInstance(instanceId) {
  return http.post(`${process.env.REACT_APP_API_URL}/v1/ec2/instance/stop/${instanceId}`)
}
