<script setup lang="ts">
import type { MenuItem, MenuItemCommandEvent } from 'primevue/menuitem';
import { v4 } from 'uuid';

export interface IControlButton {
  icon?: string | undefined;
  customIcon?: string | undefined;
  label?: string | undefined;
  tooltip?: string;
  loading?: boolean;
  class?: string;
  group?: string;
  type: 'button' | 'overlay';
  severity?: 'secondary' | 'primary' | 'success' | 'info' | 'warning' | 'danger';
  disabled?: boolean;
  hidden?: boolean;
  command: (e?: Event | MenuItemCommandEvent) => void;
}

const props = withDefaults(
  defineProps<{
    buttons?: IControlButton[];
    disabled?: boolean;
    buttonClass?: string;
    listText?: string;
    listIcon?: string;
    listIconPos?: 'left' | 'right';
    text?: boolean;
    outlined?: boolean;
  }>(),
  {
    buttons: () => [],
    disabled: false,
    buttonClass: '',
    listIcon: 'pi pi-ellipsis-h',
    listIconPos: 'left',
    listText: undefined,
    text: false,
    outlined: false,
  },
);

const menuRef = ref();
const menuId = `overlay_menu_${v4()}`;
const buttons = computed<IControlButton[]>(() => props.buttons);
const visibleButtons = computed<IControlButton[]>(() =>
  buttons.value.filter((el: IControlButton) => el.type === 'button' && el.hidden !== true),
);
const overlayButtons = computed(() => buttons.value.filter((el: IControlButton) => el.type === 'overlay' && el.hidden !== true));
const menu = computed<MenuItem[]>(() => {
  const grouped: Record<string, MenuItem> = {};
  const ungrouped: MenuItem[] = [];
  const result: MenuItem[] = [];

  for (const button of overlayButtons.value) {
    if (button.group) {
      if (!grouped[button.group]) {
        grouped[button.group] = {
          label: button.group,
          items: [],
        };
        const data = grouped[button.group];
        if (data) result.push(data);
      }
      (grouped[button.group]!.items as MenuItem[]).push(button);
    } else {
      ungrouped.push(button);
    }
  }

  if (ungrouped.length > 0) {
    if (result.length > 0) {
      result.push({ separator: true });
    }
    result.push(...ungrouped);
  }

  return result;
});

const handleToggleMenu = (event: Event) => {
  menuRef.value?.toggle(event);
};
</script>

<template>
  <div class="flex gap-1 justify-end">
    <Button
      v-for="(button, index) in visibleButtons"
      :key="index"
      v-tooltip.bottom="button.tooltip"
      :text="text"
      :outlined="outlined"
      :severity="button.severity || 'secondary'"
      :icon="button.customIcon ? undefined : button.icon"
      :label="button.label"
      :loading="button.loading"
      size="small"
      :class="`${button?.class} ${buttonClass}`"
      aria-haspopup="true"
      :disabled="disabled || button.disabled"
      :aria-controls="menuId"
      @click.stop="button.command" />

    <Button
      v-if="menu.length"
      v-tooltip.bottom="'Управление'"
      :label="listText"
      :icon="listIcon"
      :disabled="disabled"
      :text="text"
      :outlined="outlined"
      :icon-pos="listIconPos"
      size="small"
      severity="secondary"
      aria-haspopup="true"
      :class="`${buttonClass}`"
      @click.stop="handleToggleMenu" />
    <Menu :id="menuId" ref="menuRef" :model="menu" :popup="true" />
  </div>
</template>
