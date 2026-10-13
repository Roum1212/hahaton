import type { MessageInitShape } from '@bufbuild/protobuf';
import type { LoginWithFlidaRequestSchema } from '~/gen/auth/v1/auth_pb';
import type { User } from '~/gen/user/v1/user_pb';

export namespace ApiAuth {
  export type UserResponse = User;
  export type LoginWithFlidaRequest = MessageInitShape<typeof LoginWithFlidaRequestSchema>;
}
