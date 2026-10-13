<script setup lang="ts">
definePageMeta({
  layout: false,
});
</script>

<template>
  <div class="flex min-h-screen items-center justify-center overflow-hidden bg-surface-50 dark:bg-surface-950">
    <div class="flex flex-col items-center justify-center gap-2 p-8 text-center">
      <span class="text-3xl font-bold text-primary">403</span>
      <h1 class="mb-2 text-3xl font-bold text-surface-900 dark:text-surface-0 lg:text-5xl">Доступ запрещен</h1>
      <div class="text-surface-600 dark:text-surface-300">У вас недостаточно прав для доступа. Обратитесь к администратору.</div>
      <div class="mt-8">
        <i class="pi pi-fw pi-arrow-left mr-2 text-primary" />
        <NuxtLink :to="ROUTES.cabinet" class="text-primary">В кабинет</NuxtLink>
      </div>
    </div>
  </div>
</template>
