<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string;
    image?: string;
  }>(),
  {
    image: '/images/landing/events-empty.jpg',
  },
);
</script>

<template>
  <div class="app-empty">
    <img :src="image" alt="" class="app-empty__image" />
    <span class="app-empty__title">{{ title }}</span>
    <slot />
  </div>
</template>
