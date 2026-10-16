import type { MessageInitShape } from '@bufbuild/protobuf';
import type {
  CoordinatorVolunteer,
  ListVolunteerOperationsResponse,
  UpdateVolunteerRequestSchema,
  Volunteer,
  VolunteerOperation,
  VolunteerOperationType,
} from '#genproto/volunteer/v1/volunteer_service_pb';

export namespace ApiVolunteer {
  export type Response = Volunteer;
  export type CoordinatorView = CoordinatorVolunteer;
  export type Operation = VolunteerOperation;
  export type OperationsResponse = ListVolunteerOperationsResponse;
  export type OperationType = VolunteerOperationType;
  export type UpdateRequest = MessageInitShape<typeof UpdateVolunteerRequestSchema>;
}
