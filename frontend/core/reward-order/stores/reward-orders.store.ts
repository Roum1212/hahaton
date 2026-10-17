import { defineStore } from 'pinia';
import type { IRewardOrder, IRewardOrderFilter } from '../services/reward-order/reward-order.types';

interface IStore {
  pagination: IPaginationParams;
  filter: IRewardOrderFilter | null;
  total: number;
  list: IRewardOrder[];
}

export const useRewardOrdersStore = defineStore('rewardOrders', {
  state: (): IStore => ({
    pagination: { ...DEFAULT_PAGINATION },
    filter: null,
    total: 0,
    list: [],
  }),
  getters: {
    service: () => useServices().REWARD_ORDER.RewardOrderService,
  },
  actions: {
    async find(page: number = 1, filter?: IRewardOrderFilter) {
      if (page <= 1) {
        this.resetPagination();
        if (filter) this.filter = filter;
      }
      const current = this.filter;
      if (!current) return;
      const { list, total, loadAll } = await executeWithPagination(this.pagination, page, (params) => this.service.find(current, params));
      this.pagination.loadAll = loadAll;
      this.total = total ?? 0;
      if (page <= 1) this.list = list;
      else this.list.push(...list);
    },

    async purchase(rewardId: string, expectedPriceCoins: number, requestId: string) {
      const order = await this.service.create(rewardId, expectedPriceCoins, requestId);
      this.list = [order, ...this.list.filter((item) => item.id !== order.id)];
      return order;
    },

    async issue(id: string) {
      return this.replace(await this.service.issue(id));
    },

    async cancel(id: string) {
      return this.replace(await this.service.cancel(id));
    },

    replace(order: IRewardOrder) {
      const status = this.filter?.status;
      if (status !== undefined && order.status !== status) {
        const before = this.list.length;
        this.list = this.list.filter((item) => item.id !== order.id);
        this.total = Math.max(0, this.total - (before - this.list.length));
        return order;
      }
      this.list = this.list.map((item) => (item.id === order.id ? order : item));
      return order;
    },

    resetPagination() {
      this.pagination = { ...DEFAULT_PAGINATION };
    },

    $reset() {
      this.resetPagination();
      this.filter = null;
      this.total = 0;
      this.list = [];
    },
  },
});
