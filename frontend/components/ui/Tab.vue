<script setup lang="ts">
import type { MenuItem } from 'primevue/menuitem';
import Tab from 'primevue/tab';
import TabList from 'primevue/tablist';
import Tabs from 'primevue/tabs';
const emits = defineEmits(['update:modelValue']);
const props = withDefaults(
  defineProps<{
    isVisible?: boolean;
    loading?: boolean;
    menu?: MenuItem[];
    modelValue?: string | null;
  }>(),
  {
    isVisible: true,
    loading: false,
    menu: () => [],
    modelValue: null,
  },
);

const { route } = useApp();
const currentRoute = ref<string | number>(props.modelValue || route.path);
const internalValue = computed(() => props.menu.filter((el) => el?.visible !== false));
watch(
  () => currentRoute.value,
  (newValue: string | number) => {
    emits('update:modelValue', newValue);
  },
);
const onUpdated = (value: string | number) => {
  const itemMenu = props.menu.find((el: MenuItem) => el.value === value);
  if (itemMenu?.command) {
    itemMenu.command({
      originalEvent: new Event('click'),
      item: itemMenu,
    });
  } else {
    emits('update:modelValue', value);
  }
};
</script>
<template>
  <Tabs v-if="isVisible" :value="currentRoute" class="!bg-transparent" :scrollable="true" @update:value="onUpdated">
    <TabList class="!bg-transparent">
      <Tab v-for="tab in internalValue" :key="tab.value || tab.path" :value="tab.value || tab.path" class="bg-transparent !pt-2 !pb-3">
        <div class="flex gap-3 justify-center">
          <div class="flex items-center gap-1">
            <span v-if="tab?.icon" :class="tab.icon" />
            <span class="text-center">{{ tab.label }}</span>
          </div>
          <span v-if="tab.badge" class="rounded-2xl font-bold flex justify-center items-center text-xs bg-blue-100 text-blue-800 px-4 py-1">
            {{ tab.badge }}
          </span>
        </div>
      </Tab>
    </TabList>
  </Tabs>
</template>
<style>
.p-tablist-tab-list {
  background-color: transparent !important;
}
.p-tablist-nav-button {
  background-color: var(--code-background) !important;
  box-shadow: none !important;
}
</style>
