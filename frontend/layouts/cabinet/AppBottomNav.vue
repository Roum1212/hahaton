<script setup lang="ts">
import type { IMenuItem } from '~/config/menu.config';

const route = useRoute();
const { menu, isActiveItem } = useAppMenu();

const items = computed<IMenuItem[]>(() =>
  menu.value.flatMap((section) => section.groups.flatMap((group) => group.items)).filter((item) => item.to && !item.disabled),
);
</script>

<template>
  <nav class="app-bottomnav">
    <NuxtLink
      v-for="item in items"
      :key="item.label"
      :to="item.to!"
      class="app-bottomnav__tab"
      :class="{ 'app-bottomnav__tab--active': isActiveItem(item, route.path) }">
      <i class="app-bottomnav__icon" :class="item.icon" />
      <span class="app-bottomnav__label">{{ item.label }}</span>
    </NuxtLink>
  </nav>
</template>
