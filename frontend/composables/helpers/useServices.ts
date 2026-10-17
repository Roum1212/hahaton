import { AUTH_SERVICES } from '~/core/auth';
import { COOPERATION_SERVICES } from '~/core/cooperation';
import { COORDINATOR_SERVICES } from '~/core/coordinator';
import { EVENT_SERVICES } from '~/core/event';
import { EVENT_CHAT_SERVICES } from '~/core/event-chat';
import { GEO_SERVICES } from '~/core/geo';
import { IMAGE_SERVICES } from '~/core/image';
import { PARTICIPATION_SERVICES } from '~/core/participation';
import { REWARD_SERVICES } from '~/core/reward';
import { REWARD_ORDER_SERVICES } from '~/core/reward-order';
import { UPDATES_SERVICES } from '~/core/updates';
import { VOLUNTEER_SERVICES } from '~/core/volunteer';

export const useServices = () => {
  return {
    AUTH: AUTH_SERVICES(),
    GEO: GEO_SERVICES(),
    IMAGE: IMAGE_SERVICES(),
    VOLUNTEER: VOLUNTEER_SERVICES(),
    COORDINATOR: COORDINATOR_SERVICES(),
    COOPERATION: COOPERATION_SERVICES(),
    EVENT: EVENT_SERVICES(),
    PARTICIPATION: PARTICIPATION_SERVICES(),
    REWARD: REWARD_SERVICES(),
    REWARD_ORDER: REWARD_ORDER_SERVICES(),
    EVENT_CHAT: EVENT_CHAT_SERVICES(),
    UPDATES: UPDATES_SERVICES(),
  };
};
