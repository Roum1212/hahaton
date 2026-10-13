<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    classList?: string;
    containerClass?: string;
    bodyClass?: string;
    severity?: 'shadow' | 'outline' | 'default';
    height?: string;
  }>(),
  {
    classList: '',
    containerClass: '',
    bodyClass: '',
    severity: 'default',
    height: undefined,
  },
);

const containerClassList = computed<string>(() => ['flex', 'flex-col', 'gap-2', 'm-4'].join(' '));
const cardClassList = computed<string>(() => {
  const classes = ['flex', 'flex-col', 'rounded-lg', 'bg-surface-0', 'dark:bg-surface-900'];
  if (props.severity === 'shadow') {
    classes.push('shadow-md');
  } else if (props.severity === 'outline') {
    classes.push('border');
  }
  if (props.classList) {
    classes.push(props.classList);
  }

  return classes.join(' ');
});
</script>
<template>
  <div v-if="$slots.title || $slots.content" v-bind="$attrs" :class="cardClassList + ' border-gray-200 dark:!border-surface-700'">
    <div :class="containerClassList + ' ' + containerClass">
      <div v-if="$slots.title" class="font-semibold text-lg">
        <slot name="title" />
      </div>
      <div v-if="$slots.content" class="flex flex-col gap-2" :class="bodyClass">
        <slot name="content" />
      </div>
    </div>
    <div v-if="$slots.footer" class="flex flex-col gap-2 px-4 border-t border-gray-200 dark:!border-surface-700">
      <slot name="footer" />
    </div>
  </div>
</template>
