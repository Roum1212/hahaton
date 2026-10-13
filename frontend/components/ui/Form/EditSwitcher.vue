<script setup lang="ts">
const emits = defineEmits(['update:modelValue']);

const props = withDefaults(
  defineProps<{
    modelValue?: 'edit' | 'view';
    size?: 'small' | 'large';
  }>(),
  {
    modelValue: 'view',
    size: undefined,
  },
);

const { breakpoints } = useApp();
const localValue = ref(props.modelValue);
const title = computed(() => {
  if (breakpoints.md.value) {
    return props.modelValue === 'edit' ? 'В режим просмотра' : 'В режим редактирования';
  } else {
    return props.modelValue === 'edit' ? 'Просмотреть' : 'Редактировать';
  }
});
const textSize = computed(() => {
  switch (props.size) {
    case 'small':
      return 'text-sm';
    case 'large':
      return 'text-lg';
    default:
      return '';
  }
});
</script>
<template>
  <div class="flex items-center gap-2">
    <span :class="textSize">{{ title }}</span>
    <ToggleSwitch v-model="localValue" true-value="edit" false-value="view" @update:model-value="emits('update:modelValue', localValue)" />
  </div>
</template>
