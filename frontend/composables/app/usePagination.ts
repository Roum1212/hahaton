/* eslint-disable @typescript-eslint/no-explicit-any */
export const usePagination = <T>(params: { page?: number; limit?: number } = { page: 1, limit: 50 }) => {
  const searchQuery = ref<string>('');
  const page = ref<number>(params?.page ?? 1);
  const limit = ref<number>(params?.limit ?? 50);
  const loading = ref<boolean>(true);
  const fetching = ref<boolean>(false);
  const isLoadAll = ref<boolean>(false);
  const list = ref<T[]>([]);

  const load = async (fn: (limit: number, offset: number) => Promise<{ list: T[]; total?: number } | T[]>) => {
    if (isLoadAll.value) return { list: [], total: 0 };
    const offset = (page.value - 1) * limit.value;

    if (page.value <= 1) loading.value = true;
    else fetching.value = true;

    try {
      const result = await fn(limit.value, offset);
      const resList = 'list' in result ? result.list : result;
      if (page.value === 1) list.value = resList;
      else {
        for (const item of resList) {
          list.value.push(item as any);
        }
      }

      if (list.value.length < limit.value) {
        isLoadAll.value = true;
      } else {
        page.value += 1;
      }

      return { list: resList };
    } finally {
      loading.value = false;
      fetching.value = false;
    }
  };
  const reset = () => {
    page.value = 1;
    isLoadAll.value = false;
  };
  const loadAll = async (fn: (limit: number, offset: number) => Promise<{ list: T[]; total: number }>) => {
    reset();
    while (!isLoadAll.value) {
      await load(fn);
    }
    return list.value;
  };

  return {
    searchQuery,
    page,
    limit,
    loading,
    fetching,
    isLoadAll,
    list,
    load,
    reset,
    loadAll,
  };
};
