import type { MessageInitShape } from '@bufbuild/protobuf';
import type {
  CreateEventRequestSchema,
  Event,
  EventParticipationResultSchema,
  ListEventsRequestSchema,
  ListEventsResponse,
  UpdateEventRequestSchema,
} from '#genproto/event/v1/event_service_pb';

export namespace ApiEvent {
  export type Response = Event;
  export type ListResponse = ListEventsResponse;
  export type FindParams = MessageInitShape<typeof ListEventsRequestSchema>;
  export type CreateRequest = MessageInitShape<typeof CreateEventRequestSchema>;
  export type UpdateRequest = MessageInitShape<typeof UpdateEventRequestSchema>;
  export type Result = MessageInitShape<typeof EventParticipationResultSchema>;
}
