import { getItem, setItem, removeItem } from "./localStorage"

const ACCESS_TOKEN_KEY = "access_token"

function getUrlVars() {
  const vars = {}
  window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, (m, key, value) => {
    vars[key] = value
  })
  return vars
}

function getUrlAccessToken() {
  if (window.location.href.indexOf(ACCESS_TOKEN_KEY) > -1) {
    return getUrlVars().access_token
  }
  return null
}

export function getToken() {
  return getItem(ACCESS_TOKEN_KEY)
}

export function verifyToken() {
  if (getToken() != null) {
    return true
  }
  const accessToken = getUrlAccessToken()
  if (accessToken != null) {
    const path = window.location.pathname === "/login" ? "/" : window.location.pathname
    window.history.pushState({ up_plugins: true }, document.title, path)
    setItem(ACCESS_TOKEN_KEY, accessToken)
    return true
  }
  return false
}

export function getConfigBearerToken() {
  const token = getToken()
  return token
    ? {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    : null
}

export function cleanToken() {
  removeItem(ACCESS_TOKEN_KEY)
}
