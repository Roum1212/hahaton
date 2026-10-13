<script setup lang="ts">
import type { TableColumnType } from './index.vue';

const emits = defineEmits(['update:columns', 'reset-settings']);
const props = defineProps<{
  columns: TableColumnType[];
}>();
const dialogVisible = ref<boolean>(false);
const searchQuery = ref<string | undefined>(undefined);
const columnList = ref<TableColumnType[]>([]);
const filterColumnsList = computed(() =>
  columnList.value.filter((el: TableColumnType) =>
    searchQuery.value ? el.header?.toLowerCase().includes(searchQuery.value?.toLowerCase()) : true,
  ),
);
const visibleColumns = computed<TableColumnType[]>(() => columnList.value.filter((el: TableColumnType) => el.visible));
const sortableColumns = ref<TableColumnType[]>([]);

const handleOpenDialog = () => {
  columnList.value = [];
  for (const col of props.columns) {
    if (col.technical) continue;
    columnList.value.push({ ...col });
  }
  sortableColumns.value = visibleColumns.value;
  dialogVisible.value = true;
};

const handleChangeVisible = (col: TableColumnType, visible: boolean) => {
  const updateIndex = sortableColumns.value.findIndex((el) => el.key === col.key);

  if (!visible && updateIndex !== -1) {
    sortableColumns.value.splice(updateIndex, 1);
  } else {
    sortableColumns.value.push({ ...col, visible: true });
  }
};
const handleHideColumn = (key: string) => {
  const column = columnList.value.find((el: TableColumnType) => el.key === key);
  if (column) column.visible = false;
  sortableColumns.value = visibleColumns.value.filter((c) => c.visible);
};
const handleCancel = () => {
  dialogVisible.value = false;
};
const handleSave = () => {
  emits('update:columns', columnList.value);
  dialogVisible.value = false;
};
const handleUpdateOrder = () => {
  sortableColumns.value.forEach((col, index) => {
    col.order = index + 1;
  });
  const sortableColumnsMap = listToMap(sortableColumns.value, (c) => c.key);
  columnList.value = columnList.value.map((c) => {
    const s = sortableColumnsMap[c.key];
    if (s) c.order = s.order;
    return c;
  });
};
</script>
<template>
  <Button icon="pi pi-cog" severity="secondary" size="small" @click="handleOpenDialog" />
  <Dialog
    v-model:visible="dialogVisible"
    header="Настройка таблицы"
    :draggable="false"
    :style="{ 'max-width': '800px', width: '100%' }"
    :modal="true">
    <slot v-if="$slots.body" name="body" />
    <template v-else>
      <div class="grid grid-cols-2 gap-4">
        <div class="grid col-span-1">
          <div class="flex flex-col gap-2">
            <UiSearch v-model="searchQuery" />
            <UiScrollPanel :style="{ height: '300px' }">
              <div class="flex flex-col gap-2">
                <template v-for="column in filterColumnsList" :key="column.field">
                  <div class="flex gap-4 items-center">
                    <Checkbox v-model="column.visible" binary @update:model-value="(v) => handleChangeVisible(column, v)" />
                    {{ column.header }}
                  </div>
                </template>
              </div>
            </UiScrollPanel>
          </div>
        </div>
        <div class="col-span-1">
          <div class="bg-surface-50 dark:bg-surface-800 rounded-lg py-2 px-4">
            <div class="text-lg mb-4">Выбрано {{ sortableColumns.length }} колонок</div>
            <UiScrollPanel :style="{ height: '300px' }">
              <div :list="sortableColumns" :animation="200" ghost-class="ghost-card" class="flex flex-col gap-2" :dragging="false">
                <draggable
                  :list="sortableColumns"
                  :animation="200"
                  ghost-class="ghost-card"
                  :class="`grid gap-2 grid-cols-${columns}`"
                  :dragging="false"
                  @update="handleUpdateOrder">
                  <div
                    v-for="column in sortableColumns"
                    :key="column.key"
                    class="flex justify-between border dark:border-surface-900 rounded-md bg-surface-0 dark:bg-surface-800 py-2 px-[12px]">
                    <div class="flex items-center gap-4">
                      <i class="pi pi-ellipsis-v text-surface-500 hidden" />
                      <span>{{ column.header || column.key }}</span>
                    </div>
                    <Button
                      text
                      severity="contrast"
                      size="small"
                      icon="pi pi-times text-surface-500"
                      @click="handleHideColumn(column.key)" />
                  </div>
                </draggable>
              </div>
            </UiScrollPanel>
          </div>
        </div>
      </div>
    </template>
    <template v-if="!$slots.body" #footer>
      <div class="flex justify-between gap-2 w-full items-center">
        <Button label="Сбросить настройки" text severity="secondary" size="small" outlined @click="$emit('reset-settings', true)" />
        <div class="flex gap-2">
          <Button label="Отменить" severity="secondary" outlined @click="handleCancel" />
          <Button label="Сохранить" @click="handleSave" />
        </div>
      </div>
    </template>
  </Dialog>
</template>
