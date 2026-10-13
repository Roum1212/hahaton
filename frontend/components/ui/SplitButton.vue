<script setup lang="ts">
import type { ButtonProps } from 'primevue/button';
const { breakpoints } = useApp();
type UiButtonProps = ButtonProps & { smLabel?: string };
const props = withDefaults(defineProps<UiButtonProps>(), {
  smLabel: undefined,
});

const buttonProps = filterFieldsNeg(props, 'smLabel', 'label');
</script>
<template>
  <SplitButton v-bind="buttonProps" :label="breakpoints.sm.value ? props.label : props.smLabel" />
</template>
