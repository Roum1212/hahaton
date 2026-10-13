<script lang="ts" setup>
import type { SidebarProps } from 'primevue/sidebar';

const props = withDefaults(
  defineProps<{
    size?: 'small' | 'medium' | 'large' | 'full' | 'custom';
    header?: string;
    position?: SidebarProps['position'];
    headerClass?: string;
    footerClass?: string;
  }>(),
  {
    size: 'small',
    position: 'right',
    header: '',
    headerClass: '',
    footerClass: '',
  },
);

const dialogPosition = ref(props.position);
const panelStyles = computed(() => {
  const maxWidth =
    dialogPosition.value === 'full' || props.size === 'custom'
      ? '100%'
      : props.size === 'full'
        ? '100%'
        : props.size === 'large'
          ? '950px'
          : props.size === 'medium'
            ? '550px'
            : '350px';

  return `width: 100%; max-width: ${maxWidth}`;
});
</script>

<template>
  <Drawer v-bind="props" :style="panelStyles" :block-scroll="true" :dismissable="true" :position="dialogPosition" class="p-fluid">
    <template #header>
      <div class="flex justify-between items-center w-full" :class="headerClass">
        <slot v-if="$slots.header" name="header" />
        <div v-else class="p-drawer-title" data-pc-section="title">{{ header }}</div>
      </div>
    </template>
    <slot v-if="$slots.body" name="body" />
    <template v-if="$slots.footer" #footer>
      <div class="flex gap-2 justify-end" :class="footerClass">
        <slot name="footer" />
      </div>
    </template>
  </Drawer>
</template>
<style>
.p-drawer-content {
  overflow-x: hidden;
}
</style>
