import { Code, ConnectError } from '@connectrpc/connect';

const RPC_ERROR_MESSAGES: Partial<Record<Code, string>> = {
  [Code.InvalidArgument]: 'Проверьте правильность заполнения полей',
  [Code.NotFound]: 'Запрашиваемые данные не найдены',
  [Code.AlreadyExists]: 'Такая запись уже существует',
  [Code.PermissionDenied]: 'Недостаточно прав для выполнения операции',
  [Code.Unauthenticated]: 'Необходимо войти в систему',
  [Code.ResourceExhausted]: 'Слишком много запросов, попробуйте позже',
  [Code.FailedPrecondition]: 'Действие сейчас недоступно',
  [Code.Aborted]: 'Данные изменились, обновите страницу',
  [Code.DeadlineExceeded]: 'Превышено время ожидания ответа',
  [Code.Unavailable]: 'Сервер недоступен, попробуйте позже',
  [Code.Internal]: 'Произошла ошибка на сервере, попробуйте позже',
  [Code.Unknown]: 'Произошла ошибка',
};

const EVENT_DATES_MESSAGE = 'Проверьте даты: начало в будущем, длительность кратна 30 минутам';

export const RPC_METHOD_ERROR_MESSAGES: Record<string, Partial<Record<Code, string>>> = {
  AuthenticateSession: {
    [Code.FailedPrecondition]: 'Во Флида ID не хватает данных профиля для регистрации',
  },
  AcceptCooperationInvitation: {
    [Code.NotFound]: 'Приглашение не найдено или отозвано',
    [Code.FailedPrecondition]: 'Срок действия приглашения истёк',
    [Code.InvalidArgument]: 'Нельзя принять собственное приглашение',
  },
  RequestCooperation: {
    [Code.AlreadyExists]: 'Заявка уже подана или сотрудничество активно',
    [Code.FailedPrecondition]: 'Повторная заявка пока недоступна',
    [Code.InvalidArgument]: 'Нельзя сотрудничать с самим собой',
  },
  ApproveCooperation: { [Code.FailedPrecondition]: 'Заявка уже рассмотрена' },
  RejectCooperation: { [Code.FailedPrecondition]: 'Заявка уже рассмотрена' },
  TerminateCooperation: { [Code.FailedPrecondition]: 'Сотрудничество уже не активно' },
  CreateCooperationInvitation: { [Code.InvalidArgument]: 'Срок действия должен быть в будущем' },
  ListEvents: { [Code.FailedPrecondition]: 'Для этого города нет границы, выберите поиск по радиусу' },
  CreateEvent: { [Code.InvalidArgument]: EVENT_DATES_MESSAGE },
  UpdateEvent: {
    [Code.InvalidArgument]: EVENT_DATES_MESSAGE,
    [Code.FailedPrecondition]: 'Изменить нельзя: есть активные заявки или мероприятие уже началось',
  },
  CancelEvent: { [Code.FailedPrecondition]: 'Отменить можно только до начала мероприятия' },
  CompleteEvent: {
    [Code.InvalidArgument]: 'Проверьте итоги: нужна строка для каждого принятого участника, минуты кратны 30',
    [Code.FailedPrecondition]: 'Завершить можно только после окончания мероприятия и только один раз',
    [Code.AlreadyExists]: 'Итоги уже отправлены с другими данными',
  },
  RequestParticipation: {
    [Code.AlreadyExists]: 'Вы уже подавали заявку на это мероприятие',
    [Code.FailedPrecondition]: 'Запись сейчас недоступна: мероприятие началось или нет сотрудничества',
  },
  RestoreParticipation: { [Code.FailedPrecondition]: 'Восстановить заявку сейчас нельзя' },
  CancelParticipation: { [Code.FailedPrecondition]: 'Отменить заявку можно только до начала мероприятия' },
  ApproveParticipation: { [Code.FailedPrecondition]: 'Принять нельзя: нет свободных мест или мероприятие уже началось' },
  RejectParticipation: { [Code.FailedPrecondition]: 'Эту заявку уже нельзя отклонить' },
  CreateRewardOrder: {
    [Code.Aborted]: 'Цена изменилась, проверьте карточку поощрения',
    [Code.FailedPrecondition]: 'Недостаточно койнов или поощрение закончилось',
    [Code.AlreadyExists]: 'Эта покупка уже оформлена',
  },
  IssueRewardOrder: { [Code.FailedPrecondition]: 'Этот заказ нельзя выдать' },
  CancelRewardOrder: { [Code.FailedPrecondition]: 'Выданный заказ нельзя отменить' },
  UpdateReward: {
    [Code.FailedPrecondition]: 'Количество меньше выданных и зарезервированных единиц',
    [Code.NotFound]: 'Изображение не найдено, загрузите его заново',
    [Code.PermissionDenied]: 'Это изображение загружено другим участником',
  },
  SendEventChatMessage: { [Code.FailedPrecondition]: 'Чат доступен только для чтения' },
  SubscribeEventChatMessages: { [Code.PermissionDenied]: 'Доступ к чату закрыт' },
  UploadImage: {
    [Code.InvalidArgument]: 'Файл не подходит: нужен JPEG или PNG до 5 МиБ и не больше 4096 пикселей по стороне',
    [Code.Unauthenticated]: 'Войдите, чтобы загрузить изображение',
  },
  UpdateVolunteer: {
    [Code.NotFound]: 'Изображение не найдено, загрузите его заново',
    [Code.PermissionDenied]: 'Это изображение загружено другим участником',
  },
  UpdateCoordinator: {
    [Code.NotFound]: 'Изображение не найдено, загрузите его заново',
    [Code.PermissionDenied]: 'Это изображение загружено другим участником',
  },
  CreateReward: {
    [Code.NotFound]: 'Изображение не найдено, загрузите его заново',
    [Code.PermissionDenied]: 'Это изображение загружено другим участником',
  },
  SuggestAddresses: {
    [Code.Unavailable]: 'Сервис подсказок адресов недоступен',
    [Code.ResourceExhausted]: 'Лимит подсказок адресов исчерпан, попробуйте позже',
  },
};

export function getRpcErrorMessage(error: unknown, overrides: Partial<Record<Code, string>> = {}): string {
  const { code } = ConnectError.from(error);
  return overrides[code] ?? RPC_ERROR_MESSAGES[code] ?? 'Произошла ошибка';
}

export function isRpcError(error: unknown, code: Code): boolean {
  return error instanceof ConnectError && error.code === code;
}
