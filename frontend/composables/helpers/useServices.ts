import { AUTH_SERVICES } from '~/core/auth';
import { USERS_SERVICES } from '~/core/users';

export const useServices = () => {
  return {
    AUTH: AUTH_SERVICES(),
    USERS: USERS_SERVICES(),
  };
};
