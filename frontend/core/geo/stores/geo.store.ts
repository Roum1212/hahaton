import { defineStore } from 'pinia';
import type { ICity } from '../services/geo/geo.types';

interface IStore {
  cities: Record<string, ICity>;
}

export const useGeoStore = defineStore('geo', {
  state: (): IStore => ({
    cities: {},
  }),
  getters: {
    service: () => useServices().GEO.GeoService,
    getCityById: (state) => (id: string) => state.cities[id],
  },
  actions: {
    async searchCities(query: string, limit = 20) {
      const { list } = await this.service.listCities(query, { limit, offset: 0 });
      list.forEach((city) => (this.cities[city.id] = city));
      return list;
    },

    async fetchCity(id: string) {
      const city = this.cities[id] ?? (await this.service.getCity(id));
      this.cities[id] = city;
      return city;
    },
  },
});
