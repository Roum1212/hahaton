import { RewardOrderCancellationReason, RewardOrderListScope, RewardOrderStatus } from '#genproto/reward-order/v1/reward_order_service_pb';

export { RewardOrderCancellationReason, RewardOrderListScope, RewardOrderStatus };

export interface IRewardOrder {
  id: string;
  coordinatorId: string;
  rewardId: string;
  volunteerId: string;
  priceCoins: number;
  rewardTitle: string;
  status: RewardOrderStatus;
  cancellationReason: RewardOrderCancellationReason;
  createdAt: Date | null;
  updatedAt: Date | null;
}

export interface IRewardOrderFilter {
  scope: RewardOrderListScope;
  coordinatorId?: string;
  rewardId?: string;
  volunteerId?: string;
  status?: RewardOrderStatus;
}

export const REWARD_ORDER_STATUS_LABELS: Record<RewardOrderStatus, string> = {
  [RewardOrderStatus.UNSPECIFIED]: 'Неизвестно',
  [RewardOrderStatus.PENDING]: 'Ожидает выдачи',
  [RewardOrderStatus.ISSUED]: 'Выдано',
  [RewardOrderStatus.CANCELLED]: 'Отменено',
};

export const REWARD_ORDER_STATUS_SEVERITY: Record<RewardOrderStatus, 'success' | 'info' | 'warn' | 'danger' | 'secondary'> = {
  [RewardOrderStatus.UNSPECIFIED]: 'secondary',
  [RewardOrderStatus.PENDING]: 'warn',
  [RewardOrderStatus.ISSUED]: 'success',
  [RewardOrderStatus.CANCELLED]: 'secondary',
};

export const REWARD_ORDER_CANCELLATION_LABELS: Record<RewardOrderCancellationReason, string> = {
  [RewardOrderCancellationReason.UNSPECIFIED]: '',
  [RewardOrderCancellationReason.COOPERATION_TERMINATED]: 'Сотрудничество прекращено',
  [RewardOrderCancellationReason.COORDINATOR_CANCELLED]: 'Отменено координатором',
  [RewardOrderCancellationReason.VOLUNTEER_CANCELLED]: 'Отменено волонтёром',
  [RewardOrderCancellationReason.REWARD_DELETED]: 'Поощрение удалено',
};
