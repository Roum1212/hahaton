<script setup lang="ts">
definePageMeta({
  layout: false,
});

const siteName = useRuntimeConfig().public.siteName;
const { configured, initializing, authorizing, init, signIn } = useFlidaAuth();
const demo = useDemoAuth();

useHead({ title: `Вход | ${siteName}` });

onMounted(() => {
  if (configured) init({ silent: true });
});
</script>

<template>
  <div class="app-auth">
    <div class="app-auth__main">
      <NuxtLink :to="ROUTES.home" class="app-auth__brand">
        <img src="/images/common/volonterim_logo_transparent.png" :alt="siteName" />
        <img src="/images/common/wordmark.png" :alt="siteName" />
      </NuxtLink>

      <div class="app-auth__form">
        <h1 class="app-auth__title">Вход</h1>
        <p class="app-auth__note">Один аккаунт для волонтёров и координаторов. Войдите через Флида ID — регистрация не нужна.</p>
        <Message v-if="!configured" severity="warn" size="small">Не задан APP_FLIDA_CLIENT_ID</Message>
        <Button
          label="Войти через Флида ID"
          size="large"
          fluid
          :loading="initializing || authorizing"
          :disabled="!configured"
          @click="signIn" />
        <div v-if="demo.enabled" class="app-auth__demo">
          <span class="app-auth__divider">или демо-вход</span>
          <button
            v-for="account in DEMO_ACCOUNTS"
            :key="account.subject"
            type="button"
            class="app-auth__demo-account"
            @click="demo.signIn(account)">
            <AppAvatar :name="account.name" />
            <span class="app-auth__demo-text">
              <b>{{ account.label }}</b>
              <span>{{ account.name }}</span>
            </span>
            <i class="pi pi-arrow-right" />
          </button>
        </div>
        <NuxtLink :to="ROUTES.home" class="app-link">Вернуться на главную</NuxtLink>
      </div>

      <span class="app-auth__footer">© {{ new Date().getFullYear() }} {{ siteName }}</span>
    </div>

    <div class="app-auth__aside">
      <img src="/images/landing/feature-coins.webp" alt="" />
      <h2>Помогайте городу —<br />копите часы волонтёрства</h2>
      <p class="app-auth__note">1 час волонтёрства = 10 койнов, которые можно обменять на поощрения</p>
    </div>
  </div>
</template>
