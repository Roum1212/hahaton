import { ConnectError } from '@connectrpc/connect';
import type { ApiPagination, IListResult } from '~/core/app/types';
import { EVENT_CHAT_API } from '../..';
import type { ApiEventChat } from '../../api/event-chat/event-chat.api.types';
import type { IEventChat, IEventChatHistory, IEventChatMessage } from './event-chat.types';

export class EventChatService {
  private service = EVENT_CHAT_API().eventChat;

  async get(eventId: string, signal?: AbortSignal): Promise<IEventChat> {
    const chat = await this.service.get(eventId, signal);
    if (!chat) throw new ConnectError('Пустой ответ сервера');
    return this.toChat(chat);
  }

  async findMine(pagination: ApiPagination, signal?: AbortSignal): Promise<IListResult<IEventChat>> {
    const { chats, totalCount } = await this.service.findMine(toApiPage(pagination), signal);
    return { list: chats.map((item) => this.toChat(item)), total: fromBigInt(totalCount) };
  }

  async send(eventId: string, text: string, requestId: string): Promise<IEventChatMessage> {
    const message = await this.service.send(eventId, text, requestId);
    if (!message) throw new ConnectError('Пустой ответ сервера');
    return this.toMessage(message);
  }

  async findMessages(
    eventId: string,
    snapshotSequence: number | undefined,
    pagination: ApiPagination,
    signal?: AbortSignal,
  ): Promise<IEventChatHistory> {
    const response = await this.service.findMessages(
      eventId,
      snapshotSequence === undefined ? undefined : BigInt(snapshotSequence),
      toApiPage(pagination),
      signal,
    );
    return {
      list: response.messages.map((item) => this.toMessage(item)),
      total: fromBigInt(response.totalCount),
      snapshotSequence: fromBigInt(response.snapshotSequence),
      visibleAfterSequence: fromBigInt(response.visibleAfterSequence),
    };
  }

  toChat(data: ApiEventChat.Chat): IEventChat {
    return {
      eventId: data.eventId,
      eventTitle: data.eventTitle,
      readOnly: data.readOnly,
      latestSequence: fromBigInt(data.latestSequence),
      visibleAfterSequence: fromBigInt(data.visibleAfterSequence),
    };
  }

  toMessage(data: ApiEventChat.Message): IEventChatMessage {
    return {
      id: data.id,
      eventId: data.eventId,
      author: data.author.case === 'coordinatorId' ? 'coordinator' : 'volunteer',
      authorId: data.author.value ?? '',
      authorDisplayName: data.authorDisplayName,
      text: data.text,
      sequence: fromBigInt(data.sequence),
      createdAt: fromTimestamp(data.createdAt),
    };
  }
}
