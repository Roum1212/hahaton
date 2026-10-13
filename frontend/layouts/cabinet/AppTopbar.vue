<script setup lang="ts">
const { layoutState, isDarkTheme, toggleDarkMode, toggleSidebar, toggleMobileSidebar } = useAppLayout();
const { sectionLabel, pageLabel } = useAppMenu();

const collapseIcon = computed(() => (layoutState.sidebarCollapsed ? 'pi pi-angle-double-right' : 'pi pi-angle-double-left'));
</script>

<template>
  <header class="app-topbar">
    <button type="button" class="app-topbar__button app-topbar__button--mobile" aria-label="Меню" @click="toggleMobileSidebar">
      <i class="pi pi-bars" />
    </button>

    <button
      type="button"
      class="app-topbar__button app-topbar__button--desktop"
      :aria-label="layoutState.sidebarCollapsed ? 'Развернуть меню' : 'Свернуть меню'"
      @click="toggleSidebar">
      <i :class="collapseIcon" />
    </button>

    <div class="app-topbar__crumbs">
      <span class="app-topbar__crumb">{{ sectionLabel }}</span>
      <template v-if="pageLabel">
        <i class="app-topbar__sep pi pi-angle-right" />
        <span class="app-topbar__crumb app-topbar__crumb--current">{{ pageLabel }}</span>
      </template>
    </div>

    <div class="app-topbar__actions">
      <button type="button" class="app-topbar__button" :aria-label="isDarkTheme ? 'Светлая тема' : 'Тёмная тема'" @click="toggleDarkMode">
        <i :class="`pi ${isDarkTheme ? 'pi-sun' : 'pi-moon'}`" />
      </button>
    </div>
  </header>
</template>
