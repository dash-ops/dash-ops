import { getItem, setItem, removeItem } from './localStorage';

const ACCESS_TOKEN_KEY = 'access_token';

interface UrlVars {
  [key: string]: string;
}

interface BearerConfig {
  headers: {
    Authorization: string;
  };
}

function getUrlVars(): UrlVars {
  const vars: UrlVars = {};
  window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, (m, key, value) => {
    vars[key] = value;
    return m;
  });
  return vars;
}

function getUrlAccessToken(): string | null {
  if (window.location.href.indexOf(ACCESS_TOKEN_KEY) > -1) {
    return getUrlVars().access_token || null;
  }
  return null;
}

export function getToken(): string | null {
  return getItem(ACCESS_TOKEN_KEY);
}

export function verifyToken(): boolean {
  if (getToken() != null) {
    return true;
  }
  const accessToken = getUrlAccessToken();
  if (accessToken != null) {
    const path =
      window.location.pathname === '/login' ? '/' : window.location.pathname;
    window.history.pushState({}, document.title, path);
    setItem(ACCESS_TOKEN_KEY, accessToken);
    return true;
  }

  // Redirect to login if no token is found
  if (window.location.pathname !== '/login') {
    const currentPath = window.location.pathname;
    window.location.href = `/login?redirect_url=${encodeURIComponent(
      currentPath
    )}`;
  }

  return false;
}

export function getConfigBearerToken(): BearerConfig | null {
  const token = getToken();
  return token
    ? {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    : null;
}

export function cleanToken(): void {
  removeItem(ACCESS_TOKEN_KEY);
}
