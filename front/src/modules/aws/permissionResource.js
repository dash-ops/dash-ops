import http from "../../helpers/http"

export function getPermissions(filter, config) {
  return http
    .get(`/v1/aws/${filter.account}/permissions`, config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}
