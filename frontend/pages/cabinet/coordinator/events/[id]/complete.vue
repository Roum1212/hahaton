<script setup lang="ts">
import { EventStatus, type IEvent, type IEventResult } from '~/core/event/services/event/event.types';
import { ParticipationStatus } from '~/core/participation/services/participation/participation.types';

definePageMeta({
  title: 'Итоги мероприятия',
  permissions: [ROLES.coordinator],
});

interface IResultRow {
  participationId: string;
  volunteerId: string;
  attended: boolean;
  minutes: number;
}

const route = useRoute();
const eventsStore = useEventsStore();
const { list, load: loadApplications } = useEventApplications();
const { name: volunteerName, avatar: volunteerAvatar } = useVolunteerNames();
const { toast, confirm } = useApp();

const eventId = computed(() => String(route.params.id));
const loading = ref(true);
const event = ref<IEvent | null>(null);
const rows = ref<IResultRow[]>([]);
const requestId = ref(createRequestId());

const attendanceOptions = [
  { label: 'Был', value: true },
  { label: 'Не был', value: false },
];

const pendingCount = computed(() => list.value.filter((item) => item.status === ParticipationStatus.PENDING).length);
const completedResults = computed(() =>
  list.value.filter((item) => item.status === ParticipationStatus.CONFIRMED || item.status === ParticipationStatus.NOT_ATTENDED),
);
const summary = computed(() => {
  const attended = rows.value.filter((row) => row.attended);
  const minutes = attended.reduce((sum, row) => sum + row.minutes, 0);
  return {
    attended: attended.length,
    absent: rows.value.length - attended.length,
    minutes,
    coins: attended.reduce((sum, row) => sum + coinsForMinutes(row.minutes), 0),
  };
});

const load = async () => {
  loading.value = true;
  await handleAsyncOperation(
    async () => {
      event.value = await eventsStore.fetchOne(eventId.value, true);
      await loadApplications(eventId.value);
      rows.value = list.value
        .filter((item) => item.status === ParticipationStatus.APPROVED)
        .map((item) => ({
          participationId: item.id,
          volunteerId: item.volunteerId,
          attended: true,
          minutes: event.value!.durationMinutes,
        }));
    },
    () => {
      loading.value = false;
    },
  );
};

watch(
  rows,
  () => {
    requestId.value = createRequestId();
  },
  { deep: true },
);

const toResults = (): IEventResult[] =>
  rows.value.map((row) =>
    row.attended
      ? { participationId: row.participationId, confirmedMinutes: row.minutes }
      : { participationId: row.participationId, absent: true },
  );

const submit = () => {
  if (rows.value.some((row) => row.attended && row.minutes <= 0)) {
    toast.error('Укажите отработанное время для каждого присутствовавшего');
    return;
  }
  const { attended, absent, minutes, coins } = summary.value;
  confirm.accept(
    `Подтвердить итоги? Присутствовали ${attended}, не пришли ${absent}. Будет начислено ${formatDuration(minutes)} и ${formatCoins(coins)}. Изменить итоги после подтверждения нельзя.`,
    {
      header: 'Итоговая проверка',
      acceptLabel: 'Подтвердить итоги',
      accept: async () => {
        event.value = await eventsStore.complete(eventId.value, toResults(), requestId.value);
        toast.success('Итоги подведены, койны начислены');
        await loadApplications(eventId.value);
      },
    },
  );
};

onMounted(load);
</script>

<template>
  <AppPage :loading="loading" :back="`${ROUTES.coordinator.events}/${eventId}`" back-label="К мероприятию">
    <template #title>Итоги: {{ event?.title ?? '' }}</template>

    <Skeleton v-if="loading" height="260px" border-radius="20px" />

    <template v-else-if="event?.status === EventStatus.COMPLETED">
      <div class="app-card">
        <Message severity="success">Итоги зафиксированы, койны начислены участникам</Message>
        <div class="app-list">
          <div v-for="item in completedResults" :key="item.id" class="app-list__item">
            <div class="app-list__lead">
              <AppAvatar :name="volunteerName(item.volunteerId)" :image="volunteerAvatar(item.volunteerId)" />
              <NuxtLink :to="`${ROUTES.coordinator.volunteers}/${item.volunteerId}`" class="app-list__title">{{
                volunteerName(item.volunteerId)
              }}</NuxtLink>
            </div>
            <span v-if="item.status === ParticipationStatus.CONFIRMED" class="app-complete__coins">
              {{ formatDuration(item.awardedMinutes ?? 0) }} · +{{ formatCoins(item.awardedCoins ?? 0) }}
            </span>
            <AppStatus v-else label="Не пришёл" severity="danger" />
          </div>
        </div>
      </div>
    </template>

    <template v-else-if="event && !canCompleteEvent(event)">
      <Message severity="warn">Подвести итоги можно только после окончания опубликованного мероприятия</Message>
    </template>

    <template v-else-if="event">
      <div v-if="!rows.length" class="app-card">
        <AppEmpty title="Принятых участников нет" image="/images/landing/step-results.jpg" />
      </div>

      <div v-else class="app-card">
        <span class="app-card__title">Кто пришёл и сколько помогал</span>
        <div class="app-list">
          <div v-for="row in rows" :key="row.participationId" class="app-list__item">
            <div class="app-list__lead">
              <AppAvatar :name="volunteerName(row.volunteerId)" :image="volunteerAvatar(row.volunteerId)" />
              <NuxtLink :to="`${ROUTES.coordinator.volunteers}/${row.volunteerId}`" class="app-list__title">{{
                volunteerName(row.volunteerId)
              }}</NuxtLink>
            </div>
            <div class="app-complete__controls">
              <SelectButton
                v-model="row.attended"
                :options="attendanceOptions"
                option-label="label"
                option-value="value"
                :allow-empty="false"
                size="small" />
              <EventDurationInput v-model="row.minutes" :disabled="!row.attended" class="app-complete__duration" />
              <span class="app-complete__coins">{{ row.attended ? `+${formatCoins(coinsForMinutes(row.minutes))}` : '—' }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="app-card app-complete__summary">
        <div class="app-props">
          <div class="app-props__row">
            <span class="app-props__key">Присутствовали</span>
            <span class="app-props__value">{{ summary.attended }}</span>
          </div>
          <div class="app-props__row">
            <span class="app-props__key">Не пришли</span>
            <span class="app-props__value">{{ summary.absent }}</span>
          </div>
          <div class="app-props__row">
            <span class="app-props__key">Всего времени</span>
            <span class="app-props__value">{{ formatDuration(summary.minutes) }}</span>
          </div>
          <div class="app-props__row">
            <span class="app-props__key">Будет начислено</span>
            <span class="app-props__value">{{ formatCoins(summary.coins) }}</span>
          </div>
          <div v-if="pendingCount" class="app-props__row">
            <span class="app-props__key">Закроются без рассмотрения</span>
            <span class="app-props__value">{{ pendingCount }}</span>
          </div>
        </div>
        <div class="app-form__actions">
          <Button label="Проверить и завершить" @click="submit" />
        </div>
      </div>
    </template>
  </AppPage>
</template>
