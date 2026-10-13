import { UsersApi } from './api/users/users.api';
import { UsersService } from './services/users/users.service';

export const USERS_API = () => {
  const { $transport } = useNuxtApp();
  return {
    users: new UsersApi($transport),
  };
};

export const USERS_SERVICES = () => {
  return {
    UsersService: new UsersService(),
  };
};
