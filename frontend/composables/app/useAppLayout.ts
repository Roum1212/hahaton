import { computed, reactive } from 'vue';
import { ROUTES } from '~/config/routes.config';
import { MENU, type IMenuGroup, type IMenuItem, type IMenuSection } from '~/config/menu.config';

interface ILayoutConfig {
  darkTheme: boolean;
}

interface ILayoutState {
  sidebarCollapsed: boolean;
  mobileSidebarOpen: boolean;
  selectedSection: string | null;
  menu: IMenuSection[];
}

export interface IActiveMenu {
  section: IMenuSection;
  group: IMenuGroup;
  item: IMenuItem;
}

const STORAGE_KEYS = {
  darkTheme: 'darkTheme',
  sidebarCollapsed: 'sidebarCollapsed',
};

const readStorage = (key: string): boolean => {
  if (!import.meta.client) return false;
  return localStorage.getItem(key) === 'true';
};

const layoutConfig = reactive<ILayoutConfig>({
  darkTheme: readStorage(STORAGE_KEYS.darkTheme),
});

const layoutState = reactive<ILayoutState>({
  sidebarCollapsed: readStorage(STORAGE_KEYS.sidebarCollapsed),
  mobileSidebarOpen: false,
  selectedSection: null,
  menu: MENU,
});

if (import.meta.client) document.documentElement.classList.toggle('app-dark', layoutConfig.darkTheme);

export function useAppLayout() {
  const isDarkTheme = computed(() => layoutConfig.darkTheme);

  const executeDarkModeToggle = () => {
    layoutConfig.darkTheme = !layoutConfig.darkTheme;
    localStorage.setItem(STORAGE_KEYS.darkTheme, String(layoutConfig.darkTheme));
    document.documentElement.classList.toggle('app-dark', layoutConfig.darkTheme);
  };

  const toggleDarkMode = () => {
    if (!document.startViewTransition) {
      executeDarkModeToggle();
      return;
    }
    document.startViewTransition(() => executeDarkModeToggle());
  };

  const setSidebarCollapsed = (value: boolean) => {
    layoutState.sidebarCollapsed = value;
    localStorage.setItem(STORAGE_KEYS.sidebarCollapsed, String(value));
  };

  const toggleSidebar = () => setSidebarCollapsed(!layoutState.sidebarCollapsed);

  const toggleMobileSidebar = () => {
    layoutState.mobileSidebarOpen = !layoutState.mobileSidebarOpen;
  };

  const closeMobileSidebar = () => {
    layoutState.mobileSidebarOpen = false;
  };

  const setMenu = (menu: IMenuSection[]) => {
    layoutState.menu = menu;
  };

  const selectSection = (key: string | null) => {
    layoutState.selectedSection = key;
  };

  const isActiveItem = (item: IMenuItem, path: string): boolean => {
    if (!item.to) return false;
    if (item.to === path) return true;
    return item.to !== ROUTES.cabinet && path.startsWith(`${item.to}/`);
  };

  const findActive = (path: string): IActiveMenu | null => {
    let match: IActiveMenu | null = null;
    for (const section of layoutState.menu) {
      for (const group of section.groups) {
        for (const item of group.items) {
          if (!isActiveItem(item, path)) continue;
          if (!match || (item.to?.length || 0) > (match.item.to?.length || 0)) match = { section, group, item };
        }
      }
    }
    return match;
  };

  return {
    layoutConfig,
    layoutState,
    isDarkTheme,
    toggleDarkMode,
    toggleSidebar,
    setSidebarCollapsed,
    toggleMobileSidebar,
    closeMobileSidebar,
    setMenu,
    selectSection,
    isActiveItem,
    findActive,
  };
}
