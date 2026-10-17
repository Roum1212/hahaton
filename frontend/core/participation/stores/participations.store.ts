import { defineStore } from 'pinia';
import type { IParticipation, IParticipationFilter } from '../services/participation/participation.types';

interface IStore {
  pagination: IPaginationParams;
  filter: IParticipationFilter | null;
  total: number;
  list: IParticipation[];
}

export const useParticipationsStore = defineStore('participations', {
  state: (): IStore => ({
    pagination: { ...DEFAULT_PAGINATION },
    filter: null,
    total: 0,
    list: [],
  }),
  getters: {
    service: () => useServices().PARTICIPATION.ParticipationService,
    getById: (state) => (id: string) => state.list.find((item) => item.id === id),
    getByEvent: (state) => (eventId: string) => state.list.find((item) => item.eventId === eventId),
  },
  actions: {
    async find(page: number = 1, filter?: IParticipationFilter) {
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

    async request(eventId: string) {
      const participation = await this.service.request(eventId);
      this.replace(participation);
      return participation;
    },

    async restore(id: string) {
      return this.replace(await this.service.restore(id));
    },

    async approve(id: string) {
      return this.replace(await this.service.approve(id));
    },

    async reject(id: string, reason?: string) {
      return this.replace(await this.service.reject(id, reason));
    },

    async cancel(id: string) {
      return this.replace(await this.service.cancel(id));
    },

    replace(participation: IParticipation) {
      const status = this.filter?.status;
      const exists = this.list.some((item) => item.id === participation.id);
      if (status !== undefined && participation.status !== status) {
        this.list = this.list.filter((item) => item.id !== participation.id);
        if (exists) this.total = Math.max(0, this.total - 1);
        return participation;
      }
      this.list = exists ? this.list.map((item) => (item.id === participation.id ? participation : item)) : [...this.list, participation];
      return participation;
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
