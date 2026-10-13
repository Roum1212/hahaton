import { defineStore } from 'pinia';
import type { ICreateUserDto, IUser } from '../services/users/users.types';

interface IStore {
  pagination: IPaginationParams;
  search: string | undefined;
  total: number;
  list: IUser[];
  item: IUser | undefined;
}

export const useUsersStore = defineStore('users', {
  state: (): IStore => ({
    pagination: { ...DEFAULT_PAGINATION },
    search: undefined,
    total: 0,
    list: [],
    item: undefined,
  }),
  getters: {
    service: () => useServices().USERS.UsersService,
    getUserById: (state) => (id: string) => state.list.find((item) => item.id === id),
  },
  actions: {
    async find(page: number = 1, search?: string) {
      if (page <= 1) {
        this.resetPagination();
        this.search = search;
      }
      const { list, total, loadAll } = await executeWithPagination(this.pagination, page, (params) =>
        this.service.find(this.search, params),
      );
      this.pagination.loadAll = loadAll;
      this.total = total ?? 0;
      if (page <= 1) {
        this.list = list;
      } else {
        this.list.push(...list);
      }
    },

    async fetchOne(id: string) {
      this.item = await this.service.get(id);
      return this.item;
    },

    async create(dto: ICreateUserDto) {
      const model = await this.service.create(dto);
      this.list = [model, ...this.list];
      this.total++;
      return model;
    },

    async update(id: string, dto: ICreateUserDto) {
      const model = await this.service.update(id, dto);
      this.list = this.list.map((item) => (item.id === id ? model : item));
      const authStore = useAuthStore();
      if (authStore.user?.id === id) authStore.setUser(model);
      return model;
    },

    async delete(id: string) {
      await this.service.delete(id);
      this.list = this.list.filter((item) => item.id !== id);
      this.total = Math.max(0, this.total - 1);
    },

    resetPagination() {
      this.pagination = { ...DEFAULT_PAGINATION };
    },

    $reset() {
      this.resetPagination();
      this.search = undefined;
      this.total = 0;
      this.list = [];
      this.item = undefined;
    },
  },
});
