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
    menu: MenuItem[];
    modelValue?: string | null;
  }>(),
  {
    isVisible: true,
    loading: false,
    modelValue: null,
  },
);
const { route } = useApp();
const currentRoute = ref(route.path);

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
  }
};
</script>
<template>
  <Tabs :value="currentRoute" class="bg-transparent no-border" :scrollable="true" @update:value="onUpdated">
    <TabList class="bg-transparent flex items-center gap-3 steps">
      <Tab
        v-for="(tab, index) in menu"
        :key="index"
        :unstyled="true"
        :value="tab.value"
        class="flex items-center gap-3 mr-3 steps_item px-2 shrink-0"
        :class="{ last: menu.length - 1 === index }">
        <div class="flex gap-3 justify-center items-center">
          <div
            class="flex items-center text-sm justify-center shrink-0 w-7 h-7 rounded-full bg-surface-200 dark:bg-surface-600"
            :class="{ '!bg-primary text-white': currentRoute === tab.value }">
            {{ tab.number }}
          </div>
          <span class="whitespace-nowrap"> {{ tab.label }} </span>
        </div>
      </Tab>
    </TabList>
  </Tabs>
</template>
<style></style>

<style lang="scss">
.no-border .p-tablist-tab-list {
  border: none !important;
}
.bg-transparent,
.p-tablist-tab-list {
  background-color: transparent !important;
  box-shadow: none !important;
}
.p-tablist-nav-button {
  box-shadow: none !important;
  background-color: var(--code-background) !important;
}
.steps_item {
  &:after {
    content: '';
    display: block;
    width: 25px;
    height: 1px;
    background: #e5e7eb;
  }
  &.last:after {
    content: none;
  }
}
</style>
