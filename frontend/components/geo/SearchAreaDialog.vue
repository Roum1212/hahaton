<script setup lang="ts">
import { EventListScope, EventStatus } from '~/core/event/services/event/event.types';
import type { IAddressSuggestion, ICity, IGeoPoint } from '~/core/geo/services/geo/geo.types';
import { createDefaultArea, toSearchArea, type ISearchAreaState } from '~/composables/geo/useSearchArea';

const { state, dialogVisible, apply } = useSearchArea();
const authStore = useAuthStore();
const geoStore = useGeoStore();
const eventsStore = useEventsStore();
const { toast } = useApp();

const draft = ref<ISearchAreaState>(createDefaultArea());
const cityQuery = ref<ICity | string>('');
const citySuggestions = ref<ICity[]>([]);
const addressQuery = ref<IAddressSuggestion | string>('');
const addressSuggestions = ref<IAddressSuggestion[]>([]);
const count = ref<number | null>(null);
const counting = ref(false);
const { resolving, resolve, findCity } = useReverseGeocode();
let pickToken = 0;

const center = computed<IGeoPoint>(() => draft.value.center ?? draft.value.city.center ?? DEFAULT_CITY.center!);
const effectiveRadius = computed(() => draft.value.radiusKm ?? (draft.value.city.hasBoundary ? null : DEFAULT_SEARCH_RADIUS_KM));
const canReset = computed(() => !!draft.value.radiusKm && draft.value.city.hasBoundary);

const applyLabel = computed(() => {
  if (!authStore.isLoggedIn || count.value === null) return 'Применить';
  return `Показать ${count.value} ${pluralize(count.value, ['мероприятие', 'мероприятия', 'мероприятий'])}`;
});

watch(dialogVisible, (visible) => {
  if (!visible) return;
  draft.value = JSON.parse(JSON.stringify(state.value));
  cityQuery.value = draft.value.city.name;
  addressQuery.value = draft.value.address;
});

const searchCities = async ({ query }: { query: string }) => {
  citySuggestions.value = (await handleAsyncOperation(() => geoStore.searchCities(query, 10))) ?? [];
};

const selectCity = ({ value }: { value: ICity }) => {
  pickToken++;
  draft.value = {
    city: value,
    radiusKm: value.hasBoundary ? null : DEFAULT_SEARCH_RADIUS_KM,
    center: null,
    address: '',
  };
  cityQuery.value = value.name;
  addressQuery.value = '';
};

const searchAddresses = async ({ query }: { query: string }) => {
  const result = await handleAsyncOperation(() => geoStore.service.suggest(`${draft.value.city.name}, ${query}`, 'street', 'house', 7));
  addressSuggestions.value = result ?? [];
};

const selectAddress = async ({ value }: { value: IAddressSuggestion }) => {
  pickToken++;
  let point = value.point;
  if (!point && value.fiasId) {
    const [precise] = (await handleAsyncOperation(() => geoStore.service.suggest(value.fiasId!, 'street', 'house', 1))) ?? [];
    point = precise?.point ?? null;
  }
  draft.value.address = value.value;
  addressQuery.value = value.value;
  if (!point) {
    toast.warning('Для этого адреса нет координат, отметьте точку на карте');
    return;
  }
  draft.value.center = point;
  if (!draft.value.radiusKm) draft.value.radiusKm = DEFAULT_SEARCH_RADIUS_KM;
};

const setCenter = async (point: IGeoPoint) => {
  const token = ++pickToken;
  draft.value.center = point;
  draft.value.address = '';
  addressQuery.value = '';
  if (!draft.value.radiusKm) draft.value.radiusKm = DEFAULT_SEARCH_RADIUS_KM;
  const result = await resolve(point);
  if (!result || token !== pickToken) return;
  if (result.address) {
    draft.value.address = result.address;
    addressQuery.value = result.address;
  }
  const city = await findCity(result.city);
  if (!city || city.id === draft.value.city.id || token !== pickToken) return;
  draft.value.city = city;
  cityQuery.value = city.name;
};

const setRadius = (km: number) => {
  draft.value.radiusKm = km;
};

const resetRadius = () => {
  pickToken++;
  draft.value.radiusKm = null;
  draft.value.center = null;
  draft.value.address = '';
  addressQuery.value = '';
};

watchDebounced(
  () => [dialogVisible.value, JSON.stringify(toSearchArea(draft.value))],
  async () => {
    if (!dialogVisible.value || !authStore.isLoggedIn) {
      count.value = null;
      return;
    }
    counting.value = true;
    const result = await handleAsyncOperation(
      () =>
        eventsStore.service.find(
          { scope: EventListScope.CATALOG, searchArea: toSearchArea(draft.value), status: EventStatus.PUBLISHED, startsFrom: new Date() },
          { limit: 1, offset: 0 },
        ),
      () => {
        counting.value = false;
      },
    );
    count.value = result?.total ?? null;
  },
  { debounce: 400 },
);
</script>

<template>
  <Dialog
    v-model:visible="dialogVisible"
    modal
    dismissable-mask
    header="Город или регион"
    :draggable="false"
    class="geo-area-dialog"
    :style="{ width: '920px' }"
    :breakpoints="{ '767px': '100vw' }">
    <div class="geo-area">
      <AutoComplete
        v-model="cityQuery"
        :suggestions="citySuggestions"
        option-label="name"
        placeholder="Город"
        :delay="300"
        fluid
        @complete="searchCities"
        @option-select="selectCity">
        <template #option="{ option }">
          <div class="flex flex-col">
            <span>{{ option.name }}</span>
            <span class="geo-area__hint">{{ option.regionName }}</span>
          </div>
        </template>
      </AutoComplete>

      <span class="geo-area__section">Радиус поиска</span>

      <AutoComplete
        v-model="addressQuery"
        :suggestions="addressSuggestions"
        option-label="value"
        :placeholder="resolving ? 'Определяем адрес…' : 'Улица и номер дома или точка на карте'"
        :delay="300"
        :loading="resolving"
        fluid
        @complete="searchAddresses"
        @option-select="selectAddress" />

      <GeoMap
        v-if="dialogVisible"
        :center="center"
        :radius-meters="effectiveRadius ? effectiveRadius * 1000 : null"
        class="geo-area__map"
        @update:center="setCenter" />

      <span class="geo-area__label">Радиус поиска, км</span>
      <div class="geo-area__footer">
        <div class="geo-area__radius">
          <button
            v-for="km in SEARCH_RADIUS_OPTIONS"
            :key="km"
            type="button"
            class="geo-area__chip"
            :class="{ 'geo-area__chip--active': effectiveRadius === km }"
            @click="setRadius(km)">
            {{ km }}
          </button>
        </div>
        <div class="geo-area__actions">
          <Button v-if="canReset" label="Сбросить фильтр" severity="secondary" @click="resetRadius" />
          <Button :label="applyLabel" :loading="counting" @click="apply(draft)" />
        </div>
      </div>
    </div>
  </Dialog>
</template>
