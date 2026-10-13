<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
const props = defineProps<{ error: any }>();

const statusCode = computed<number>(() => props.error?.statusCode || 500);

const message = computed(() => {
  switch (statusCode.value) {
    case 500:
      return { title: 'Ошибка на сервере', description: 'Произошла ошибка на сервере, попробуйте позже' };
    case 403:
      return { title: 'Доступ запрещен', description: 'Доступ к запрашиваемой странице запрещен' };
    case 404:
      return { title: 'Страница не найдена', description: 'Запрашиваемая страница не найдена' };
    default:
      return { title: 'Неизвестная ошибка', description: 'Произошла неизвестная ошибка, попробуйте позже' };
  }
});

const textColor = computed(() => (statusCode.value >= 500 ? 'text-red-500' : 'text-primary'));

const authStore = useAuthStore();
</script>

<template>
  <div class="flex min-h-screen items-center justify-center overflow-hidden bg-surface-50 dark:bg-surface-950">
    <div class="flex flex-col items-center justify-center gap-2 p-8 text-center">
      <span class="text-3xl font-bold" :class="textColor">{{ statusCode }}</span>
      <h1 class="mb-2 text-3xl font-bold text-surface-900 dark:text-surface-0 lg:text-5xl">{{ message.title }}</h1>
      <div class="text-surface-600 dark:text-surface-300">{{ message.description }}</div>
      <div class="mt-8">
        <i class="pi pi-fw pi-arrow-left mr-2" :class="textColor" />
        <NuxtLink :to="authStore.isLoggedIn ? ROUTES.cabinet : ROUTES.login" :class="textColor">
          {{ authStore.isLoggedIn ? 'В кабинет' : 'Ко входу' }}
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
