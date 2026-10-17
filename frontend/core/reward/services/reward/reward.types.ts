import { RewardListScope } from '#genproto/reward/v1/reward_service_pb';

export { RewardListScope };

export interface IReward {
  id: string;
  coordinatorId: string;
  title: string;
  description: string;
  priceCoins: number;
  availableQuantity: number;
  quantityLimit: number;
  imageUrl: string | null;
}

export interface IRewardDto {
  title: string;
  description: string;
  priceCoins: number;
  quantityLimit: number;
  imagePath?: string;
}

export interface IRewardFilter {
  scope: RewardListScope;
  coordinatorId?: string;
}
