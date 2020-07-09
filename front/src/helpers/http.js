import axios from "axios"
import { notification } from "antd"
import { getConfigBearerToken, cleanToken } from "./oauth"

const http = axios.create({
  baseURL: process.env.REACT_APP_API_URL,
})

http.interceptors.request.use(
  (config) => ({ ...config, ...getConfigBearerToken() }),
  (error) => Promise.reject(error),
)

http.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.message === "Network Error") {
      return Promise.reject(error.message)
    }
    if (axios.isCancel(error)) {
      return Promise.reject("Request canceled")
    }
    if (error.response.status === 401) {
      notification.error({
        message: "Unauthorized request",
        description: error.response.data.error,
      })
      cleanToken()
    }
    return Promise.reject(error.response)
  },
)

export const cancelToken = axios.CancelToken

export default http
