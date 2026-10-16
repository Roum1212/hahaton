export const useUpcomingEvents = () => {
  const store = useEventsStore();
  const loading = ref(true);

  const load = async (city: string) => {
    loading.value = true;
    await handleAsyncOperation(
      () => store.findUpcoming(city),
      () => {
        loading.value = false;
      },
    );
  };

  return { store, loading, load };
};
