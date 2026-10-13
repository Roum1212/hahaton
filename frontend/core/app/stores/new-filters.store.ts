/* eslint-disable @typescript-eslint/no-explicit-any */
import { merge } from 'lodash-es';
import { defineStore } from 'pinia';

export type FilterKey =
  | 'catalogs'
  | 'categories'
  | 'brands'
  | 'collections'
  | 'properties'
  | 'propertyGroups'
  | 'relations'
  | 'products'
  | 'accounts'
  | 'accountWarehouses'
  | 'accountProducts'
  | 'accountPromotions'
  | 'accountTasks'
  | 'accountReports'
  | 'accountLogs'
  | 'accountContentQueuesItems'
  | 'accountContentQueuesGroups'
  | 'accountPricesQueuesItems'
  | 'accountPricesQueuesGroups'
  | 'accountStocksQueuesItems'
  | 'accountStocksQueuesGroups'
  | 'users'
  | 'userGroups';

export const filterKeys: FilterKey[] = [
  'catalogs',
  'categories',
  'brands',
  'collections',
  'properties',
  'propertyGroups',
  'relations',
  'products',
  'accounts',
  'accountWarehouses',
  'accountProducts',
  'accountPromotions',
  'accountTasks',
  'accountReports',
  'accountLogs',
  'accountContentQueuesItems',
  'accountContentQueuesGroups',
  'accountPricesQueuesItems',
  'accountPricesQueuesGroups',
  'accountStocksQueuesItems',
  'accountStocksQueuesGroups',
  'users',
  'userGroups',
];

const initFilters = <T>(defaultValue: T) => Object.fromEntries(filterKeys.map((k) => [k, defaultValue])) as Record<FilterKey, T>;

export type FilterValue = Record<string, any>;

export interface FilterItem {
  key: FilterKey;
  value: FilterValue;
  catalogId?: string | null;
  accountId?: string | null;
  marketplaceId?: string | null;
  global?: boolean;
  name?: string;
}

interface FiltersState {
  current: Record<FilterKey, FilterValue>;
  saved: Record<FilterKey, FilterItem[]>;
}

const LOCAL_STORAGE_KEY = 'filters_saved';
const LOCAL_STORAGE_VERSION = 1;

export const useNewFiltersStore = defineStore('new-filters', {
  state: (): FiltersState => ({
    current: initFilters<FilterValue>({}),
    saved: initFilters<FilterItem[]>([]),
  }),

  actions: {
    setFilter(key: FilterKey, value: FilterValue) {
      this.current[key] = value;
    },

    applyFilters(key: FilterKey, filters: FilterValue[]) {
      this.current[key] = merge({}, ...filters);
    },

    resetFilter(key: FilterKey) {
      this.current[key] = {};
    },

    resetAll() {
      filterKeys.forEach((k) => (this.current[k] = {}));
    },

    saveFilter(key: FilterKey, name: string, options?: Partial<Pick<FilterItem, 'catalogId' | 'accountId' | 'marketplaceId' | 'global'>>) {
      const item: FilterItem = {
        key,
        value: this.current[key],
        name,
        catalogId: options?.catalogId ?? null,
        accountId: options?.accountId ?? null,
        marketplaceId: options?.marketplaceId ?? null,
        global: options?.global ?? false,
      };
      this.saved[key].push(item);
      this.persistSavedFilters();
    },

    loadSavedFilter(key: FilterKey, name: string) {
      const item = this.saved[key].find((f) => f.name === name);
      if (item) this.setFilter(key, item.value);
    },

    deleteSavedFilter(key: FilterKey, name: string) {
      this.saved[key] = this.saved[key].filter((f) => f.name !== name);
      this.persistSavedFilters();
    },

    getApplicableSavedFilters(
      key: FilterKey,
      context: { catalogId?: string; accountId?: string; marketplaceId?: string; global?: boolean },
    ) {
      return this.saved[key].filter(
        (f) =>
          (f.catalogId === undefined || f.catalogId === context.catalogId) &&
          (f.accountId === undefined || f.accountId === context.accountId) &&
          (f.marketplaceId === undefined || f.marketplaceId === context.marketplaceId) &&
          (f.global || context.global),
      );
    },

    persistSavedFilters() {
      try {
        const data = { version: LOCAL_STORAGE_VERSION, saved: this.saved };
        localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(data));
      } catch (e) {
        console.error('Error saving filters to localStorage', e);
      }
    },

    loadSavedFilters() {
      try {
        const data = localStorage.getItem(LOCAL_STORAGE_KEY);
        if (data) {
          const parsed = JSON.parse(data);
          if (parsed.version === LOCAL_STORAGE_VERSION) {
            this.saved = parsed.saved ?? initFilters<FilterItem[]>([]);
          } else {
            console.warn('Filters version mismatch. Resetting saved filters.');
            this.saved = initFilters<FilterItem[]>([]);
          }
        }
      } catch (e) {
        console.error('Error loading filters from localStorage', e);
      }
    },
  },
});
