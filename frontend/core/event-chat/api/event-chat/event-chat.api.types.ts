import type {
  EventChat,
  EventChatMessage,
  ListEventChatMessagesResponse,
  ListMyEventChatsResponse,
} from '#genproto/event-chat/v1/event_chat_service_pb';

export namespace ApiEventChat {
  export type Chat = EventChat;
  export type ChatsResponse = ListMyEventChatsResponse;
  export type Message = EventChatMessage;
  export type MessagesResponse = ListEventChatMessagesResponse;
}
