import type { RouteLocationNormalized } from 'vue-router';

export default defineNuxtRouteMiddleware(async (to: RouteLocationNormalized) => {
  const authStore = useAuthStore();
  const loader = useLoader();
  const path = to.path.toLowerCase();

  loader.show();
  try {
    await authStore.init();

    if (!GUEST_PAGES.includes(path) && !authStore.isLoggedIn) {
      rememberLoginRedirect(to.fullPath);
      return navigateTo({ path: ROUTES.login, query: { returnTo: to.fullPath } });
    }

    if (AUTH_PAGES.includes(path) && typeof to.query.returnTo === 'string') {
      rememberLoginRedirect(to.query.returnTo);
    }

    if (AUTH_PAGES.includes(path) && authStore.isLoggedIn) {
      return useAuthActions().redirectAfterLogin();
    }

    return true;
  } finally {
    loader.hide();
  }
});
