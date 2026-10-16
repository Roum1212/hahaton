import type { MessageInitShape } from '@bufbuild/protobuf';
import type {
  Coordinator,
  CoordinatorSettings,
  ListCoordinatorsResponse,
  UpdateCoordinatorRequestSchema,
} from '#genproto/coordinator/v1/coordinator_service_pb';

export namespace ApiCoordinator {
  export type Response = Coordinator;
  export type ListResponse = ListCoordinatorsResponse;
  export type Settings = CoordinatorSettings;
  export type UpdateRequest = MessageInitShape<typeof UpdateCoordinatorRequestSchema>;
}
