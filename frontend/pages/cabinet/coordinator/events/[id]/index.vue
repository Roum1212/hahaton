<script setup lang="ts">
import {
  EVENT_REQUIREMENT_LABELS,
  EVENT_STATUS_LABELS,
  EVENT_STATUS_SEVERITY,
  EventStatus,
  type IEvent,
} from '~/core/event/services/event/event.types';
import {
  PARTICIPATION_CANCELLATION_LABELS,
  PARTICIPATION_STATUS_LABELS,
  PARTICIPATION_STATUS_SEVERITY,
  ParticipationStatus,
} from '~/core/participation/services/participation/participation.types';

definePageMeta({
  title: 'Мероприятие',
  permissions: [ROLES.coordinator],
});

const route = useRoute();
const eventsStore = useEventsStore();
const { list, loading: applicationsLoading, load: loadApplications, approve, reject } = useEventApplications();
const { name: volunteerName, view: volunteerView, avatar: volunteerAvatar } = useVolunteerNames();
const { cancel } = useEventActions();
const history = useParticipationHistory();
const { active: section } = useProfileTabs(OWNED_EVENT_TABS);

const eventId = computed(() => String(route.params.id));
const loading = ref(true);
const event = ref<IEvent | null>(null);
const tab = ref<'pending' | 'approved' | 'other'>('pending');
const rejectTarget = ref<string | null>(null);
const rejectReason = ref('');
const participantStatuses = [
  ParticipationStatus.APPROVED,
  ParticipationStatus.CONFIRMED,
  ParticipationStatus.NOT_ATTENDED,
];

const tabs = computed(() => [
  {
    value: 'pending',
    label: `Новые ${list.value.filter((item) => item.status === ParticipationStatus.PENDING).length}`,
  },
  {
    value: 'approved',
    label: `Участники ${list.value.filter((item) => participantStatuses.includes(item.status)).length}`,
  },
  { value: 'other', label: 'Остальные' },
]);

const confirmedCount = computed(
  () => list.value.filter((item) => item.status === ParticipationStatus.CONFIRMED).length,
);

const visible = computed(() =>
  list.value.filter((item) => {
    if (tab.value === 'pending') return item.status === ParticipationStatus.PENDING;
    if (tab.value === 'approved') return participantStatuses.includes(item.status);
    return [ParticipationStatus.REJECTED, ParticipationStatus.CANCELLED, ParticipationStatus.CLOSED].includes(
      item.status,
    );
  }),
);

const canDecide = computed(() => !!event.value && canCancelEvent(event.value));
const canReject = computed(() => event.value?.status === EventStatus.PUBLISHED);
const rejectingApproved = computed(() =>
  list.value.some((item) => item.id === rejectTarget.value && item.status === ParticipationStatus.APPROVED),
);

const refresh = async () => {
  event.value = await eventsStore.fetchOne(eventId.value, true);
};

const load = async () => {
  loading.value = true;
  await handleAsyncOperation(refresh, () => {
    loading.value = false;
  });
  await loadApplications(eventId.value);
};

const onApprove = async (id: string) => {
  await approve(id);
  await handleAsyncOperation(refresh);
};

const openReject = (id: string) => {
  rejectTarget.value = id;
  rejectReason.value = '';
};

const confirmReject = async () => {
  if (!rejectTarget.value) return;
  const ok = await reject(rejectTarget.value, rejectReason.value);
  if (ok) {
    rejectTarget.value = null;
    await handleAsyncOperation(refresh);
  }
};

const onCancelEvent = () => {
  if (event.value) cancel(event.value, () => load());
};

onMounted(load);
</script>

<template>
  <AppPage :loading="loading" :back="ROUTES.coordinator.events" back-label="Мои мероприятия">
    <div v-if="loading" class="app-columns app-columns--left">
      <Skeleton height="360px" border-radius="20px" />
      <Skeleton height="420px" border-radius="20px" />
    </div>

    <div v-else-if="!event" class="app-card">
      <AppEmpty title="Мероприятие недоступно" />
    </div>

    <div v-else class="app-columns app-columns--left">
      <aside class="app-columns__aside">
        <div class="app-card app-event-aside">
          <div class="app-event-aside__image" :class="{ 'app-event-aside__image--photo': event.imageUrl }">
            <img :src="eventImage(event, event.id)" alt="" />
          </div>
          <span v-if="event.status === EventStatus.PUBLISHED" class="app-event-aside__count">
            Участвуют <b>{{ event.approvedCount }}</b> из {{ event.volunteerLimit }}
          </span>
          <span
            v-else-if="event.status === EventStatus.COMPLETED && !applicationsLoading"
            class="app-event-aside__count">
            Подтверждённых участников: {{ confirmedCount }}
          </span>
          <div class="app-event-aside__actions">
            <NuxtLink
              v-if="event.status === EventStatus.COMPLETED"
              :to="`${ROUTES.coordinator.events}/${event.id}/complete`">
              <Button label="Посмотреть итоги" severity="secondary" outlined fluid />
            </NuxtLink>
            <NuxtLink v-if="canCompleteEvent(event)" :to="`${ROUTES.coordinator.events}/${event.id}/complete`">
              <Button label="Подвести итоги" size="large" fluid />
            </NuxtLink>
            <Button
              label="Чат мероприятия"
              icon="pi pi-comments"
              severity="secondary"
              outlined
              fluid
              @click="section = 'chat'" />
            <NuxtLink v-if="canEditEvent(event)" :to="`${ROUTES.coordinator.events}/${event.id}/edit`">
              <Button label="Редактировать" icon="pi pi-pencil" severity="secondary" outlined fluid />
            </NuxtLink>
            <Button
              v-if="canCancelEvent(event)"
              label="Отменить мероприятие"
              severity="danger"
              text
              fluid
              @click="onCancelEvent" />
          </div>
        </div>
      </aside>

      <div class="app-columns__main">
        <div class="app-card">
          <div class="app-hero__chips">
            <AppStatus :label="EVENT_STATUS_LABELS[event.status]" :severity="EVENT_STATUS_SEVERITY[event.status]" />
            <span v-if="event.status === EventStatus.PUBLISHED" class="app-chip app-chip--accent">
              <i class="pi pi-star" />План: +{{ formatCoins(coinsForMinutes(event.durationMinutes)) }}
            </span>
          </div>
          <h1 class="app-page__title">{{ event.title }}</h1>
          <div class="app-info">
            <div class="app-info__item">
              <i class="pi pi-calendar" />
              <span><b>Когда</b>{{ formatEventPeriod(event.startsAt, event.endsAt) }}</span>
            </div>
            <div class="app-info__item">
              <i class="pi pi-clock" />
              <span><b>Плановая длительность</b>{{ formatDuration(event.durationMinutes) }}</span>
            </div>
            <div class="app-info__item">
              <i class="pi pi-inbox" />
              <span><b>Новых заявок</b>{{ event.pendingCount }}</span>
            </div>
          </div>
          <p v-if="event.description" class="app-text">{{ event.description }}</p>
          <div v-if="event.requirements.length" class="app-hero__chips">
            <span v-for="item in event.requirements" :key="item" class="app-chip">{{
              EVENT_REQUIREMENT_LABELS[item]
            }}</span>
          </div>
        </div>

        <EventLocation :event="event" />

        <AppTabs v-model="section" :tabs="OWNED_EVENT_TABS" :counts="{ applications: event.pendingCount }" />

        <EventChat v-if="section === 'chat'" :event-id="event.id" />

        <div v-else class="app-card">
          <AppFilter v-model="tab" :options="tabs" />

          <div v-if="applicationsLoading" class="flex flex-col gap-3">
            <Skeleton v-for="n in 2" :key="n" height="64px" border-radius="16px" />
          </div>

          <AppEmpty v-else-if="!visible.length" title="Заявок нет" image="/images/landing/step-request.jpg" />

          <div v-else class="app-list">
            <div v-for="item in visible" :key="item.id" class="app-list__item">
              <div class="app-list__lead">
                <NuxtLink :to="`${ROUTES.coordinator.volunteers}/${item.volunteerId}`"
                  ><AppAvatar :name="volunteerName(item.volunteerId)" :image="volunteerAvatar(item.volunteerId)"
                /></NuxtLink>
                <div class="app-list__main">
                  <span class="app-list__title">
                    <NuxtLink :to="`${ROUTES.coordinator.volunteers}/${item.volunteerId}`">{{
                      volunteerName(item.volunteerId)
                    }}</NuxtLink>
                    <AppStatus
                      :label="PARTICIPATION_STATUS_LABELS[item.status]"
                      :severity="PARTICIPATION_STATUS_SEVERITY[item.status]" />
                  </span>
                  <div class="app-list__meta">
                    <span v-if="volunteerView(item.volunteerId)?.phone">
                      <i class="pi pi-phone" />{{ volunteerView(item.volunteerId)?.phone }}
                    </span>
                    <span v-if="volunteerView(item.volunteerId)?.email">
                      <i class="pi pi-envelope" />{{ volunteerView(item.volunteerId)?.email }}
                    </span>
                    <span v-if="item.awardedCoins !== null"
                      ><i class="pi pi-star" />+{{ formatCoins(item.awardedCoins) }}</span
                    >
                  </div>
                  <span v-if="item.cancellationReason" class="app-list__comment">
                    {{ PARTICIPATION_CANCELLATION_LABELS[item.cancellationReason] }}
                  </span>
                  <span v-if="item.rejectionReason" class="app-list__comment">{{ item.rejectionReason }}</span>
                </div>
              </div>
              <div class="app-list__actions">
                <Button
                  v-if="canReject && item.status === ParticipationStatus.APPROVED"
                  label="Снять с мероприятия"
                  severity="secondary"
                  text
                  size="small"
                  @click="openReject(item.id)" />
                <template v-if="item.status === ParticipationStatus.PENDING">
                  <Button
                    v-if="canReject"
                    label="Отклонить"
                    severity="secondary"
                    text
                    size="small"
                    @click="openReject(item.id)" />
                  <Button
                    v-if="canDecide"
                    label="Принять"
                    size="small"
                    :disabled="event.freePlaces === 0"
                    @click="onApprove(item.id)" />
                </template>
                <Button
                  v-tooltip.bottom="'История заявки'"
                  icon="pi pi-history"
                  severity="secondary"
                  text
                  rounded
                  size="small"
                  aria-label="История заявки"
                  @click="history.open(item.id)" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <Dialog
      :visible="!!rejectTarget"
      modal
      :header="rejectingApproved ? 'Снятие участника' : 'Отклонение заявки'"
      :draggable="false"
      :style="{ width: '460px' }"
      @update:visible="(value: boolean) => !value && (rejectTarget = null)">
      <p class="app-text mb-4">
        {{
          rejectingApproved
            ? 'Волонтёр потеряет место и доступ к чату, койны за мероприятие не начислятся. ' +
              'Подать заявку на это мероприятие снова он не сможет.'
            : 'Подать заявку на это мероприятие снова волонтёр не сможет.'
        }}
      </p>
      <label class="app-field">
        <span class="app-field__label">Причина (видит только волонтёр)</span>
        <Textarea v-model="rejectReason" rows="3" maxlength="2000" auto-resize fluid />
      </label>
      <template #footer>
        <Button label="Отмена" severity="secondary" text @click="rejectTarget = null" />
        <Button
          :label="rejectingApproved ? 'Снять с мероприятия' : 'Отклонить'"
          severity="danger"
          @click="confirmReject" />
      </template>
    </Dialog>

    <ParticipationHistoryDialog
      v-model:visible="history.visible.value"
      :loading="history.loading.value"
      :changes="history.changes.value" />
  </AppPage>
</template>
