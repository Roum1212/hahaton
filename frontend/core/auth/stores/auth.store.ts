import { defineStore } from 'pinia';
import { accessToken } from '~/core/app/rpc/token';
import { UserRole } from '~/core/users/services/users/users.types';
import type { IAuthUser } from '../services/auth/auth.types';

interface IAuthStore {
  user: IAuthUser | null;
  isLoggedIn: boolean;
}

export const useAuthStore = defineStore('auth', {
  state: (): IAuthStore => ({
    user: null,
    isLoggedIn: false,
  }),
  getters: {
    service: () => useServices().AUTH.AuthService,
    hasToken: (): boolean => !!accessToken.get(),
    isAdmin: (state): boolean => state.user?.role === UserRole.ADMIN,
    roles: (state): string[] => (state.user?.role === UserRole.ADMIN ? ['admin', 'user'] : state.user ? ['user'] : []),
  },
  actions: {
    async loginWithFlida(flidaAccessToken: string) {
      accessToken.clear();
      accessToken.set(await this.service.loginWithFlida(flidaAccessToken));
      await this.checkAuth();
    },

    async checkAuth() {
      if (!accessToken.get()) {
        this.resetState();
        return;
      }
      try {
        this.setUser(await this.service.me());
        this.isLoggedIn = true;
      } catch (e) {
        accessToken.clear();
        this.resetState();
        throw e;
      }
    },

    async logout() {
      try {
        await this.service.logout();
      } finally {
        accessToken.clear();
        this.resetState();
        this.resetUserScopedStores();
      }
    },

    setUser(user: IAuthUser) {
      this.user = { ...user };
    },

    resetUserScopedStores() {
      useUsersStore().$reset();
    },

    resetState() {
      this.user = null;
      this.isLoggedIn = false;
    },
  },
});
