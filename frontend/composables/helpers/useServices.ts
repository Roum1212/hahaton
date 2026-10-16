import { AUTH_SERVICES } from '~/core/auth';
import { EVENTS_SERVICES } from '~/core/events';
import { USERS_SERVICES } from '~/core/users';

export const useServices = () => {
  return {
    AUTH: AUTH_SERVICES(),
    USERS: USERS_SERVICES(),
    EVENTS: EVENTS_SERVICES(),
  };
};
