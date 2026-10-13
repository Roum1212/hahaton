<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts" generic="T">
import type { CSSProperties } from 'vue';

interface Props<T> {
  items: T[];
  itemHeight?: number;
  autoHeight?: boolean;
  buffer?: number;
  height?: string;
  loading?: boolean;
  allowInnerScroll?: boolean;
}

const props = withDefaults(defineProps<Props<T>>(), {
  itemHeight: 40,
  autoHeight: false,
  buffer: 5,
  height: '',
  loading: false,
  allowInnerScroll: true,
});

defineExpose({
  scrollToPos,
  scrollToTop,
});
const emit = defineEmits(['pre-scroll-bottom', 'scroll-bottom', 'scroll-top', 'scroll']);

const container = ref<HTMLElement | null>(null);
const scrollTop = ref(0);
const itemHeights = ref<number[]>([]);

const getItemKey = (item: T, index: number) => {
  return (item as any)?.id ?? `item-${index}`;
};
const totalHeight = computed(() =>
  props.autoHeight ? itemHeights.value.reduce((a, b) => a + b, 0) : props.items.length * props.itemHeight,
);
const cumHeights = computed(() => {
  const arr: number[] = [];
  let sum = 0;
  for (const h of itemHeights.value) {
    sum += h;
    arr.push(sum);
  }
  return arr;
});

const startIndex = computed(() => {
  if (!props.autoHeight) {
    const indexFloat = scrollTop.value / props.itemHeight;
    const indexBase = Math.floor(indexFloat);
    const fraction = indexFloat - indexBase;
    const threshold = 0.4;
    const adjustedIndex = fraction > threshold ? indexBase : Math.max(0, indexBase - 1);
    return Math.max(0, adjustedIndex - props.buffer);
  }

  const scroll = scrollTop.value;
  if (cumHeights.value.length === 0) return 0;

  let low = 0;
  let high = cumHeights.value.length - 1;
  let result = 0;

  while (low <= high) {
    const mid = Math.floor((low + high) / 2);
    if (cumHeights.value[mid] && cumHeights.value[mid] <= scroll) {
      result = mid + 1;
      low = mid + 1;
    } else {
      high = mid - 1;
    }
  }

  return Math.max(0, result - props.buffer);
});

const offsetTop = computed(() => {
  if (!props.autoHeight) {
    return startIndex.value * props.itemHeight;
  }
  return startIndex.value === 0 ? 0 : cumHeights.value[startIndex.value - 1];
});

const endIndex = computed(() => {
  if (!props.autoHeight) {
    const visibleCount = Math.ceil((props.height ? parseInt(props.height) : window.innerHeight) / props.itemHeight);
    return Math.min(props.items.length, startIndex.value + visibleCount + props.buffer * 2);
  }

  const containerHeight = props.height ? parseInt(props.height) : window.innerHeight;
  const scrollBottom = scrollTop.value + containerHeight;

  let end = startIndex.value;
  let currentHeight = offsetTop.value || 0;

  while (end < props.items.length && currentHeight < scrollBottom + props.buffer * 100) {
    currentHeight += itemHeights.value[end] || props.itemHeight;
    end++;
  }

  return Math.min(props.items.length, end + props.buffer);
});

const visibleItems = computed(() => props.items.slice(startIndex.value, endIndex.value));

const containerStyle = computed<CSSProperties>(() => ({
  height: props.height || '100%',
  overflowY: props.allowInnerScroll ? 'auto' : 'hidden',
  overflowX: 'auto',
  position: 'relative',
}));
const contentStyle = computed<CSSProperties>(() => ({
  transform: `translate3d(0, ${offsetTop.value}px, 0)`,
  position: 'relative',
  top: 0,
  left: 0,
  willChange: 'transform',
  width: 'fit-content',
}));
const itemStyle = computed<CSSProperties>(() => ({
  height: !props.autoHeight ? `${props.itemHeight}px` : 'auto',
}));

let scrollRaf: number | null = null;

const handleScroll = () => {
  if (props.loading || !container.value) return;

  const el = container.value;
  const currentScrollTop = el.scrollTop;

  if (Math.abs(currentScrollTop - scrollTop.value) < 0.5) return;

  if (scrollRaf) {
    cancelAnimationFrame(scrollRaf);
  }

  scrollRaf = requestAnimationFrame(() => {
    scrollTop.value = currentScrollTop;

    const scrollHeight = el.scrollHeight;
    const clientHeight = el.clientHeight;
    const scrollBottom = scrollHeight - (currentScrollTop + clientHeight);

    if (scrollBottom <= 150) emit('pre-scroll-bottom');
    if (scrollBottom <= 10) emit('scroll-bottom');
    if (currentScrollTop <= 5) emit('scroll-top');

    scrollRaf = null;
  });
  emit('scroll', scrollTop.value);
};

const isDraggingH = ref(false);
const dragStartX = ref(0);
const scrollStartX = ref(0);
const dragStartY = ref(0);
const dragThreshold = 5;

function onMouseDownH(e: MouseEvent) {
  if (!container.value) return;
  const el = e.target as HTMLElement;
  const isSelect = !el.classList.contains('table-cell') && !el.classList.contains('table-row') && !el.classList.contains('vs-item');
  if (isSelect || e.button !== 0) {
    container.value.style.userSelect = 'text';
    return;
  }

  dragStartX.value = e.clientX;
  dragStartY.value = e.clientY;
  scrollStartX.value = container.value.scrollLeft;

  window.addEventListener('mousemove', onMouseMoveH);
  window.addEventListener('mouseup', onMouseUpH);
  e.preventDefault();
}

let raf: number | null = null;
function onMouseMoveH(e: MouseEvent) {
  const dx = e.clientX - dragStartX.value;
  const dy = e.clientY - dragStartY.value;

  if (!isDraggingH.value && Math.abs(dy) > Math.abs(dx) && Math.abs(dy) > dragThreshold) {
    return;
  }

  if (!isDraggingH.value && Math.abs(dx) > dragThreshold) {
    isDraggingH.value = true;
    if (container.value) {
      container.value.style.cursor = 'grabbing';
      container.value.style.userSelect = 'none';
    }
  }

  if (isDraggingH.value && container.value) {
    if (raf) cancelAnimationFrame(raf);
    raf = requestAnimationFrame(() => {
      container.value!.scrollLeft = scrollStartX.value - dx;
      raf = null;
    });
  }
}

function onMouseUpH(e: MouseEvent) {
  window.removeEventListener('mousemove', onMouseMoveH);
  window.removeEventListener('mouseup', onMouseUpH);

  if (isDraggingH.value) {
    isDraggingH.value = false;

    if (container.value) {
      container.value.style.cursor = '';
      container.value.style.removeProperty('user-select');
    }

    e.preventDefault();
  }
}

const updateItemHeights = () => {
  if (!container.value || !props.autoHeight) return;

  requestAnimationFrame(() => {
    const children = container.value!.querySelectorAll('.vs-item:not(.vs-loading)');
    const newHeights = Array.from(children).map((el) => (el as HTMLElement).offsetHeight);

    if (newHeights.length !== itemHeights.value.length || !newHeights.every((h, i) => h === itemHeights.value[i])) {
      itemHeights.value = newHeights;
    }
  });
};

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

let resizeObserver: ResizeObserver | null = null;

onMounted(() => {
  const el = container.value;
  if (!el) return;
  el.addEventListener('scroll', handleScroll, { passive: true });
  resizeObserver = new ResizeObserver(() => {
    if (props.autoHeight) updateItemHeights();
  });
  resizeObserver.observe(el);
  if (props.autoHeight) setTimeout(updateItemHeights, 100);
});

watch(
  visibleItems,
  () => {
    if (props.autoHeight) {
      nextTick(updateItemHeights);
    }
  },
  { flush: 'post' },
);

onBeforeUnmount(() => {
  const el = container.value;
  if (el) el.removeEventListener('scroll', handleScroll);
  if (resizeObserver) resizeObserver.disconnect();

  window.removeEventListener('mousemove', onMouseMoveH);
  window.removeEventListener('mouseup', onMouseUpH);

  if (scrollRaf) cancelAnimationFrame(scrollRaf);
  if (raf) cancelAnimationFrame(raf);
});
</script>

<template>
  <div ref="container" class="virtual-scroll" :style="containerStyle" @mousedown="onMouseDownH">
    <div class="virtual-scroll__container" :style="{ height: totalHeight + 'px', position: 'relative' }">
      <div class="virtual-scroll__content" :style="contentStyle">
        <div v-for="(item, index) in visibleItems" :key="getItemKey(item, startIndex + index)" class="vs-item" :style="itemStyle">
          <slot :item="item" :index="startIndex + index" />
        </div>
      </div>
    </div>
    <div v-show="loading" class="vs-loading">
      <slot name="loading">
        <div class="default-loader">Загрузка...</div>
      </slot>
    </div>
  </div>
</template>

<style scoped lang="scss">
.virtual-scroll {
  width: 100%;
  scrollbar-gutter: stable;
  transform: translateZ(0);
  will-change: transform;
  -webkit-overflow-scrolling: touch;
  overflow-anchor: none;

  -webkit-backface-visibility: hidden;
  backface-visibility: hidden;
  -webkit-perspective: 1000;
  perspective: 1000;
  contain: layout style;
  overflow-y: auto;
  overflow-x: auto;

  &::-webkit-scrollbar {
    height: 4px;
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

.virtual-scroll__container {
  will-change: height;
  width: fit-content;
}

.virtual-scroll__content {
  contain: layout style;
  will-change: transform;
  width: fit-content;
}

.vs-item {
  width: 100%;
  transform: translateZ(0);
  contain: layout style paint;
  will-change: transform;
  user-select: text;

  -webkit-transform: translate3d(0, 0, 0);
  -webkit-backface-visibility: hidden;
  backface-visibility: hidden;
}

.vs-loading {
  text-align: center;
  padding: 12px;
  contain: layout style paint;
}

.default-loader {
  color: #666;
  font-size: 14px;
}
</style>
