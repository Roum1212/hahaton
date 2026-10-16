import type { RouteLocationNormalized } from 'vue-router';

export default defineNuxtRouteMiddleware(async (to: RouteLocationNormalized) => {
  const authStore = useAuthStore();
  const loader = useLoader();
  const path = to.path.toLowerCase();

  loader.show();
  try {
    await authStore.init();

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
