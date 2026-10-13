/* eslint-disable @typescript-eslint/no-explicit-any */
export interface ListState<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  loading: boolean;
  loadMore: boolean;
  error: string | null;
  filters: Record<string, any>;
}

export const useListStore = <T>(storeKey: string) =>
  defineStore(storeKey, {
    state: (): ListState<T> => ({
      items: [],
      total: 0,
      page: 1,
      pageSize: 20,
      loading: false,
      loadMore: true,
      error: null,
      filters: {},
    }),
    actions: {
      async fetch(
        fetchFn: (params: { page: number; pageSize: number; filters: Record<string, any> }) => Promise<{ items: T[]; total: number }>,
      ) {
        if (!this.loadMore) return;
        this.loading = true;
        this.error = null;
        try {
          const { items, total } = await fetchFn({ page: this.page, pageSize: this.pageSize, filters: this.filters });
          if (this.page === 0) {
            this.items = items as typeof this.items;
          } else {
            this.items.push(...(items as typeof this.items));
          }
          this.total = total;
        } catch (err: any) {
          this.error = err?.message || 'Unknown error';
        } finally {
          this.loading = false;
        }
      },
      setPage(page: number) {
        this.page = page;
      },
      setPageSize(size: number) {
        this.pageSize = size;
      },
      setFilters(filters: Record<string, any>) {
        this.filters = filters;
      },
      reset() {
        this.items = [];
        this.total = 0;
        this.page = 1;
        this.filters = {};
        this.error = null;
      },
    },
  });
