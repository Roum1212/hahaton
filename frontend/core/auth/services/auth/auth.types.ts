import type { IUser } from '~/core/users/services/users/users.types';
import type { Flida } from '../../api/flida/flida.types';

export type IAuthUser = IUser;

export type IFlidaTokens = Flida.TokenResponse;
export type IFlidaError = Flida.Error;
