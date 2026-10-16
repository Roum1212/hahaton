import { ConnectError } from '@connectrpc/connect';
import type * as YMaps from 'yandex-maps';
import type { ApiPagination, IListResult } from '~/core/app/types';
import { GEO_API } from '../..';
import type { ApiGeo } from '../../api/geo/geo.api.types';
import type { AddressBound, IAddressSuggestion, ICity, IGeoPoint, IReverseAddress, IReverseGeocode, ISearchArea } from './geo.types';

export class GeoService {
  private service = GEO_API().geo;
  private yandexMapsApiKey = String(useRuntimeConfig().public.yandexMapsApiKey || '');

  async listCities(query: string | undefined, pagination: ApiPagination, signal?: AbortSignal): Promise<IListResult<ICity>> {
    const { cities, totalCount } = await this.service.listCities(query?.trim() || undefined, toApiPage(pagination), signal);
    return { list: cities.map((item) => this.toCity(item)), total: fromBigInt(totalCount) };
  }

  async getCity(id: string): Promise<ICity> {
    const city = await this.service.getCity(id);
    if (!city) throw new ConnectError('Пустой ответ сервера');
    return this.toCity(city);
  }

  async suggest(query: string, from: AddressBound, to: AddressBound, count = 10, signal?: AbortSignal): Promise<IAddressSuggestion[]> {
    const suggestions = await this.service.suggest({ query, fromBound: { value: from }, toBound: { value: to }, count }, signal);
    return suggestions.map((item) => this.toSuggestion(item));
  }

  async reverse(point: IGeoPoint, signal?: AbortSignal): Promise<IReverseGeocode | null> {
    return this.yandexMapsApiKey ? this.reverseWithYandex(point) : this.reverseWithOsm(point, signal);
  }

  async reverseWithYandex(point: IGeoPoint): Promise<IReverseGeocode | null> {
    const maps = await loadYandexMaps(this.yandexMapsApiKey);
    const { geoObjects } = await maps.geocode(toYandexCoords(point), { results: 1 });
    const item = geoObjects.get(0) as YMaps.GeocodeResult | undefined;
    if (!item) return null;
    const street = item.getThoroughfare();
    const address = street ? [street, item.getPremiseNumber()].filter(Boolean).join(', ') : String(item.properties.get('name', {}) || '');
    return { address: address || null, city: item.getLocalities()[0] ?? null };
  }

  async reverseWithOsm(point: IGeoPoint, signal?: AbortSignal): Promise<IReverseGeocode | null> {
    const params = new URLSearchParams({
      format: 'jsonv2',
      lat: String(point.latitude),
      lon: String(point.longitude),
      zoom: '18',
      addressdetails: '1',
      'accept-language': 'ru',
    });
    const response = await fetch(`${REVERSE_GEOCODE_URL}?${params}`, { signal });
    if (!response.ok) return null;
    const data: { address?: IReverseAddress; name?: string } = await response.json();
    if (!data.address) return null;
    return {
      address: this.toReverseLabel(data.address, data.name),
      city: data.address.city ?? data.address.town ?? data.address.village ?? null,
    };
  }

  toReverseLabel(address: IReverseAddress, name?: string): string | null {
    const street = address.road ?? address.pedestrian ?? address.footway;
    if (street) return [street, address.house_number].filter(Boolean).join(', ');
    const place = name || address.amenity || address.building || address.park || address.leisure;
    const area = address.neighbourhood ?? address.quarter ?? address.suburb ?? address.city ?? address.town ?? address.village;
    return [place, area].filter(Boolean).join(', ') || null;
  }

  toApiPoint(point: IGeoPoint): ApiGeo.PointInit {
    return { latitude: point.latitude, longitude: point.longitude };
  }

  toApiSearchArea(area: ISearchArea): ApiGeo.SearchArea {
    if ('cityId' in area) return { area: { case: 'cityId', value: area.cityId } };
    return { area: { case: 'radius', value: { center: this.toApiPoint(area.center), radiusMeters: area.radiusMeters } } };
  }

  toPoint(point: ApiGeo.Point | undefined): IGeoPoint | null {
    if (point?.latitude === undefined || point.longitude === undefined) return null;
    return { latitude: point.latitude, longitude: point.longitude };
  }

  toCity(data: ApiGeo.CityResponse): ICity {
    return {
      id: data.id,
      name: data.name,
      regionName: data.regionName,
      center: this.toPoint(data.center),
      hasBoundary: data.hasBoundary,
    };
  }

  toSuggestion(data: ApiGeo.Suggestion): IAddressSuggestion {
    return {
      value: data.value,
      unrestrictedValue: data.unrestrictedValue,
      fiasId: data.data?.fiasId ?? null,
      cityFiasId: data.data?.cityFiasId ?? null,
      settlementFiasId: data.data?.settlementFiasId ?? null,
      cityWithType: data.data?.cityWithType ?? null,
      regionWithType: data.data?.regionWithType ?? null,
      settlementWithType: data.data?.settlementWithType ?? null,
      point: this.toPoint(data.data?.point),
      qcGeo: data.data?.qcGeo ?? null,
    };
  }
}
