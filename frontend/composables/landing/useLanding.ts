export const useLanding = () => {
  const authStore = useAuthStore();
  const siteName = useRuntimeConfig().public.siteName;

  const isLoggedIn = computed(() => authStore.isLoggedIn);
  const primaryRoute = computed(() => (authStore.isLoggedIn ? ROUTES.cabinet : ROUTES.login));

  const scrollTo = (anchor: string) => {
    if (!import.meta.client) return;
    document.getElementById(anchor)?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  };

  return { siteName, isLoggedIn, primaryRoute, scrollTo };
};
