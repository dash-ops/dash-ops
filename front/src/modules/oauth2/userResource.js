import http from "../../helpers/http"

export function getUserData() {
  return http.get(`${process.env.REACT_APP_API_URL}/v1/me`)
}
