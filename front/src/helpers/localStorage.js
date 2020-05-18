export function setItem(key, value) {
  localStorage.setItem(`dash-ops:${key}`, value)
}

export function getItem(key) {
  return localStorage.getItem(`dash-ops:${key}`)
}

export function removeItem(key) {
  localStorage.removeItem(`dash-ops:${key}`)
}
