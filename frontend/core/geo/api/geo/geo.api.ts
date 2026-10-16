import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { GeoService } from '#genproto/geo/v1/geo_service_pb';
import type { ApiGeo } from './geo.api.types';

export class GeoApi extends BaseRpcApi<typeof GeoService> {
  constructor(transport: Transport) {
    super(GeoService, transport);
  }

  async listCities(
    query: string | undefined,
    page: { limit: number; offset: bigint },
    signal?: AbortSignal,
  ): Promise<ApiGeo.CitiesResponse> {
    return this.client.listCities({ query, ...page }, this.options(signal));
  }

  async getCity(id: string): Promise<ApiGeo.CityResponse | undefined> {
    const { city } = await this.client.getCity({ id });
    return city;
  }

  async suggest(data: ApiGeo.SuggestRequest, signal?: AbortSignal): Promise<ApiGeo.Suggestion[]> {
    const { suggestions } = await this.client.suggestAddresses(data, this.options(signal));
    return suggestions;
  }
}
