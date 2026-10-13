<script lang="ts" setup>
withDefaults(
  defineProps<{
    loading?: boolean;
    loaded?: number;
    count?: number;
  }>(),
  {
    loading: true,
    count: 0,
    loaded: 0,
  },
);
</script>
<template>
  <div v-if="loading" class="text-surface-500 dark:text-gray-300 text-sm mb-2">
    <Skeleton width="13rem" />
  </div>
  <div v-else class="text-surface-600 dark:text-surface-50 text-sm mb-2">Найдено результатов: {{ count }}</div>
</template>
