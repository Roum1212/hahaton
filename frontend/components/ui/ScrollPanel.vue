<script setup lang="ts">
import type { CSSProperties } from 'vue';

const emit = defineEmits(['update:scroll', 'scroll-bottom']);
const props = withDefaults(defineProps<{ style?: CSSProperties }>(), {
  style: undefined,
});

const scrollPanel = ref<HTMLElement | null>(null);
const isScrolling = ref<boolean>(false);
const startX = ref<number>(0);
const scrollLeft = ref<number>(0);

const startScroll = (event: MouseEvent) => {
  if (!scrollPanel.value) return;
  isScrolling.value = true;
  startX.value = event.pageX - scrollPanel.value.offsetLeft;
  scrollLeft.value = scrollPanel.value.scrollLeft;
};
const onScroll = (event: MouseEvent) => {
  if (!isScrolling.value || !scrollPanel.value) return;
  event.preventDefault();
  const x = event.pageX - scrollPanel.value.offsetLeft;
  const walk = (x - startX.value) * 2;
  scrollPanel.value.scrollLeft = scrollLeft.value - walk;
};

const stopScroll = () => {
  isScrolling.value = false;
};
const handleScroll = () => {
  if (scrollPanel.value) {
    const scrollTop = scrollPanel.value.scrollTop;
    const scrollHeight = scrollPanel.value.scrollHeight;
    const clientHeight = scrollPanel.value.clientHeight;

    emit('update:scroll', { scrollTop, scrollHeight, clientHeight });
    if (scrollTop + clientHeight + 30 >= scrollHeight) {
      emit('scroll-bottom', true);
    }
  }
};
onMounted(() => {
  if (scrollPanel.value) {
    scrollPanel.value.addEventListener('scroll', handleScroll);
  }
});
onBeforeUnmount(() => {
  if (scrollPanel.value) {
    scrollPanel.value.removeEventListener('scroll', handleScroll);
  }
});
</script>
<template>
  <div
    ref="scrollPanel"
    class="ui-scroll-panel"
    :class="{ grabbing: isScrolling }"
    :style="props.style"
    @mousedown="startScroll"
    @mousemove="onScroll"
    @mouseup="stopScroll"
    @mouseleave="stopScroll">
    <div>
      <slot />
    </div>
  </div>
</template>
<style scoped>
.ui-scroll-panel {
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 4px;
}

.ui-scroll-panel::-webkit-scrollbar {
  width: 4px;
}

.ui-scroll-panel::-webkit-scrollbar-track {
  background: transparent;
}

.ui-scroll-panel::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.396);
  border-radius: 10px;
}

.ui-scroll-panel.grabbing {
  cursor: grabbing;
}
</style>
