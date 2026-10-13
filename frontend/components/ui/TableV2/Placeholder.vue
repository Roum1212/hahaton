<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import type { CSSProperties } from 'vue';
import type { TableColumnType } from './index.vue';

const props = withDefaults(
  defineProps<{
    selectMode: 'checkbox' | 'fill' | undefined;
    columns: TableColumnType[];
    rows?: number;
  }>(),
  {
    rows: 20,
  },
);
const selectCol = { key: '__select', width: 60, pin: true, pinRelation: 'left' };
const visibleCols = computed<TableColumnType[]>(() => props.columns.filter((c) => c.visible));

function colPinnedStyle(col: TableColumnType, _: number) {
  const base: CSSProperties = {
    width: (col.width || 120) + 'px',
    minWidth: (col.width || 120) + 'px',
    maxWidth: (col.width || 120) + 'px',
  };
  if (!col.pin) return base;
  const leftOffset = computeLeftOffset(col);
  const rightOffset = computeRightOffset(col);
  if (col.pinRelation === 'left') base.left = leftOffset + 'px';
  else base.right = rightOffset + 'px';
  base.position = 'sticky';
  base.zIndex = 6;
  return base;
}
function computeLeftOffset(col: TableColumnType) {
  let offset = 0;
  offset += selectCol.width;
  for (const c of visibleCols.value) {
    if (c.key === col.key) break;
    if (c.pin && c.pinRelation === 'left') offset += c.width || 120;
  }
  return offset;
}

function computeRightOffset(col: TableColumnType) {
  let offset = 0;
  const rightPinned = visibleCols.value.filter((c) => c.pin && c.pinRelation === 'right');
  for (const c of rightPinned.slice().reverse()) {
    if (c.key === col.key) break;
    offset += c.width || 120;
  }
  return offset;
}
</script>

<template>
  <div v-for="(row, rowIndex) of rows" :key="rowIndex" class="table-row border-b">
    <div v-if="selectMode === 'checkbox'" class="table-cell">
      <Skeleton height="17px" width="17px" />
    </div>
    <div
      v-for="(col, colIndex) in visibleCols"
      :key="col.key"
      class="table-cell"
      :style="col.temp ? {} : colPinnedStyle(col, colIndex)"
      :class="[{ pinned: col.pin }, col.pin ? `pinned-${col.pinRelation}` : '', col.temp ? 'table-cell-temp' : '']">
      <Skeleton v-if="!col.temp" height="14px" width="10rem" />
    </div>
  </div>
</template>
<style scoped>
.table-cell-temp {
  width: var(--temp-col-width);
  min-width: var(--temp-col-width);
  max-width: var(--temp-col-width);
}
.table-row {
  display: flex;
  align-items: center;
  width: max-content;
  min-width: 100%;
}
.table-cell {
  flex-shrink: 0;
  padding: 8px 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: flex;
  height: -webkit-fill-available;
}
.pinned-right {
  position: sticky;
  background-color: #fff;
  border-left: 1px solid var(--surface-border);
}
.pinned-left {
  position: sticky;
  background-color: #fff;
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
