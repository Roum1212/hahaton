/* eslint-disable @typescript-eslint/no-explicit-any */
export const DEFAULT_PAGINATION = {
  page: 1,
  limit: 50,
  loadAll: false,
};
export interface IPaginationParams {
  page: number;
  limit: number;
  loadAll: boolean;
}
export async function executeWithPagination<RESPONSE>(
  pagination: IPaginationParams,
  _: any,
  asyncFn: (params: { limit: number; offset: number }) => Promise<{ list: RESPONSE[]; total?: number }>,
): Promise<{ list: RESPONSE[]; total?: number; loadAll: boolean }> {
  const { limit, loadAll } = pagination;
  const defaultResponse = {
    list: [],
    total: undefined,
    loadAll: true,
  };
  if (loadAll) return defaultResponse;
  if (pagination.page === 0) pagination.page = 1;
  const offset = limit * (pagination.page - 1);
  const { list, total }: { list: RESPONSE[]; total?: number } = await asyncFn({ limit, offset });
  pagination.page++;
  return {
    list,
    total,
    loadAll: list.length < limit,
  };
}
