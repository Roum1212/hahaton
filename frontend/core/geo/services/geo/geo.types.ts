export interface IGeoPoint {
  latitude: number;
  longitude: number;
}

export interface ICity {
  id: string;
  name: string;
  regionName: string;
  center: IGeoPoint | null;
  hasBoundary: boolean;
}

export type ISearchArea = { cityId: string } | { center: IGeoPoint; radiusMeters: number };

export type AddressBound = 'country' | 'region' | 'area' | 'city' | 'settlement' | 'street' | 'house';

export interface IAddressSuggestion {
  value: string;
  unrestrictedValue: string;
  fiasId: string | null;
  cityFiasId: string | null;
  settlementFiasId: string | null;
  cityWithType: string | null;
  regionWithType: string | null;
  settlementWithType: string | null;
  point: IGeoPoint | null;
  qcGeo: number | null;
}

export const SEARCH_RADIUS_LIMITS = { min: 1, max: 200000 } as const;

export interface IReverseAddress {
  road?: string;
  pedestrian?: string;
  footway?: string;
  house_number?: string;
  amenity?: string;
  building?: string;
  park?: string;
  leisure?: string;
  neighbourhood?: string;
  suburb?: string;
  quarter?: string;
  city?: string;
  town?: string;
  village?: string;
}

export interface IReverseGeocode {
  address: string | null;
  city: string | null;
}
