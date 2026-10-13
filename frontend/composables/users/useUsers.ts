import type { UiTableColumn } from '~/components/ui/Table/index.vue';
import { UserRole, USER_ROLE_LABELS, type ICreateUserDto, type IUser } from '~/core/users/services/users/users.types';

interface IState {
  loading: boolean;
  fetching: boolean;
  page: number;
  searchQuery: string | undefined;
}

const manager = useDialogManager();

export const useUsersForm = () => {
  const { formRef, submitting, submitForm } = useEntityForm();

  const generateInitialData = (user?: IUser): ICreateUserDto => ({
    email: user?.email ?? '',
    firstName: user?.firstName ?? '',
    lastName: user?.lastName ?? '',
    role: user?.role ?? UserRole.USER,
    active: user?.active ?? true,
  });

  const roleOptions = [UserRole.USER, UserRole.ADMIN].map((value) => ({ value, label: USER_ROLE_LABELS[value] }));

  return { formRef, submitting, submitForm, generateInitialData, roleOptions };
};

export const useUsersActions = () => {
  const state = reactive<IState>({
    loading: true,
    fetching: false,
    page: 1,
    searchQuery: undefined,
  });
  const { toast, confirm, loader, debounce } = useApp();
  const store = useUsersStore();

  const resetList = () => {
    store.resetPagination();
    state.page = 1;
  };

  const load = async () => {
    if (state.fetching || (state.page > 1 && store.pagination.loadAll)) return;
    await handleAsyncOperation(
      async () => {
        state.fetching = true;
        await store.find(state.page, state.searchQuery);
        state.page++;
      },
      () => {
        state.fetching = false;
      },
    );
  };

  const reload = async () => {
    state.loading = true;
    resetList();
    await load();
    state.loading = false;
  };

  const search = debounce(reload, 500);

  const save = async (userId: string | undefined, dto: ICreateUserDto): Promise<IUser | false> => {
    loader.show();
    const result = await handleAsyncOperation(
      async () => {
        if (userId) {
          const model = await store.update(userId, dto);
          toast.success('Информация успешно обновлена');
          return model;
        }
        const model = await store.create(dto);
        toast.success('Пользователь успешно создан');
        return model;
      },
      () => loader.hide(),
    );
    return result ?? false;
  };

  const deleteMany = (ids: string[], fn?: () => void) => {
    confirm.delete('Вы действительно хотите удалить пользователей?', {
      acceptLabel: 'Удалить',
      accept: async () => {
        await Promise.all(ids.map((id) => store.delete(id)));
        toast.success('Пользователи успешно удалены');
        if (fn) fn();
      },
    });
  };

  return { state, load, reload, search, resetList, save, deleteMany };
};

export const useUsersTable = () => {
  const columns: UiTableColumn[] = [
    { key: 'name', header: 'Пользователь' },
    { key: 'role', header: 'Роль', width: '170px' },
    { key: 'status', header: 'Статус', width: '150px', hideOnMobile: true },
    { key: 'createdAt', header: 'Создан', width: '170px', hideOnMobile: true },
    { key: 'actions', header: '', width: '110px', align: 'right' },
  ];
  return { columns };
};

export const useUsersDialogs = () => {
  const { loader } = useApp();
  const store = useUsersStore();

  const openUserFormDialog = async (userId?: string) => {
    const { default: UsersForm } = await import('~/components/users/Form.vue');

    return manager.openDialog<{ userId?: string }, { user?: IUser }>(UsersForm, {
      id: `user-form-${userId ?? 'new'}`,
      options: {
        type: 'sidebar',
        header: 'Добавить пользователя',
        size: 'medium',
      },
      props: { userId },
      onOpen: async ({ userId }, options) => {
        if (!userId) return { user: undefined };
        loader.show();
        const user = await handleAsyncOperation(
          async () => store.getUserById(userId) ?? (await store.fetchOne(userId)),
          () => loader.hide(),
        );
        if (user) options.header = user.fullName;
        return { user };
      },
      footerContent: (modal) => manager.defaultFooter(modal),
    });
  };

  return { openUserFormDialog };
};
