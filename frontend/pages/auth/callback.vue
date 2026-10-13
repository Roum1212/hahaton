<script setup lang="ts">
definePageMeta({
  layout: false,
});

const { init } = useFlidaAuth();

onMounted(async () => {
  const success = await init({
    onReady: () => window.history.replaceState(window.history.state, '', ROUTES.authCallback),
  });
  if (!success) await navigateTo(ROUTES.login, { replace: true });
});
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center gap-4 bg-surface-0 p-4 dark:bg-surface-950">
    <ProgressSpinner />
    <span class="text-surface-500 dark:text-surface-400">Выполняется вход…</span>
  </div>
</template>
