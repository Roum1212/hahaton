import type { MessageInitShape } from '@bufbuild/protobuf';
import type {
  CreateUserRequestSchema,
  ListUsersRequestSchema,
  ListUsersResponse,
  UpdateUserRequestSchema,
  User,
} from '~/gen/user/v1/user_pb';

export namespace ApiUser {
  export type Response = User;
  export type ListResponse = ListUsersResponse;
  export type FindParams = MessageInitShape<typeof ListUsersRequestSchema>;
  export type CreateRequest = MessageInitShape<typeof CreateUserRequestSchema>;
  export type UpdateRequest = MessageInitShape<typeof UpdateUserRequestSchema>;
}
