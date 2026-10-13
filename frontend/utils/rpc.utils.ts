import { Code, ConnectError } from '@connectrpc/connect';

const RPC_ERROR_MESSAGES: Partial<Record<Code, string>> = {
  [Code.InvalidArgument]: 'Некорректные данные запроса',
  [Code.NotFound]: 'Запрашиваемые данные не найдены',
  [Code.AlreadyExists]: 'Запись уже существует',
  [Code.PermissionDenied]: 'Недостаточно прав для выполнения операции',
  [Code.Unauthenticated]: 'Необходимо авторизоваться',
  [Code.ResourceExhausted]: 'Слишком много запросов, попробуйте позже',
  [Code.FailedPrecondition]: 'Операция не может быть выполнена',
  [Code.DeadlineExceeded]: 'Превышено время ожидания ответа',
  [Code.Unavailable]: 'Сервер недоступен, попробуйте позже',
  [Code.Internal]: 'Произошла ошибка на сервере, попробуйте позже',
  [Code.Unknown]: 'Произошла ошибка',
};

const SERVER_ERROR_CODES = [Code.Internal, Code.Unknown, Code.Unavailable, Code.DataLoss];

export function getRpcErrorMessage(error: unknown): string {
  const rpcError = ConnectError.from(error);
  if (rpcError.rawMessage && !SERVER_ERROR_CODES.includes(rpcError.code)) return rpcError.rawMessage;
  return RPC_ERROR_MESSAGES[rpcError.code] ?? 'Произошла ошибка';
}

export function isRpcError(error: unknown, code: Code): boolean {
  return error instanceof ConnectError && error.code === code;
}
