import { RewardApi } from './api/reward/reward.api';
import { RewardService } from './services/reward/reward.service';

export const REWARD_API = () => {
  const { $transport } = useNuxtApp();
  return {
    reward: new RewardApi($transport),
  };
};

export const REWARD_SERVICES = () => {
  return {
    RewardService: new RewardService(),
  };
};
