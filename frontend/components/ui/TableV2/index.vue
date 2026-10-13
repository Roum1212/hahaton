<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<!-- eslint-disable @typescript-eslint/no-invalid-void-type -->
<script setup lang="ts" generic="T extends object extends {id?:string | number}">
import type { CSSProperties } from 'vue';

export interface TableColumnType {
  key: string;
  header?: string;
  width: number;
  visible?: boolean;
  pin?: boolean;
  pinRelation?: 'left' | 'right';
  order?: number;
  sortable?: boolean;
  technical?: boolean;
  temp?: boolean;
  [key: string]: any;
}
type RowType = T;
type KeyOfRow = keyof RowType;

const emit = defineEmits<{
  'scroll-bottom': [flag: boolean];
  'pre-scroll-bottom': [flag: boolean];
  'update:select': [selected: RowType[KeyOfRow][], entities: RowType[]];
  'update:columns': [columns: TableColumnType[]];
  'column-resize': [columns: TableColumnType[]];
  'update:select-all': [flag: boolean];
  resetSettings: [];
}>();
const props = withDefaults(
  defineProps<{
    values?: RowType[];
    loading?: boolean;
    fetching?: boolean;
    select?: RowType[KeyOfRow][];
    selectMode?: 'checkbox' | 'fill' | undefined;
    hoverMode?: boolean;
    valueKey?: string;
    searchFilter?: string | undefined;
    isShowSelectPanel?: boolean;
    isShowSettings?: boolean;
    scrollContainer?: string | 'auto';
    isShowResult?: boolean;
    isShowScrollTopButton?: boolean;
    columns?: TableColumnType[];
    virtualize?: boolean;
    virtualizeRowHeight?: number;
    virtualizeRowAutoHeight?: boolean;
    clickAction?: (row: RowType, e: Event) => void;
    isResizable?: boolean;
    placeholderLines?: number;
    allowInnerScroll?: boolean;
  }>(),
  {
    loading: false,
    fetching: false,
    values: () => [],
    columns: () => [],
    select: undefined,
    selectMode: undefined,
    isShowSelectPanel: true,
    isShowScrollTopButton: true,
    hoverMode: false,
    valueKey: 'id',
    searchFilter: '',
    clickAction: () => void 0,
    isShowSettings: false,
    scrollContainer: undefined,
    isShowResult: false,
    virtualize: false,
    virtualizeRowHeight: 58,
    virtualizeRowAutoHeight: false,
    isResizable: true,
    placeholderLines: 20,
    allowInnerScroll: true,
  },
);

const tableRef = ref<HTMLDivElement | null>(null);
const headerRef = ref<HTMLDivElement | null>(null);
const headerOuterRef = ref<HTMLDivElement | null>(null);
const bodyOuterRef = ref<HTMLDivElement | null>(null);
const scrollContainerRef = ref<HTMLDivElement | null>(null);

const localValueKey: keyof RowType = (props.valueKey as KeyOfRow) || ('id' as keyof T);

const selectedKeys = ref<T[KeyOfRow][]>(Array.isArray(props.select) ? [...props.select] : []) as Ref<T[keyof T][]>;
const lastClickWasShift = ref(false);
const columnsState = ref<TableColumnType[]>(initializeColumns());

onMounted(() => {
  nextTick(autoResizeColumns);
  observeEvents();
  syncHeaderEvents();
});
const observeEvents = () => {
  const resizeObserver: any = null;

  if (resizeObserver && scrollContainerRef.value) {
    resizeObserver.unobserve(scrollContainerRef.value);
    resizeObserver.disconnect();
  }
  onBeforeUnmount(() => {
    if (resizeObserver && scrollContainerRef.value) {
      resizeObserver.unobserve(scrollContainerRef.value);
      resizeObserver.disconnect();
    }
  });
};
const syncHeaderEvents = () => {
  let raf: number | null = null;

  const header = headerOuterRef.value;
  const bodyContainer =
    scrollContainerRef.value?.querySelector('.virtual-scroll') || scrollContainerRef.value?.querySelector('.base-table-body');
  if (!header || !bodyContainer) return;

  const syncHeaderScroll = () => {
    if (raf) cancelAnimationFrame(raf);
    raf = requestAnimationFrame(() => {
      header.scrollLeft = bodyContainer.scrollLeft;
      raf = null;
    });
  };
  bodyContainer.addEventListener('scroll', syncHeaderScroll);

  onBeforeUnmount(() => {
    bodyContainer.removeEventListener('scroll', syncHeaderScroll);
    if (raf) cancelAnimationFrame(raf);
  });
};

watch(
  () => props.select,
  (v) => {
    if (Array.isArray(v)) selectedKeys.value = [...v];
  },
);
watch(
  () => props.columns.map((el) => el.visible),
  () => {
    columnsState.value = initializeColumns();
    if (props.loading) return;
    nextTick(autoResizeColumns);
  },
);
watch(
  () => props.loading,
  () => {
    if (props.loading) return;
    nextTick(autoResizeColumns);
  },
);

function initializeColumns(): TableColumnType[] {
  if (props.columns && props.columns.length) {
    return props.columns.map((c, i) => ({
      key: c.key || `col_${i}`,
      header: c.header !== undefined ? c.header : c.key || `Column ${i}`,
      width: c.width,
      visible: c.visible !== false,
      pin: !!c.pin,
      pinRelation: c.pinRelation || 'left',
      order: typeof c.order === 'number' ? c.order : i,
      sortable: !!c.sortable,
      technical: !!c.technical,
    }));
  }
  const first = props.values && props.values[0];
  if (first)
    return Object.keys(first).map((k, i) => ({
      key: k,
      header: k,
      width: 160,
      visible: true,
      pin: false,
      pinRelation: 'left',
      order: i,
      sortable: false,
    }));
  return [];
}

const visibleCols = computed<TableColumnType[]>(() => columnsState.value.filter((c) => c.visible));

const orderedCols = computed<TableColumnType[]>(() => visibleCols.value.slice().sort((a, b) => (a.order || 0) - (b.order || 0)));
const orderedColsWithTemp = computed<TableColumnType[]>(() => {
  const result: TableColumnType[] = [];
  let isAddTempCol = false;
  const tempCol: TableColumnType = { key: 'temp', width: 0, visible: true, sortable: false, technical: true, temp: true };
  orderedCols.value.forEach((col) => {
    if (col.pin && col.pinRelation === 'right' && !isAddTempCol) {
      isAddTempCol = true;
      result.push(tempCol);
    }
    result.push(col);
  });

  if (!isAddTempCol) result.push(tempCol);
  return result;
});
const selectCol = { key: '__select', width: 60, pin: true, pinRelation: 'left' };
const lastSelectedIndex = ref<number | null>(null);

const totalTableWidth = computed(() => {
  let w = 0;
  if (props.selectMode === 'checkbox') w += selectCol.width;
  for (const c of orderedColsWithTemp.value) w += c.width || 120;
  return w;
});

const filteredData = computed<RowType[]>(() => {
  if (!props.searchFilter || !props.searchFilter.length) return props.values;
  const q = props.searchFilter.toLowerCase();
  return props.values.filter((r) => Object.values(r).some((v) => String(v).toLowerCase().includes(q)));
});

const sort = reactive<{ key: string | null; order: 'asc' | 'desc' }>({ key: null, order: 'asc' });

const displayedData = computed(() => {
  const res = filteredData.value.slice();

  if (sort.key) {
    const key = sort.key as keyof (typeof res)[0];
    res.sort((a, b) => {
      const av = a[key];
      const bv = b[key];

      if (av == null) return 1;
      if (bv == null) return -1;
      if (av === bv) return 0;
      return ((av as string | number) > (bv as string | number) ? 1 : -1) * (sort.order === 'asc' ? 1 : -1);
    });
  }

  return res;
});
const allSelected = computed<boolean>(
  () => !!filteredAndSortedData.value.length && filteredAndSortedData.value.every((r) => selectedKeys.value.includes(r[localValueKey])),
);
const pinnedOffsets = computed(() => {
  let left = props.selectMode === 'checkbox' ? selectCol.width : 0;
  let right = 0;
  const leftMap: Record<string, number> = {};
  const rightMap: Record<string, number> = {};

  for (const col of orderedCols.value) {
    if (col.pin && col.pinRelation === 'left') {
      leftMap[col.key] = left;
      left += col.width || 120;
    }
  }

  const rightPinned = orderedCols.value.filter((c) => c.pin && c.pinRelation === 'right').reverse();
  for (const col of rightPinned) {
    rightMap[col.key] = right;
    right += col.width || 120;
  }

  return { leftMap, rightMap };
});
const containerHeight = computed(() => {
  if (props.scrollContainer === 'auto') return 'initial';
  if (!props.scrollContainer) {
    if (!tableRef.value) return undefined;

    const windowHeight = window.innerHeight;
    const headerHeight = 40;
    const gap = 40;
    const top = tableRef.value.getBoundingClientRect().top;
    return windowHeight - top - headerHeight - gap + 'px';
  }
  return typeof props.scrollContainer === 'number' ? props.scrollContainer + 'px' : props.scrollContainer;
});
const filteredAndSortedData = displayedData;

function isSelected(id: RowType[KeyOfRow]) {
  return selectedKeys.value.includes(id);
}
function onRowCheckboxClickCapture(e: MouseEvent) {
  lastClickWasShift.value = e.shiftKey === true;
}
function onRowCheckboxChange(checked: boolean, row: RowType) {
  const id = row[localValueKey] as unknown as T[KeyOfRow];

  if (lastClickWasShift.value && lastSelectedIndex.value !== null) {
    const currentIndex = filteredAndSortedData.value.findIndex((r) => r[localValueKey] === id);
    const start = Math.min(lastSelectedIndex.value, currentIndex);
    const end = Math.max(lastSelectedIndex.value, currentIndex);
    const rangeIds = filteredAndSortedData.value.slice(start, end + 1).map((r) => r[localValueKey]);

    if (checked) {
      selectedKeys.value = Array.from(new Set([...selectedKeys.value, ...rangeIds]));
    } else {
      selectedKeys.value = selectedKeys.value.filter((k) => !rangeIds.includes(k));
    }

    lastSelectedIndex.value = currentIndex;
  } else {
    const i = selectedKeys.value.indexOf(id);
    if (checked) {
      if (i === -1) selectedKeys.value.push(id);
    } else {
      if (i !== -1) selectedKeys.value.splice(i, 1);
    }

    lastSelectedIndex.value = filteredAndSortedData.value.findIndex((r) => r[localValueKey] === id);
  }
  const clearSelectedKeys = selectedKeys.value.slice();
  const entities: RowType[] = [];
  for (const id of clearSelectedKeys) {
    entities.push(filteredAndSortedData.value.find((r) => r[localValueKey] === id)!);
  }

  emit('update:select', clearSelectedKeys, entities);
}
function toggleSelectAll(flag: boolean) {
  if (flag) selectedKeys.value = filteredAndSortedData.value.map((r) => r[localValueKey]);
  else selectedKeys.value = [];
  const clearSelectedKeys = selectedKeys.value.slice();
  const entities: RowType[] = [];
  for (const id of clearSelectedKeys) {
    entities.push(filteredAndSortedData.value.find((r) => r[localValueKey] === id)!);
  }

  emit('update:select', clearSelectedKeys, entities);
  emit('update:select-all', flag);
}
function clearSelection() {
  selectedKeys.value = [];
  emit('update:select', [], []);
}

function onRowClick(row: RowType, ev: Event) {
  if (props.selectMode === 'fill') {
    const id = row[localValueKey];
    if (selectedKeys.value.includes(id)) selectedKeys.value = selectedKeys.value.filter((x) => x !== id);
    else selectedKeys.value.push(id);
    const clearSelectedKeys = selectedKeys.value.slice();
    const entities: RowType[] = [];
    for (const id of clearSelectedKeys) {
      entities.push(filteredAndSortedData.value.find((r) => r[localValueKey] === id)!);
    }

    emit('update:select', clearSelectedKeys, entities);
  }
  if (typeof props.clickAction === 'function') props?.clickAction(row, ev);
}

function rowClass(row: RowType, index: number) {
  return {
    'is-selected': selectedKeys.value.includes(row[localValueKey]),
    'is-odd': index % 2 === 1,
  };
}
function getCellValue(row: RowType, col: TableColumnType) {
  const key = col.key;
  if (key.includes('.')) {
    return key.split('.').reduce<any>((acc, p) => (acc && typeof acc === 'object' ? (acc as Record<string, any>)[p] : null), row);
  }
  return (row as Record<string, any>)[key];
}

function colPinnedStyle(col: TableColumnType) {
  const base: CSSProperties = {
    width: (col.width || 120) + 'px',
    minWidth: (col.width || 120) + 'px',
    maxWidth: (col.width || 120) + 'px',
  };
  if (!col.pin) return base;

  if (col.pinRelation === 'left') base.left = pinnedOffsets.value.leftMap[col.key] + 'px';
  else base.right = pinnedOffsets.value.rightMap[col.key] + 'px';

  base.position = 'sticky';
  base.zIndex = 6;
  return base;
}

function handleUpdateColumns(cols: TableColumnType[]) {
  const newColsMap = listToMap(cols, (c) => c.key);

  columnsState.value.forEach((c) => {
    const newSettings = newColsMap[c.key];
    if (newSettings) {
      c.order = newSettings.order;
      c.visible = newSettings.visible;
    }
  });
  emit('update:columns', columnsState.value);
  nextTick(autoResizeColumns);
}
function startColumnResize(e: MouseEvent, col: TableColumnType) {
  if (!props.isResizable) return;
  e.stopPropagation();
  e.preventDefault();

  const startX = e.clientX;
  const startWidth = col.width || 120;

  const onMouseMove = (moveEvent: MouseEvent) => {
    const delta = moveEvent.clientX - startX;
    const newWidth = Math.max(60, startWidth + delta);
    col.width = newWidth;
  };

  const onMouseUp = () => {
    document.removeEventListener('mousemove', onMouseMove);
    document.removeEventListener('mouseup', onMouseUp);
    emit('update:columns', columnsState.value);
    emit('column-resize', columnsState.value);
  };

  document.addEventListener('mousemove', onMouseMove);
  document.addEventListener('mouseup', onMouseUp);
}

const tempColWidth = ref(0);
function autoResizeColumns() {
  const containerWidth = scrollContainerRef.value?.clientWidth;

  if (!containerWidth) return;

  const row = scrollContainerRef.value?.querySelector('.table-row');
  const otherCols = row?.querySelectorAll('.table-cell:not(.table-cell-temp)');
  let otherColsWidth = 0;
  otherCols?.forEach((el) => (otherColsWidth += el.clientWidth));
  const gap = 4;
  const width = containerWidth - otherColsWidth - gap;
  tempColWidth.value = width < 20 ? 20 : width;
}

const showScrollTop = ref(false);
const lastScrollTop = ref(0);
let ticking = false;

function onScroll(position: number) {
  if (!ticking) {
    requestAnimationFrame(() => {
      const isScrollingDown = position > lastScrollTop.value;

      if (isScrollingDown && position > 300) {
        showScrollTop.value = true;
      } else if (!isScrollingDown) {
        showScrollTop.value = false;
      }

      lastScrollTop.value = position;
      ticking = false;
    });
    ticking = true;
  }
}
function handleScrollToTop() {
  if (!bodyOuterRef.value || !(bodyOuterRef.value as any)?.scrollToTop) return;
  (bodyOuterRef.value as any).scrollToTop();
}
</script>
<template>
  <UiTableV2Result v-if="isShowResult" :loading="loading" :count="filteredData.length" />
  <div
    ref="tableRef"
    class="table-v2-wrapper bg-white dark:bg-surface-900"
    :class="{ 'table-v2-wrapper--loading': loading }"
    :style="{ '--row-height': virtualize ? virtualizeRowHeight + 'px' : 'auto' }">
    <div ref="headerOuterRef" class="table-header-outer">
      <div ref="headerRef" class="table-header" :style="{ width: totalTableWidth + 'px' }">
        <div class="table-row table-row--head">
          <div v-if="selectMode === 'checkbox'" class="table-cell table-cell--head pinned pinned-left left-0 table-cell--select">
            <Checkbox :model-value="allSelected" binary size="small" @update:model-value="toggleSelectAll" />
          </div>
          <div
            v-for="col in orderedColsWithTemp"
            :key="`header-${col.key}`"
            class="table-cell table-cell--head relative gap-2 items-center"
            :class="[{ pinned: col.pin }, col.pin ? `pinned-${col.pinRelation}` : '', col.temp ? 'table-cell-temp' : '']"
            :style="colPinnedStyle(col)">
            <div class="head-content inline-flex">
              <slot :name="`header-${col.key}`" :column="col">
                <UiTableV2HeaderColumn
                  v-model:pin="col.pin"
                  v-model:visible="col.visible"
                  :title="col.header"
                  :technical="col.technical"
                  :style="colPinnedStyle(col)" />
              </slot>
            </div>

            <div v-if="isResizable && !col.temp" class="resize-handle" @mousedown="(e) => startColumnResize(e, col)" />
          </div>
        </div>
      </div>
    </div>
    <div
      ref="scrollContainerRef"
      class="table-scroll-container"
      :style="{ '--temp-col-width': tempColWidth + 'px', height: !virtualize && containerHeight ? containerHeight : 'initial' }">
      <UiTableV2Placeholder
        v-if="loading"
        ref="bodyOuterRef"
        :rows="placeholderLines"
        :columns="orderedColsWithTemp"
        :select-mode="selectMode" />

      <slot v-else-if="!filteredData.length" name="empty">
        <UiTableV2Empty />
      </slot>
      <template v-if="virtualize">
        <UiVirtualScroll
          v-show="!loading && filteredData.length"
          ref="bodyOuterRef"
          v-model:items="filteredData"
          :item-height="virtualizeRowHeight"
          :buffer="10"
          :height="containerHeight"
          :loading="fetching"
          :allow-inner-scroll="allowInnerScroll"
          @scroll="onScroll"
          @pre-scroll-bottom="emit('pre-scroll-bottom', true)"
          @scroll-bottom="emit('scroll-bottom', true)">
          <template #default="{ item: row, index }">
            <div class="table-row" :class="rowClass(row, index)" @click="onRowClick(row, $event)">
              <div v-if="selectMode === 'checkbox'" class="table-cell table-cell--select pinned pinned-left left-0">
                <Checkbox
                  :model-value="isSelected(row[localValueKey as KeyOfRow])"
                  binary
                  size="small"
                  @click.stop="onRowCheckboxClickCapture($event)"
                  @update:model-value="(checked) => onRowCheckboxChange(checked, row)" />
              </div>
              <div
                v-for="(col, colIndex) in orderedColsWithTemp"
                :key="`${index}_${col.key}_${colIndex}`"
                class="table-cell"
                :class="[{ pinned: col.pin }, col.pin ? `pinned-${col.pinRelation}` : '', col.temp ? 'table-cell-temp' : '']"
                :style="col.temp ? {} : colPinnedStyle(col)">
                <template v-if="!col.temp">
                  <slot :name="`col-${col.key}`" :row="row" :column="col">
                    <slot name="column" :row="row" :column="col">{{ getCellValue(row, col) }}</slot>
                  </slot>
                </template>
              </div>
            </div>
          </template>
          <template #loading>
            <slot name="loading">
              <div class="flex items-center justify-center gap-2">
                <i class="pi pi-spin pi-spinner" />
                Загрузка...
              </div>
            </slot>
          </template>
        </UiVirtualScroll>
      </template>
      <template v-else>
        <UiTableV2Body
          v-show="!loading && filteredData.length"
          ref="bodyOuterRef"
          v-model:items="filteredData"
          :height="containerHeight"
          :loading="fetching"
          :allow-inner-scroll="allowInnerScroll"
          @scroll="onScroll"
          @pre-scroll-bottom="emit('pre-scroll-bottom', true)"
          @scroll-bottom="emit('scroll-bottom', true)">
          <template #default="{ item: row, index }">
            <div class="table-row" :class="rowClass(row, index)" @click="onRowClick(row, $event)">
              <div v-if="selectMode === 'checkbox'" class="table-cell table-cell--select pinned pinned-left left-0">
                <Checkbox
                  :model-value="isSelected(row[localValueKey as KeyOfRow])"
                  binary
                  size="small"
                  @click.stop="onRowCheckboxClickCapture($event)"
                  @update:model-value="(checked) => onRowCheckboxChange(checked, row)" />
              </div>
              <div
                v-for="(col, colIndex) in orderedColsWithTemp"
                :key="`${index}_${col.key}_${colIndex}`"
                class="table-cell"
                :class="[{ pinned: col.pin }, col.pin ? `pinned-${col.pinRelation}` : '', col.temp ? 'table-cell-temp' : '']"
                :style="col.temp ? {} : colPinnedStyle(col)">
                <template v-if="!col.temp">
                  <slot :name="`col-${col.key}`" :row="row" :column="col">
                    <slot name="column" :row="row" :column="col">{{ getCellValue(row, col) }}</slot>
                  </slot>
                </template>
              </div>
            </div>
          </template>
          <template #loading>
            <slot name="loading">
              <div class="flex items-center justify-center gap-2">
                <i class="pi pi-spin pi-spinner" />
                Загрузка...
              </div>
            </slot>
          </template>
        </UiTableV2Body>
      </template>
    </div>
    <Button
      v-if="isShowScrollTopButton && showScrollTop"
      icon="pi pi-arrow-up"
      class="!absolute bottom-1 !right-4 z-50"
      rounded
      severity="contrast"
      @click.stop="handleScrollToTop" />
    <UiTableV2SelectPanel v-if="selectedKeys.length && isShowSelectPanel" :selected="selectedKeys.length" @clear-selected="clearSelection">
      <template #selectCount>
        <slot v-if="$slots.selectCount" name="selectCount" />
        <template v-else>
          {{ selectedKeys.length }}
        </template>
      </template>
      <slot name="selectPanel" :selected="selectedKeys" :count="selectedKeys.length" />
    </UiTableV2SelectPanel>

    <div v-if="isShowSettings" class="settings-drawer">
      <UiTableV2Settings :columns="columnsState" @update:columns="handleUpdateColumns" />
    </div>
  </div>
</template>

<style scoped lang="scss">
.table-v2-wrapper {
  width: 100%;
  overflow: hidden;
  position: relative;
  contain: layout paint;
  &--loading .table-scroll-container {
    overflow: hidden;
  }
}
.table-header-outer {
  overflow: hidden;
  border-bottom: 2px solid var(--surface-border);
  min-height: 40px;
}
.table-header {
  display: block;
  padding-right: 4px;
}
.table-row {
  display: flex;
  align-items: normal;
  width: max-content;

  &.is-selected {
    background: var(--app-accent-soft);
  }
}

.table-cell {
  flex-shrink: 0;
  padding: 8px 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: flex;
}

.table-cell--select {
  width: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.table-scroll-container {
  width: 100%;
  overflow: auto;
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
  .table-row {
    border-bottom: 1px solid var(--surface-border);
  }
  .table-row:hover,
  .table-row:hover .table-cell.pinned,
  .table-row.is-selected .table-cell.pinned {
    background: var(--p-surface-100);
  }

  .table-cell {
    height: calc(var(--row-height, 100%) - 2px);
    align-items: center;
  }
}
.table-cell-temp {
  width: var(--temp-col-width);
  min-width: var(--temp-col-width);
  max-width: var(--temp-col-width);
}
.table-body {
  min-width: 100%;
}
.pinned-left {
  z-index: 10;
  position: sticky;
  background-color: var(--p-surface-0);
}

.pinned-right {
  position: sticky;
  background-color: var(--p-surface-0);
  border-left: 1px solid var(--surface-border);
}

.pinned-left,
.pinned-right {
  will-change: left, right;
}
.settings-drawer {
  position: absolute;
  right: 12px;
  top: 2px;
  z-index: 10;
}

.table-header .pinned-right {
  border-left: none;
}
.resize-handle {
  position: absolute;
  right: 3px;
  top: 0;
  width: 3px;
  height: 100%;
  cursor: col-resize;
  user-select: none;
  z-index: 50;
  transition: background 0.2s;
}

.resize-handle:hover {
  background: rgba(0, 0, 0, 0.1);
}
.app-dark .table-v2-wrapper {
  .table-row:hover,
  .table-row:hover .table-cell.pinned,
  .table-row.is-selected .table-cell.pinned {
    background: var(--p-surface-900);
  }
  .pinned-left,
  .pinned-right {
    background-color: var(--p-surface-900);
  }
}
</style>
