<script setup lang="ts">
import {
  COOPERATION_STATUS_LABELS,
  COOPERATION_STATUS_SEVERITY,
  CooperationStatus,
} from '~/core/cooperation/services/cooperation/cooperation.types';
import type { ICoordinator } from '~/core/coordinator/services/coordinator/coordinator.types';

definePageMeta({
  title: 'Координатор',
});

const route = useRoute();
const coordinatorsStore = useCoordinatorsStore();
const volunteerStore = useVolunteerStore();
const { store: cooperationsStore, load: loadCooperations, request, terminate } = useVolunteerCooperations();
const { area } = useSearchArea();
const { store: eventsStore, loading: eventsLoading, fetching: eventsFetching, load: loadEvents, loadMore } = useEventCatalog();
const { store: rewardsStore, loading: rewardsLoading, load: loadRewards } = useRewardCatalog();
const { purchase } = useRewardPurchase();
const { active: tab } = useProfileTabs(COORDINATOR_PUBLIC_TABS);

const coordinatorId = computed(() => String(route.params.id));
const loading = ref(true);
const coordinator = ref<ICoordinator | null>(null);

const cooperation = computed(() => cooperationsStore.withCoordinator(coordinatorId.value));
const isMember = computed(() => cooperation.value?.status === CooperationStatus.ACTIVE);
const canJoin = computed(() => canRequestCooperation(cooperation.value));
const balance = computed(() => volunteerStore.profile?.coinBalance ?? 0);
const stats = computed(() =>
  isMember.value
    ? [
        { label: 'Мероприятия', value: eventsStore.total },
        { label: 'Поощрения', value: rewardsStore.list.length },
      ]
    : [],
);

const loadContent = async () => {
  if (!isMember.value) return;
  await Promise.all([loadEvents(area.value, 'upcoming', 1, coordinatorId.value), loadRewards(coordinatorId.value)]);
};

onMounted(async () => {
  loading.value = true;
  await handleAsyncOperation(
    async () => {
      coordinator.value = await coordinatorsStore.fetchOne(coordinatorId.value);
      await Promise.all([loadCooperations(), volunteerStore.fetchProfile()]);
    },
    () => {
      loading.value = false;
    },
  );
  await loadContent();
});

watch(isMember, () => {
  if (!loading.value) loadContent();
});
watch(area, loadContent, { deep: true });
</script>

<template>
  <AppPage :back="ROUTES.volunteer.coordinators" back-label="Все координаторы">
    <Skeleton v-if="loading" height="260px" border-radius="20px" />

    <div v-else-if="!coordinator" class="app-card">
      <AppEmpty title="Координатор не найден" />
    </div>

    <template v-else>
      <AppProfileHero
        :name="coordinator.displayName"
        :avatar="coordinator.avatarUrl || null"
        :cover="PROFILE_COVER_COORDINATOR"
        :about="coordinator.description"
        :stats="stats">
        <template #chips>
          <span class="app-chip app-chip--violet"><i class="pi pi-briefcase" />Координатор</span>
          <AppStatus
            v-if="cooperation"
            :label="COOPERATION_STATUS_LABELS[cooperation.status]"
            :severity="COOPERATION_STATUS_SEVERITY[cooperation.status]" />
        </template>
        <template #actions>
          <Button v-if="canJoin" label="Предложить сотрудничество" @click="request(coordinatorId)" />
          <span v-else-if="cooperation?.status === CooperationStatus.REJECTED" class="app-card__hint">
            Повторная заявка с {{ formatDateTime(cooperation.retryAfter) }}
          </span>
          <span v-else-if="cooperation?.status === CooperationStatus.PENDING" class="app-card__hint">
            Координатор рассмотрит заявку и откроет доступ к мероприятиям
          </span>
          <Button
            v-if="isMember"
            label="Прекратить сотрудничество"
            severity="secondary"
            outlined
            @click="terminate(coordinatorId, coordinator.displayName)" />
        </template>
      </AppProfileHero>

      <template v-if="isMember">
        <AppTabs v-model="tab" :tabs="COORDINATOR_PUBLIC_TABS" />

        <template v-if="tab === 'events'">
          <div v-if="eventsLoading" class="app-cards">
            <Skeleton v-for="n in 3" :key="n" height="360px" border-radius="20px" />
          </div>
          <div v-else-if="!eventsStore.list.length" class="app-card">
            <AppEmpty title="Предстоящих мероприятий в выбранном районе нет" />
          </div>
          <template v-else>
            <div class="app-cards">
              <EventCard v-for="event in eventsStore.list" :key="event.id" :event="event" :to="`${ROUTES.volunteer.events}/${event.id}`">
                <template #actions>
                  <NuxtLink :to="`${ROUTES.volunteer.events}/${event.id}`">
                    <Button label="Подробнее" size="small" />
                  </NuxtLink>
                </template>
              </EventCard>
            </div>
            <div v-if="!eventsStore.pagination.loadAll" class="app-more">
              <Button label="Показать ещё" severity="secondary" outlined :loading="eventsFetching" @click="loadMore(area, 'upcoming')" />
            </div>
          </template>
        </template>

        <template v-else>
          <div v-if="rewardsLoading" class="app-cards">
            <Skeleton v-for="n in 3" :key="n" height="340px" border-radius="20px" />
          </div>
          <div v-else-if="!rewardsStore.list.length" class="app-card">
            <AppEmpty title="Поощрений пока нет" image="/images/landing/step-reward.jpg" />
          </div>
          <div v-else class="app-cards">
            <RewardCard
              v-for="reward in rewardsStore.list"
              :key="reward.id"
              :reward="reward"
              :balance="balance"
              @purchase="purchase(reward)" />
          </div>
        </template>
      </template>

      <div v-else class="app-card">
        <AppEmpty
          title="Мероприятия и поощрения координатора доступны волонтёрам, которые с ним сотрудничают"
          image="/images/landing/step-request.jpg" />
      </div>
    </template>
  </AppPage>
</template>
