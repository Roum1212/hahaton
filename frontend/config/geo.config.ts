import type { ICity } from '~/core/geo/services/geo/geo.types';

export const DEFAULT_CITY: ICity = {
  id: '05a647ff-b859-4c93-8c8e-3590110a0085',
  name: 'Обнинск',
  regionName: 'Калужская обл',
  center: { latitude: 55.094411, longitude: 36.612196 },
  hasBoundary: false,
};

export const CITY_LOCATIVES: Record<string, string> = {
  [DEFAULT_CITY.id]: 'Обнинске',
};

export const SEARCH_RADIUS_OPTIONS = [1, 2, 3, 5, 10, 15, 25, 50, 100];

export const DEFAULT_SEARCH_RADIUS_KM = 5;

export const YANDEX_MAPS_URL = 'https://api-maps.yandex.ru/2.1/';

export const YANDEX_MAPS_WEB_URL = 'https://yandex.ru/maps/';

export const MAP_ACCENT_COLOR = '#f97316';

export const REVERSE_GEOCODE_URL = 'https://nominatim.openstreetmap.org/reverse';
