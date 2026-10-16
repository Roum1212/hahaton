import {
  EVENT_AGE_REQUIREMENTS,
  EventListScope,
  EventStatus,
  type IEvent,
  type IEventDto,
  type IEventFilter,
} from '~/core/event/services/event/event.types';
import type { ISearchArea } from '~/core/geo/services/geo/geo.types';

export type CatalogPeriod = 'upcoming' | 'all';

export const useEventLookup = () => {
  const store = useEventsStore();

  const ensure = async (ids: string[]) => {
    const missing = [...new Set(ids)].filter((id) => id && !store.getById(id));
    await Promise.all(missing.map((id) => handleAsyncOperation(() => store.fetchOne(id))));
  };

  const event = (id: string) => store.getById(id);

  return { ensure, event };
};

export const useEventCatalog = () => {
  const store = useEventsStore();
  const { ensure: ensureCoordinators } = useCoordinatorNames();
  const loading = ref(true);
  const fetching = ref(false);

  const load = async (area: ISearchArea, period: CatalogPeriod, page = 1, coordinatorId?: string) => {
    if (page === 1) loading.value = true;
    else fetching.value = true;
    await handleAsyncOperation(
      async () => {
        await store.find(page, {
          scope: EventListScope.CATALOG,
          searchArea: area,
          coordinatorId,
          status: period === 'upcoming' ? EventStatus.PUBLISHED : undefined,
          startsFrom: period === 'upcoming' ? new Date() : undefined,
        });
        await ensureCoordinators(store.list.map((item) => item.coordinatorId));
      },
      () => {
        loading.value = false;
        fetching.value = false;
      },
    );
  };

  const loadMore = async (area: ISearchArea, period: CatalogPeriod) => {
    if (store.pagination.loadAll || fetching.value) return;
    await load(area, period, store.pagination.page);
  };

  return { store, loading, fetching, load, loadMore };
};

export const canApplyToEvent = (event: IEvent) => event.status === EventStatus.PUBLISHED && !isPast(event.startsAt);

export type OwnedEventsFilter = 'upcoming' | 'awaiting' | 'completed' | 'cancelled' | 'all';

export const OWNED_EVENTS_FILTERS: { label: string; value: OwnedEventsFilter }[] = [
  { label: 'Предстоящие', value: 'upcoming' },
  { label: 'Ожидают итогов', value: 'awaiting' },
  { label: 'Завершённые', value: 'completed' },
  { label: 'Отменённые', value: 'cancelled' },
  { label: 'Все', value: 'all' },
];

export const toOwnedFilter = (filter: OwnedEventsFilter): IEventFilter => {
  const now = new Date();
  if (filter === 'upcoming') return { scope: EventListScope.OWNED, status: EventStatus.PUBLISHED, startsFrom: now };
  if (filter === 'awaiting') return { scope: EventListScope.OWNED, status: EventStatus.PUBLISHED, endsTo: now };
  if (filter === 'completed') return { scope: EventListScope.OWNED, status: EventStatus.COMPLETED };
  if (filter === 'cancelled') return { scope: EventListScope.OWNED, status: EventStatus.CANCELLED };
  return { scope: EventListScope.OWNED };
};

export const useOwnedEvents = () => {
  const store = useEventsStore();
  const loading = ref(true);
  const fetching = ref(false);

  const load = async (filter: OwnedEventsFilter, page = 1) => {
    if (page === 1) loading.value = true;
    else fetching.value = true;
    await handleAsyncOperation(
      () => store.find(page, toOwnedFilter(filter)),
      () => {
        loading.value = false;
        fetching.value = false;
      },
    );
  };

  return { store, loading, fetching, load };
};

export const canEditEvent = (event: IEvent) =>
  event.status === EventStatus.PUBLISHED && !isPast(event.startsAt) && event.approvedCount + event.pendingCount === 0;

export const canCancelEvent = (event: IEvent) => event.status === EventStatus.PUBLISHED && !isPast(event.startsAt);

export const canCompleteEvent = (event: IEvent) => event.status === EventStatus.PUBLISHED && isPast(event.endsAt);

export const useEventForm = () => {
  const generateInitialData = (event?: IEvent): IEventDto & { durationMinutes: number } => {
    const start = event?.startsAt ?? nextHalfHour(new Date(Date.now() + 24 * 60 * 60 * 1000));
    return {
      title: event?.title ?? '',
      description: event?.description ?? '',
      location: event?.location ?? '',
      point: event?.point ?? null,
      startsAt: start,
      endsAt: event?.endsAt ?? new Date(start.getTime() + 180 * 60000),
      durationMinutes: event?.durationMinutes ?? 180,
      requirements: [...(event?.requirements ?? [])],
      volunteerLimit: event?.volunteerLimit ?? 10,
    };
  };

  const validate = (dto: IEventDto & { durationMinutes: number }): string | null => {
    if (!dto.title.trim()) return 'Укажите название';
    if (!dto.point) return 'Отметьте место проведения на карте';
    if (!dto.startsAt || isPast(dto.startsAt)) return 'Начало должно быть в будущем';
    if (dto.durationMinutes <= 0) return 'Укажите длительность';
    if (dto.durationMinutes % MINUTES_STEP !== 0) return 'Длительность должна быть кратна 30 минутам';
    if (dto.volunteerLimit <= 0) return 'Укажите количество мест';
    if (dto.requirements.filter((item) => EVENT_AGE_REQUIREMENTS.includes(item)).length > 1) return 'Выберите один возрастной порог';
    return null;
  };

  const toDto = (form: IEventDto & { durationMinutes: number }): IEventDto => {
    const startsAt = form.startsAt ? new Date(form.startsAt) : null;
    if (startsAt) startsAt.setSeconds(0, 0);
    return {
      title: form.title.trim(),
      description: form.description.trim(),
      location: form.location.trim(),
      point: form.point,
      startsAt,
      endsAt: startsAt ? new Date(startsAt.getTime() + form.durationMinutes * 60000) : null,
      requirements: form.requirements,
      volunteerLimit: form.volunteerLimit,
    };
  };

  return { generateInitialData, validate, toDto };
};

export const nextHalfHour = (value: Date) => {
  const result = new Date(value);
  result.setSeconds(0, 0);
  const minutes = result.getMinutes();
  result.setMinutes(minutes <= 30 ? 30 : 60);
  if (minutes === 0) result.setMinutes(0);
  return result;
};

export const useEventActions = () => {
  const store = useEventsStore();
  const { toast, confirm } = useApp();

  const cancel = (event: IEvent, onDone?: () => void) =>
    confirm.delete(`Отменить мероприятие «${event.title}»? Все активные заявки будут отменены.`, {
      header: 'Отмена мероприятия',
      acceptLabel: 'Отменить мероприятие',
      rejectLabel: 'Не отменять',
      accept: async () => {
        await store.cancel(event.id);
        toast.success('Мероприятие отменено');
        onDone?.();
      },
    });

  return { cancel };
};

export const useEventCreateDialog = () => {
  const visible = useState('eventCreateDialog', () => false);

  const open = () => {
    visible.value = true;
  };

  return { visible, open };
};
