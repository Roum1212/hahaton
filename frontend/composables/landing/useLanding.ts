import type { ISearchArea } from '~/core/geo/services/geo/geo.types';

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

export const useLandingEvents = () => {
  const authStore = useAuthStore();
  const store = useEventsStore();
  const loading = ref(true);

  const load = async (area: ISearchArea) => {
    loading.value = true;
    if (!authStore.isLoggedIn) {
      store.upcoming = [];
      loading.value = false;
      return;
    }
    await handleAsyncOperation(
      () => store.findUpcoming(area),
      () => {
        loading.value = false;
      },
    );
  };

  return { store, loading, load };
};
