export const useAppMenu = () => {
  const route = useRoute();
  const { layoutState, findActive, isActiveItem } = useAppLayout();

  const menu = computed(() => layoutState.menu);
  const active = computed(() => findActive(route.path));

  const currentSection = computed(
    () => menu.value.find((section) => section.key === layoutState.selectedSection) ?? active.value?.section ?? menu.value[0] ?? null,
  );

  const sectionLabel = computed(() => active.value?.section.label || 'Кабинет');
  const pageLabel = computed(() => (route.meta.title as string | undefined) || active.value?.item.label || '');

  return { menu, active, currentSection, sectionLabel, pageLabel, isActiveItem };
};
