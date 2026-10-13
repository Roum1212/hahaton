<script setup lang="ts">
defineProps<{
  count: number;
}>();

const emit = defineEmits<{
  clear: [];
}>();
</script>

<template>
  <Teleport to="body">
    <Transition name="app-select-panel">
      <div v-if="count" class="app-select-panel">
        <button type="button" class="app-select-panel__close" aria-label="Снять выделение" @click="emit('clear')">
          <i class="pi pi-times" />
        </button>

        <span class="app-select-panel__count">Выбрано: {{ count }}</span>

        <div class="app-select-panel__actions">
          <slot />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
