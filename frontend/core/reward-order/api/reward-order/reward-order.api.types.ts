import type { MessageInitShape } from '@bufbuild/protobuf';
import type {
  ListRewardOrdersRequestSchema,
  ListRewardOrdersResponse,
  RewardOrder,
} from '#genproto/reward-order/v1/reward_order_service_pb';

export namespace ApiRewardOrder {
  export type Response = RewardOrder;
  export type ListResponse = ListRewardOrdersResponse;
  export type FindParams = MessageInitShape<typeof ListRewardOrdersRequestSchema>;
}
