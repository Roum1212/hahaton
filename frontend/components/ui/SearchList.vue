<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script lang="ts" setup generic="T">
interface SelectOption {
  select?: boolean;
  [key: string]: any;
}

const emit = defineEmits<{
  'update:modelValue': [value: any[]];
  'scroll-bottom': [];
  search: [value: string];
}>();

const props = withDefaults(
  defineProps<{
    modelValue: any[];
    multiple?: boolean;
    options?: any[];
    optionValue?: string;
    optionLabel?: string;
    fetching?: boolean;
    loading?: boolean;
    searchFetching?: boolean;
    badges?: boolean;
    height?: string;
  }>(),
  {
    multiple: false,
    options: () => [],
    optionValue: 'id',
    optionLabel: 'label',
    fetching: false,
    loading: false,
    searchFetching: false,
    badges: false,
    height: '200px',
  },
);

const localValue = ref<any[]>(props.modelValue);
const optionsList = ref<Array<any>>([...props.options]);
const selectOptions = ref<Array<SelectOption>>([]);
const search = ref<string>('');

const filterOptions = (searchTerm: string) => {
  const searchLower = searchTerm.toLowerCase().trim();
  optionsList.value =
    searchLower === ''
      ? [...props.options]
      : props.options.filter((item: any) => item[props.optionLabel].toLowerCase().includes(searchLower));
};

const generateSelectOptions = () => {
  optionsList.value = [...props.options];

  if (props.multiple) {
    selectOptions.value = props.options.map((item: any) => ({
      ...item,
      select: !!localValue.value?.find((el: any) => el[props.optionValue] === item[props.optionValue]),
    }));
  } else {
    selectOptions.value = props.modelValue;
  }
};

const updateLocalValue = () => {
  const selectedItems = selectOptions.value.filter((item: any) => item.select);
  localValue.value = [...localValue.value, ...selectedItems];

  localValue.value = localValue.value.filter(
    (item: any, index: number, self: any[]) => index === self.findIndex((t: any) => t[props.optionValue] === item[props.optionValue]),
  );
};

const changeSelected = (el: any) => {
  if (!localValue.value) return;

  if (props.multiple) {
    const existingItem = localValue.value.find((item: any) => item[props.optionValue] === el[props.optionValue]);

    if (existingItem) {
      localValue.value = localValue.value.map((item: any) =>
        item[props.optionValue] === el[props.optionValue] ? { ...item, select: el.select } : item,
      );
    } else {
      localValue.value.push(el);
    }

    localValue.value = localValue.value.filter((item: any) => item.select !== false);
  } else {
    localValue.value = [el];
  }

  emit('update:modelValue', localValue.value);
};

const removeSelected = (el: any) => {
  if (!localValue.value) return;

  if (props.multiple) {
    selectOptions.value = selectOptions.value.map((item: any) => ({
      ...item,
      select: item[props.optionValue] !== el[props.optionValue] && item.select,
    }));

    localValue.value = localValue.value.filter((item: any) => item[props.optionValue] !== el[props.optionValue]);
  } else {
    selectOptions.value = [];
  }

  emit(
    'update:modelValue',
    localValue.value.filter((item: any) => item.select !== false),
  );
};

onMounted(() => {
  generateSelectOptions();
});

watch(search, (newValue: string) => {
  if (!props.searchFetching) {
    filterOptions(newValue);
  } else {
    emit('search', newValue);
  }
});

watch(() => props.options, generateSelectOptions);

watch(() => selectOptions.value, updateLocalValue, { deep: true });

watch(
  () => props.modelValue,
  () => {
    localValue.value = props.modelValue;
    generateSelectOptions();
  },
);

const scrollPanel = () => emit('scroll-bottom');
</script>

<template>
  <div class="flex items-center gap-2 w-full">
    <UiSearch v-model="search" />
  </div>

  <div v-if="props.badges" class="flex flex-wrap gap-1 my-2">
    <template v-if="loading">
      <Skeleton v-for="n in 2" :key="n" width="57px" height="18px" border-radius="16px" />
    </template>
    <template v-else>
      <Chip
        v-for="(item, index) in localValue"
        :key="index"
        class="text-xs"
        :value="item"
        :label="item[props.optionLabel]"
        removable
        @remove="removeSelected(item)" />
    </template>
  </div>
  <div class="p-2">
    <UiScrollPanel :style="{ width: '100%', height: props.height }" @scroll-bottom="scrollPanel">
      <template v-if="loading">
        <div class="flex flex-col gap-2">
          <div v-for="n in 10" :key="n" class="flex gap-2 items-center text-sm">
            <Skeleton width="17px" height="17px" />
            <Skeleton width="11rem" height="17px" />
          </div>
        </div>
      </template>
      <template v-else>
        <div class="flex flex-col gap-2">
          <template v-for="(item, index) in optionsList" :key="index">
            <div v-if="multiple" class="flex gap-2 items-center">
              <Checkbox
                v-if="selectOptions[index]?.select !== undefined"
                v-model="selectOptions[index].select"
                :binary="true"
                :input-id="`${item[optionValue]}_${item[optionLabel]}`"
                @change="changeSelected(selectOptions[index])" />
              <slot v-if="$slots.option" name="option" :data="item" />
              <label v-else :for="`${item[optionValue]}_${item[optionLabel]}`">
                {{ item[optionLabel] }}
              </label>
            </div>
            <div v-else class="flex gap-2 items-center text-sm">
              <RadioButton
                v-model="selectOptions[0]"
                :input-id="`${item[optionValue]}_${item[optionLabel]}`"
                :value="item"
                @change="changeSelected(item)" />
              <slot v-if="$slots.option" name="option" :data="item" />
              <label v-else :for="`${item[optionValue]}_${item[optionLabel]}`" class="ml-2">
                {{ item[optionLabel] }}
              </label>
            </div>
          </template>
        </div>
        <div v-if="fetching" class="flex justify-center">
          <ProgressSpinner style="width: 30px; height: 30px" stroke-width="2" />
        </div>
      </template>
    </UiScrollPanel>
  </div>
</template>
