<script setup lang="ts">
const emits = defineEmits(['update:modelValue', 'showEditMode']);

const props = withDefaults(
  defineProps<{
    modelValue?: string | number | boolean | undefined;
    overlay?: boolean;
  }>(),
  {
    modelValue: undefined,
    overlay: true,
  },
);

const localValue = ref(props.modelValue);
const isEdit = ref(false);
const popoverRef = ref();

function handleShowEdit(event: Event) {
  localValue.value = props.modelValue;
  emits('showEditMode', true);

  if (props.overlay && popoverRef.value) {
    popoverRef.value.toggle(event);
  } else {
    isEdit.value = true;
  }
}

function handleCancel() {
  isEdit.value = false;
  emits('showEditMode', false);
  if (props.overlay) popoverRef.value?.hide();
}

function handleSave() {
  isEdit.value = false;
  emits('update:modelValue', localValue.value);
  emits('showEditMode', false);
  if (props.overlay) popoverRef.value?.hide();
}

watch(
  () => props.modelValue,
  (value) => {
    localValue.value = value;
  },
);
</script>

<template>
  <div class="flex justify-between gap-2 group relative items-center w-full">
    <template v-if="overlay || !isEdit">
      <slot v-if="$slots.value" name="value" :value="localValue" />
      <span v-else>{{ localValue }}</span>
    </template>
    <Button
      v-tooltip.top="'Редактировать'"
      icon="pi pi-pencil text-xs"
      severity="secondary"
      size="small"
      text
      rounded
      class="!p-0 !w-7 !h-7 opacity-0 group-hover:opacity-100 transition-opacity"
      @click="handleShowEdit" />
    <Popover v-if="overlay" ref="popoverRef" :dismissable="false">
      <div class="w-fit-content min-w-[230px] flex flex-col gap-2">
        <InputGroup class="w-full">
          <slot v-if="$slots.edit" name="edit" :value="localValue" />
          <template v-else>
            <InputText
              v-if="typeof localValue === 'string'"
              v-model="localValue"
              placeholder="Введите значение"
              size="small"
              class="w-full" />
            <InputNumber
              v-else-if="typeof localValue === 'number'"
              v-model="localValue"
              placeholder="Введите значение"
              size="small"
              class="w-full" />
            <InputSwitch v-else-if="typeof localValue === 'boolean'" v-model="localValue" size="small" class="w-full" />
          </template>

          <InputGroupAddon>
            <Button icon="pi pi-times" severity="secondary" size="small" class="h-[35px] !w-full" text @click="handleCancel" />
          </InputGroupAddon>
          <InputGroupAddon>
            <Button icon="pi pi-check" severity="primary" size="small" class="h-[35px] !w-full" @click="handleSave" />
          </InputGroupAddon>
        </InputGroup>
      </div>
    </Popover>
    <template v-else>
      <InputGroup v-if="isEdit" class="w-full">
        <slot v-if="$slots.edit" name="edit" :value="localValue" />
        <template v-else>
          <InputText
            v-if="typeof localValue === 'string'"
            v-model="localValue"
            placeholder="Введите значение"
            size="small"
            class="w-full" />
          <InputNumber
            v-else-if="typeof localValue === 'number'"
            v-model="localValue"
            placeholder="Введите значение"
            size="small"
            class="w-full" />
          <InputSwitch v-else-if="typeof localValue === 'boolean'" v-model="localValue" size="small" class="w-full" />
        </template>
        <InputGroupAddon>
          <Button icon="pi pi-times" severity="secondary" size="small" class="h-[35px]" text @click="handleCancel" />
        </InputGroupAddon>
        <InputGroupAddon>
          <Button icon="pi pi-check" severity="primary" size="small" class="h-[35px]" @click="handleSave" />
        </InputGroupAddon>
      </InputGroup>
    </template>
  </div>
</template>
