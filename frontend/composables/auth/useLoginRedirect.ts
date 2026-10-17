const LOGIN_REDIRECT_KEY = 'flida_login_redirect';

export const normalizeLoginRedirect = (path: string | null | undefined): string | null => {
  if (!path?.startsWith('/')) return null;
  const base = 'https://app.invalid';
  try {
    const target = new URL(path, base);
    if (target.origin !== base || AUTH_PAGES.includes(target.pathname.toLowerCase())) return null;
    return `${target.pathname}${target.search}${target.hash}`;
  } catch {
    return null;
  }
};

export const rememberLoginRedirect = (path: string) => {
  const target = normalizeLoginRedirect(path);
  if (target) sessionStorage.setItem(LOGIN_REDIRECT_KEY, target);
};

export const takeLoginRedirect = (): string | null => {
  const target = normalizeLoginRedirect(sessionStorage.getItem(LOGIN_REDIRECT_KEY));
  sessionStorage.removeItem(LOGIN_REDIRECT_KEY);
  return target;
};
