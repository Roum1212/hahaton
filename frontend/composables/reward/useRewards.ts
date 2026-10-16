import type { IReward, IRewardDto } from '~/core/reward/services/reward/reward.types';
import { RewardListScope } from '~/core/reward/services/reward/reward.types';
import { RewardOrderListScope, RewardOrderStatus } from '~/core/reward-order/services/reward-order/reward-order.types';

export const useRewardCatalog = () => {
  const store = useRewardsStore();
  const { ensure: ensureCoordinators } = useCoordinatorNames();
  const loading = ref(true);

  const load = async (coordinatorId?: string) => {
    loading.value = true;
    await handleAsyncOperation(
      async () => {
        await store.find(1, { scope: RewardListScope.CATALOG, coordinatorId });
        while (!store.pagination.loadAll) await store.find(store.pagination.page);
        await ensureCoordinators(store.list.map((item) => item.coordinatorId));
      },
      () => {
        loading.value = false;
      },
    );
  };

  return { store, loading, load };
};

export const useRewardPurchase = () => {
  const ordersStore = useRewardOrdersStore();
  const rewardsStore = useRewardsStore();
  const volunteerStore = useVolunteerStore();
  const { toast, confirm } = useApp();

  const purchase = (reward: IReward) => {
    const requestId = createRequestId();
    confirm.accept(`Обменять ${formatCoins(reward.priceCoins)} на «${reward.title}»?`, {
      header: 'Обмен коинов',
      acceptLabel: 'Обменять',
      accept: async () => {
        await ordersStore.purchase(reward.id, reward.priceCoins, requestId);
        await Promise.all([volunteerStore.fetchProfile(), rewardsStore.refresh(reward.id)]);
        toast.success('Запрос на поощрение оформлен');
      },
    });
  };

  return { purchase };
};

export const useMyPurchases = () => {
  const store = useRewardOrdersStore();
  const volunteerStore = useVolunteerStore();
  const { ensure: ensureCoordinators } = useCoordinatorNames();
  const { toast, confirm } = useApp();
  const loading = ref(true);
  const fetching = ref(false);

  const load = async (page = 1) => {
    if (page === 1) loading.value = true;
    else fetching.value = true;
    await handleAsyncOperation(
      async () => {
        await store.find(page, { scope: RewardOrderListScope.MY_PURCHASES });
        await ensureCoordinators(store.list.map((item) => item.coordinatorId));
      },
      () => {
        loading.value = false;
        fetching.value = false;
      },
    );
  };

  const cancel = (id: string) =>
    confirm.delete('Отменить заказ? Коины вернутся на баланс.', {
      header: 'Отмена заказа',
      acceptLabel: 'Отменить заказ',
      rejectLabel: 'Не отменять',
      accept: async () => {
        await store.cancel(id);
        await volunteerStore.fetchProfile();
        toast.success('Заказ отменён, коины возвращены');
      },
    });

  const canCancel = (status: RewardOrderStatus) => status === RewardOrderStatus.PENDING;

  return { store, loading, fetching, load, cancel, canCancel };
};

export const useOwnedRewards = () => {
  const store = useRewardsStore();
  const { toast, confirm, loader } = useApp();
  const loading = ref(true);

  const load = async () => {
    loading.value = true;
    await handleAsyncOperation(
      async () => {
        await store.find(1, { scope: RewardListScope.OWNED });
        while (!store.pagination.loadAll) await store.find(store.pagination.page);
      },
      () => {
        loading.value = false;
      },
    );
  };

  const save = async (id: string | null, dto: IRewardDto) => {
    loader.show();
    const result = await handleAsyncOperation(
      () => (id ? store.update(id, dto) : store.create(dto)),
      () => loader.hide(),
    );
    if (result) toast.success(id ? 'Поощрение обновлено' : 'Поощрение создано');
    return !!result;
  };

  const remove = (reward: IReward) =>
    confirm.delete(`Удалить «${reward.title}»? Невыданные заказы будут отменены с возвратом коинов.`, {
      header: 'Удаление поощрения',
      acceptLabel: 'Удалить',
      rejectLabel: 'Отмена',
      accept: async () => {
        await store.delete(reward.id);
        toast.success('Поощрение удалено');
      },
    });

  return { store, loading, load, save, remove };
};

export const useIncomingOrders = () => {
  const store = useRewardOrdersStore();
  const { ensure: ensureVolunteers } = useVolunteerNames();
  const { toast, confirm, loader } = useApp();
  const loading = ref(true);
  const fetching = ref(false);

  const load = async (status?: RewardOrderStatus, page = 1) => {
    if (page === 1) loading.value = true;
    else fetching.value = true;
    await handleAsyncOperation(
      async () => {
        await store.find(page, { scope: RewardOrderListScope.INCOMING, status });
        await ensureVolunteers(store.list.map((item) => item.volunteerId));
      },
      () => {
        loading.value = false;
        fetching.value = false;
      },
    );
  };

  const issue = async (id: string) => {
    loader.show();
    const result = await handleAsyncOperation(
      () => store.issue(id),
      () => loader.hide(),
    );
    if (result) toast.success('Поощрение выдано');
  };

  const cancel = (id: string) =>
    confirm.delete('Отменить заказ? Коины вернутся волонтёру.', {
      header: 'Отмена заказа',
      acceptLabel: 'Отменить заказ',
      rejectLabel: 'Не отменять',
      accept: async () => {
        await store.cancel(id);
        toast.success('Заказ отменён');
      },
    });

  return { store, loading, fetching, load, issue, cancel };
};
