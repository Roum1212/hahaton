<script setup lang="ts">
import type { IMenuSection } from '~/config/menu.config';
import AppNav from './AppNav.vue';
import AppUserMenu from './AppUserMenu.vue';

const { layoutState, closeMobileSidebar, selectSection, setSidebarCollapsed } = useAppLayout();
const { menu, currentSection } = useAppMenu();
const siteName = useRuntimeConfig().public.siteName;

const onSectionClick = (section: IMenuSection) => {
  selectSection(section.key);
  if (layoutState.sidebarCollapsed) setSidebarCollapsed(false);
};
</script>

<template>
  <aside class="app-sidebar">
    <button type="button" class="app-sidebar__close" aria-label="Закрыть меню" @click="closeMobileSidebar">
      <i class="pi pi-times" />
    </button>

    <div class="app-rail">
      <NuxtLink :to="ROUTES.cabinet" class="app-rail__brand" :aria-label="siteName">
        <span class="app-rail__mark">{{ siteName.slice(0, 1) }}</span>
      </NuxtLink>

      <div class="app-rail__items">
        <button
          v-for="section in menu"
          :key="section.key"
          v-tooltip.right="section.label"
          type="button"
          class="app-rail__item"
          :class="{ 'app-rail__item--active': currentSection?.key === section.key }"
          :aria-label="section.label"
          @click="onSectionClick(section)">
          <component :is="section.lucide" :size="18" :stroke-width="1.75" />
        </button>
      </div>

      <div v-if="layoutState.sidebarCollapsed" class="app-rail__footer">
        <app-user-menu compact />
      </div>
    </div>

    <div class="app-panel">
      <div class="app-panel__head">
        <span class="app-panel__title">{{ siteName }}</span>
        <span class="app-panel__subtitle">{{ currentSection?.label }}</span>
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
