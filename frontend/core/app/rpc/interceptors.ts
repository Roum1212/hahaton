import { Code, ConnectError, type Interceptor } from '@connectrpc/connect';

export const createAuthInterceptor =
  (getToken: () => string | null): Interceptor =>
  (next) =>
  async (req) => {
    const token = getToken();
    if (token) req.header.set('Authorization', `Bearer ${token}`);
    return next(req);
  };

export const createRefreshInterceptor =
  (refresh: () => Promise<string | null>): Interceptor =>
  (next) =>
  async (req) => {
    try {
      return await next(req);
    } catch (e) {
      if (req.stream || ConnectError.from(e).code !== Code.Unauthenticated) throw e;
      const token = await refresh();
      if (!token) throw e;
      return next(req);
    }
  };

export const createErrorInterceptor =
  (onError: (error: ConnectError) => void): Interceptor =>
  (next) =>
  async (req) => {
    try {
      return await next(req);
    } catch (e) {
      const error = ConnectError.from(e);
      if (error.code !== Code.Canceled) onError(error);
      throw error;
    }
  };
