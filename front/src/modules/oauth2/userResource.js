import http from "../../helpers/http"

export function getUserData() {
  return http.get("/v1/me")
}
