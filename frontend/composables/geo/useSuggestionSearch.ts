export const useSuggestionSearch = <T>(find: (query: string) => Promise<T[]>, prompt: string) => {
  const suggestions = shallowRef<T[]>([]);
  const searching = ref(false);
  const searched = ref(false);
  const emptyMessage = computed(() => (searching.value ? 'Ищем…' : searched.value ? 'Ничего не найдено' : prompt));

  let requestId = 0;

  const reset = () => {
    requestId++;
    searching.value = false;
    suggestions.value = [];
    searched.value = false;
  };

  const search = async ({ query }: { query: string }) => {
    if (!query.trim()) {
      reset();
      return;
    }
    const currentRequest = ++requestId;
    searching.value = true;
    searched.value = false;
    const result = await handleAsyncOperation(() => find(query.trim()));
    if (currentRequest !== requestId) return;
    searching.value = false;
    suggestions.value = result ?? [];
    searched.value = result !== undefined;
  };

  return { suggestions, searching, searched, emptyMessage, search, reset };
};
