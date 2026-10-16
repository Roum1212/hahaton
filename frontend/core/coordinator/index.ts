import { CoordinatorApi } from './api/coordinator/coordinator.api';
import { CoordinatorService } from './services/coordinator/coordinator.service';

export const COORDINATOR_API = () => {
  const { $transport } = useNuxtApp();
  return {
    coordinator: new CoordinatorApi($transport),
  };
};

export const COORDINATOR_SERVICES = () => {
  return {
    CoordinatorService: new CoordinatorService(),
  };
};
