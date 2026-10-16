import { defineStore } from 'pinia';
import type { ICoordinator, ICoordinatorSettings, IUpdateCoordinatorDto } from '../services/coordinator/coordinator.types';

interface IStore {
  pagination: IPaginationParams;
  total: number;
  list: ICoordinator[];
  byId: Record<string, ICoordinator>;
  me: ICoordinator | null;
  settings: ICoordinatorSettings | null;
}

export const useCoordinatorsStore = defineStore('coordinators', {
  state: (): IStore => ({
    pagination: { ...DEFAULT_PAGINATION },
    total: 0,
    list: [],
    byId: {},
    me: null,
    settings: null,
  }),
  getters: {
    service: () => useServices().COORDINATOR.CoordinatorService,
    getById: (state) => (id: string) => state.byId[id],
  },
  actions: {
    async find(page: number = 1) {
      if (page <= 1) this.resetPagination();
      const { list, total, loadAll } = await executeWithPagination(this.pagination, page, (params) => this.service.find(params));
      this.pagination.loadAll = loadAll;
      this.total = total ?? 0;
      list.forEach((item) => (this.byId[item.id] = item));
      if (page <= 1) this.list = list;
      else this.list.push(...list);
    },

    async fetchOne(id: string) {
      const coordinator = this.byId[id] ?? (await this.service.get(id));
      this.byId[id] = coordinator;
      return coordinator;
    },

    async fetchMe() {
      this.me = await this.service.get();
      this.byId[this.me.id] = this.me;
      return this.me;
    },

    async updateMe(dto: IUpdateCoordinatorDto) {
      this.me = await this.service.update(dto);
      this.byId[this.me.id] = this.me;
      return this.me;
    },

    async fetchSettings() {
      this.settings = await this.service.getSettings();
      return this.settings;
    },

    async updateSettings(settings: ICoordinatorSettings) {
      this.settings = await this.service.updateSettings(settings);
      return this.settings;
    },

    resetPagination() {
      this.pagination = { ...DEFAULT_PAGINATION };
    },

    $reset() {
      this.resetPagination();
      this.total = 0;
      this.list = [];
      this.byId = {};
      this.me = null;
      this.settings = null;
    },
  },
});
