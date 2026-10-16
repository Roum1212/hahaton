import { SessionApi } from './api/session/session.api';
import { AuthSessionService } from './services/auth/auth.service';

export const AUTH_API = () => {
  const { $transport } = useNuxtApp();
  return {
    session: new SessionApi($transport),
  };
};

export const AUTH_SERVICES = () => {
  return {
    AuthService: new AuthSessionService(),
  };
};
