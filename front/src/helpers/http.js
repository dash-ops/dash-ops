import axios from "axios"
import { getConfigBearerToken, cleanToken } from "./oauth"

const http = axios.create({
  baseURL: process.env.REACT_APP_API_URL,
})

http.interceptors.request.use(
  config => ({ ...config, ...getConfigBearerToken() }),
  error => Promise.reject(error),
)

http.interceptors.response.use(
  response => response,
  error => {
    if (error.response.status === 401) {
      cleanToken()
    }
    return Promise.reject(error.response)
  },
)

export default http
