<script setup lang="ts">
import {
  REWARD_ORDER_CANCELLATION_LABELS,
  REWARD_ORDER_STATUS_LABELS,
  REWARD_ORDER_STATUS_SEVERITY,
  RewardOrderStatus,
} from '~/core/reward-order/services/reward-order/reward-order.types';

definePageMeta({
  title: 'Заказы',
  permissions: [ROLES.coordinator],
});

const { store, loading, fetching, load, issue, cancel } = useIncomingOrders();
const { name: volunteerName } = useVolunteerNames();

const status = ref<RewardOrderStatus | null>(RewardOrderStatus.PENDING);
const options = [
  ...[RewardOrderStatus.PENDING, RewardOrderStatus.ISSUED, RewardOrderStatus.CANCELLED].map((value) => ({
    value,
    label: REWARD_ORDER_STATUS_LABELS[value],
  })),
  { value: null, label: 'Все' },
];

watch(status, (value) => load(value ?? undefined), { immediate: true });
</script>

<template>
  <AppPage>
    <template #title>Заказы</template>

    <div class="app-card">
      <AppFilter v-model="status" :options="options" />

      <div v-if="loading" class="flex flex-col gap-3">
        <Skeleton v-for="n in 3" :key="n" height="64px" border-radius="16px" />
      </div>

      <AppEmpty v-else-if="!store.list.length" title="Заказов нет" image="/images/landing/step-reward.jpg" />

      <template v-else>
        <div class="app-list">
          <div v-for="order in store.list" :key="order.id" class="app-list__item">
            <div class="app-list__lead">
              <img :src="pickIllustration(REWARD_ILLUSTRATIONS, order.rewardId)" alt="" class="app-thumb" />
              <div class="app-list__main">
                <span class="app-list__title">
                  {{ order.rewardTitle }}
                  <AppStatus :label="REWARD_ORDER_STATUS_LABELS[order.status]" :severity="REWARD_ORDER_STATUS_SEVERITY[order.status]" />
                </span>
                <div class="app-list__meta">
                  <NuxtLink :to="`${ROUTES.coordinator.volunteers}/${order.volunteerId}`" class="app-link"
                    ><i class="pi pi-user" />{{ volunteerName(order.volunteerId) }}</NuxtLink
                  >
                  <span><i class="pi pi-wallet" />{{ formatCoins(order.priceCoins) }}</span>
                  <span><i class="pi pi-clock" />{{ formatDateTime(order.createdAt) }}</span>
                </div>
                <span v-if="order.cancellationReason" class="app-list__comment">{{
                  REWARD_ORDER_CANCELLATION_LABELS[order.cancellationReason]
                }}</span>
              </div>
            </div>
            <div v-if="order.status === RewardOrderStatus.PENDING" class="app-list__actions">
              <Button label="Отменить" severity="secondary" text size="small" @click="cancel(order.id)" />
              <Button label="Выдано" icon="pi pi-check" size="small" @click="issue(order.id)" />
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
    </div>
  </AppPage>
</template>
