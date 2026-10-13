<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
export interface IComponentClasses {
  container: string;
  label: string;
  input: string;
}

const emits = defineEmits(['update:modelValue', 'updated']);
const props = withDefaults(
  defineProps<{
    modelValue?: (string | number | bigint)[] | undefined;
    label?: string | null;
    placeholder?: string;
    disabled?: boolean;
    type?: 'string' | 'number';
    factoryDigits?: number;
    classList?: IComponentClasses;
  }>(),
  {
    modelValue: () => [],
    label: null,
    placeholder: 'Введите несколько значений',
    disabled: false,
    type: 'string',
    factoryDigits: 0,
    classList: () => ({
      container: '',
      label: '',
      input: '',
    }),
  },
);
const { debounce } = useApp();
const inputRef = ref();
const containerRef = ref();
const inputValue = ref<string | null>(null);
const isFocus = ref(false);
const internalValue = ref<(string | number | bigint)[]>([...props.modelValue]);
const internalValueChips = computed<string[]>(() => internalValue.value.map((el: string | number | bigint) => (el ? el.toString() : '')));

const isShowPlaceholder = computed(() => !isFocus.value && !internalValueChips.value.length);
const onSetValue = () => {
  if (!inputValue.value) {
    isFocus.value = false;
    return;
  }

  const newValue = checkInputValue(inputValue.value);
  if (newValue !== null) {
    internalValue.value = [...internalValue.value, newValue];

    const emitValue =
      props.type === 'number'
        ? internalValue.value.map((el) => (typeof el === 'number' ? el : parseFloat(el as string)))
        : internalValue.value;

    emits('update:modelValue', emitValue);
  }

  inputValue.value = null;
  isFocus.value = false;
};
const onRemoveValue = (index: number) => {
  internalValue.value = internalValue.value.filter((_, i) => i !== index);
  emits('updated', true);
  emits('update:modelValue', internalValue.value);
};
const onInputFocus = () => {
  if (inputRef.value) {
    inputRef.value.focus();
  }
  isFocus.value = true;
};

const onInputUnfocus = () => {
  isFocus.value = false;
};

const handleClickOutside = (event: Event) => {
  if (containerRef.value && !containerRef.value.contains(event.target as Node)) {
    onInputUnfocus();
  }
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside, { passive: true });
});

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside);
});

const inputWidth = computed<number>(() => {
  const length = inputValue.value?.length ?? 0;
  return Math.max(1, length * 0.75);
});
const filterInputValue = (event: KeyboardEvent) => {
  if (event.key === 'Enter') {
    event.preventDefault();
    onSetValue();
    isFocus.value = false;
    return;
  }
  if (props.type === 'number') {
    if (
      !/^[0-9.,]$/.test(event.key) &&
      !['Backspace', 'Delete', 'ArrowLeft', 'ArrowRight', 'Meta', 'Control', 'v', 'м'].includes(event.key)
    ) {
      event.preventDefault();
    }
    if (event.key === ',' && inputValue.value?.includes(',')) {
      event.preventDefault();
    }
  }
};
const handlePaste = (event: ClipboardEvent) => {
  event.preventDefault();
  const pastedText = event.clipboardData?.getData('text');

  if (!pastedText) return;

  const values = pastedText
    .split(',')
    .map((v) => v.trim())
    .filter((v) => v.length > 0);

  if (values.length === 0) return;

  const validValues: (string | number)[] = [];

  for (const value of values) {
    const newValue = checkInputValue(value);
    if (newValue !== null) {
      validValues.push(newValue);
    }
  }

  if (validValues.length > 0) {
    internalValue.value = [...internalValue.value, ...validValues];
    emits('update:modelValue', internalValue.value);
  }

  inputValue.value = null;
  isFocus.value = false;
};

const onUpdate = debounce(() => {
  emits('updated', true);
}, 100);
const handleClear = () => {
  inputValue.value = null;
  internalValue.value = [];
  emits('update:modelValue', internalValue.value);
};

const checkInputValue = (value: string | undefined) => {
  if (!value) return null;
  if (props.type === 'number') {
    const num = Number(value);
    if (isNaN(num)) {
      return null;
    }
    return num;
  }
  return value;
};

watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue.length !== internalValue.value.length || !newValue.every((val, index) => val === internalValue.value[index])) {
      internalValue.value = [...newValue];
    }
  },
  { flush: 'post' },
);
</script>
<template>
  <div class="flex flex-col items-end">
    <div
      ref="containerRef"
      class="multyinput gap-2 cursor-pointer w-full overflow-y-auto relative"
      :class="{ disabled: disabled }"
      @click.stop="onInputFocus">
      <span v-if="isShowPlaceholder" class="dark:text-gray-400 text-gray-500">
        {{ placeholder }}
      </span>
      <Chip
        v-for="(label, index) in internalValueChips"
        :key="`${label}-${index}`"
        :label="label"
        removable
        class="!py-1 text-xs items-center break-all"
        @remove="() => onRemoveValue(index)" />
      <UiForm class="flex items-center justify-start gap-2" @submit="onSetValue">
        <input
          ref="inputRef"
          v-model="inputValue"
          unstyled
          :disabled="disabled"
          class="unsyled-input string"
          :class="classList.input"
          :style="{
            'max-width': '100%',
            'min-width': '10px',
            width: `${inputWidth}rem`,
          }"
          @keydown="filterInputValue"
          @keyup="onUpdate"
          @focus="isFocus = true"
          @blur="onSetValue"
          @paste="handlePaste" />
      </UiForm>
      <Button
        v-if="internalValueChips.length > 3"
        text
        icon="pi pi-times"
        size="small"
        severity="secondary"
        class="!py-0 !px-0 !h-5 !w-5 text-xs items-center break-all !absolute right-2 top-[8px]"
        @click="handleClear" />
    </div>
  </div>
</template>
<style lang="scss">
.multyinput.disabled {
  opacity: 1;
  background: var(--p-inputtext-disabled-background);
  color: var(--p-inputtext-disabled-color);
}
.multyinput {
  display: flex;
  flex-wrap: wrap;

  border-start-end-radius: var(--p-inputgroup-addon-border-radius);
  border-end-end-radius: var(--p-inputgroup-addon-border-radius);
  font-size: var(--p-inputtext-sm-font-size);
  padding-block: var(--p-inputtext-sm-padding-y);
  padding-inline: var(--p-inputtext-sm-padding-x);
  font-family: inherit;
  font-feature-settings: inherit;
  font-size: 1rem;
  color: var(--p-inputtext-color);
  background: var(--p-inputtext-background);
  padding-block: var(--p-inputtext-padding-y);
  padding-inline: var(--p-inputtext-padding-x);
  border: 1px solid var(--p-inputtext-border-color);
  transition:
    background var(--p-inputtext-transition-duration),
    color var(--p-inputtext-transition-duration),
    border-color var(--p-inputtext-transition-duration),
    outline-color var(--p-inputtext-transition-duration),
    box-shadow var(--p-inputtext-transition-duration);
  appearance: none;
  border-radius: var(--p-inputtext-border-radius);
  outline-color: transparent;
  box-shadow: var(--p-inputtext-shadow);

  .unsyled-input.string {
    width: 20px;
    background: transparent;
    &:focus {
      border: none;
      outline: none;
    }
  }
  .unsyled-input.number {
    background: transparent;
    &:focus {
      border: none;
      outline: none;
    }
  }
}
.p-inputgroup > .multyinput:not(:last-child) {
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
}
.p-chip-remove-icon {
  flex-shrink: 0;
}
</style>
