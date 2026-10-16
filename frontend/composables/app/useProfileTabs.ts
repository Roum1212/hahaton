import type { IProfileTab } from '~/config/cabinet.config';

export const useProfileTabs = (tabs: IProfileTab[]) => {
  const route = useRoute();
  const router = useRouter();

  const active = computed<string>({
    get: () => {
      const value = String(route.query.tab ?? '');
      return tabs.some((tab) => tab.key === value) ? value : tabs[0]!.key;
    },
    set: (value) => {
      router.replace({ query: { ...route.query, tab: value } });
    },
  });

  return { tabs, active };
};
