import type { ICity, IGeoPoint, ISearchArea } from '~/core/geo/services/geo/geo.types';

export interface ISearchAreaState {
  city: ICity;
  radiusKm: number | null;
  center: IGeoPoint | null;
  address: string;
}

export const createDefaultArea = (): ISearchAreaState => ({
  city: DEFAULT_CITY,
  radiusKm: null,
  center: null,
  address: '',
});

export const toSearchArea = (state: ISearchAreaState): ISearchArea => {
  const center = state.center ?? state.city.center;
  if ((state.radiusKm || !state.city.hasBoundary) && center) {
    return { center, radiusMeters: Math.round((state.radiusKm ?? DEFAULT_SEARCH_RADIUS_KM) * 1000) };
  }
  return { cityId: state.city.id };
};

export const useSearchArea = () => {
  const state = useLocalStorage<ISearchAreaState>('searchArea', createDefaultArea(), { mergeDefaults: true });
  const confirmed = useLocalStorage('cityConfirmed', false);
  const dialogVisible = useState('searchAreaDialogVisible', () => false);

  const area = computed(() => toSearchArea(state.value));
  const city = computed(() => state.value.city);
  const label = computed(() => (state.value.radiusKm ? `${state.value.city.name} · ${state.value.radiusKm} км` : state.value.city.name));
  const inCityTitle = computed(() => {
    const locative = CITY_LOCATIVES[state.value.city.id];
    return locative ? `в ${locative}` : `в городе ${state.value.city.name}`;
  });

  const confirm = () => {
    confirmed.value = true;
  };

  const openDialog = () => {
    confirmed.value = true;
    dialogVisible.value = true;
  };

  const apply = (value: ISearchAreaState) => {
    state.value = { ...value };
    confirmed.value = true;
    dialogVisible.value = false;
  };

  return { state, area, city, label, inCityTitle, confirmed, dialogVisible, confirm, openDialog, apply };
};
