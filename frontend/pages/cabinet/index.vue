<script setup lang="ts">
definePageMeta({
  title: 'Дашборд',
});

const authStore = useAuthStore();
const usersStore = useUsersStore();

onMounted(() => {
  if (authStore.isAdmin && !usersStore.list.length) handleAsyncOperation(() => usersStore.find(1));
});
</script>

<template>
  <AppPage>
    <template #title>Добро пожаловать, {{ authStore.user?.firstName || authStore.user?.email }}</template>
    <template #subtitle>{{ authStore.user?.email }}</template>

    <div class="app-grid">
      <NuxtLink :to="`${ROUTES.cabinet}/${ROUTES.settings}/${ROUTES.profile}`" class="app-card app-card--action">
        <div class="flex items-start justify-between gap-3">
          <div class="flex flex-col gap-1">
            <span class="app-card__title">Профиль</span>
            <span class="app-card__hint">Имя, фамилия и email вашей учётной записи</span>
          </div>
          <i class="pi pi-user text-xl text-primary" />
        </div>
      </NuxtLink>

      <NuxtLink v-if="authStore.isAdmin" :to="`${ROUTES.cabinet}/${ROUTES.settings}/${ROUTES.users}`" class="app-card app-card--action">
        <div class="flex items-start justify-between gap-3">
          <div class="flex flex-col gap-1">
            <span class="app-card__title">Пользователи</span>
            <span class="app-card__hint">Учётные записи и роли</span>
          </div>
          <span class="app-metric">{{ usersStore.total }}</span>
        </div>
      </NuxtLink>
    </div>
  </AppPage>
</template>
