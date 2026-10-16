import type { MessageInitShape } from '@bufbuild/protobuf';
import type {
  Cooperation,
  CooperationInvitation,
  ListCooperationInvitationsResponse,
  ListCooperationsRequestSchema,
  ListCooperationsResponse,
} from '#genproto/cooperation/v1/cooperation_service_pb';

export namespace ApiCooperation {
  export type Response = Cooperation;
  export type ListResponse = ListCooperationsResponse;
  export type FindParams = MessageInitShape<typeof ListCooperationsRequestSchema>;
  export type Invitation = CooperationInvitation;
  export type InvitationsResponse = ListCooperationInvitationsResponse;
}
