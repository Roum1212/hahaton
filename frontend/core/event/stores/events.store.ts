import { defineStore } from 'pinia';
import type { ISearchArea } from '~/core/geo/services/geo/geo.types';
import {
  EventListScope,
  EventStatus,
  type IEvent,
  type IEventDto,
  type IEventFilter,
  type IEventResult,
} from '../services/event/event.types';

interface IStore {
  pagination: IPaginationParams;
  filter: IEventFilter | null;
  total: number;
  list: IEvent[];
  byId: Record<string, IEvent>;
  upcoming: IEvent[];
}

export const useEventsStore = defineStore('events', {
  state: (): IStore => ({
    pagination: { ...DEFAULT_PAGINATION },
    filter: null,
    total: 0,
    list: [],
    byId: {},
    upcoming: [],
  }),
  getters: {
    service: () => useServices().EVENT.EventService,
    getById: (state) => (id: string) => state.byId[id],
  },
  actions: {
    async find(page: number = 1, filter?: IEventFilter) {
      if (page <= 1) {
        this.resetPagination();
        if (filter) this.filter = filter;
      }
      const current = this.filter;
      if (!current) return;
      const { list, total, loadAll } = await executeWithPagination(this.pagination, page, (params) => this.service.find(current, params));
      this.pagination.loadAll = loadAll;
      this.total = total ?? 0;
      this.remember(list);
      if (page <= 1) this.list = list;
      else this.list.push(...list);
    },

    async findUpcoming(searchArea: ISearchArea, limit = 3) {
      const { list } = await this.service.find(
        { scope: EventListScope.CATALOG, searchArea, status: EventStatus.PUBLISHED, startsFrom: new Date() },
        { limit, offset: 0 },
      );
      this.remember(list);
      this.upcoming = list;
      return list;
    },

    async fetchOne(id: string, force = false) {
      if (!force && this.byId[id]) return this.byId[id];
      const event = await this.service.get(id);
      this.replace(event);
      return event;
    },

    async create(dto: IEventDto) {
      const event = await this.service.create(dto);
      this.replace(event);
      this.list = [event, ...this.list];
      this.total++;
      return event;
    },

    async update(id: string, dto: Partial<IEventDto>) {
      const event = await this.service.update(id, dto);
      this.replace(event);
      return event;
    },

    async cancel(id: string) {
      const event = await this.service.cancel(id);
      this.replace(event);
      return event;
    },

    async complete(id: string, results: IEventResult[], requestId: string) {
      const event = await this.service.complete(id, results, requestId);
      this.replace(event);
      return event;
    },

    remember(list: IEvent[]) {
      list.forEach((item) => (this.byId[item.id] = item));
    },

    replace(event: IEvent) {
      this.byId[event.id] = event;
      this.list = this.list.map((item) => (item.id === event.id ? event : item));
      this.upcoming = this.upcoming.map((item) => (item.id === event.id ? event : item));
    },

    resetPagination() {
      this.pagination = { ...DEFAULT_PAGINATION };
    },

    $reset() {
      this.resetPagination();
      this.filter = null;
      this.total = 0;
      this.list = [];
      this.byId = {};
      this.upcoming = [];
    },
  },
});
