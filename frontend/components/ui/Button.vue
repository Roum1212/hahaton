<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<!-- eslint-disable vue/require-default-prop -->
<script setup lang="ts">
const { breakpoints } = useApp();
interface UiButtonProps {
  style?: any;
  class?: any;
  label?: string | undefined;
  icon?: string | undefined;
  iconPos?: 'left' | 'right' | 'top' | 'bottom';
  iconClass?: string | object | undefined;
  badge?: string | undefined;
  badgeClass?: string | object | undefined;
  badgeSeverity?: 'secondary' | 'info' | 'success' | 'warn' | 'danger' | 'contrast' | null | undefined;
  loading?: boolean | undefined;
  loadingIcon?: string | undefined;
  as?: string | Component | undefined;
  asChild?: boolean | undefined;
  link?: boolean | undefined;
  severity?: 'primary' | 'secondary' | 'success' | 'info' | 'warn' | 'help' | 'danger' | 'contrast' | undefined;
  raised?: boolean | undefined;
  rounded?: boolean | undefined;
  text?: boolean | undefined;
  outlined?: boolean | undefined;
  size?: 'small' | 'large' | undefined;
  variant?: 'outlined' | 'text' | 'link' | undefined;
  plain?: boolean | undefined;
  fluid?: boolean | undefined;
  unstyled?: boolean;
  smLabel?: string;
}
const props = withDefaults(defineProps<UiButtonProps>(), {
  smLabel: undefined,
});

const buttonProps = computed(() => filterFieldsNeg(props, 'smLabel', 'label'));
const classList = 'shrink-0';
</script>
<template>
  <Button
    v-bind="buttonProps"
    :class="classList"
    :label="breakpoints.sm.value ? props.label : props.smLabel === undefined ? props.label : props.smLabel" />
</template>
