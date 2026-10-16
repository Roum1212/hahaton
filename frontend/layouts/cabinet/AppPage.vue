<script setup lang="ts">
withDefaults(
  defineProps<{
    loading?: boolean;
    back?: string;
    backLabel?: string;
  }>(),
  {
    loading: false,
    back: undefined,
    backLabel: 'Назад',
  },
);
</script>

<template>
  <div class="app-page">
    <NuxtLink v-if="back" :to="back" class="app-page__back">
      <i class="pi pi-arrow-left" />
      {{ backLabel }}
    </NuxtLink>

    <div v-if="$slots.title || $slots.actions" class="app-page__header">
      <Skeleton v-if="loading" width="16rem" height="2.25rem" />
      <h1 v-else-if="$slots.title" class="app-page__title">
        <slot name="title" />
      </h1>

      <div v-if="$slots.actions && !loading" class="app-page__actions">
        <slot name="actions" />
      </div>
    </div>

    <slot />
  </div>
</template>
