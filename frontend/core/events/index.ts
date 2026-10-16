import { EventsService } from './services/events/events.service';

export const EVENTS_SERVICES = () => {
  return {
    EventsService: new EventsService(),
  };
};
