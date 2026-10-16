import { defineStore } from 'pinia';
import type { ISession } from '../services/auth/auth.types';

interface IAuthStore {
  session: ISession | null;
  initialized: boolean;
}

export const useAuthStore = defineStore('auth', {
  state: (): IAuthStore => ({
    session: null,
    initialized: false,
  }),
  getters: {
    service: () => useServices().AUTH.AuthService,
    isLoggedIn: (state): boolean => !!state.session?.userId,
    hasCoordinatorProfile: (state): boolean => !!state.session?.hasCoordinatorProfile,
    roles: (state): string[] => {
      if (!state.session?.userId) return [];
      return state.session.hasCoordinatorProfile ? [ROLES.volunteer, ROLES.coordinator] : [ROLES.volunteer];
    },
  },
  actions: {
    async init() {
      if (this.initialized) return;
      try {
        this.session = await this.service.open();
        if (this.isLoggedIn) await useVolunteerStore().fetchProfile();
      } catch {
        this.resetState();
      } finally {
        this.initialized = true;
      }
    },

    async loginWithFlida(idToken: string) {
      await this.service.open();
      this.session = await this.service.authenticate(idToken);
      await useVolunteerStore().fetchProfile();
    },

    async refreshSession() {
      this.session = await this.service.get();
    },

    async logout() {
      try {
        await this.service.close();
      } finally {
        this.resetState();
        this.resetUserScopedStores();
      }
    },

    expire() {
      if (!this.isLoggedIn) return;
      this.resetState();
      this.resetUserScopedStores();
      const route = useRouter().currentRoute.value;
      if (GUEST_PAGES.includes(route.path.toLowerCase())) return;
      useCookie('redirectAfterLogin').value = route.fullPath;
      navigateTo(ROUTES.login);
    },

    resetUserScopedStores() {
      useVolunteerStore().$reset();
      useCoordinatorsStore().$reset();
      useCooperationsStore().$reset();
      useEventsStore().$reset();
      useParticipationsStore().$reset();
      useRewardsStore().$reset();
      useRewardOrdersStore().$reset();
      useEventChatsStore().$reset();
    },

    resetState() {
      this.session = null;
    },
  },
});
