<script setup lang="ts">
import AppHeader from './AppHeader.vue';
import AppBottomNav from './AppBottomNav.vue';

const authStore = useAuthStore();
</script>

<template>
  <div class="app-shell">
    <app-header />
    <main class="app-shell__content">
      <div class="app-container">
        <NuxtPage />
      </div>
    </main>
    <app-bottom-nav />
    <GeoSearchAreaDialog />
    <EventCreateDialog v-if="authStore.hasCoordinatorProfile" />
  </div>
</template>
