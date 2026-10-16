<script setup lang="ts">
import type { IMenuItem } from '~/config/menu.config';

const route = useRoute();
const { menu, isActiveItem } = useAppMenu();

const itemClass = (item: IMenuItem) => ({
  'app-nav__item--active': isActiveItem(item, route.path),
  'app-nav__item--disabled': item.disabled,
});
</script>

<template>
  <nav class="app-nav">
    <div v-for="section in menu" :key="section.key" class="app-nav__group">
      <span class="app-nav__label">{{ section.label }}</span>

      <template v-for="group in section.groups" :key="group.key">
        <template v-for="item in group.items" :key="item.label">
          <span v-if="item.disabled || !item.to" class="app-nav__item" :class="itemClass(item)" aria-disabled="true">
            <i v-if="item.icon" class="app-nav__icon" :class="item.icon" />
            <span class="app-nav__text">{{ item.label }}</span>
            <span v-if="item.badge" class="app-nav__badge">{{ item.badge }}</span>
          </span>

          <NuxtLink v-else :to="item.to" class="app-nav__item" :class="itemClass(item)">
            <i v-if="item.icon" class="app-nav__icon" :class="item.icon" />
            <span class="app-nav__text">{{ item.label }}</span>
            <span v-if="item.badge" class="app-nav__badge">{{ item.badge }}</span>
          </NuxtLink>
        </template>
      </template>
    </div>
  </nav>
</template>
