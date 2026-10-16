import type { RouteLocationNormalized } from 'vue-router';

export const AUTH_PAGES: string[] = [ROUTES.login, ROUTES.authCallback];
export const GUEST_PAGES: string[] = [ROUTES.home, ...AUTH_PAGES, ROUTES.access];

export default defineNuxtRouteMiddleware(async (to: RouteLocationNormalized) => {
  const authStore = useAuthStore();
  const loader = useLoader();
  const path = to.path.toLowerCase();

  loader.show();
  try {
    if (authStore.hasToken && !authStore.isLoggedIn) {
      try {
        await authStore.checkAuth();
      } catch {
        if (!GUEST_PAGES.includes(path)) {
          useCookie('redirectAfterLogin').value = to.fullPath;
          return navigateTo(ROUTES.login);
        }
      }
    }

    if (!GUEST_PAGES.includes(path) && !authStore.isLoggedIn) {
      useCookie('redirectAfterLogin').value = to.fullPath;
      return navigateTo(ROUTES.login);
    }

    if (AUTH_PAGES.includes(path) && authStore.isLoggedIn) {
      return navigateTo(ROUTES.cabinet);
    }

    return true;
  } finally {
    loader.hide();
  }
});
