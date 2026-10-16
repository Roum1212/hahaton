import { CooperationApi } from './api/cooperation/cooperation.api';
import { CooperationService } from './services/cooperation/cooperation.service';

export const COOPERATION_API = () => {
  const { $transport } = useNuxtApp();
  return {
    cooperation: new CooperationApi($transport),
  };
};

export const COOPERATION_SERVICES = () => {
  return {
    CooperationService: new CooperationService(),
  };
};
