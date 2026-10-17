import { CooperationListScope, CooperationStatus, type ICooperation } from '~/core/cooperation/services/cooperation/cooperation.types';

export const canRequestCooperation = (cooperation?: ICooperation) => {
  if (!cooperation) return true;
  if (cooperation.status === CooperationStatus.TERMINATED) return true;
  return cooperation.status === CooperationStatus.REJECTED && (!cooperation.retryAfter || isPast(cooperation.retryAfter));
};

export const useVolunteerCooperations = () => {
  const store = useCooperationsStore();
  const { toast, confirm, loader } = useApp();
  const loading = ref(true);

  const load = async (silent = false) => {
    if (!silent) loading.value = true;
    await handleAsyncOperation(
      () => store.fetchAll(CooperationListScope.AS_VOLUNTEER),
      () => {
        loading.value = false;
      },
    );
  };

  const request = async (coordinatorId: string) => {
    loader.show();
    const result = await handleAsyncOperation(
      () => store.request(coordinatorId),
      () => loader.hide(),
    );
    if (result) toast.success('Заявка на сотрудничество отправлена');
  };

  const terminate = (coordinatorId: string, name: string) =>
    confirm.delete(`Прекратить сотрудничество с координатором ${name}? Активные заявки и невыданные заказы будут отменены.`, {
      header: 'Прекращение сотрудничества',
      acceptLabel: 'Прекратить',
      rejectLabel: 'Отмена',
      accept: async () => {
        await store.terminateAsVolunteer(coordinatorId);
        toast.success('Сотрудничество прекращено');
      },
    });

  return { store, loading, load, request, terminate };
};

export const useCoordinatorCooperations = () => {
  const store = useCooperationsStore();
  const { ensure: ensureVolunteers } = useVolunteerNames();
  const { toast, confirm, loader } = useApp();
  const loading = ref(true);

  const load = async (silent = false) => {
    if (!silent) loading.value = true;
    await handleAsyncOperation(
      async () => {
        await Promise.all([store.fetchAll(CooperationListScope.AS_COORDINATOR), store.fetchInvitations()]);
        await ensureVolunteers(store.asCoordinator.map((item) => item.volunteerId));
      },
      () => {
        loading.value = false;
      },
    );
  };

  const run = async (fn: () => Promise<unknown>, message: string) => {
    loader.show();
    const ok = await handleAsyncOperation(
      async () => {
        await fn();
        return true;
      },
      () => loader.hide(),
    );
    if (ok) toast.success(message);
  };

  const approve = (volunteerId: string) => run(() => store.approve(volunteerId), 'Сотрудничество установлено');

  const reject = (volunteerId: string, name: string) =>
    confirm.delete(`Отклонить заявку волонтёра ${name}?`, {
      header: 'Отклонение заявки',
      acceptLabel: 'Отклонить',
      rejectLabel: 'Отмена',
      accept: async () => {
        await store.reject(volunteerId);
        toast.success('Заявка отклонена');
      },
    });

  const terminate = (volunteerId: string, name: string) =>
    confirm.delete(`Прекратить сотрудничество с волонтёром ${name}? Его активные заявки и невыданные заказы будут отменены.`, {
      header: 'Прекращение сотрудничества',
      acceptLabel: 'Прекратить',
      rejectLabel: 'Отмена',
      accept: async () => {
        await store.terminateAsCoordinator(volunteerId);
        toast.success('Сотрудничество прекращено');
      },
    });

  const createInvitation = (expiresAt: Date | null) => run(() => store.createInvitation(expiresAt), 'Приглашение создано');

  const revokeInvitation = (id: string) =>
    confirm.delete('Отозвать приглашение? Ссылка перестанет работать, текущие сотрудничества сохранятся.', {
      header: 'Отзыв приглашения',
      acceptLabel: 'Отозвать',
      rejectLabel: 'Отмена',
      accept: async () => {
        await store.revokeInvitation(id);
        toast.success('Приглашение отозвано');
      },
    });

  return { store, loading, load, approve, reject, terminate, createInvitation, revokeInvitation };
};
