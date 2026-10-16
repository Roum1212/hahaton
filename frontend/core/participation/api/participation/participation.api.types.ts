import type { MessageInitShape } from '@bufbuild/protobuf';
import type {
  ListParticipationChangesResponse,
  ListParticipationsRequestSchema,
  ListParticipationsResponse,
  Participation,
  ParticipationChange,
} from '#genproto/participation/v1/participation_service_pb';

export namespace ApiParticipation {
  export type Response = Participation;
  export type ListResponse = ListParticipationsResponse;
  export type FindParams = MessageInitShape<typeof ListParticipationsRequestSchema>;
  export type Change = ParticipationChange;
  export type ChangesResponse = ListParticipationChangesResponse;
}
