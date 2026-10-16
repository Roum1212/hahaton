import { EventChatApi } from './api/event-chat/event-chat.api';
import { EventChatService } from './services/event-chat/event-chat.service';

export const EVENT_CHAT_API = () => {
  const { $transport } = useNuxtApp();
  return {
    eventChat: new EventChatApi($transport),
  };
};

export const EVENT_CHAT_SERVICES = () => {
  return {
    EventChatService: new EventChatService(),
  };
};
