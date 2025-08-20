export function setItem(key: string, value: string): void {
  localStorage.setItem(`dash-ops:${key}`, value);
}

export function getItem(key: string): string | null {
  return localStorage.getItem(`dash-ops:${key}`);
}

export function removeItem(key: string): void {
  localStorage.removeItem(`dash-ops:${key}`);
}
