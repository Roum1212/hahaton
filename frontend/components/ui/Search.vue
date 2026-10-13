<script setup lang="ts">
const emits = defineEmits(['update:modelValue', 'reset']);
const props = withDefaults(
  defineProps<{
    modelValue: string | undefined;
    placeholder?: string;
    size?: 'small' | 'large' | undefined | null;
    iconPos?: 'left' | 'right';
    disabled?: boolean;
  }>(),
  {
    placeholder: 'Поиск...',
    size: 'small',
    iconPos: 'left',
    disabled: false,
  },
);

const { debounce } = useApp();
const localValue = ref<string | undefined>(props.modelValue);

const resetSearch = () => {
  localValue.value = '';
  emits('update:modelValue', localValue.value);
  emits('reset', localValue.value);
};
const debounceUpdate = debounce(() => {
  emits('update:modelValue', localValue.value);
}, 500);
</script>
<template>
  <IconField class="w-full">
    <InputIcon v-if="iconPos === 'left'" class="pi pi-search" />
    <InputText
      v-model="localValue"
      :size="size"
      class="w-full"
      :disabled="disabled"
      :placeholder="placeholder"
      @update:model-value="debounceUpdate" />
    <InputIcon v-if="localValue?.length" class="pi pi-fw pi-times cursor-pointer" @click="resetSearch" />
    <InputIcon v-if="iconPos === 'right'" class="pi pi-search" />
  </IconField>
</template>
