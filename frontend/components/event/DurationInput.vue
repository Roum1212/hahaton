<script setup lang="ts">
import type { ComponentPublicInstance } from 'vue';

const minutes = defineModel<number>({ required: true });

defineProps<{
  disabled?: boolean;
}>();

const MAX_HOURS = 23;

const minuteOptions = [
  { value: 0, label: '00 мин' },
  { value: MINUTES_STEP, label: '30 мин' },
];

const hours = computed<number>({
  get: () => Math.floor(minutes.value / 60),
  set: (value) => {
    minutes.value = Math.min(Math.max(value ?? 0, 0), MAX_HOURS) * 60 + (minutes.value % 60);
  },
});

const root = ref<ComponentPublicInstance | null>(null);

const onKeydown = (event: KeyboardEvent) => {
  const input = (root.value?.$el as HTMLElement | undefined)?.querySelector('input');
  if (event.target !== input || !input || !/^\d$/.test(event.key)) return;
  if (input.selectionStart !== input.selectionEnd || Number.parseInt(input.value, 10) !== 0) return;
  input.setSelectionRange(0, 1);
};

onMounted(() => (root.value?.$el as HTMLElement | undefined)?.addEventListener('keydown', onKeydown, true));
onBeforeUnmount(() => (root.value?.$el as HTMLElement | undefined)?.removeEventListener('keydown', onKeydown, true));

const rest = computed<number>({
  get: () => (minutes.value % 60 >= MINUTES_STEP ? MINUTES_STEP : 0),
  set: (value) => {
    minutes.value = hours.value * 60 + value;
  },
});
</script>

<template>
  <InputGroup ref="root">
    <InputNumber v-model="hours" :min="0" :max="MAX_HOURS" show-buttons suffix=" ч" :disabled="disabled" :allow-empty="false" fluid />
    <Select
      v-model="rest"
      :options="minuteOptions"
      option-label="label"
      option-value="value"
      :disabled="disabled"
      class="app-duration__minutes" />
  </InputGroup>
</template>
