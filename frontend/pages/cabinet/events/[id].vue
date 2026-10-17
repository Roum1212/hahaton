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
  type IParticipation,
} from '~/core/participation/services/participation/participation.types';

definePageMeta({
  title: 'Мероприятие',
});

const route = useRoute();
const eventsStore = useEventsStore();
const { ensure: ensureCoordinators, name: coordinatorName } = useCoordinatorNames();
const { findForEvent } = useMyParticipations();
const { request, cancel, restore } = useParticipationActions();
const history = useParticipationHistory();
const { active: tab } = useProfileTabs(EVENT_TABS);

const eventId = computed(() => String(route.params.id));
const loading = ref(true);
const event = ref<IEvent | null>(null);
const participation = ref<IParticipation | null>(null);

const CHAT_STATUSES = [ParticipationStatus.APPROVED, ParticipationStatus.CONFIRMED, ParticipationStatus.NOT_ATTENDED];

const result = computed(() => {
  const value = participation.value;
  if (event.value?.status !== EventStatus.COMPLETED || value?.awardedCoins == null) return null;
  return { coins: value.awardedCoins, minutes: value.awardedMinutes ?? 0 };
});

const canApply = computed(() => !!event.value && canApplyToEvent(event.value));
const canCancel = computed(() => canApply.value && isParticipationActive(participation.value));
const canRestore = computed(() => canApply.value && participation.value?.status === ParticipationStatus.CANCELLED);
const canChat = computed(() => !!participation.value && CHAT_STATUSES.includes(participation.value.status));

const load = async () => {
  loading.value = true;
  await handleAsyncOperation(
    async () => {
      event.value = await eventsStore.fetchOne(eventId.value, true);
      await ensureCoordinators([event.value.coordinatorId]);
      participation.value = await findForEvent(eventId.value);
    },
    () => {
      loading.value = false;
    },
  );
};

const refresh = async () => {
  event.value = await eventsStore.fetchOne(eventId.value, true);
  participation.value = await findForEvent(eventId.value);
};

const onRequest = async () => {
  if (await request(eventId.value)) await refresh();
};

const onRestore = async () => {
  if (participation.value && (await restore(participation.value.id))) await refresh();
};

const onCancel = () => {
  if (participation.value) cancel(participation.value.id, () => refresh());
};

onMounted(load);
useLiveRefresh(() => handleAsyncOperation(refresh));
</script>

<template>
  <AppPage :loading="loading" :back="ROUTES.volunteer.events" back-label="Все мероприятия">
    <div v-if="loading" class="app-columns app-columns--left">
      <Skeleton height="360px" border-radius="20px" />
      <Skeleton height="420px" border-radius="20px" />
    </div>

    <div v-else-if="!event" class="app-card">
      <AppEmpty title="Мероприятие недоступно">
        <NuxtLink :to="ROUTES.volunteer.events">
          <Button label="К мероприятиям" />
        </NuxtLink>
      </AppEmpty>
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
          <div v-if="participation" class="flex justify-center">
            <AppStatus
              :label="PARTICIPATION_STATUS_LABELS[participation.status]"
              :severity="PARTICIPATION_STATUS_SEVERITY[participation.status]" />
          </div>
          <div class="app-event-aside__actions">
            <Button v-if="!participation && canApply" label="Подать заявку" size="large" fluid @click="onRequest" />
            <Button v-if="canRestore" label="Подать заявку снова" size="large" fluid @click="onRestore" />
            <Button
              v-if="canChat"
              label="Чат мероприятия"
              icon="pi pi-comments"
              severity="secondary"
              outlined
              fluid
              @click="tab = 'chat'" />
            <Button v-if="canCancel" label="Отменить заявку" severity="secondary" text fluid @click="onCancel" />
            <Button
              v-if="participation"
              label="История заявки"
              severity="secondary"
              text
              fluid
              @click="history.open(participation.id)" />
          </div>
          <span v-if="!participation && !canApply" class="app-card__hint text-center">Приём заявок закрыт</span>
        </div>
      </aside>

      <div class="app-columns__main">
        <div class="app-card">
          <div class="app-hero__chips">
            <AppStatus
              v-if="event.status !== EventStatus.PUBLISHED"
              :label="EVENT_STATUS_LABELS[event.status]"
              :severity="EVENT_STATUS_SEVERITY[event.status]" />
            <span v-if="result" class="app-chip app-chip--accent">
              <i class="pi pi-star" />Начислено {{ formatCoins(result.coins) }}
            </span>
            <span v-else-if="event.status === EventStatus.PUBLISHED" class="app-chip app-chip--accent">
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
              <i class="pi pi-user" />
              <span>
                <b>Организатор</b>
                <NuxtLink :to="`${ROUTES.volunteer.coordinators}/${event.coordinatorId}`" class="app-link">
                  {{ coordinatorName(event.coordinatorId) }}
                </NuxtLink>
              </span>
            </div>
            <div class="app-info__item">
              <i class="pi pi-clock" />
              <span><b>Плановая длительность</b>{{ formatDuration(event.durationMinutes) }}</span>
            </div>
            <div v-if="result" class="app-info__item">
              <i class="pi pi-check-circle" />
              <span><b>Подтверждено</b>{{ formatDuration(result.minutes) }}</span>
            </div>
          </div>
        </div>

        <AppTabs v-if="canChat" v-model="tab" :tabs="EVENT_TABS" />

        <EventChat v-if="canChat && tab === 'chat'" :event-id="event.id" />

        <template v-else>
          <EventLocation :event="event" />

          <div
            v-if="
              participation &&
              (participation.rejectionReason || participation.cancellationReason || participation.awardedCoins !== null)
            "
            class="app-card">
            <span class="app-card__title">Моя заявка</span>
            <span v-if="participation.cancellationReason" class="app-card__hint">
              {{ PARTICIPATION_CANCELLATION_LABELS[participation.cancellationReason] }}
            </span>
            <Message v-if="participation.rejectionReason" severity="warn" size="small">{{
              participation.rejectionReason
            }}</Message>
            <span v-if="participation.awardedCoins !== null" class="app-text">
              Подтверждено {{ formatDuration(participation.awardedMinutes ?? 0) }}, начислено +{{
                formatCoins(participation.awardedCoins)
              }}
            </span>
          </div>

          <div v-if="event.description" class="app-card">
            <span class="app-card__title">О мероприятии</span>
            <p class="app-text">{{ event.description }}</p>
          </div>

          <div v-if="event.requirements.length" class="app-card">
            <span class="app-card__title">Условия участия</span>
            <div class="app-hero__chips">
              <span v-for="item in event.requirements" :key="item" class="app-chip">{{
                EVENT_REQUIREMENT_LABELS[item]
              }}</span>
            </div>
          </div>
        </template>
      </div>
    </div>

    <ParticipationHistoryDialog
      v-model:visible="history.visible.value"
      :loading="history.loading.value"
      :changes="history.changes.value" />
  </AppPage>
</template>
