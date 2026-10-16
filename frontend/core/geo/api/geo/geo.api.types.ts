import type { MessageInitShape } from '@bufbuild/protobuf';
import type {
  AddressSuggestion,
  City,
  GeoPoint,
  GeoPointSchema,
  ListCitiesResponse,
  SearchAreaSchema,
  SuggestAddressesRequestSchema,
} from '#genproto/geo/v1/geo_service_pb';

export namespace ApiGeo {
  export type Point = GeoPoint;
  export type PointInit = MessageInitShape<typeof GeoPointSchema>;
  export type CityResponse = City;
  export type CitiesResponse = ListCitiesResponse;
  export type Suggestion = AddressSuggestion;
  export type SuggestRequest = MessageInitShape<typeof SuggestAddressesRequestSchema>;
  export type SearchArea = MessageInitShape<typeof SearchAreaSchema>;
}
