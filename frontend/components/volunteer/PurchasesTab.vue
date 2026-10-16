<script setup lang="ts">
import {
  REWARD_ORDER_CANCELLATION_LABELS,
  REWARD_ORDER_STATUS_LABELS,
  REWARD_ORDER_STATUS_SEVERITY,
} from '~/core/reward-order/services/reward-order/reward-order.types';

const { store, loading, fetching, load, cancel, canCancel } = useMyPurchases();
const { name: coordinatorName } = useCoordinatorNames();

onMounted(() => load());
</script>

<template>
  <div class="app-card">
    <div v-if="loading" class="flex flex-col gap-3">
      <Skeleton v-for="n in 3" :key="n" height="72px" border-radius="16px" />
    </div>

    <AppEmpty v-else-if="!store.list.length" title="Вы ещё ничего не обменяли" image="/images/landing/step-reward.jpg">
      <NuxtLink :to="ROUTES.volunteer.rewards">
        <Button label="Выбрать поощрение" icon="pi pi-gift" />
      </NuxtLink>
    </AppEmpty>

    <template v-else>
      <div class="app-list">
        <div v-for="order in store.list" :key="order.id" class="app-list__item">
          <div class="app-list__lead">
            <img :src="pickIllustration(REWARD_ILLUSTRATIONS, order.rewardId)" alt="" class="app-thumb" />
            <div class="app-list__main">
              <span class="app-list__title">{{ order.rewardTitle }}</span>
              <div class="app-list__meta">
                <span><i class="pi pi-wallet" />{{ formatCoins(order.priceCoins) }}</span>
                <span><i class="pi pi-user" />{{ coordinatorName(order.coordinatorId) }}</span>
                <span><i class="pi pi-clock" />{{ formatDateTime(order.createdAt) }}</span>
              </div>
              <span v-if="order.cancellationReason" class="app-list__comment">
                {{ REWARD_ORDER_CANCELLATION_LABELS[order.cancellationReason] }}
              </span>
            </div>
          </div>
          <div class="app-list__actions">
            <AppStatus :label="REWARD_ORDER_STATUS_LABELS[order.status]" :severity="REWARD_ORDER_STATUS_SEVERITY[order.status]" />
            <Button v-if="canCancel(order.status)" label="Отменить" severity="secondary" text size="small" @click="cancel(order.id)" />
          </div>
        </div>
      </div>
      <div v-if="!store.pagination.loadAll" class="app-more">
        <Button label="Показать ещё" severity="secondary" outlined :loading="fetching" @click="load(store.pagination.page)" />
      </div>
    </template>
  </div>
</template>
