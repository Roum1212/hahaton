<script setup lang="ts">
type Position = 'left' | 'right' | 'center' | 'topLeft' | 'topCenter' | 'topRight' | 'bottomLeft' | 'bottomCenter' | 'bottomRight';
const props = withDefaults(
  defineProps<{
    relationEl?: HTMLElement | null;
    position?: Position;
    offsetTop?: number;
    offsetBottom?: number;
    offsetRight?: number;
    offsetLeft?: number;
  }>(),
  {
    relationEl: undefined,
    position: 'bottomCenter',
    offsetBottom: undefined,
    offsetTop: undefined,
    offsetLeft: undefined,
    offsetRight: undefined,
  },
);
const style = ref<Record<string, string>>({});
const panelRef = ref<HTMLElement | null>(null);
onMounted(() => {
  updatePosition();
  window.addEventListener('scroll', updatePosition);
  window.addEventListener('resize', updatePosition);
});

onBeforeUnmount(() => {
  window.removeEventListener('scroll', updatePosition);
  window.removeEventListener('resize', updatePosition);
});
function updatePosition() {
  const offsetTop = props.offsetTop ?? 50;
  const offsetBottom = props.offsetBottom ?? 50;
  const offsetLeft = props.offsetLeft ?? 50;
  const offsetRight = props.offsetRight ?? 50;

  if (!props.relationEl) {
    style.value = getFixedStyle(props.position, offsetTop, offsetBottom, offsetLeft, offsetRight);
    return;
  }

  const rect = props.relationEl.getBoundingClientRect();

  let top = 0;
  let left = 0;

  switch (props.position) {
    case 'topLeft':
      top = rect.top - offsetTop;
      left = rect.left + offsetLeft;
      break;
    case 'topCenter':
      top = rect.top - offsetTop;
      left = rect.left + rect.width / 2;
      break;
    case 'topRight':
      top = rect.top - offsetTop;
      left = rect.right - offsetRight;
      break;
    case 'bottomLeft':
      top = rect.bottom + offsetTop;
      left = rect.left + offsetLeft;
      break;
    case 'bottomCenter':
      top = rect.bottom + offsetTop;
      left = rect.left + rect.width / 2;
      break;
    case 'bottomRight':
      top = rect.bottom + offsetTop;
      left = rect.right - offsetRight - (panelRef.value?.offsetWidth || 0);
      break;
  }

  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;

  const panelWidth = 320;
  const panelHeight = 100;

  top = Math.min(Math.max(top, 0), viewportHeight - panelHeight);
  left = Math.min(Math.max(left, 0), viewportWidth - panelWidth);

  style.value = {
    position: 'fixed',
    top: `${top}px`,
    left: `${left}px`,
    transform: props.position?.includes('Center') ? 'translateX(-50%)' : '',
  };
}
function getFixedStyle(position: string, offsetTop: number, offsetBottom: number, offsetLeft: number, offsetRight: number) {
  const styleObj: Record<string, string> = { position: 'fixed' };

  switch (position) {
    case 'topLeft':
      styleObj.top = `${offsetTop}px`;
      styleObj.left = `${offsetLeft}px`;
      break;
    case 'topCenter':
      styleObj.top = `${offsetTop}px`;
      styleObj.left = '50%';
      styleObj.transform = 'translateX(-50%)';
      break;
    case 'topRight':
      styleObj.top = `${offsetTop}px`;
      styleObj.right = `${offsetRight}px`;
      break;
    case 'bottomLeft':
      styleObj.bottom = `${offsetBottom}px`;
      styleObj.left = `${offsetLeft}px`;
      break;
    case 'bottomCenter':
      styleObj.bottom = `${offsetBottom}px`;
      styleObj.left = '50%';
      styleObj.transform = 'translateX(-50%)';
      break;
    case 'bottomRight':
      styleObj.bottom = `${offsetBottom}px`;
      styleObj.right = `${offsetRight}px`;
      break;
  }

  return styleObj;
}
watch(
  () => props.relationEl,
  () => updatePosition,
);
</script>
<template>
  <div ref="panelRef" class="floating-panel" :style="style">
    <slot />
  </div>
</template>

<style scoped lang="scss">
.floating-panel {
  width: fit-content;
  max-width: 300px;
  padding: 10px 24px;
  border-radius: 18px;
  background: #000000ed;
  color: #ffffff;
  z-index: 500;

  @media (min-width: 560px) {
    max-width: 480px;
  }
}
</style>
