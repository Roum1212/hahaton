<script setup lang="ts">
withDefaults(
  defineProps<{
    title?: string;
    description?: string;
  }>(),
  {
    title: 'Здесь пока пусто',
    description: 'Необходимо добавить данные и они появятся здесь',
  },
);
</script>
<template>
  <div class="flex flex-col items-center justify-center gap-2 py-5">
    <div class="text-xl font-semibold">{{ title }}</div>
    <div class="">{{ description }}</div>
  </div>
</template>
