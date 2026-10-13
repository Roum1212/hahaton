<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts" generic="T extends Object extends { id: any }">
export interface IOnLoadEvent {
  rows: number;
  search?: string;
}
type OptionType = T & { id: any };

const emits = defineEmits(['update:modelValue', 'scroll-bottom', 'search']);

const props = withDefaults(
  defineProps<{
    modelValue?: any | any[];
    searchQuery?: string;
    filter?: boolean;
    loading?: boolean;
    fetching?: boolean;
    display?: 'chip' | 'list' | 'none';
    selectMode?: 'checkbox' | 'item';
    multiple?: boolean;
    options?: OptionType[];
    optionLabel?: string;
    optionValue?: string;
    filterValues?: string[];
    height?: string;
    class?: string;
  }>(),
  {
    modelValue: undefined,
    searchQuery: undefined,
    filter: false,
    loading: false,
    fetching: false,
    display: 'list',
    multiple: false,
    selectMode: 'checkbox',
    options: () => [],
    optionLabel: undefined,
    optionValue: undefined,
    filterValues: () => [],
    height: '100%',
    class: '',
  },
);
const chipsPopoverRef = ref();
const { debounce, breakpoints } = useApp();

const searchQuery = ref<string | undefined>(props.searchQuery);
const internalValue = ref<any[]>(props.multiple ? props.modelValue : [props.modelValue]);

const containerRef = ref<HTMLElement>();
const headerRef = ref<HTMLElement>();
const optionsListRef = ref<HTMLElement>();
const containerHeight = ref(0);
const headerHeight = ref(0);
const optionsListHeight = ref(0);
const lastClickWasShift = ref(false);
const lastSelectedIndex = ref<number | null>(null);
const labelSliceCount = computed(() => (breakpoints.md.value ? 3 : 0));

const optionList = computed(() => {
  const options = clearDuplicates(props.options);
  if (!searchQuery.value?.length) return options;

  return options.filter((option: OptionType) => {
    return props.filterValues.some((path: any) => {
      const val = getValueByPath(option, path);
      return String(val)
        .toLowerCase()
        .includes(searchQuery.value?.toLowerCase() ?? '');
    });
  });
});
const optionsDictionary = computed(() => generateOptionsDictionary(clearDuplicates(props.options)));

function generateOptionsDictionary(options: OptionType[]) {
  return listToMap(options, (el) => generateKey(el));
}
function generateKey(el: OptionType) {
  return props.optionValue
    ? getValueByPath(el, props.optionValue)
    : JSON.stringify(el)
        .replace(/["\\\s]/g, '')
        .replace(/[{}]/g, '')
        .replace(/,/g, ':');
}
function prepareInternalValue(): any[] {
  return Array.isArray(props.modelValue) ? props.modelValue : isValuable(props.modelValue) ? [props.modelValue] : [];
}
function getItemValue(item: any) {
  return props.optionValue ? getValueByPath(item, props.optionValue) : item;
}
function getValueByPath(obj: any, path?: string): any {
  if (!path) return obj;
  return path.split('.').reduce((acc, key) => acc?.[key], obj);
}
function getSelectedLabel(item: any): string | OptionType | undefined {
  const key = props.optionValue ? item : generateKey(item);
  const option = optionsDictionary.value[key];
  return props.optionLabel ? option?.[props.optionLabel] : option;
}
function getLabel(item: any): string | OptionType | undefined {
  const key = generateKey(item);
  const option = optionsDictionary.value[key];
  return props.optionLabel ? option?.[props.optionLabel] : option;
}

function getValue(value: any | any[]): string[] {
  const val = Array.isArray(value) ? value : [value];
  return val.map((el) => (typeof el === 'object' ? JSON.stringify(el) : el));
}

function deepEqual(a: any, b: any): boolean {
  if (a === b) return true;
  if (a === null || b === null || a === undefined || b === undefined) return false;

  return JSON.stringify(a) === JSON.stringify(b);
}

async function handleSelect(item: any) {
  const selectedValue = getItemValue(item);

  const currentIndex = optionList.value.findIndex((el) => deepEqual(getItemValue(el), selectedValue));

  if (props.multiple && lastClickWasShift.value && lastSelectedIndex.value !== null) {
    const start = Math.min(lastSelectedIndex.value, currentIndex);
    const end = Math.max(lastSelectedIndex.value, currentIndex);
    const rangeItems = optionList.value.slice(start, end + 1);

    const newValues = [...internalValue.value];

    rangeItems.forEach((rangeItem) => {
      const val = getItemValue(rangeItem);

      if (!newValues.some((v) => deepEqual(v, val))) {
        newValues.push(val);
      }
    });

    internalValue.value = newValues;
  } else {
    const index = internalValue.value?.findIndex((val) => deepEqual(val, selectedValue));

    if (props.multiple) {
      if (index >= 0) {
        internalValue.value.splice(index, 1);
        internalValue.value = internalValue.value.filter((el: any) => !deepEqual(el, selectedValue));
      } else {
        internalValue.value.push(selectedValue);
      }
    } else {
      internalValue.value = [selectedValue];
    }
  }

  lastSelectedIndex.value = currentIndex;

  await nextTick();

  containerHeight.value = containerRef.value?.closest('.p-drawer-content')?.getBoundingClientRect().height || 0;
  headerHeight.value = headerRef.value?.getBoundingClientRect().height || 0;
  optionsListHeight.value = containerHeight.value - headerHeight.value - 40;

  const entities = internalValue.value.map((v) => optionsDictionary.value[v]);

  emits('update:modelValue', props.multiple ? internalValue.value : internalValue.value?.[0], props.multiple ? entities : entities?.[0]);
}

const handleScroll = debounce(async () => {
  if (!optionsListRef.value) return;
  emits('scroll-bottom');
}, 100);

function handleRemove(index: number) {
  if (!props.multiple) {
    internalValue.value = [];
    emits('update:modelValue', undefined);
  } else {
    internalValue.value.splice(index, 1);
    emits('update:modelValue', internalValue.value);
  }
}

function clearDuplicates(items: OptionType[]): OptionType[] {
  const reduceItems = items.reduce((acc: Record<any, any>, item) => {
    acc[item.id] = item;
    return acc;
  }, {});

  return Object.values(reduceItems);
}
const displayValues = computed(() => getValue(internalValue.value));

function isSelected(item: any): boolean {
  const itemValue = getValueByPath(item, props.optionValue);
  return internalValue.value.some((val: any) => deepEqual(val, itemValue));
}
function handleOpenChips(e: Event) {
  chipsPopoverRef.value?.toggle(e);
}
function onRowCheckboxClickCapture(e: MouseEvent) {
  lastClickWasShift.value = e.shiftKey === true;
}

const handleClearAll = debounce(() => {
  if (!Array.isArray(internalValue.value)) {
    emits('update:modelValue', undefined, undefined);
  } else {
    emits('update:modelValue', [], []);
  }
}, 100);

onMounted(async () => {
  containerHeight.value = containerRef.value?.closest('.p-drawer-content')?.getBoundingClientRect().height || 0;
  headerHeight.value = headerRef.value?.getBoundingClientRect().height || 0;
  optionsListHeight.value = containerHeight.value - headerHeight.value - 40;
});

watch(
  () => props.modelValue,
  async () => {
    internalValue.value = prepareInternalValue() as any;
  },
);
</script>
<template>
  <div ref="containerRef">
    <div v-if="filter || props.display !== 'none'" ref="headerRef" class="flex flex-col gap-2 mb-4">
      <div v-if="filter" class="flex items-center gap-2 w-full">
        <UiSearch v-model="searchQuery" @update:model-value="emits('search', searchQuery)" />
      </div>

      <div v-if="display === 'chip'" class="flex gap-1 items-center flex-wrap">
        <Button
          v-show="displayValues.slice(labelSliceCount).length"
          severity="secondary"
          icon="pi pi-list"
          size="small"
          :label="`+${displayValues.slice(labelSliceCount).length}`"
          class="rounded-2xl"
          @click.stop="handleOpenChips" />
        <Popover v-if="displayValues.slice(labelSliceCount).length" ref="chipsPopoverRef" class="w-[250px]">
          <UiScrollPanel :style="{ height: '200px' }">
            <div class="flex gap-1 items-center flex-wrap">
              <template v-for="(item, i) in displayValues.slice(labelSliceCount)" :key="`chip_${item}`">
                <Chip v-if="getSelectedLabel(item)" :label="getSelectedLabel(item)" removable class="text-xs" @remove="handleRemove(i)" />
                <Skeleton v-else width="5rem" height="18px" />
              </template>
            </div>
          </UiScrollPanel>
        </Popover>
        <template v-for="(item, i) in displayValues.slice(0, labelSliceCount)" :key="`chip_${item}`">
          <Chip v-if="getSelectedLabel(item)" :label="getSelectedLabel(item)" removable class="text-xs" @remove="handleRemove(i)" />
          <Skeleton v-else width="5rem" height="18px" />
        </template>
        <Button
          v-if="displayValues.slice(labelSliceCount).length"
          label="Очистить все"
          class="text-xs rounded-2xl"
          severity="secondary"
          size="small"
          rounded
          @click.stop="handleClearAll" />
      </div>
      <template v-else-if="props.display === 'list'">
        <template v-for="(item, i) in displayValues" :key="`chip_${i}`">
          <template v-if="getSelectedLabel(item)">
            {{ getSelectedLabel(item) }}
          </template>
          <Skeleton v-else width="5rem" height="18px" />
          <span v-if="i < displayValues.length - 1">, </span>
        </template>
      </template>
    </div>

    <UiScrollPanel ref="optionsListRef" :style="{ 'min-height': '300px', height: props.height || '500px' }" @scroll-bottom="handleScroll">
      <div v-if="loading" class="flex justify-center py-2">
        <slot v-if="$slots.loading" name="loading" />
        <i v-else class="pi pi-spin pi-spinner" style="font-size: 2rem" />
      </div>

      <div
        v-for="item in optionList"
        v-else
        :key="getLabel(item)"
        class="cursor-pointer py-1 px-2 rounded-xl"
        :class="isSelected(item) && props.selectMode === 'item' ? 'bg-surface-100 font-semibold' : ''">
        <div class="flex gap-2 items-center" @click="selectMode === 'item' ? handleSelect(item) : $event">
          <Checkbox
            v-if="selectMode === 'checkbox'"
            :model-value="isSelected(item)"
            binary
            @click.stop="onRowCheckboxClickCapture($event)"
            @update:model-value="handleSelect(item)" />
          <slot v-if="$slots.option" name="option" :option="{ label: getLabel(item), item }" />
          <span v-else>{{ getLabel(item) }}</span>
        </div>
      </div>

      <div v-show="!loading && !optionList.length">
        <div class="flex flex-col items-center justify-center gap-2 py-5">
          <span>{{ searchQuery?.length ? 'По вашему запросу ничего не найдено' : 'Список пуст' }}</span>
        </div>
      </div>
      <div v-if="fetching">
        <slot v-if="$slots.fetching" name="fetching" />
        <div v-else class="flex justify-center py-2">
          <i class="pi pi-spin pi-spinner" />
        </div>
      </div>
    </UiScrollPanel>
  </div>
</template>
