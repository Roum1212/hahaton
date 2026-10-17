import type { IListResult } from '~/core/app/types';
import { CooperationListScope } from '~/core/cooperation/services/cooperation/cooperation.types';
import {
  ParticipationListScope,
  ParticipationStatus,
  type IParticipation,
} from '~/core/participation/services/participation/participation.types';
import { RewardOrderListScope, type IRewardOrder } from '~/core/reward-order/services/reward-order/reward-order.types';

export const useVolunteerNames = () => {
  const store = useVolunteerStore();

  const ensure = async (ids: string[]) => {
    const missing = [...new Set(ids)].filter((id) => id && !store.getCoordinatorView(id));
    await Promise.all(missing.map((id) => handleAsyncOperation(() => store.fetchCoordinatorView(id))));
  };

  const name = (id: string) => store.getCoordinatorView(id)?.fullName || 'Волонтёр';
  const view = (id: string) => store.getCoordinatorView(id);
  const avatar = (id: string) => store.getCoordinatorView(id)?.avatarUrl ?? null;

  return { ensure, name, view, avatar };
};

const PROFILE_PAGE_LIMIT = 100;

const fetchEvery = async <T>(fetchPage: (limit: number, offset: number) => Promise<IListResult<T>>): Promise<T[]> => {
  const items: T[] = [];
  for (let offset = 0; ; offset += PROFILE_PAGE_LIMIT) {
    const { list, total } = await fetchPage(PROFILE_PAGE_LIMIT, offset);
    items.push(...list);
    if (!list.length || items.length >= (total ?? 0)) return items;
  }
};

export const useCoordinatorVolunteerProfile = (volunteerId: Ref<string>) => {
  const volunteerStore = useVolunteerStore();
  const participationsStore = useParticipationsStore();
  const ordersStore = useRewardOrdersStore();
  const cooperationsStore = useCooperationsStore();
  const { ensure: ensureEvents } = useEventLookup();

  const loading = ref(true);
  const participations = ref<IParticipation[]>([]);
  const orders = ref<IRewardOrder[]>([]);

  const volunteer = computed(() => volunteerStore.getCoordinatorView(volunteerId.value) ?? null);
  const cooperation = computed(() => cooperationsStore.withVolunteer(volunteerId.value));
  const rejections = computed(() => participations.value.filter((item) => item.status === ParticipationStatus.REJECTED));
  const history = computed(() => participations.value.filter((item) => item.status !== ParticipationStatus.REJECTED));
  const confirmed = computed(() => participations.value.filter((item) => item.status === ParticipationStatus.CONFIRMED));

  const stats = computed(() => ({
    events: confirmed.value.length,
    minutes: confirmed.value.reduce((sum, item) => sum + (item.awardedMinutes ?? 0), 0),
    coins: confirmed.value.reduce((sum, item) => sum + (item.awardedCoins ?? 0), 0),
    rejections: rejections.value.length,
  }));

  const load = async (silent = false) => {
    if (!silent) loading.value = true;
    await handleAsyncOperation(
      async () => {
        const id = volunteerId.value;
        const [, , participationList, orderList] = await Promise.all([
          volunteerStore.fetchCoordinatorView(id, true),
          cooperationsStore.fetchAll(CooperationListScope.AS_COORDINATOR),
          fetchEvery((limit, offset) =>
            participationsStore.service.find({ scope: ParticipationListScope.MY_EVENTS, volunteerId: id }, { limit, offset }),
          ),
          fetchEvery((limit, offset) =>
            ordersStore.service.find({ scope: RewardOrderListScope.INCOMING, volunteerId: id }, { limit, offset }),
          ),
        ]);
        participations.value = [...participationList].sort((a, b) => (b.updatedAt?.getTime() ?? 0) - (a.updatedAt?.getTime() ?? 0));
        orders.value = orderList;
        await ensureEvents(participationList.map((item) => item.eventId));
      },
      () => {
        loading.value = false;
      },
    );
  };

  return { loading, volunteer, cooperation, participations, history, rejections, orders, stats, load };
};
