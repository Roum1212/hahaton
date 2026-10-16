<script setup lang="ts">
import AppSidebar from './AppSidebar.vue';
import AppTopbar from './AppTopbar.vue';
import AppBottomNav from './AppBottomNav.vue';

const route = useRoute();
const { layoutState, closeMobileSidebar, selectSection } = useAppLayout();

const shellClass = computed(() => ({
  'app-shell--collapsed': layoutState.sidebarCollapsed,
  'app-shell--mobile-open': layoutState.mobileSidebarOpen,
}));

watch(
  () => route.fullPath,
  () => {
    closeMobileSidebar();
    selectSection(null);
  },
);

watch(
  () => layoutState.mobileSidebarOpen,
  (open) => document.body.classList.toggle('app-scroll-lock', open),
);

onUnmounted(() => document.body.classList.remove('app-scroll-lock'));
</script>

<template>
  <div class="app-shell" :class="shellClass">
    <app-sidebar />
    <div class="app-shell__backdrop" @click="closeMobileSidebar" />

    <div class="app-shell__main">
      <app-topbar />
      <main class="app-shell__content">
        <div class="app-shell__inner">
          <NuxtPage />
        </div>
      </main>
    </div>

    <app-bottom-nav />
  </div>
</template>
