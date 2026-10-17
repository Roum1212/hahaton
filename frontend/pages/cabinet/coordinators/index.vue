<script setup lang="ts">
import {
  COOPERATION_STATUS_LABELS,
  COOPERATION_STATUS_SEVERITY,
  CooperationStatus,
  type ICooperation,
} from '~/core/cooperation/services/cooperation/cooperation.types';

definePageMeta({
  title: 'Координаторы',
});

const authStore = useAuthStore();
const coordinatorsStore = useCoordinatorsStore();
const { store: cooperationsStore, loading: cooperationsLoading, load: loadCooperations, request, terminate } = useVolunteerCooperations();

const loading = ref(true);

const coordinators = computed(() => coordinatorsStore.list.filter((item) => item.id !== authStore.session?.userId));

const cooperationOf = (coordinatorId: string): ICooperation | undefined => cooperationsStore.withCoordinator(coordinatorId);

onMounted(async () => {
  loading.value = true;
  await handleAsyncOperation(
    async () => {
      await coordinatorsStore.find(1);
      while (!coordinatorsStore.pagination.loadAll) await coordinatorsStore.find(coordinatorsStore.pagination.page);
      await loadCooperations();
    },
    () => {
      loading.value = false;
    },
  );
});
</script>

<template>
  <AppPage>
    <template #title>Координаторы</template>

    <div v-if="loading || cooperationsLoading" class="app-grid">
      <Skeleton v-for="n in 4" :key="n" height="180px" border-radius="20px" />
    </div>

    <div v-else-if="!coordinators.length" class="app-card">
      <AppEmpty title="Координаторов пока нет" image="/images/landing/step-request.jpg" />
    </div>

    <div v-else class="app-grid">
      <article v-for="coordinator in coordinators" :key="coordinator.id" class="app-card app-person">
        <div class="app-person__head">
          <NuxtLink :to="`${ROUTES.volunteer.coordinators}/${coordinator.id}`">
            <AppAvatar :name="coordinator.displayName" :image="coordinator.avatarUrl || null" />
          </NuxtLink>
          <div class="app-list__main">
            <NuxtLink :to="`${ROUTES.volunteer.coordinators}/${coordinator.id}`" class="app-list__title">
              {{ coordinator.displayName }}
            </NuxtLink>
            <AppStatus
              v-if="cooperationOf(coordinator.id)"
              :label="COOPERATION_STATUS_LABELS[cooperationOf(coordinator.id)!.status]"
              :severity="COOPERATION_STATUS_SEVERITY[cooperationOf(coordinator.id)!.status]" />
          </div>
        </div>
        <p v-if="coordinator.description" class="app-person__text">{{ coordinator.description }}</p>
        <span
          v-if="
            cooperationOf(coordinator.id)?.status === CooperationStatus.REJECTED && !canRequestCooperation(cooperationOf(coordinator.id))
          "
          class="app-card__hint">
          Повторная заявка с {{ formatDateTime(cooperationOf(coordinator.id)?.retryAfter) }}
        </span>
        <div class="app-person__actions">
          <Button
            v-if="canRequestCooperation(cooperationOf(coordinator.id))"
            label="Предложить сотрудничество"
            size="small"
            @click="request(coordinator.id)" />
          <template v-if="cooperationOf(coordinator.id)?.status === CooperationStatus.ACTIVE">
            <NuxtLink :to="`${ROUTES.volunteer.coordinators}/${coordinator.id}`">
              <Button label="Профиль" severity="secondary" outlined size="small" />
            </NuxtLink>
            <Button
              label="Прекратить сотрудничество"
              severity="secondary"
              text
              size="small"
              @click="terminate(coordinator.id, coordinator.displayName)" />
          </template>
        </div>
      </article>
    </div>
  </AppPage>
</template>
