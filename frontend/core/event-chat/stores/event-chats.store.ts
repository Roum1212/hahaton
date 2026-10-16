import { defineStore } from 'pinia';
import type { IEventChat } from '../services/event-chat/event-chat.types';

interface IStore {
  pagination: IPaginationParams;
  total: number;
  list: IEventChat[];
}

export const useEventChatsStore = defineStore('eventChats', {
  state: (): IStore => ({
    pagination: { ...DEFAULT_PAGINATION },
    total: 0,
    list: [],
  }),
  getters: {
    service: () => useServices().EVENT_CHAT.EventChatService,
    getByEvent: (state) => (eventId: string) => state.list.find((item) => item.eventId === eventId),
  },
  actions: {
    async find(page: number = 1) {
      if (page <= 1) this.resetPagination();
      const { list, total, loadAll } = await executeWithPagination(this.pagination, page, (params) => this.service.findMine(params));
      this.pagination.loadAll = loadAll;
      this.total = total ?? 0;
      if (page <= 1) this.list = list;
      else this.list.push(...list);
    },

    async refresh(eventId: string) {
      const chat = await this.service.get(eventId);
      const exists = this.list.some((item) => item.eventId === eventId);
      this.list = exists ? this.list.map((item) => (item.eventId === eventId ? chat : item)) : [chat, ...this.list];
      return chat;
    },

    resetPagination() {
      this.pagination = { ...DEFAULT_PAGINATION };
    },

    $reset() {
      this.resetPagination();
      this.total = 0;
      this.list = [];
    },
  },
});
