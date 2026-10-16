<script setup lang="ts" generic="T">
defineProps<{
  options: { label: string; value: T }[];
}>();

const model = defineModel<T>({ required: true });
</script>

<template>
  <div class="app-filter">
    <button
      v-for="option in options"
      :key="String(option.value)"
      type="button"
      class="app-filter__item"
      :class="{ 'app-filter__item--active': model === option.value }"
      @click="model = option.value">
      {{ option.label }}
    </button>
  </div>
</template>
