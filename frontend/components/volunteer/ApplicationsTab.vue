<script setup lang="ts">
import {
  PARTICIPATION_CANCELLATION_LABELS,
  PARTICIPATION_STATUS_LABELS,
  PARTICIPATION_STATUS_SEVERITY,
  ParticipationStatus,
  type IParticipation,
} from '~/core/participation/services/participation/participation.types';

const { store, loading, fetching, load } = useMyParticipations();
const { event } = useEventLookup();
const { cancel, restore } = useParticipationActions();
const history = useParticipationHistory();

const status = ref<ParticipationStatus | null>(null);
const statusOptions = [
  { label: 'Все', value: null },
  ...[ParticipationStatus.PENDING, ParticipationStatus.APPROVED, ParticipationStatus.CONFIRMED, ParticipationStatus.CANCELLED].map(
    (value) => ({
      label: PARTICIPATION_STATUS_LABELS[value],
      value,
    }),
  ),
];

const isUpcoming = (item: IParticipation) => {
  const target = event(item.eventId);
  return !!target && canApplyToEvent(target);
};

const reload = () => load(status.value ?? undefined);

watch(status, reload, { immediate: true });
useLiveRefresh(() => load(status.value ?? undefined, 1, true));
</script>

<template>
  <div class="app-card">
    <AppFilter v-model="status" :options="statusOptions" />

    <div v-if="loading" class="flex flex-col gap-3">
      <Skeleton v-for="n in 3" :key="n" height="72px" border-radius="16px" />
    </div>

    <AppEmpty v-else-if="!store.list.length" title="Здесь появятся ваши заявки">
      <NuxtLink :to="ROUTES.volunteer.events">
        <Button label="Найти мероприятие" icon="pi pi-search" />
      </NuxtLink>
    </AppEmpty>

    <template v-else>
      <div class="app-list">
        <div v-for="item in store.list" :key="item.id" class="app-list__item">
          <div class="app-list__lead">
            <img :src="eventImage(event(item.eventId), item.eventId)" alt="" class="app-thumb" />
            <div class="app-list__main">
              <NuxtLink :to="`${ROUTES.volunteer.events}/${item.eventId}`" class="app-list__title">
                {{ event(item.eventId)?.title ?? 'Мероприятие' }}
              </NuxtLink>
              <div class="app-list__meta">
                <span>
                  <i class="pi pi-calendar" />
                  {{ formatEventPeriod(event(item.eventId)?.startsAt ?? null, event(item.eventId)?.endsAt ?? null) }}
                </span>
                <span v-if="item.awardedCoins !== null">
                  <i class="pi pi-star" />+{{ formatCoins(item.awardedCoins) }} за {{ formatDuration(item.awardedMinutes ?? 0) }}
                </span>
              </div>
              <span v-if="item.cancellationReason" class="app-list__comment">
                {{ PARTICIPATION_CANCELLATION_LABELS[item.cancellationReason] }}
              </span>
              <span v-if="item.rejectionReason" class="app-list__comment">{{ item.rejectionReason }}</span>
            </div>
          </div>
          <div class="app-list__actions">
            <AppStatus :label="PARTICIPATION_STATUS_LABELS[item.status]" :severity="PARTICIPATION_STATUS_SEVERITY[item.status]" />
            <Button
              v-if="isUpcoming(item) && isParticipationActive(item)"
              label="Отменить"
              severity="secondary"
              text
              size="small"
              @click="cancel(item.id)" />
            <Button
              v-if="isUpcoming(item) && item.status === ParticipationStatus.CANCELLED"
              label="Подать снова"
              size="small"
              @click="restore(item.id)" />
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
      <div v-if="!store.pagination.loadAll" class="app-more">
        <Button
          label="Показать ещё"
          severity="secondary"
          outlined
          :loading="fetching"
          @click="load(status ?? undefined, store.pagination.page)" />
      </div>
    </template>

    <ParticipationHistoryDialog v-model:visible="history.visible.value" :loading="history.loading.value" :changes="history.changes.value" />
  </div>
</template>
