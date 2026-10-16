<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    name: string;
    image?: string | null;
    size?: 'sm' | 'md' | 'lg';
  }>(),
  {
    image: null,
    size: 'md',
  },
);

const letter = computed(() => props.name.trim().slice(0, 1).toLocaleUpperCase() || '?');
</script>

<template>
  <span class="app-avatar" :class="size !== 'md' && `app-avatar--${size}`">
    <img v-if="image" :src="image" :alt="name" />
    <template v-else>{{ letter }}</template>
  </span>
</template>
