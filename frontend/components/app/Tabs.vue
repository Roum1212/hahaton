<script setup lang="ts">
import type { IProfileTab } from '~/config/cabinet.config';

defineProps<{
  tabs: IProfileTab[];
  counts?: Record<string, number>;
}>();

const active = defineModel<string>({ required: true });
</script>

<template>
  <div class="app-tabs" role="tablist">
    <button
      v-for="tab in tabs"
      :key="tab.key"
      type="button"
      role="tab"
      class="app-tabs__tab"
      :class="{ 'app-tabs__tab--active': active === tab.key }"
      :aria-selected="active === tab.key"
      @click="active = tab.key">
      {{ tab.label }}<b v-if="counts?.[tab.key]">{{ counts[tab.key] }}</b>
    </button>
  </div>
</template>
