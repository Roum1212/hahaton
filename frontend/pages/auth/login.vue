<script setup lang="ts">
definePageMeta({
  layout: false,
});

const siteName = useRuntimeConfig().public.siteName;
const { configured, initializing, authorizing, init, signIn } = useFlidaAuth();

useHead({ title: `Вход | ${siteName}` });

onMounted(() => {
  if (configured) init({ silent: true });
});
</script>

<template>
  <div class="box-border flex h-[100vh] w-full flex-col items-center justify-between bg-surface-0 p-[2rem] dark:bg-surface-950">
    <NuxtLink to="/" class="flex items-center gap-3">
      <img src="/images/common/volonterim_logo_transparent.png" :alt="siteName" class="h-9 w-9" />
      <img src="/images/common/wordmark.png" :alt="siteName" class="h-6 w-auto" />
    </NuxtLink>

    <div class="flex w-full max-w-[420px] flex-col gap-6">
      <div class="flex flex-col gap-2">
        <div class="text-[32px] font-bold text-surface-900 dark:text-surface-0">Вход в систему</div>
        <span class="text-surface-500 dark:text-surface-400">Войдите с помощью аккаунта Флида ID</span>
      </div>
      <Message v-if="!configured" severity="warn" size="small">Не задан APP_FLIDA_CLIENT_ID</Message>
      <Button
        label="ВОЙТИ ЧЕРЕЗ ФЛИДА ID"
        icon="pi pi-sign-in"
        class="w-full !font-semibold !tracking-[0.08em]"
        :loading="initializing || authorizing"
        :disabled="!configured"
        @click="signIn" />
    </div>

    <div class="text-center text-sm text-surface-500 dark:text-surface-400">© {{ new Date().getFullYear() }} {{ siteName }}</div>
  </div>
</template>
