import { VolunteerApi } from './api/volunteer/volunteer.api';
import { VolunteerService } from './services/volunteer/volunteer.service';

export const VOLUNTEER_API = () => {
  const { $transport } = useNuxtApp();
  return {
    volunteer: new VolunteerApi($transport),
  };
};

export const VOLUNTEER_SERVICES = () => {
  return {
    VolunteerService: new VolunteerService(),
  };
};
