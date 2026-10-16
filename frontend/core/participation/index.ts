import { ParticipationApi } from './api/participation/participation.api';
import { ParticipationService } from './services/participation/participation.service';

export const PARTICIPATION_API = () => {
  const { $transport } = useNuxtApp();
  return {
    participation: new ParticipationApi($transport),
  };
};

export const PARTICIPATION_SERVICES = () => {
  return {
    ParticipationService: new ParticipationService(),
  };
};
