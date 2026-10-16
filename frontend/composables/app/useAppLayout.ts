import { computed, reactive } from 'vue';

interface ILayoutConfig {
  darkTheme: boolean;
}

const layoutConfig = reactive<ILayoutConfig>({
  darkTheme: false,
});

if (import.meta.client) document.documentElement.classList.remove('app-dark');

export function useAppLayout() {
  const isDarkTheme = computed(() => layoutConfig.darkTheme);

  const isActivePath = (to: string, path: string, exact = false): boolean => {
    if (to === path) return true;
    return !exact && path.startsWith(`${to}/`);
  };

  return { layoutConfig, isDarkTheme, isActivePath };
}
