import { defineStore } from 'pinia';
import type { IReward, IRewardDto, IRewardFilter } from '../services/reward/reward.types';

interface IStore {
  pagination: IPaginationParams;
  filter: IRewardFilter | null;
  total: number;
  list: IReward[];
}

export const useRewardsStore = defineStore('rewards', {
  state: (): IStore => ({
    pagination: { ...DEFAULT_PAGINATION },
    filter: null,
    total: 0,
    list: [],
  }),
  getters: {
    service: () => useServices().REWARD.RewardService,
    getById: (state) => (id: string) => state.list.find((item) => item.id === id),
  },
  actions: {
    async find(page: number = 1, filter?: IRewardFilter) {
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

    async refresh(id: string) {
      return this.replace(await this.service.get(id));
    },

    async create(dto: IRewardDto) {
      const reward = await this.service.create(dto);
      this.list = [reward, ...this.list];
      this.total++;
      return reward;
    },

    async update(id: string, dto: Partial<IRewardDto>) {
      return this.replace(await this.service.update(id, dto));
    },

    async delete(id: string) {
      await this.service.delete(id);
      this.list = this.list.filter((item) => item.id !== id);
      this.total = Math.max(0, this.total - 1);
    },

    replace(reward: IReward) {
      this.list = this.list.map((item) => (item.id === reward.id ? reward : item));
      return reward;
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
