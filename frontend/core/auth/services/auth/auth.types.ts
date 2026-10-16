import type { Flida } from '../../api/flida/flida.types';

export interface ISession {
  userId: string | null;
  hasCoordinatorProfile: boolean;
  expiresAt: Date | null;
}

export type IFlidaTokens = Flida.TokenResponse;
export type IFlidaError = Flida.Error;
