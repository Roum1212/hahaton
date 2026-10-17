import type { IParticipation, IParticipationChange } from '~/core/participation/services/participation/participation.types';
import { ParticipationListScope, ParticipationStatus } from '~/core/participation/services/participation/participation.types';

const ACTIVE_STATUSES = [ParticipationStatus.PENDING, ParticipationStatus.APPROVED];

export const useParticipationActions = () => {
  const store = useParticipationsStore();
  const { toast, confirm, loader } = useApp();

  const run = async <T>(fn: () => Promise<T>, message: string) => {
    loader.show();
    const result = await handleAsyncOperation(fn, () => loader.hide());
    if (result) toast.success(message);
    return result;
  };

  const request = (eventId: string) => run(() => store.request(eventId), 'Заявка отправлена координатору');

  const restore = (id: string) => run(() => store.restore(id), 'Заявка снова на рассмотрении');

  const cancel = (id: string, onDone?: (participation: IParticipation) => void) =>
    confirm.delete('Отменить заявку на мероприятие?', {
      header: 'Отмена заявки',
      acceptLabel: 'Отменить заявку',
      rejectLabel: 'Не отменять',
      accept: async () => {
        const participation = await store.cancel(id);
        toast.success('Заявка отменена');
        onDone?.(participation);
      },
    });

  return { request, restore, cancel };
};

export const useMyParticipations = () => {
  const store = useParticipationsStore();
  const { ensure: ensureEvents } = useEventLookup();
  const loading = ref(true);
  const fetching = ref(false);

  const load = async (status?: ParticipationStatus, page = 1, silent = false) => {
    if (page > 1) fetching.value = true;
    else if (!silent) loading.value = true;
    await handleAsyncOperation(
      async () => {
        await store.find(page, { scope: ParticipationListScope.MY_APPLICATIONS, status });
        await ensureEvents(store.list.map((item) => item.eventId));
      },
      () => {
        loading.value = false;
        fetching.value = false;
      },
    );
  };

  const findForEvent = async (eventId: string) => {
    const { list } = await store.service.find({ scope: ParticipationListScope.MY_APPLICATIONS, eventId }, { limit: 1, offset: 0 });
    return list[0] ?? null;
  };

  return { store, loading, fetching, load, findForEvent };
};

export const useParticipationHistory = () => {
  const store = useParticipationsStore();
  const visible = ref(false);
  const loading = ref(false);
  const changes = ref<IParticipationChange[]>([]);

  const open = async (participationId: string) => {
    visible.value = true;
    loading.value = true;
    const result = await handleAsyncOperation(
      () => store.service.findChanges(participationId, { limit: 100, offset: 0 }),
      () => {
        loading.value = false;
      },
    );
    changes.value = result?.list ?? [];
  };

  return { visible, loading, changes, open };
};

export const isParticipationActive = (participation: IParticipation | null | undefined) =>
  !!participation && ACTIVE_STATUSES.includes(participation.status);

export const useEventApplications = () => {
  const service = useParticipationsStore().service;
  const { ensure: ensureVolunteers } = useVolunteerNames();
  const { toast, loader } = useApp();
  const loading = ref(true);
  const list = ref<IParticipation[]>([]);

  const load = async (eventId: string, silent = false) => {
    if (!silent) loading.value = true;
    await handleAsyncOperation(
      async () => {
        const result: IParticipation[] = [];
        for (let offset = 0; ; offset += 100) {
          const page = await service.find({ scope: ParticipationListScope.MY_EVENTS, eventId }, { limit: 100, offset });
          result.push(...page.list);
          if (result.length >= page.total || !page.list.length) break;
        }
        await ensureVolunteers(result.map((item) => item.volunteerId));
        list.value = result;
      },
      () => {
        loading.value = false;
      },
    );
  };

  const replace = (participation: IParticipation) => {
    list.value = list.value.map((item) => (item.id === participation.id ? participation : item));
  };

  const approve = async (id: string) => {
    loader.show();
    const result = await handleAsyncOperation(
      () => service.approve(id),
      () => loader.hide(),
    );
    if (!result) return;
    replace(result);
    toast.success('Заявка принята');
  };

  const reject = async (id: string, reason?: string) => {
    loader.show();
    const result = await handleAsyncOperation(
      () => service.reject(id, reason),
      () => loader.hide(),
    );
    if (!result) return false;
    const wasApproved = list.value.some((item) => item.id === id && item.status === ParticipationStatus.APPROVED);
    replace(result);
    toast.success(wasApproved ? 'Участник снят с мероприятия' : 'Заявка отклонена');
    return true;
  };

  const byStatus = (status: ParticipationStatus) => list.value.filter((item) => item.status === status);

  return { list, loading, load, approve, reject, byStatus };
};
