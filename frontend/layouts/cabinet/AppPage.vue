<script setup lang="ts">
withDefaults(
  defineProps<{
    loading?: boolean;
  }>(),
  {
    loading: false,
  },
);
</script>

<template>
  <div class="app-page">
    <div v-if="$slots.title || $slots.subtitle || $slots.actions" class="app-page__header">
      <div class="app-page__heading">
        <Skeleton v-if="loading" width="14rem" height="1.75rem" />
        <template v-else>
          <h1 v-if="$slots.title" class="app-page__title">
            <slot name="title" />
          </h1>
          <p v-if="$slots.subtitle" class="app-page__subtitle">
            <slot name="subtitle" />
          </p>
        </template>
      </div>

      <div v-if="$slots.actions" class="app-page__actions">
        <Skeleton v-if="loading" width="8rem" height="2.25rem" />
        <slot v-else name="actions" />
      </div>
    </div>

    <div class="app-page__body">
      <slot />
    </div>
  </div>
</template>
