<script setup lang="ts">
defineEmits(['clear-selected']);
defineProps({
  selected: {
    type: Number as PropType<number | undefined>,
    required: false,
    default: undefined,
  },
  showClear: {
    type: Boolean as PropType<boolean>,
    required: false,
    default: true,
  },
});
const selectPanel = ref<HTMLElement | null>(null);
const tableWrapper = ref<HTMLElement | null>(null);
const leftOffset = ref(0);

function updatePanelPosition() {
  if (!tableWrapper.value || !selectPanel.value) return;

  const rect = tableWrapper.value.getBoundingClientRect();
  const panelWidth = selectPanel.value.offsetWidth;

  leftOffset.value = rect.left + rect.width / 2 - panelWidth / 2;
}

onMounted(() => {
  tableWrapper.value = document.querySelector('.table-v2-wrapper') as HTMLElement;
  nextTick(updatePanelPosition);

  window.addEventListener('resize', updatePanelPosition);
});
onBeforeUnmount(() => window.removeEventListener('resize', updatePanelPosition));
</script>
<template>
  <Teleport to="body">
    <div ref="selectPanel" class="table-select-panel" :style="{ left: leftOffset + 'px' }">
      <div class="flex gap-2 justify-between items-center">
        <div class="flex gap-4 items-center">
          <div v-if="showClear" class="w-8 h-8 flex items-center justify-center cursor-pointer" @click="$emit('clear-selected', true)">
            <i class="pi pi-times text-white hover:opacity-80" />
          </div>
          <span v-if="selected" class="whitespace-nowrap"
            >Выбрано:
            <slot v-if="$slots.selectCount" name="selectCount" />
            <template v-else>
              {{ selected }}
            </template>
          </span>
        </div>
        <slot />
      </div>
    </div>
  </Teleport>
</template>
<style scoped lang="scss">
.table-select-panel {
  position: fixed;
  width: 320px;
  padding: 10px;
  padding: 10px 24px 10px 24px;
  gap: 10px;
  bottom: 25px;
  border-radius: 18px;
  background: var(--app-black);
  color: #ffffff;
  z-index: 1200;
  @media (min-width: 560px) {
    width: 480px;
  }
}
</style>
