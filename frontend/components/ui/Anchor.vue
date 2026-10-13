<script lang="ts" setup>
import type { MenuItem } from 'primevue/menuitem';

export interface IAnchorMenu {
  label: string;
  slug: string;
  items?: MenuItem[];
}
const emits = defineEmits(['update:activeSlug']);
const props = withDefaults(
  defineProps<{
    menu: IAnchorMenu[];
    activeSlug?: string;
  }>(),
  {
    activeSlug: '',
  },
);

const internalActiveSlug = ref<string>(props.activeSlug.length > 0 ? props.activeSlug : (props.menu[0]?.slug as string));
const internalMenu = computed<MenuItem[]>(() => props.menu.map((el: MenuItem) => ({ ...el, command: () => scrollToSection(el.slug) })));

const scrollToSection = (slug: string) => {
  internalActiveSlug.value = slug;
  emits('update:activeSlug', slug);
  const section = document.getElementById(slug);
  if (section) {
    const offset = 50;
    const elementPosition = section.getBoundingClientRect().top;
    const offsetPosition = elementPosition + window.scrollY - offset;
    window.scrollTo({
      top: offsetPosition,
      behavior: 'smooth',
    });
  }
};
</script>
<template>
  <UiAccordionCard severity="outline">
    <template #title>
      <div class="flex justify-between items-center gap-2">
        <span class="text-md">Навигация</span>
      </div>
    </template>
    <template #content>
      <Menu v-if="internalMenu" :model="internalMenu" class="!border-0 !p-0 !ml-[-10px]" />
    </template>
  </UiAccordionCard>
</template>
