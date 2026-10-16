import type * as YMaps from 'yandex-maps';
import type { IGeoPoint } from '~/core/geo/services/geo/geo.types';

export type YandexMaps = typeof YMaps;

let loading: Promise<YandexMaps> | null = null;

const globalMaps = () => (window as unknown as { ymaps?: YandexMaps }).ymaps;

export function loadYandexMaps(apiKey = ''): Promise<YandexMaps> {
  if (loading) return loading;
  loading = new Promise<YandexMaps>((resolve, reject) => {
    const ready = () => globalMaps()!.ready(() => resolve(globalMaps()!));
    if (globalMaps()) {
      ready();
      return;
    }
    const params = new URLSearchParams({ lang: 'ru_RU' });
    if (apiKey) params.set('apikey', apiKey);
    const script = document.createElement('script');
    script.src = `${YANDEX_MAPS_URL}?${params}`;
    script.async = true;
    script.onload = ready;
    script.onerror = () => {
      loading = null;
      script.remove();
      reject(new Error('Не удалось загрузить Яндекс Карты'));
    };
    document.head.appendChild(script);
  });
  return loading;
}

export function toYandexCoords(point: IGeoPoint): [number, number] {
  return [point.latitude, point.longitude];
}

export function yandexMapsLink(point: IGeoPoint, zoom = 17): string {
  const ll = `${point.longitude},${point.latitude}`;
  return `${YANDEX_MAPS_WEB_URL}?${new URLSearchParams({ pt: ll, ll, z: String(zoom), l: 'map' })}`;
}
