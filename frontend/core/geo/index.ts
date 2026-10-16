import { GeoApi } from './api/geo/geo.api';
import { GeoService } from './services/geo/geo.service';

export const GEO_API = () => {
  const { $transport } = useNuxtApp();
  return {
    geo: new GeoApi($transport),
  };
};

export const GEO_SERVICES = () => {
  return {
    GeoService: new GeoService(),
  };
};
