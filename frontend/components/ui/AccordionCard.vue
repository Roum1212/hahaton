<script lang="ts" setup>
defineExpose({ updateContentHeight });

const emits = defineEmits(['update:expanded']);
const props = withDefaults(
  defineProps<{
    severity?: 'shadow' | 'outline' | 'default';
    classList?: string;
    height?: string;
    expanded?: boolean;
  }>(),
  {
    expanded: false,
    severity: 'default',
    classList: '',
    height: undefined,
  },
);

const containerClassList = computed<string>(() => ['flex', 'flex-col', 'm-4'].join(' '));
const cardClassList = computed<string>(() => {
  const classes = ['flex', 'flex-col', 'rounded-lg', 'bg-surface-0', 'dark:bg-surface-900'];
  if (props.severity === 'shadow') {
    classes.push('shadow-md');
  } else if (props.severity === 'outline') {
    classes.push('border border-gray-200 dark:border-surface-700');
  }
  if (props.classList) {
    classes.push(props.classList);
  }
  return classes.join(' ');
});

const contentRef = ref<HTMLElement | null>();
const iconClass = computed(() => (isExpanded.value ? 'pi pi-chevron-up' : 'pi pi-chevron-down'));
const isExpanded = ref(props.expanded);
const contentHeight = ref(0);

function updateContentHeight() {
  if (isExpanded.value && contentRef.value) {
    contentHeight.value = contentRef.value.scrollHeight;
  } else {
    contentHeight.value = 0;
  }
}

const toggleCard = () => {
  isExpanded.value = !isExpanded.value;
  emits('update:expanded', isExpanded.value);
};

onMounted(() => {
  if (isExpanded.value) {
    updateContentHeight();
  }
});

useResizeObserver(contentRef, () => {
  if (isExpanded.value && contentRef.value) {
    contentHeight.value = contentRef.value.scrollHeight;
  }
});

watch(isExpanded, () => {
  updateContentHeight();
});
</script>
<template>
  <div v-if="$slots.title || $slots.content" v-bind="$attrs" :class="cardClassList">
    <div :class="containerClassList">
      <div class="flex justify-between items-center gap-1">
        <div v-if="$slots.title" class="font-semibold text-lg w-full">
          <slot name="title" />
        </div>
        <Button severity="secondary" text rounded size="small" class="shrink-0 p-0 !h-6 !w-6" :icon="iconClass" @click.stop="toggleCard" />
      </div>
      <div
        v-if="$slots.content"
        ref="contentRef"
        class="flex flex-col gap-2 overflow-hidden transition-all duration-300 ease-in-out"
        :style="{ 'max-height': `${contentHeight}px` }">
        <slot name="content" />
      </div>
    </div>
    <div v-if="$slots.footer" class="flex flex-col gap-2 px-4 border-t border-gray-200 dark:border-surface-700">
      <slot name="footer" />
    </div>
  </div>
</template>
