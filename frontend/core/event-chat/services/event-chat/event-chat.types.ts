import type { IListResult } from '~/core/app/types';

export interface IEventChat {
  eventId: string;
  eventTitle: string;
  readOnly: boolean;
  latestSequence: number;
  visibleAfterSequence: number;
}

export interface IEventChatMessage {
  id: string;
  eventId: string;
  author: 'coordinator' | 'volunteer';
  authorId: string;
  authorDisplayName: string;
  text: string;
  sequence: number;
  createdAt: Date | null;
}

export interface IEventChatHistory extends IListResult<IEventChatMessage> {
  snapshotSequence: number;
  visibleAfterSequence: number;
}

export const EVENT_CHAT_MESSAGE_MAX_LENGTH = 2000;
