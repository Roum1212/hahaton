import type { Transport } from '@connectrpc/connect';
import { BaseRpcApi } from '~/core/app/rpc/base-api';
import { EventChatService } from '#genproto/event-chat/v1/event_chat_service_pb';
import type { ApiEventChat } from './event-chat.api.types';

export class EventChatApi extends BaseRpcApi<typeof EventChatService> {
  constructor(transport: Transport) {
    super(EventChatService, transport);
  }

  async get(eventId: string, signal?: AbortSignal): Promise<ApiEventChat.Chat | undefined> {
    const { chat } = await this.client.getEventChat({ eventId }, this.options(signal));
    return chat;
  }

  async findMine(page: { limit: number; offset: bigint }, signal?: AbortSignal): Promise<ApiEventChat.ChatsResponse> {
    return this.client.listMyEventChats(page, this.options(signal));
  }

  async send(eventId: string, text: string, requestId: string): Promise<ApiEventChat.Message | undefined> {
    const { message } = await this.client.sendEventChatMessage({ eventId, text, requestId });
    return message;
  }

  async findMessages(
    eventId: string,
    snapshotSequence: bigint | undefined,
    page: { limit: number; offset: bigint },
    signal?: AbortSignal,
  ): Promise<ApiEventChat.MessagesResponse> {
    return this.client.listEventChatMessages({ eventId, snapshotSequence, ...page }, this.options(signal));
  }
}
