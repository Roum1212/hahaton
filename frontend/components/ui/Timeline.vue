<script setup lang="ts">
import { addDays } from 'date-fns';

export type TimelineValue = {
  id: string | number;
  name: string;
  startDate: string;
  endDate?: string;
};

const props = withDefaults(
  defineProps<{
    value: TimelineValue[];
    search?: string | undefined;
    from?: Date;
    loading?: boolean;
  }>(),
  {
    from: () => new Date(),
    search: undefined,
  },
);

const emit = defineEmits<{
  (e: 'click:value', value: TimelineValue): void;
}>();
const { isDarkTheme } = useAppLayout();
const search = ref<string | undefined>(props.search);
const now = new Date();
const startDate = computed(() =>
  props.from ? props.from : valuesSorted.value[0]?.startDate ? new Date(valuesSorted.value[0].startDate) : now,
);
const endDate = computed(() => {
  const max = valuesSorted.value.reduce((acc, promo) => {
    const d = new Date(promo.endDate ?? promo.startDate);
    return d > acc ? d : acc;
  }, new Date(startDate.value));
  return max;
});

const totalUnits = computed(() => {
  const diff = endDate.value.getTime() - startDate.value.getTime();
  return Math.floor(diff / 86400000) + 1;
});
const nowUnitIndex = computed(() => {
  const diff = now.getTime() - startDate.value.getTime();
  return Math.floor(diff / 86400000) + 1;
});

const valuesSorted = computed(() => [...props.value].sort((a, b) => new Date(a.startDate).getTime() - new Date(b.startDate).getTime()));
const filteredValues = computed(() =>
  valuesSorted.value.filter((p) => (search.value ? p.name.toLowerCase().includes(search.value.toLowerCase()) : true)),
);

const monthSeparators = computed(() => {
  const result: { label: string; left: number; date: Date }[] = [];

  const date = new Date(startDate.value);
  date.setDate(1);

  while (date <= endDate.value) {
    const label = formatDate(date, 'MMMM yyyy');
    let left = 0;
    const diffDays = (date.getTime() - startDate.value.getTime()) / 86400000;
    left = (diffDays / totalUnits.value) * 100;
    result.push({ label, left, date });
    date.setMonth(date.getMonth() + 1);
  }

  return result;
});

function getColor(name: string, saturation: number = 50, lightness: number = 95) {
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }
  const hue = Math.abs(hash) % 360;
  return `hsl(${hue}, ${saturation}%, ${isDarkTheme.value ? 120 - lightness : lightness}%)`;
}

function unitToLeft(dateStr: string) {
  const date = new Date(dateStr);
  const diff = date.getTime() - startDate.value.getTime();

  return (diff / 86400000 / totalUnits.value) * 100;
}
const unitWidth = computed(() => {
  return 100 / totalUnits.value;
});
function getWidth(start: string, end?: string) {
  const fromD = props.from || new Date(startDate.value);
  const startD = new Date(start);
  const endD = end ? new Date(end) : new Date(start);

  const startOffset = (startD.getTime() - fromD.getTime()) / 86400000;
  const endOffset = (endD.getTime() - fromD.getTime()) / 86400000;
  const widthInDays = endOffset - startOffset + 1;
  return (widthInDays / totalUnits.value) * 100;
}
function isSameMonthStart(date1: Date, date2: Date) {
  return date1.getFullYear() === date2.getFullYear() && date1.getMonth() === date2.getMonth();
}
function getTooltipText(promo: TimelineValue) {
  return `${promo.name} (${formatDate(promo.startDate)} - ${promo.endDate ? formatDate(promo.endDate) : 'бессрочно'})`;
}
function handleClick(item: TimelineValue) {
  emit('click:value', item);
}

const scrollContainer = ref<HTMLElement>();
const unitsRef = ref<HTMLElement>();
const listRef = ref<HTMLElement>();

let isDragging = false;
let dragStartX = 0;
let dragStartY = 0;
let scrollLeftStart = 0;
let scrollTopStart = 0;
let dragTimeout: NodeJS.Timeout | null = null;

function startDrag(e: MouseEvent) {
  dragTimeout = setTimeout(() => {
    isDragging = true;
    dragStartX = e.pageX;
    dragStartY = e.pageY;
    scrollLeftStart = scrollContainer.value?.scrollLeft ?? 0;
    scrollTopStart = scrollContainer.value?.scrollTop ?? 0;
  }, 150);
}

function onDrag(e: MouseEvent) {
  if (!isDragging || !scrollContainer.value) return;
  const dx = e.pageX - dragStartX;
  const dy = e.pageY - dragStartY;
  scrollContainer.value.scrollLeft = scrollLeftStart - dx;
  scrollContainer.value.scrollTop = scrollTopStart - dy;
}

function stopDrag() {
  if (dragTimeout) clearTimeout(dragTimeout);
  isDragging = false;
}
watch(
  () => props.search,
  (newValue) => {
    search.value = newValue;
  },
);
</script>
<template>
  <div class="relative w-full h-full overflow-hidden">
    <div v-if="loading" class="w-full flex items-center justify-center h-10">
      <i class="pi pi-spinner pi-spin" />
    </div>
    <div
      v-else
      ref="scrollContainer"
      class="timeline overflow-x-auto overflow-y-auto h-[calc(100%-40px)]"
      @mousedown="startDrag"
      @mousemove="onDrag"
      @mouseup="stopDrag"
      @mouseleave="stopDrag">
      <div class="min-w-[100vw] w-fit relative select-none">
        <div class="relative h-4 z-20 border-t border-surface-200 dark:border-surface-950 dark:bg-surface-950">
          <div
            v-for="(separator, idx) in monthSeparators"
            :key="'label-' + idx"
            class="absolute text-[10px] text-surface-500 dark:text-surface-50 whitespace-nowrap border-l border-surface-200 dark:border-surface-950 pl-1"
            :style="{
              top: '0px',
              left: isSameMonthStart(separator.date, startDate) ? '0%' : `calc(${separator.left}% + 2px)`,
            }">
            {{ separator.label }}
          </div>
        </div>

        <div ref="unitsRef" class="flex sticky top-0 bg-white dark:bg-surface-950 z-10 border-b border-surface-200 -full">
          <div class="flex flex-1">
            <div
              v-for="n in totalUnits"
              :key="n"
              class="px-1 text-center py-2 text-surface-600 dark:text-surface-100 text-xs relative"
              :style="{ width: unitWidth + '%', minWidth: unitWidth + '%' }">
              {{ formatDate(addDays(startDate, n - 1), 'dd') }}
              <div v-if="n === nowUnitIndex" class="h-full w-full absolute top-0 left-0 bg-surface-400 opacity-30 z-0" />
            </div>
          </div>
        </div>

        <div ref="listRef" class="flex flex-col">
          <div v-for="promo in filteredValues" :key="promo.id" class="flex group hover:bg-surface-50 dark:hover:bg-surface-950">
            <div class="flex-1 relative h-[30px]">
              <div
                class="absolute top-1 left-0 h-7 rounded-lg text-xs text-white px-2 flex items-center cursor-pointer"
                :style="{
                  left: unitToLeft(promo.startDate) + '%',
                  width: getWidth(promo.startDate, promo.endDate) + '%',
                  backgroundColor: getColor(promo.name),
                }"
                @click="handleClick(promo)">
                <div
                  v-tooltip="getTooltipText(promo)"
                  class="px-2 py-1 text-sm font-semibold truncate z-20 sticky left-0 text-white"
                  :style="{
                    color: getColor(promo.name, 60, 40),
                  }">
                  {{ promo.name }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<style scoped lang="scss">
.timeline {
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
</style>
