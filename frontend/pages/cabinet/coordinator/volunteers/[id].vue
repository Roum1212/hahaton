<script setup lang="ts">
import {
  COOPERATION_STATUS_LABELS,
  COOPERATION_STATUS_SEVERITY,
  CooperationStatus,
} from '~/core/cooperation/services/cooperation/cooperation.types';
import {
  PARTICIPATION_CANCELLATION_LABELS,
  PARTICIPATION_STATUS_LABELS,
  PARTICIPATION_STATUS_SEVERITY,
} from '~/core/participation/services/participation/participation.types';
import { REWARD_ORDER_STATUS_LABELS, REWARD_ORDER_STATUS_SEVERITY } from '~/core/reward-order/services/reward-order/reward-order.types';

definePageMeta({
  title: 'Волонтёр',
  permissions: [ROLES.coordinator],
});

const route = useRoute();
const volunteerId = computed(() => String(route.params.id));
const { loading, volunteer, cooperation, history, rejections, orders, stats, load } = useCoordinatorVolunteerProfile(volunteerId);
const { approve, reject, terminate } = useCoordinatorCooperations();
const { event } = useEventLookup();
const participationHistory = useParticipationHistory();
const { active: tab } = useProfileTabs(VOLUNTEER_PUBLIC_TABS);

const name = computed(() => volunteer.value?.fullName || 'Волонтёр');
const heroStats = computed(() => [
  { label: 'Койны', value: stats.value.coins },
  { label: 'Часы', value: formatHours(stats.value.minutes) },
  { label: 'Мероприятия', value: stats.value.events },
]);
const counts = computed(() => ({ participations: history.value.length, rejections: rejections.value.length, orders: orders.value.length }));

const onApprove = async () => {
  await approve(volunteerId.value);
  await load();
};

onMounted(load);
</script>

<template>
  <AppPage :back="ROUTES.coordinator.volunteers" back-label="Все волонтёры">
    <Skeleton v-if="loading && !volunteer" height="260px" border-radius="20px" />

    <div v-else-if="!volunteer" class="app-card">
      <AppEmpty title="Волонтёр недоступен" />
    </div>

    <template v-else>
      <AppProfileHero :name="name" :avatar="volunteer.avatarUrl" :cover="PROFILE_COVER_VOLUNTEER" :stats="heroStats">
        <template #chips>
          <span class="app-chip app-chip--accent"><i class="pi pi-heart-fill" />Волонтёр</span>
          <AppStatus
            v-if="cooperation"
            :label="COOPERATION_STATUS_LABELS[cooperation.status]"
            :severity="COOPERATION_STATUS_SEVERITY[cooperation.status]" />
        </template>
        <template #meta>
          <span v-if="volunteer.email"><i class="pi pi-envelope" />{{ volunteer.email }}</span>
          <span v-if="volunteer.phone"><i class="pi pi-phone" />{{ volunteer.phone }}</span>
          <span v-if="!volunteer.email && !volunteer.phone">
            Контакты скрыты
            <i v-tooltip.bottom="'Контакты видны, пока с волонтёром действует сотрудничество'" class="pi pi-question-circle" />
          </span>
        </template>
        <template v-if="cooperation?.status === CooperationStatus.PENDING || cooperation?.status === CooperationStatus.ACTIVE" #actions>
          <template v-if="cooperation.status === CooperationStatus.PENDING">
            <Button label="Отклонить" severity="secondary" text @click="reject(volunteerId, name)" />
            <Button label="Принять заявку" @click="onApprove" />
          </template>
          <Button v-else label="Прекратить сотрудничество" severity="secondary" outlined @click="terminate(volunteerId, name)" />
        </template>
      </AppProfileHero>

      <AppTabs v-model="tab" :tabs="VOLUNTEER_PUBLIC_TABS" :counts="counts" />

      <div class="app-card">
        <div v-if="loading" class="flex flex-col gap-3">
          <Skeleton v-for="n in 3" :key="n" height="64px" border-radius="16px" />
        </div>

        <template v-else-if="tab === 'participations'">
          <AppEmpty v-if="!history.length" title="В ваших мероприятиях волонтёр ещё не участвовал" />
          <div v-else class="app-list">
            <div v-for="item in history" :key="item.id" class="app-list__item">
              <div class="app-list__lead">
                <img :src="eventImage(event(item.eventId), item.eventId)" alt="" class="app-thumb" />
                <div class="app-list__main">
                  <NuxtLink :to="`${ROUTES.coordinator.events}/${item.eventId}`" class="app-list__title">
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
                </div>
              </div>
              <div class="app-list__actions">
                <AppStatus :label="PARTICIPATION_STATUS_LABELS[item.status]" :severity="PARTICIPATION_STATUS_SEVERITY[item.status]" />
                <Button
                  v-tooltip.bottom="'История заявки'"
                  icon="pi pi-history"
                  severity="secondary"
                  text
                  rounded
                  size="small"
                  aria-label="История заявки"
                  @click="participationHistory.open(item.id)" />
              </div>
            </div>
          </div>
        </template>

        <template v-else-if="tab === 'rejections'">
          <AppEmpty v-if="!rejections.length" title="Отказов не было" image="/images/landing/step-results.jpg" />
          <div v-else class="app-list">
            <div v-for="item in rejections" :key="item.id" class="app-list__item">
              <div class="app-list__lead">
                <img :src="eventImage(event(item.eventId), item.eventId)" alt="" class="app-thumb" />
                <div class="app-list__main">
                  <NuxtLink :to="`${ROUTES.coordinator.events}/${item.eventId}`" class="app-list__title">
                    {{ event(item.eventId)?.title ?? 'Мероприятие' }}
                  </NuxtLink>
                  <div class="app-list__meta">
                    <span><i class="pi pi-clock" />Отказ {{ formatDateTime(item.updatedAt) }}</span>
                  </div>
                  <span class="app-list__comment">{{ item.rejectionReason || 'Причина не указана' }}</span>
                </div>
              </div>
              <div class="app-list__actions">
                <AppStatus :label="PARTICIPATION_STATUS_LABELS[item.status]" :severity="PARTICIPATION_STATUS_SEVERITY[item.status]" />
              </div>
            </div>
          </div>
        </template>

        <template v-else>
          <AppEmpty v-if="!orders.length" title="Заказов на ваши поощрения нет" image="/images/landing/step-reward.jpg" />
          <div v-else class="app-list">
            <div v-for="order in orders" :key="order.id" class="app-list__item">
              <div class="app-list__lead">
                <img :src="pickIllustration(REWARD_ILLUSTRATIONS, order.rewardId)" alt="" class="app-thumb" />
                <div class="app-list__main">
                  <span class="app-list__title">{{ order.rewardTitle }}</span>
                  <div class="app-list__meta">
                    <span><i class="pi pi-wallet" />{{ formatCoins(order.priceCoins) }}</span>
                    <span><i class="pi pi-clock" />{{ formatDateTime(order.createdAt) }}</span>
                  </div>
                </div>
              </div>
              <div class="app-list__actions">
                <AppStatus :label="REWARD_ORDER_STATUS_LABELS[order.status]" :severity="REWARD_ORDER_STATUS_SEVERITY[order.status]" />
              </div>
            </div>
          </div>
        </template>
      </div>
    </template>

    <ParticipationHistoryDialog
      v-model:visible="participationHistory.visible.value"
      :loading="participationHistory.loading.value"
      :changes="participationHistory.changes.value" />
  </AppPage>
</template>
