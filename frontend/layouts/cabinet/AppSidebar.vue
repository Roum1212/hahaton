<script setup lang="ts">
import AppNav from './AppNav.vue';
import AppUserMenu from './AppUserMenu.vue';

const { closeMobileSidebar } = useAppLayout();
const siteName = useRuntimeConfig().public.siteName;
</script>

<template>
  <aside class="app-sidebar">
    <button type="button" class="app-sidebar__close" aria-label="Закрыть меню" @click="closeMobileSidebar">
      <i class="pi pi-times" />
    </button>

    <div class="app-panel">
      <div class="app-panel__head">
        <NuxtLink :to="ROUTES.cabinet" class="app-brand" :aria-label="siteName">
          <img src="/images/common/volonterim_logo_transparent.png" :alt="siteName" class="app-brand__mark" />
          <img src="/images/common/wordmark.png" :alt="siteName" class="app-brand__name" />
        </NuxtLink>
      </div>

      <div class="app-panel__body app-scroll">
        <app-nav />
      </div>

      <div class="app-panel__footer">
        <app-user-menu />
      </div>
    </div>
  </aside>
</template>
