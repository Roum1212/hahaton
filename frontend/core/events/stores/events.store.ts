import { defineStore } from 'pinia';
import type { IEvent } from '../services/events/events.types';

interface IStore {
  city: string | undefined;
  upcoming: IEvent[];
}

export const useEventsStore = defineStore('events', {
  state: (): IStore => ({
    city: undefined,
    upcoming: [],
  }),
  getters: {
    service: () => useServices().EVENTS.EventsService,
  },
  actions: {
    async findUpcoming(city: string) {
      this.city = city;
      this.upcoming = await this.service.findUpcoming(city);
      return this.upcoming;
    },

    $reset() {
      this.city = undefined;
      this.upcoming = [];
    },
  },
});
