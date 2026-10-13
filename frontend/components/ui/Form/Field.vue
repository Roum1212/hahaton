<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    id?: string | number;
    label?: string | number;
    relation?: 'horizontal' | 'vertical';
    labelRelation?: 'right' | 'left' | undefined;
    classList?: string;
  }>(),
  {
    id: '',
    label: '',
    relation: 'vertical',
    labelRelation: undefined,
    classList: '',
  },
);

const localClassList = computed(() => {
  const classes = [props.classList];
  if (props.relation === 'horizontal') {
    classes.push(...['flex', 'gap-2', 'mb-2', 'items-center']);
  } else {
    classes.push('field');
  }
  if (props.labelRelation && props.relation === 'horizontal') {
    classes.push('items-center');
  }
  if (props.labelRelation === 'right') {
    classes.push('flex-row-reverse justify-end');
  }
  return classes;
});
</script>
<template>
  <div :class="localClassList">
    <slot v-if="$slots.label" name="label" :label="label" />
    <label v-else-if="label" :for="`${id}`" class="font-medium">{{ label }}</label>
    <slot />
  </div>
</template>
