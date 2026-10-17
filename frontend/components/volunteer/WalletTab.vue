<script setup lang="ts">
import type { IVolunteerOperation } from '~/core/volunteer/services/volunteer/volunteer.types';
import { VOLUNTEER_OPERATION_TYPE_LABELS } from '~/core/volunteer/services/volunteer/volunteer.types';

const volunteerStore = useVolunteerStore();
const { list, loading, fetching, isLoadAll, load, reset } = usePagination<IVolunteerOperation>({ page: 1, limit: 30 });

const fetchPage = (limit: number, offset: number) => volunteerStore.service.listOperations(undefined, { limit, offset });

onMounted(() => handleAsyncOperation(() => load(fetchPage)));

useLiveRefresh(() => {
  reset();
  return handleAsyncOperation(() => load(fetchPage, true));
});
</script>

<template>
  <div class="app-card">
    <div v-if="loading" class="flex flex-col gap-3">
      <Skeleton v-for="n in 4" :key="n" height="56px" border-radius="16px" />
    </div>

    <AppEmpty v-else-if="!list.length" title="Койны появятся после первого мероприятия" image="/images/landing/step-accrual.jpg" />

    <template v-else>
      <div class="app-list">
        <div v-for="operation in list" :key="operation.id" class="app-list__item app-list__item--row">
          <div class="app-list__lead">
            <span class="app-op-icon" :class="operation.coinDelta >= 0 ? 'app-op-icon--plus' : 'app-op-icon--minus'">
              <i :class="operation.coinDelta >= 0 ? 'pi pi-arrow-down-left' : 'pi pi-arrow-up-right'" />
            </span>
            <div class="app-list__main">
              <span class="app-list__title">{{ operation.description }}</span>
              <div class="app-list__meta">
                <span>{{ VOLUNTEER_OPERATION_TYPE_LABELS[operation.type] }}</span>
                <span>{{ formatDateTime(operation.createdAt) }}</span>
                <span v-if="operation.minuteDelta"><i class="pi pi-clock" />+{{ formatDuration(operation.minuteDelta) }}</span>
              </div>
            </div>
          </div>
          <div class="app-op-amount">
            <span :class="operation.coinDelta >= 0 ? 'app-amount--plus' : 'app-amount--minus'">
              {{ formatSignedCoins(operation.coinDelta) }}
            </span>
            <span class="app-card__hint">{{ formatCoins(operation.coinBalanceAfter) }}</span>
          </div>
        </div>
      </div>
      <div v-if="!isLoadAll" class="app-more">
        <Button label="Показать ещё" severity="secondary" outlined :loading="fetching" @click="load(fetchPage)" />
      </div>
    </template>
  </div>
</template>
