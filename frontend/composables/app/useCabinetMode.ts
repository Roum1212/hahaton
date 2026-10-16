import type { CabinetMode } from '~/config/menu.config';

const MODE_HOMES: Record<CabinetMode, string> = {
  volunteer: ROUTES.volunteer.home,
  coordinator: ROUTES.coordinator.home,
};

export const modeFromPath = (path: string): CabinetMode | null => {
  if (path === ROUTES.coordinator.home || path.startsWith(`${ROUTES.coordinator.home}/`)) return 'coordinator';
  if (path === ROUTES.cabinet || path.startsWith(`${ROUTES.cabinet}/`)) return 'volunteer';
  return null;
};

export const useCabinetMode = () => {
  const authStore = useAuthStore();
  const stored = useLocalStorage<CabinetMode>('cabinetMode', 'volunteer');

  const canCoordinate = computed(() => authStore.hasCoordinatorProfile);
  const mode = computed<CabinetMode>(() => (stored.value === 'coordinator' && canCoordinate.value ? 'coordinator' : 'volunteer'));

  const setMode = (value: CabinetMode) => {
    stored.value = value === 'coordinator' && !canCoordinate.value ? 'volunteer' : value;
  };

  const switchMode = async (value: CabinetMode) => {
    setMode(value);
    await navigateTo(MODE_HOMES[mode.value]);
  };

  const homeFor = (value: CabinetMode = mode.value) => MODE_HOMES[value];

  return { mode, canCoordinate, setMode, switchMode, homeFor };
};
