<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts" generic="T extends { id: string | number }">
export interface UiTableColumn {
  key: string;
  header: string;
  width?: string;
  align?: 'left' | 'center' | 'right';
  hideOnMobile?: boolean;
}

const props = withDefaults(
  defineProps<{
    items: T[];
    columns: UiTableColumn[];
    selected?: (string | number)[];
    loading?: boolean;
    fetching?: boolean;
    total?: number;
    loadAll?: boolean;
    selectable?: boolean;
    emptyTitle?: string;
    emptyText?: string;
    emptyIcon?: string;
  }>(),
  {
    selected: () => [],
    loading: false,
    fetching: false,
    total: 0,
    loadAll: true,
    selectable: true,
    emptyTitle: 'Данных нет',
    emptyText: 'Измените условия поиска или добавьте запись',
    emptyIcon: 'pi pi-inbox',
  },
);

const emit = defineEmits<{
  'update:selected': [ids: (string | number)[]];
  'load-more': [];
}>();

const { isBottom } = useScrollBottomListener();

const columnsCount = computed(() => props.columns.length + (props.selectable ? 1 : 0));
const isEmpty = computed(() => !props.loading && !props.items.length);
const shown = computed(() => props.items.length);

const allSelected = computed(() => !!props.items.length && props.items.every((item) => props.selected.includes(item.id)));
const someSelected = computed(() => !allSelected.value && props.items.some((item) => props.selected.includes(item.id)));

const isSelected = (item: T) => props.selected.includes(item.id);

const toggleRow = (item: T, value: boolean) => {
  const next = value ? [...new Set([...props.selected, item.id])] : props.selected.filter((id) => id !== item.id);

  emit('update:selected', next);
};

const toggleAll = (value: boolean) => {
  emit('update:selected', value ? props.items.map((item) => item.id) : []);
};

const canLoadMore = computed(() => !props.loadAll && !props.loading && !props.fetching && shown.value > 0);

watch(isBottom, (value) => {
  if (value && canLoadMore.value) emit('load-more');
});
</script>

<template>
  <div class="app-table-wrapper">
    <div v-if="$slots.toolbar" class="app-table-toolbar">
      <slot name="toolbar" />
    </div>

    <div class="app-table-scroll">
      <table class="app-table">
        <thead>
          <tr>
            <th v-if="selectable" class="app-table__check">
              <Checkbox
                :model-value="allSelected"
                :indeterminate="someSelected"
                binary
                :disabled="!items.length"
                @update:model-value="toggleAll" />
            </th>
            <th
              v-for="column in columns"
              :key="column.key"
              :style="column.width ? { width: column.width } : undefined"
              :class="[`app-table__cell--${column.align || 'left'}`, { 'app-table__cell--mobile-hidden': column.hideOnMobile }]">
              {{ column.header }}
            </th>
          </tr>
        </thead>

        <tbody v-if="loading && !items.length">
          <tr v-for="n in 6" :key="`skeleton-${n}`">
            <td :colspan="columnsCount"><Skeleton height="1.5rem" /></td>
          </tr>
        </tbody>

        <tbody v-else-if="isEmpty">
          <tr>
            <td :colspan="columnsCount">
              <div class="app-empty">
                <span class="app-empty__icon"><i :class="emptyIcon" /></span>
                <span class="app-empty__title">{{ emptyTitle }}</span>
                <span class="app-empty__text">{{ emptyText }}</span>
              </div>
            </td>
          </tr>
        </tbody>

        <tbody v-else>
          <tr v-for="item in items" :key="item.id" :class="{ 'app-table__row--selected': isSelected(item) }">
            <td v-if="selectable" class="app-table__check">
              <Checkbox :model-value="isSelected(item)" binary @update:model-value="(value: boolean) => toggleRow(item, value)" />
            </td>
            <td
              v-for="column in columns"
              :key="column.key"
              :class="[`app-table__cell--${column.align || 'left'}`, { 'app-table__cell--mobile-hidden': column.hideOnMobile }]">
              <slot :name="`cell-${column.key}`" :row="item" />
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="app-table-footer">
      <span class="app-table-footer__result">
        {{ shown ? `Показано ${shown} из ${total || shown}` : '' }}
      </span>

      <div v-if="fetching" class="app-table-footer__loader">
        <i class="pi pi-spin pi-spinner" />
        <span>Загрузка…</span>
      </div>

      <UiButton
        v-else-if="canLoadMore"
        label="Показать ещё"
        icon="pi pi-angle-down"
        severity="secondary"
        outlined
        size="small"
        @click="emit('load-more')" />
    </div>

    <UiTableSelectPanel :count="selected.length" @clear="emit('update:selected', [])">
      <slot name="bulk" />
    </UiTableSelectPanel>
  </div>
</template>
