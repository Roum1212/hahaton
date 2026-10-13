export namespace Flida {
  export type FedCMMode = 'active' | 'passive';
  export type AuthorizationMode = 'redirect' | 'popup';

  export type ErrorCode =
    | 'invalid_config'
    | 'invalid_state'
    | 'transaction_exists'
    | 'crypto_unavailable'
    | 'fedcm_unavailable'
    | 'authorization_in_progress'
    | 'popup_blocked'
    | 'popup_closed'
    | 'timeout'
    | 'invalid_response'
    | 'http_error'
    | 'oauth_error';

  export interface SDKConfig {
    clientId: string;
    redirectUri: string;
    fedcm?: FedCMMode[];
    mode?: AuthorizationMode;
    popupWidth?: number;
    popupHeight?: number;
  }

  export interface AuthorizeOptions {
    signal?: AbortSignal;
  }

  export interface TokenResponse {
    access_token: string;
    token_type: string;
    expires_in?: number;
    refresh_token?: string;
    id_token?: string;
  }

  export interface UserInfoResponse {
    sub: string;
    given_name?: string;
    family_name?: string;
    middle_name?: string;
    nickname?: string;
    date_of_birth?: string;
    sex?: 'male' | 'female';
    email?: string;
    email_verified?: boolean;
    phone_number?: string;
    phone_number_verified?: boolean;
  }

  export interface SDK {
    readonly autoAuthorization: Promise<TokenResponse | undefined>;
    authorize(options?: AuthorizeOptions): Promise<TokenResponse | undefined>;
    getUser(accessToken: string, signal?: AbortSignal): Promise<UserInfoResponse>;
    refreshToken(refreshToken: string, signal?: AbortSignal): Promise<TokenResponse>;
  }

  export interface Error extends globalThis.Error {
    code: ErrorCode;
    status?: number;
  }

  export interface Module {
    FlidaSDK: new (config: SDKConfig) => SDK;
    FlidaError: new (...args: never[]) => Error;
  }
}
