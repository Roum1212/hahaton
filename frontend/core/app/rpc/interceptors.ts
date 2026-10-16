import { Code, ConnectError, type Interceptor } from '@connectrpc/connect';

const SILENT_CODES = [Code.Canceled, Code.Unauthenticated];

export const createErrorInterceptor =
  (onError: (error: ConnectError, method: string) => void, onUnauthenticated: () => void): Interceptor =>
  (next) =>
  async (req) => {
    try {
      return await next(req);
    } catch (e) {
      const error = ConnectError.from(e);
      if (error.code === Code.Unauthenticated) onUnauthenticated();
      if (!SILENT_CODES.includes(error.code)) onError(error, req.method.name);
      throw error;
    }
  };
