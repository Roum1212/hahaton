<script setup lang="ts">
import type { ICoordinator } from '~/core/coordinator/services/coordinator/coordinator.types';

definePageMeta({
  title: 'Приглашение',
});

const route = useRoute();
const cooperationsStore = useCooperationsStore();
const coordinatorsStore = useCoordinatorsStore();
const { toast } = useApp();

const token = computed(() => String(route.params.token ?? ''));
const accepting = ref(false);
const coordinator = ref<ICoordinator | null>(null);

const onAccept = async () => {
  accepting.value = true;
  const cooperation = await handleAsyncOperation(
    () => cooperationsStore.acceptInvitation(token.value),
    () => {
      accepting.value = false;
    },
  );
  if (!cooperation) return;
  coordinator.value = (await handleAsyncOperation(() => coordinatorsStore.fetchOne(cooperation.coordinatorId))) ?? null;
  toast.success('Сотрудничество установлено');
};
</script>

<template>
  <AppPage>
    <template #title>Приглашение к сотрудничеству</template>

    <div class="app-card app-invite">
      <img src="/images/landing/step-handover.jpg" alt="" class="app-promo__image" />
      <template v-if="coordinator">
        <div class="flex items-center gap-3">
          <AppAvatar :name="coordinator.displayName" :image="coordinator.avatarUrl || null" />
          <div class="flex flex-col">
            <span class="app-card__title">{{ coordinator.displayName }}</span>
            <span v-if="coordinator.description" class="app-card__hint">{{ coordinator.description }}</span>
          </div>
        </div>
        <div class="flex flex-wrap gap-2">
          <NuxtLink :to="ROUTES.volunteer.events">
            <Button label="Мероприятия" />
          </NuxtLink>
          <NuxtLink :to="ROUTES.volunteer.rewards">
            <Button label="Поощрения" severity="secondary" outlined />
          </NuxtLink>
        </div>
      </template>

      <template v-else>
        <span class="app-card__title">Координатор приглашает вас к сотрудничеству</span>
        <Button label="Принять приглашение" size="large" :loading="accepting" :disabled="!token" @click="onAccept" />
      </template>
    </div>
  </AppPage>
</template>
