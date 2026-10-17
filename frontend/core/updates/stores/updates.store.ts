import { defineStore } from 'pinia';

interface IUpdatesStore {
  revision: number;
  connected: boolean;
}

export const useUpdatesStore = defineStore('updates', {
  state: (): IUpdatesStore => ({
    revision: 0,
    connected: false,
  }),
  getters: {
    service: () => useServices().UPDATES.UpdatesService,
  },
  actions: {
    notify() {
      this.revision += 1;
    },
  },
});
