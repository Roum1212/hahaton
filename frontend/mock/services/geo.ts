import { create } from '@bufbuild/protobuf';
import { Code, type ConnectRouter } from '@connectrpc/connect';
import {
  AddressSuggestionDataSchema,
  AddressSuggestionSchema,
  GeoService,
  type GeoPoint,
  type SearchArea,
} from '../../genproto/geo/v1/geo_service_pb';
import { requireUser } from '../auth';
import { compareStrings, db, fail, paginate, requireFound, type CityRec, type Point } from '../db';
import { toCity, toPoint } from '../mappers';

const DADATA_URL = 'https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address';
const CATALOG_VERSION = 'mock-2026-10';
const BOUNDS = ['country', 'region', 'area', 'city', 'settlement', 'street', 'house'];

interface DadataSuggestion {
  value: string;
  unrestricted_value: string;
  data: {
    fias_id?: string | null;
    city_fias_id?: string | null;
    settlement_fias_id?: string | null;
    city?: string | null;
    settlement?: string | null;
    settlement_type?: string | null;
    city_with_type?: string | null;
    region_with_type?: string | null;
    settlement_with_type?: string | null;
    geo_lat?: string | null;
    geo_lon?: string | null;
    qc_geo?: string | number | null;
  };
}

const parsePoint = (lat?: string | null, lon?: string | null): Point | undefined => {
  const latitude = lat ? Number(lat) : NaN;
  const longitude = lon ? Number(lon) : NaN;
  if (!Number.isFinite(latitude) || !Number.isFinite(longitude)) return undefined;
  return { latitude, longitude };
};

export const callDadata = async (query: string, from: string, to: string, count: number): Promise<DadataSuggestion[]> => {
  const key = process.env.MOCK_DADATA_API_KEY;
  if (!key) return fail(Code.Unavailable, 'DaData не настроен');
  const response = await fetch(DADATA_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Accept: 'application/json', Authorization: `Token ${key}` },
    body: JSON.stringify({
      query,
      count,
      from_bound: { value: from },
      to_bound: { value: to },
      locations: [{ country_iso_code: 'RU' }],
    }),
    signal: AbortSignal.timeout(5000),
  }).catch(() => null);
  if (!response) return fail(Code.Unavailable, 'DaData недоступен');
  if (response.status === 429 || response.status === 403) return fail(Code.ResourceExhausted, 'Лимит подсказок DaData исчерпан');
  if (!response.ok) return fail(Code.Unavailable, 'DaData недоступен');
  const body = (await response.json()) as { suggestions?: DadataSuggestion[] };
  return body.suggestions ?? [];
};

const SETTLEMENT_TYPES = ['пгт', 'рп', 'кп', 'дп', 'п', 'с', 'д', 'ст-ца', 'х', 'аул', 'сл', 'нп', 'г'];

const isPopulatedPlace = (item: DadataSuggestion) =>
  !item.data.settlement ? !!item.data.city : SETTLEMENT_TYPES.includes(item.data.settlement_type ?? '');

const rememberCity = (item: DadataSuggestion): CityRec | undefined => {
  const id = item.data.settlement_fias_id || item.data.city_fias_id;
  const name = item.data.settlement || item.data.city;
  if (!id || !name) return undefined;
  const existing = db.cities.get(id);
  if (existing) return existing;
  const city: CityRec = {
    id,
    name,
    regionName: item.data.region_with_type ?? '',
    center: parsePoint(item.data.geo_lat, item.data.geo_lon),
  };
  db.cities.set(id, city);
  return city;
};

const validatePoint = (point: GeoPoint | undefined): Point => {
  const latitude = point?.latitude;
  const longitude = point?.longitude;
  if (latitude === undefined || longitude === undefined || !Number.isFinite(latitude) || !Number.isFinite(longitude)) {
    return fail(Code.InvalidArgument, 'Нужны обе координаты точки');
  }
  if (latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180) return fail(Code.InvalidArgument, 'Координаты вне диапазона');
  return { latitude, longitude };
};

export { validatePoint };

const distanceMeters = (a: Point, b: Point) => {
  const rad = Math.PI / 180;
  const dLat = (b.latitude - a.latitude) * rad;
  const dLon = (b.longitude - a.longitude) * rad;
  const h = Math.sin(dLat / 2) ** 2 + Math.cos(a.latitude * rad) * Math.cos(b.latitude * rad) * Math.sin(dLon / 2) ** 2;
  return 2 * 6371000 * Math.asin(Math.sqrt(h));
};

const insidePolygon = (point: Point, polygon: Point[]) => {
  let inside = false;
  for (let i = 0, j = polygon.length - 1; i < polygon.length; j = i++) {
    const a = polygon[i]!;
    const b = polygon[j]!;
    const crosses =
      a.latitude > point.latitude !== b.latitude > point.latitude &&
      point.longitude <= ((b.longitude - a.longitude) * (point.latitude - a.latitude)) / (b.latitude - a.latitude) + a.longitude;
    if (crosses) inside = !inside;
  }
  return inside;
};

export const buildAreaFilter = (area: SearchArea | undefined): ((point: Point) => boolean) => {
  const value = area?.area;
  if (!value?.case) return fail(Code.InvalidArgument, 'Нужна область поиска');
  if (value.case === 'cityId') {
    const city = requireFound(db.cities.get(value.value), 'Город не найден');
    if (!city.boundary?.length) return fail(Code.FailedPrecondition, 'Для города нет границы, выберите радиус');
    const boundary = city.boundary;
    return (point) => insidePolygon(point, boundary);
  }
  const center = validatePoint(value.value.center);
  const radius = value.value.radiusMeters;
  if (!Number.isInteger(radius) || radius < 1 || radius > 200000) return fail(Code.InvalidArgument, 'Радиус от 1 до 200000 метров');
  return (point) => distanceMeters(center, point) <= radius;
};

export const registerGeoService = (router: ConnectRouter) =>
  router.service(GeoService, {
    async suggestAddresses(req, ctx) {
      requireUser(ctx);
      const query = req.query.trim();
      if (!query || [...query].length > 300) return fail(Code.InvalidArgument, 'Запрос должен быть от 1 до 300 символов');
      const from = req.fromBound?.value ?? '';
      const to = req.toBound?.value ?? '';
      const fromIndex = BOUNDS.indexOf(from);
      const toIndex = BOUNDS.indexOf(to);
      if (fromIndex < 0 || toIndex < 0 || fromIndex > toIndex) return fail(Code.InvalidArgument, 'Некорректные границы подсказок');
      const count = req.count === 0 ? 10 : req.count;
      if (count < 1 || count > 20) return fail(Code.InvalidArgument, 'count должен быть от 1 до 20');
      const suggestions = await callDadata(query, from, to, count);
      return {
        suggestions: suggestions.map((item) => {
          const point = parsePoint(item.data.geo_lat, item.data.geo_lon);
          const qc = item.data.qc_geo === null || item.data.qc_geo === undefined ? undefined : Number(item.data.qc_geo);
          return create(AddressSuggestionSchema, {
            value: item.value,
            unrestrictedValue: item.unrestricted_value,
            data: create(AddressSuggestionDataSchema, {
              fiasId: item.data.fias_id ?? undefined,
              cityFiasId: item.data.city_fias_id ?? undefined,
              settlementFiasId: item.data.settlement_fias_id ?? undefined,
              cityWithType: item.data.city_with_type ?? undefined,
              regionWithType: item.data.region_with_type ?? undefined,
              settlementWithType: item.data.settlement_with_type ?? undefined,
              point: point ? toPoint(point) : undefined,
              qcGeo: Number.isFinite(qc) ? qc : undefined,
            }),
          });
        }),
      };
    },
    async listCities(req) {
      const query = req.query?.trim() ?? '';
      const foundIds = new Set<string>();
      if (query && process.env.MOCK_DADATA_API_KEY) {
        const found = await callDadata(query, 'city', 'settlement', 20);
        found
          .filter(isPopulatedPlace)
          .map(rememberCity)
          .forEach((city) => city && foundIds.add(city.id));
      }
      const needle = query.toLowerCase();
      const items = [...db.cities.values()]
        .filter((item) => !needle || item.name.toLowerCase().includes(needle) || foundIds.has(item.id))
        .sort((a, b) => compareStrings(a.name, b.name) || compareStrings(a.regionName, b.regionName) || compareStrings(a.id, b.id));
      const { page, totalCount } = paginate(items, req.limit, req.offset, 100, 1000);
      return { cities: page.map(toCity), catalogVersion: CATALOG_VERSION, totalCount };
    },
    getCity(req) {
      if (!req.id) return fail(Code.InvalidArgument, 'Не указан город');
      const city = requireFound(db.cities.get(req.id), 'Город не найден');
      return { city: toCity(city), catalogVersion: CATALOG_VERSION };
    },
  });
