<script setup lang="ts">
const emits = defineEmits(['submit']);
const formEl = ref<HTMLFormElement | null>(null);
defineExpose({
  el: formEl,
});
const props = withDefaults(
  defineProps<{
    classList?: string;
    loading?: boolean;
    submitted?: boolean;
    severity?: 'skeleton' | 'default';
    placeholderRows?: number;
    placeholderCount?: number;
    fakeSubmit?: boolean;
  }>(),
  {
    classList: '',
    loading: false,
    submitted: false,
    severity: 'default',
    placeholderRows: 1,
    placeholderCount: 3,
    fakeSubmit: false,
  },
);

const isLoading = computed(() => props.loading);
const isSubmitted = ref(props.submitted);
const handleSubmit = async () => {
  try {
    emits('submit');
    isSubmitted.value = true;
  } catch (error) {
    console.error('Form submission error:', error);
  }
};
</script>
<template>
  <form ref="formEl" class="p-fluid max-w-full" :class="classList" @submit.prevent="handleSubmit">
    <template v-if="isLoading">
      <UiFormPlaceholder v-if="props.severity === 'skeleton'" :rows="placeholderRows" :count="placeholderCount" />
      <slot v-if="$slots.loading" name="loading" />
    </template>
    <template v-if="!isLoading">
      <slot :loading="isLoading" :submitted="isSubmitted" />
      <button v-if="fakeSubmit" submit class="hidden" />
    </template>
  </form>
</template>
