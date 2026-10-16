import { RewardOrderApi } from './api/reward-order/reward-order.api';
import { RewardOrderService } from './services/reward-order/reward-order.service';

export const REWARD_ORDER_API = () => {
  const { $transport } = useNuxtApp();
  return {
    rewardOrder: new RewardOrderApi($transport),
  };
};

export const REWARD_ORDER_SERVICES = () => {
  return {
    RewardOrderService: new RewardOrderService(),
  };
};
