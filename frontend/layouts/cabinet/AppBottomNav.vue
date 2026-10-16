<script setup lang="ts">
import { NAV_BY_MODE } from '~/config/menu.config';

const route = useRoute();
const { mode } = useCabinetMode();
const { isActivePath } = useAppLayout();

const items = computed(() => NAV_BY_MODE[mode.value]);
</script>

<template>
  <nav class="app-bottomnav">
    <NuxtLink
      v-for="item in items"
      :key="item.to"
      :to="item.to"
      class="app-bottomnav__tab"
      :class="{ 'app-bottomnav__tab--active': isActivePath(item.to, route.path, item.exact) }">
      <i class="app-bottomnav__icon" :class="item.icon" />
      <span class="app-bottomnav__label">{{ item.label }}</span>
    </NuxtLink>
  </nav>
</template>
