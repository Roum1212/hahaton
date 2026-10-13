<script lang="ts" setup>
import { Eye, EyeOff, Pin, PinOff } from 'lucide-vue-next';
import type { CSSProperties } from 'vue';
const emits = defineEmits(['update:pin', 'update:sortRelation', 'update:visible']);
const props = withDefaults(
  defineProps<{
    title?: string;
    technical?: boolean;
    pin?: boolean;
    visible?: boolean;
    style?: CSSProperties | undefined;
  }>(),
  {
    title: '',
    pin: false,
    visible: false,
    style: undefined,
    technical: false,
  },
);
const th = ref();
const dropdownRef = ref();
const localPin = ref(props.pin);
const localVisible = ref(props.visible);
const dropdownVisible = ref(false);
const dropdownPosition = ref({ top: 0, left: 0 });

function updateDropdownPosition() {
  if (th.value) {
    const rect = th.value.getBoundingClientRect();
    dropdownPosition.value = {
      top: rect.bottom + window.scrollY + 5,
      left: rect.left + window.scrollX,
    };
  }
}

function handleDropdownToggle() {
  dropdownVisible.value = !dropdownVisible.value;
  if (dropdownVisible.value) {
    nextTick(() => {
      updateDropdownPosition();
    });
  }
}

function handlePin() {
  localPin.value = !localPin.value;
  emits('update:pin', localPin.value);
}
function handleVisible() {
  localVisible.value = !localPin.value;
  dropdownVisible.value = false;
  emits('update:visible', !props.visible);
}

function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node) && th.value && !th.value.contains(event.target as Node)) {
    dropdownVisible.value = false;
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
  window.addEventListener('scroll', updateDropdownPosition);
  window.addEventListener('resize', updateDropdownPosition);
});

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside);
  window.removeEventListener('scroll', updateDropdownPosition);
  window.removeEventListener('resize', updateDropdownPosition);
});
</script>
<template>
  <div ref="th" class="gap-2 head-content" :style="style">
    <span class="text-sm text-gray-500 dark:text-gray-300">{{ title }}</span>
    <template v-if="!technical">
      <Button
        text
        :icon="dropdownVisible ? 'pi pi-angle-up' : 'pi pi-angle-down'"
        severity="secondary"
        size="small"
        class="!p-0 w-5"
        @click="handleDropdownToggle" />
      <Teleport to="body">
        <div
          v-if="dropdownVisible"
          ref="dropdownRef"
          class="th-layout absolute bg-white dark:bg-surface-800 border p-2 border-surface-300 dark:border-surface-700 rounded-xl shadow-md max-w-[200px] flex flex-col gap-1"
          :style="{
            zIndex: 9999,
            top: `${dropdownPosition.top}px`,
            left: `${dropdownPosition.left}px`,
            position: 'absolute',
          }">
          <Button severity="secondary" text size="small" class="w-full !justify-start" @click="handlePin">
            <PinOff v-if="pin" :size="14" /><Pin v-else :size="14" />
            <span class="text-sm">{{ pin ? 'Открепить столбец' : 'Закрепить столбец' }}</span>
          </Button>
          <Button severity="secondary" text size="small" class="w-full !justify-start" @click="handleVisible">
            <EyeOff v-if="visible" :size="14" /><Eye v-else :size="14" />
            <span class="text-sm">{{ visible ? 'Скрыть столбец' : 'Показать столбец' }}</span>
          </Button>
        </div>
      </Teleport>
    </template>
  </div>
</template>
<style scoped>
.head-content {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  white-space: nowrap;
}
</style>
