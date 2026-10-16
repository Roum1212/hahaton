import { defineStore } from 'pinia';
import type { ICoordinatorVolunteer, IUpdateVolunteerDto, IVolunteer } from '../services/volunteer/volunteer.types';

interface IStore {
  profile: IVolunteer | null;
  coordinatorViews: Record<string, ICoordinatorVolunteer>;
}

export const useVolunteerStore = defineStore('volunteer', {
  state: (): IStore => ({
    profile: null,
    coordinatorViews: {},
  }),
  getters: {
    service: () => useServices().VOLUNTEER.VolunteerService,
    getCoordinatorView: (state) => (id: string) => state.coordinatorViews[id],
  },
  actions: {
    async fetchProfile() {
      this.profile = await this.service.getMe();
      return this.profile;
    },

    async update(dto: IUpdateVolunteerDto) {
      this.profile = await this.service.update(dto);
      return this.profile;
    },

    async fetchCoordinatorView(volunteerId: string, force = false) {
      if (!force && this.coordinatorViews[volunteerId]) return this.coordinatorViews[volunteerId];
      const view = await this.service.getForCoordinator(volunteerId);
      this.coordinatorViews[volunteerId] = view;
      return view;
    },

    $reset() {
      this.profile = null;
      this.coordinatorViews = {};
    },
  },
});
