import type { MessageInitShape } from '@bufbuild/protobuf';
import type {
  CreateRewardRequestSchema,
  ListRewardsRequestSchema,
  ListRewardsResponse,
  Reward,
  UpdateRewardRequestSchema,
} from '#genproto/reward/v1/reward_service_pb';

export namespace ApiReward {
  export type Response = Reward;
  export type ListResponse = ListRewardsResponse;
  export type FindParams = MessageInitShape<typeof ListRewardsRequestSchema>;
  export type CreateRequest = MessageInitShape<typeof CreateRewardRequestSchema>;
  export type UpdateRequest = MessageInitShape<typeof UpdateRewardRequestSchema>;
}
