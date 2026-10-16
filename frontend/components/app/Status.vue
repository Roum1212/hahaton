<script setup lang="ts">
const props = defineProps<{
  label: string;
  severity?: 'success' | 'info' | 'warn' | 'danger' | 'secondary';
}>();

const MODIFIERS = {
  success: 'app-chip--success',
  info: 'app-chip--accent',
  warn: 'app-chip--warn',
  danger: 'app-chip--danger',
  secondary: '',
};

const modifier = computed(() => MODIFIERS[props.severity ?? 'secondary']);
</script>

<template>
  <span class="app-chip" :class="modifier">{{ label }}</span>
</template>
