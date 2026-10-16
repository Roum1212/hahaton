import { EventApi } from './api/event/event.api';
import { EventService } from './services/event/event.service';

export const EVENT_API = () => {
  const { $transport } = useNuxtApp();
  return {
    event: new EventApi($transport),
  };
};

export const EVENT_SERVICES = () => {
  return {
    EventService: new EventService(),
  };
};
