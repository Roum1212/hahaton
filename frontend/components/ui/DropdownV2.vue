<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts" generic="T extends Object extends { id: any }">
export interface IOnLoadEvent {
  rows: number;
  search?: string;
}
type OptionType = T & { id: any };
defineExpose({
  close,
});
const emits = defineEmits(['update:modelValue', 'add-input-option', 'add-option', 'updated', 'close', 'open', 'scroll-bottom', 'search']);

interface Props {
  modelValue?: any | any[];
  id?: string;
  size?: 'small' | 'medium' | 'large';
  placeholder?: string;
  disabled?: boolean;
  filter?: boolean;
  loading?: boolean;
  fetching?: boolean;
  display?: 'chip' | 'list';
  multiple?: boolean;
  editable?: boolean;
  options?: OptionType[];
  optionLabel?: string;
  optionValue?: string;
  filterValues?: string[];
  emptyMessage?: string;
  emptyFilterMessage?: string;
  showClear?: boolean;
  classList?: string;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: undefined,
  id: '',
  size: 'medium',
  placeholder: 'Выберите',
  disabled: false,
  filter: false,
  loading: false,
  fetching: false,
  display: 'list',
  multiple: false,
  editable: false,
  options: () => [],
  optionLabel: undefined,
  optionValue: undefined,
  filterValues: () => [],
  emptyMessage: 'Список пуст',
  emptyFilterMessage: 'По вашему запросу ничего не найдено',
  showClear: true,
  classList: '',
});
interface IComponentState {
  isOpen: boolean;
  hoverOption: number | undefined;
  editInputValue: string;
  editInputFocus: boolean;
}
const localState = reactive<IComponentState>({
  isOpen: false,
  hoverOption: undefined as number | undefined,
  editInputValue: '',
  editInputFocus: false,
});

const optionsListRef = ref<HTMLElement>();
const dropdownRef = ref<HTMLElement>();
const dropdownValueRef = ref<HTMLElement>();

const sizeClass = computed(() => {
  switch (props.size) {
    case 'small':
      return 'p-select-sm p-inputfield-sm';
    case 'large':
      return 'p-select-lg p-inputfield-lg';
    default:
      return '';
  }
});

const dropdownPosition = ref({ top: 0, left: 0 });
const dropdownObserver = ref<MutationObserver | null>(null);

const { debounce } = useApp();

const searchQuery = ref<string | undefined>(undefined);
const internalValue = ref<any[]>(prepareInternalValue());

const isShowPlaceholder = computed(() => !internalValue.value.length && !localState.editInputFocus);
const isPositionReady = ref(false);
const lastSelectedIndex = ref<number | null>(null);

const displayValues = computed(() => getValue(internalValue.value));

const isShowClear = computed(
  () => props.showClear && ((props.multiple && internalValue.value.length > 0) || (!props.multiple && internalValue.value)),
);

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

function isSelected(item: any): boolean {
  const itemValue = getValueByPath(item, props.optionValue);
  return internalValue.value.some((val: any) => deepEqual(val, itemValue));
}

async function handleSelect(item: any, event?: MouseEvent) {
  const selectedValue = getItemValue(item);

  const currentIndex = optionList.value.findIndex((el) => deepEqual(getItemValue(el), selectedValue));

  if (props.multiple && event?.shiftKey && lastSelectedIndex.value !== null) {
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
  const entities = internalValue.value.map((v) => optionsDictionary.value[v]);

  emits('update:modelValue', props.multiple ? internalValue.value : internalValue.value?.[0], props.multiple ? entities : entities?.[0]);
  emits('updated', true);
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
  emits('updated', true);
}
function handleEditInputFocus(event: boolean) {
  localState.editInputFocus = event;
}
function handleEditInputEnter() {
  const value = localState.editInputValue.trim();
  localState.editInputFocus = false;
  localState.editInputValue = '';
  emits('add-input-option', value);
}

function updateDropdownPosition() {
  if (dropdownValueRef.value) {
    const rect = dropdownValueRef.value.getBoundingClientRect();
    dropdownPosition.value = {
      top: rect.bottom + window.scrollY + 5,
      left: rect.left + window.scrollX,
    };
    isPositionReady.value = true;
  }
}
async function close() {
  localState.isOpen = false;
}
async function toggleDropdown() {
  localState.isOpen = !localState.isOpen;
  if (!localState.isOpen) {
    searchQuery.value = '';
    isPositionReady.value = false;
    emits('close', true);
  }
  if (localState.isOpen) {
    isPositionReady.value = false;

    emits('open', true);
    await nextTick();
    updateDropdownPosition();
  }
}

function handleClickOutside(event: MouseEvent) {
  nextTick(() => {
    if (!dropdownRef.value || !dropdownValueRef.value) return;

    const path = event.composedPath();

    if (!path.includes(dropdownRef.value) && !path.includes(dropdownValueRef.value)) {
      localState.isOpen = false;
    }
  });
}
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function clearDuplicates(items: any[]): any[] {
  const reduceItems = items.reduce((acc, item) => {
    acc[JSON.stringify(item)] = item;
    return acc;
  }, {});

  return Object.values(reduceItems);
}
function handleClear() {
  internalValue.value = [];
  if (props.multiple) {
    emits('update:modelValue', [], []);
  } else {
    emits('update:modelValue', undefined, undefined);
  }
  emits('updated', true);
}

onMounted(async () => {
  document.addEventListener('click', handleClickOutside);
  window.addEventListener('scroll', updateDropdownPosition);
  window.addEventListener('resize', updateDropdownPosition);
});

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside);
  window.removeEventListener('scroll', updateDropdownPosition);
  window.removeEventListener('resize', updateDropdownPosition);
  if (dropdownObserver.value) {
    dropdownObserver.value.disconnect();
  }
});

watch(
  () => props.modelValue,
  async () => {
    internalValue.value = prepareInternalValue();
  },
);
</script>
<template>
  <div class="relative" :class="classList">
    <div
      :id="props.id"
      class="p-select p-component p-inputwrapper w-full min-h-[30px] !relative"
      :class="`${sizeClass} ${$props.classList} ${props.disabled ? 'p-disabled !bg-surface-300 !border-surface-400' : ''}`"
      @click.stop="toggleDropdown">
      <span ref="dropdownValueRef" class="p-select-label">
        <template v-if="isShowPlaceholder">
          <slot v-if="$slots.placeholder" name="placeholder" />
          <span v-else class="dark:text-gray-400 text-gray-500">
            {{ placeholder }}
          </span>
        </template>

        <template v-else>
          <slot v-if="$slots.value" name="value" :value="displayValues" />
          <template v-else>
            <div v-if="props.display === 'chip'" class="flex gap-1 items-center">
              <template v-for="(item, i) in displayValues" :key="`chip_${i}`">
                <Tag v-if="getSelectedLabel(item)" class="h-6" severity="secondary" :value="item" removable>
                  {{ getSelectedLabel(item) }}
                  <i class="pi pi-times text-xs cursor-pointer" @click.stop="handleRemove(i)" />
                </Tag>
                <Skeleton v-else width="5rem" height="18px" />
              </template>
            </div>
            <template v-else>
              <template v-for="(item, i) in displayValues" :key="`chip_${i}`">
                <template v-if="getSelectedLabel(item)">
                  {{ getSelectedLabel(item) }}
                </template>
                <Skeleton v-else width="5rem" height="18px" />
                <span v-if="i < displayValues.length - 1">, </span>
              </template>
            </template>
          </template>
        </template>
        <input
          v-if="props.editable"
          v-model="localState.editInputValue"
          class="ml-3 unsyled-input cursor-pointer"
          @keyup.enter="handleEditInputEnter"
          @focusin="handleEditInputFocus(false)"
          @focusout="handleEditInputFocus(true)" />
      </span>
      <InputIcon
        v-if="isShowClear && internalValue.length"
        class="pi pi-times cursor-pointer absolute right-10 top-5"
        @click.stop="handleClear" />
      <div class="p-select-dropdown">
        <svg
          width="14"
          height="14"
          viewBox="0 0 14 14"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          class="p-icon p-select-dropdown-icon">
          <path
            d="M7.01744 10.398C6.91269 10.3985 6.8089 10.378 6.71215 10.3379C6.61541 10.2977 6.52766 10.2386 6.45405 10.1641L1.13907 4.84913C1.03306 4.69404 0.985221 4.5065 1.00399 4.31958C1.02276 4.13266 1.10693 3.95838 1.24166 3.82747C1.37639 3.69655 1.55301 3.61742 1.74039 3.60402C1.92777 3.59062 2.11386 3.64382 2.26584 3.75424L7.01744 8.47394L11.769 3.75424C11.9189 3.65709 12.097 3.61306 12.2748 3.62921C12.4527 3.64535 12.6199 3.72073 12.7498 3.84328C12.8797 3.96582 12.9647 4.12842 12.9912 4.30502C13.0177 4.48162 12.9841 4.662 12.8958 4.81724L7.58083 10.1322C7.50996 10.2125 7.42344 10.2775 7.32656 10.3232C7.22968 10.3689 7.12449 10.3944 7.01744 10.398Z"
            fill="currentColor"></path>
        </svg>
      </div>
    </div>
    <Teleport to="body">
      <Transition name="dropdown">
        <div
          v-show="localState.isOpen && isPositionReady"
          ref="dropdownRef"
          class="p-select-overlay p-component"
          :style="{
            zIndex: 9999,
            top: `${dropdownPosition.top}px`,
            left: `${dropdownPosition.left}px`,
            position: 'absolute',
            minWidth: `${dropdownValueRef?.offsetWidth}px`,
          }">
          <slot v-if="$slots.header" name="header"></slot>
          <div v-if="props.filter" class="p-select-header" data-pc-section="header">
            <UiSearch v-model="searchQuery" @update:model-value="emits('search', searchQuery)" />
          </div>
          <div
            ref="optionsListRef"
            class="p-select-list-container overflow-y-auto"
            style="max-height: 14rem; will-change: transform"
            @scroll="handleScroll">
            <ul class="p-select-list" role="listbox">
              <li v-if="loading" class="flex justify-center py-2">
                <slot v-if="$slots.loading" name="loading" />
                <i v-else class="pi pi-spin pi-spinner" style="font-size: 2rem" />
              </li>
              <li v-else-if="!optionList.length" class="p-select-empty-message">
                {{ !searchQuery?.length ? emptyMessage : emptyFilterMessage }}
              </li>
              <li
                v-for="(item, i) in optionList"
                v-else
                :key="i"
                class="p-select-option"
                :class="{ 'p-focus': i === localState.hoverOption, 'p-select-option-selected': isSelected(item) }"
                style="overflow: hidden; position: relative"
                @mouseover="localState.hoverOption = i"
                @click="handleSelect(item)">
                <span class="p-select-option-label">
                  <slot v-if="$slots.option" name="option" :item="item" />
                  <template v-else>
                    {{ getLabel(item) }}
                  </template>
                </span>
              </li>
              <li v-if="fetching && !loading" class="p-select-empty-message justify-center flex w-full">
                <slot v-if="$slots.fetching" name="fetching" />
                <i v-else class="pi pi-spin pi-spinner" style="font-size: 1rem" />
              </li>
            </ul>
          </div>
          <slot v-if="$slots.footer" name="footer"></slot>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>
<style lang="css" scoped>
* {
  --p-select-background: var(--p-form-field-background);
  --p-select-disabled-background: var(--p-form-field-disabled-background);
  --p-select-filled-background: var(--p-form-field-filled-background);
  --p-select-filled-hover-background: var(--p-form-field-filled-hover-background);
  --p-select-filled-focus-background: var(--p-form-field-filled-focus-background);
  --p-select-border-color: var(--p-form-field-border-color);
  --p-select-hover-border-color: var(--p-form-field-hover-border-color);
  --p-select-focus-border-color: var(--p-form-field-focus-border-color);
  --p-select-invalid-border-color: var(--p-form-field-invalid-border-color);
  --p-select-color: var(--p-form-field-color);
  --p-select-disabled-color: var(--p-form-field-disabled-color);
  --p-select-placeholder-color: var(--p-form-field-placeholder-color);
  --p-select-invalid-placeholder-color: var(--p-form-field-invalid-placeholder-color);
  --p-select-shadow: var(--p-form-field-shadow);
  --p-select-padding-x: var(--p-form-field-padding-x);
  --p-select-padding-y: var(--p-form-field-padding-y);
  --p-select-border-radius: var(--p-form-field-border-radius);
  --p-select-focus-ring-width: var(--p-form-field-focus-ring-width);
  --p-select-focus-ring-style: var(--p-form-field-focus-ring-style);
  --p-select-focus-ring-color: var(--p-form-field-focus-ring-color);
  --p-select-focus-ring-offset: var(--p-form-field-focus-ring-offset);
  --p-select-focus-ring-shadow: var(--p-form-field-focus-ring-shadow);
  --p-select-transition-duration: var(--p-form-field-transition-duration);
  --p-select-sm-font-size: var(--p-form-field-sm-font-size);
  --p-select-sm-padding-x: var(--p-form-field-sm-padding-x);
  --p-select-sm-padding-y: var(--p-form-field-sm-padding-y);
  --p-select-lg-font-size: var(--p-form-field-lg-font-size);
  --p-select-lg-padding-x: var(--p-form-field-lg-padding-x);
  --p-select-lg-padding-y: var(--p-form-field-lg-padding-y);
  --p-select-dropdown-width: 2.5rem;
  --p-select-dropdown-color: var(--p-form-field-icon-color);
  --p-select-overlay-background: var(--p-overlay-select-background);
  --p-select-overlay-border-color: var(--p-overlay-select-border-color);
  --p-select-overlay-border-radius: var(--p-overlay-select-border-radius);
  --p-select-overlay-color: var(--p-overlay-select-color);
  --p-select-overlay-shadow: var(--p-overlay-select-shadow);
  --p-select-list-padding: var(--p-list-padding);
  --p-select-list-gap: var(--p-list-gap);
  --p-select-list-header-padding: var(--p-list-header-padding);
  --p-select-option-focus-background: var(--p-list-option-focus-background);
  --p-select-option-selected-background: var(--p-list-option-selected-background);
  --p-select-option-selected-focus-background: var(--p-list-option-selected-focus-background);
  --p-select-option-color: var(--p-list-option-color);
  --p-select-option-focus-color: var(--p-list-option-focus-color);
  --p-select-option-selected-color: var(--p-list-option-selected-color);
  --p-select-option-selected-focus-color: var(--p-list-option-selected-focus-color);
  --p-select-option-padding: var(--p-list-option-padding);
  --p-select-option-border-radius: var(--p-list-option-border-radius);
  --p-select-option-group-background: var(--p-list-option-group-background);
  --p-select-option-group-color: var(--p-list-option-group-color);
  --p-select-option-group-font-weight: var(--p-list-option-group-font-weight);
  --p-select-option-group-padding: var(--p-list-option-group-padding);
  --p-select-clear-icon-color: var(--p-form-field-icon-color);
  --p-select-checkmark-color: var(--p-list-option-color);
  --p-select-checkmark-gutter-start: -0.375rem;
  --p-select-checkmark-gutter-end: 0.375rem;
  --p-select-empty-message-padding: var(--p-list-option-padding);
}
.p-select {
  display: inline-flex;
  align-items: center;
  cursor: pointer;
  position: relative;
  user-select: none;
  background: var(--p-select-background);
  border: 1px solid var(--p-select-border-color);
  width: 100%;
  transition:
    background var(--p-select-transition-duration),
    color var(--p-select-transition-duration),
    border-color var(--p-select-transition-duration),
    outline-color var(--p-select-transition-duration),
    box-shadow var(--p-select-transition-duration);
  border-radius: var(--p-select-border-radius);
  outline-color: transparent;
  box-shadow: var(--p-select-shadow);
}
.p-select-dropdown {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: transparent;
  color: var(--p-select-dropdown-color);
  width: var(--p-select-dropdown-width);
  border-start-end-radius: var(--p-select-border-radius);
  border-end-end-radius: var(--p-select-border-radius);
}
.p-select-label-empty {
  overflow: hidden;
  opacity: 0;
}
.p-select-label {
  display: block;
  white-space: nowrap;
  overflow: hidden;
  flex: 1 1 auto;
  width: 1%;
  padding: var(--p-select-padding-y) var(--p-select-padding-x);
  text-overflow: ellipsis;
  cursor: pointer;
  color: var(--p-select-color);
  background: transparent;
  border: 0 none;
  outline: 0 none;
  margin-right: 20px;
}
.p-select-overlay {
  will-change: transform, opacity;
  position: absolute;
  top: 0;
  left: 0;
  background: var(--p-select-overlay-background);
  color: var(--p-select-overlay-color);
  border: 1px solid var(--p-select-overlay-border-color);
  border-radius: var(--p-select-overlay-border-radius);
  box-shadow: var(--p-select-overlay-shadow);
}
.p-select-header {
  padding: var(--p-select-list-header-padding);
}
.p-select-list-container {
  overflow: auto;
}
.p-select-list {
  margin: 0;
  padding: 0;
  list-style-type: none;
  padding: var(--p-select-list-padding);
  gap: var(--p-select-list-gap);
  display: flex;
  flex-direction: column;
}
.p-select-option {
  cursor: pointer;
  font-weight: normal;
  white-space: nowrap;
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  padding: var(--p-select-option-padding);
  border: 0 none;
  color: var(--p-select-option-color);
  background: transparent;
  transition:
    background var(--p-select-transition-duration),
    color var(--p-select-transition-duration),
    border-color var(--p-select-transition-duration),
    box-shadow var(--p-select-transition-duration),
    outline-color var(--p-select-transition-duration);
  border-radius: var(--p-select-option-border-radius);
}
.p-select-option:not(.p-select-option-selected):not(.p-disabled).p-focus {
  background: var(--p-select-option-focus-background);
  color: var(--p-select-option-focus-color);
}
.p-select-option.p-select-option-selected,
.p-select-option.p-select-option-selected.p-focus {
  background: var(--p-select-option-selected-focus-background);
  color: var(--p-select-option-selected-focus-color);
}
.unsyled-input {
  background: transparent;
}
.unsyled-input:focus {
  border: none;
  outline: none;
}

.p-select-sm .p-select-label {
  font-size: var(--p-select-sm-font-size);
  padding-block: var(--p-select-sm-padding-y);
  padding-inline: var(--p-select-sm-padding-x);
}
.p-select-sm .p-select-dropdown .p-icon {
  font-size: var(--p-select-sm-font-size);
  width: var(--p-select-sm-font-size);
  height: var(--p-select-sm-font-size);
}

.p-inputgroup > .relative:not(:last-child) .p-select {
  border-top-right-radius: 0 !important;
  border-bottom-right-radius: 0 !important;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.3s ease;
  transform-origin: top center;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px) scaleY(0.95);
}

.dropdown-enter-to,
.dropdown-leave-from {
  opacity: 1;
  transform: translateY(0) scaleY(1);
}
</style>
