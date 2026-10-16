<script setup lang="ts">
import { NAV_BY_MODE } from '~/config/menu.config';
import AppUserMenu from './AppUserMenu.vue';

const route = useRoute();
const siteName = useRuntimeConfig().public.siteName;
const { mode } = useCabinetMode();
const { isActivePath } = useAppLayout();
const { label, openDialog } = useSearchArea();

const items = computed(() => NAV_BY_MODE[mode.value]);
</script>

<template>
  <header class="app-header">
    <div class="app-container app-header__inner">
      <NuxtLink :to="ROUTES.home" class="app-header__brand">
        <img src="/images/common/volonterim_logo_transparent.png" :alt="siteName" class="app-header__logo" />
        <img src="/images/common/wordmark.png" :alt="siteName" class="app-header__wordmark" />
      </NuxtLink>

      <button v-if="mode === 'volunteer'" type="button" class="app-header__city" @click="openDialog">
        <i class="pi pi-map-marker" />
        <span class="app-header__city-label">{{ label }}</span>
      </button>

      <nav class="app-header__nav">
        <NuxtLink
          v-for="item in items"
          :key="item.to"
          :to="item.to"
          class="app-header__link"
          :class="{ 'app-header__link--active': isActivePath(item.to, route.path, item.exact) }">
          {{ item.label }}
        </NuxtLink>
      </nav>

      <div class="app-header__actions">
        <app-user-menu />
      </div>
    </div>
  </header>
</template>
