import http from "../../helpers/http"

export function getAccounts(config) {
  return http
    .get("/v1/aws/accounts", config)
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }))
}
