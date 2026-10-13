import type { TableColumnType } from '~/components/ui/TableV2/index.vue';
export const useTableWithSettings = (defaultColumns: TableColumnType[], settingsKey?: string) => {
  const tableSettingsStore = useNewTableSettingsStore();

  if (!settingsKey) {
    return defaultColumns;
  }

  tableSettingsStore.load(settingsKey);
  const saved = tableSettingsStore.tables[settingsKey];

  if (saved?.columns?.length) {
    return saved.columns;
  }

  return defaultColumns;
};
