import { UpdatesApi } from './api/updates/updates.api';
import { UpdatesService } from './services/updates/updates.service';

export const UPDATES_API = () => {
  const { $transport } = useNuxtApp();
  return {
    updates: new UpdatesApi($transport),
  };
};

export const UPDATES_SERVICES = () => {
  return {
    UpdatesService: new UpdatesService(),
  };
};
