import { defineStore } from 'pinia';
import type { TableColumnType } from '~/components/ui/TableV2/index.vue';

export interface TableSettings {
  columns: TableColumnType[];
}

const LOCAL_STORAGE_PREFIX = 'table_settings_';
const LOCAL_STORAGE_VERSION = 1;

export const useNewTableSettingsStore = defineStore('new-table-settings', {
  state: () => ({
    tables: {} as Record<string, TableSettings>,
  }),
  getters: {},
  actions: {
    getTableSettings(tableKey: string) {
      this.load(tableKey);
      return this.tables[tableKey];
    },
    load(tableKey: string) {
      try {
        const data = localStorage.getItem(`${LOCAL_STORAGE_PREFIX}${tableKey}`);
        if (data) {
          const parsed = JSON.parse(data) as { version: number; settings: TableSettings };
          if (parsed.version === LOCAL_STORAGE_VERSION) {
            this.tables[tableKey] = parsed.settings;
          } else {
            console.warn(`Table settings version mismatch for ${tableKey}. Resetting.`);
            this.tables[tableKey] = { columns: [] };
          }
        } else {
          this.tables[tableKey] = { columns: [] };
        }
      } catch (e) {
        console.error('Failed to load table settings', e);
        this.tables[tableKey] = { columns: [] };
      }
    },

    save(tableKey: string) {
      const settings = this.tables[tableKey];
      if (!settings) return;
      try {
        const data = { version: LOCAL_STORAGE_VERSION, settings };
        localStorage.setItem(`${LOCAL_STORAGE_PREFIX}${tableKey}`, JSON.stringify(data));
      } catch (e) {
        console.error('Failed to save table settings', e);
      }
    },

    set(tableKey: string, settings: Partial<TableSettings>) {
      const existing = this.tables[tableKey] ?? { columns: [] };
      this.tables[tableKey] = { ...existing, ...settings };
      this.save(tableKey);
    },

    reset(tableKey: string) {
      this.tables[tableKey] = { columns: [] };
      localStorage.removeItem(`${LOCAL_STORAGE_PREFIX}${tableKey}`);
    },
  },
});
