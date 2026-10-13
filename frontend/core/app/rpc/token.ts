const ACCESS_TOKEN_KEY = 'access_token';

export const accessToken = {
  get(): string | null {
    if (!import.meta.client) return null;
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  },
  set(token: string) {
    if (import.meta.client) localStorage.setItem(ACCESS_TOKEN_KEY, token);
  },
  clear() {
    if (import.meta.client) localStorage.removeItem(ACCESS_TOKEN_KEY);
  },
};
