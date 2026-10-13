<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts" generic="T extends Object extends { id?: string | number }">
interface Props<T> {
  items: T[];
  loading?: boolean;
  height?: string;
  allowInnerScroll?: boolean;
}
defineExpose({
  scrollToPos,
  scrollToTop,
});
const emit = defineEmits(['pre-scroll-bottom', 'scroll-bottom', 'scroll']);

withDefaults(defineProps<Props<T>>(), {
  loading: false,
  height: '700px',
  allowInnerScroll: true,
});
const container = ref<HTMLElement | null>(null);
const scrollTop = ref(0);

const isDraggingH = ref(false);
const dragStartX = ref(0);
const scrollStartX = ref(0);
const dragThreshold = 5;

const getItemKey = (item: T, index: number) => {
  return (item as any)?.id ?? `item-${index}`;
};
const handleScroll = () => {
  const el = container.value || document.documentElement;
  scrollTop.value = el.scrollTop;

  const scrollBottom = el.scrollHeight - (el.scrollTop + el.clientHeight);
  if (scrollBottom <= 10) emit('scroll-bottom');
  if (scrollBottom <= 150) emit('pre-scroll-bottom');
  emit('scroll', scrollTop.value);
};

let raf: number | null = null;
function onMouseMoveH(e: MouseEvent) {
  const dx = e.clientX - dragStartX.value;

  if (!isDraggingH.value && Math.abs(dx) > dragThreshold) {
    isDraggingH.value = true;
    container.value!.style.cursor = 'grabbing';
    container.value!.style.userSelect = 'none';
  }

  if (isDraggingH.value) {
    if (raf) cancelAnimationFrame(raf);
    raf = requestAnimationFrame(() => {
      if (container.value) container.value.scrollLeft = scrollStartX.value - dx;
      raf = null;
    });
  }
}

function onMouseDownH(e: MouseEvent) {
  if (!container.value) return;
  const el = e.target as HTMLElement;
  const isSelect = !el.classList.contains('table-cell') && !el.classList.contains('table-row') && !el.classList.contains('vs-item');
  if (isSelect || e.button !== 0) {
    container.value.style.userSelect = 'text';
    return;
  }
  dragStartX.value = e.clientX;
  scrollStartX.value = container.value!.scrollLeft;

  window.addEventListener('mousemove', onMouseMoveH);
  window.addEventListener('mouseup', onMouseUpH);
}

function onMouseUpH() {
  if (isDraggingH.value) isDraggingH.value = false;

  if (container.value) {
    container.value.style.cursor = 'default';
    container.value.style.removeProperty('user-select');
  }

  window.removeEventListener('mousemove', onMouseMoveH);
  window.removeEventListener('mouseup', onMouseUpH);
}
function scrollToPos(position: number, behavior: ScrollBehavior = 'auto') {
  if (container.value) {
    container.value.scrollTo({
      top: position,
      behavior,
    });
  }
}
function scrollToTop(behavior: ScrollBehavior = 'smooth') {
  scrollToPos(0, behavior);
}

onMounted(() => {
  const el = container.value || window;
  el.addEventListener('scroll', handleScroll, { passive: true });

  if (container.value) {
    container.value.style.cursor = 'default';
    container.value.addEventListener('mousedown', onMouseDownH);
    window.addEventListener('mousemove', onMouseMoveH);
    window.addEventListener('mouseup', onMouseUpH);
  }
});

onBeforeUnmount(() => {
  const el = container.value || window;
  el.removeEventListener('scroll', handleScroll);
  if (!container.value) return;
  container.value.removeEventListener('mousedown', onMouseDownH);
  window.removeEventListener('mousemove', onMouseMoveH);
  window.removeEventListener('mouseup', onMouseUpH);
});
</script>
<template>
  <div
    v-if="height"
    ref="container"
    class="base-table-body"
    :style="{ height, overflowY: allowInnerScroll ? 'auto' : 'hidden', position: 'relative' }">
    <div>
      <div v-for="(item, index) in items" :key="getItemKey(item, index)" class="vs-item">
        <slot :item="item" :index="index" />
      </div>
      <div v-if="loading" class="vs-loading">
        <slot name="loading">
          <div class="default-loader">Загрузка...</div>
        </slot>
      </div>
    </div>
  </div>
</template>
<style scoped>
.base-table-body {
  will-change: scroll-position;
  width: 100%;
  scrollbar-gutter: stable;
  scroll-behavior: smooth;
  &::-webkit-scrollbar {
    height: 0px;
    width: 4px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.396);
    border-radius: 10px;
  }
}
</style>
