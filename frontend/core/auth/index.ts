import { AuthApi } from './api/auth/auth.api';
import { AuthSessionService } from './services/auth/auth.service';

export const AUTH_API = () => {
  const { $transport, $publicTransport } = useNuxtApp();
  return {
    auth: new AuthApi($transport, $publicTransport),
  };
};

export const AUTH_SERVICES = () => {
  return {
    AuthService: new AuthSessionService(),
  };
};
