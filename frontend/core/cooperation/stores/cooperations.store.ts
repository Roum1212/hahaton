import { defineStore } from 'pinia';
import {
  CooperationListScope,
  CooperationStatus,
  type ICooperation,
  type ICooperationInvitation,
} from '../services/cooperation/cooperation.types';

interface IStore {
  asVolunteer: ICooperation[];
  asCoordinator: ICooperation[];
  invitations: ICooperationInvitation[];
}

const LOAD_LIMIT = 100;

export const useCooperationsStore = defineStore('cooperations', {
  state: (): IStore => ({
    asVolunteer: [],
    asCoordinator: [],
    invitations: [],
  }),
  getters: {
    service: () => useServices().COOPERATION.CooperationService,
    activeCoordinatorIds: (state) =>
      state.asVolunteer.filter((item) => item.status === CooperationStatus.ACTIVE).map((item) => item.coordinatorId),
    withCoordinator: (state) => (coordinatorId: string) => state.asVolunteer.find((item) => item.coordinatorId === coordinatorId),
    withVolunteer: (state) => (volunteerId: string) => state.asCoordinator.find((item) => item.volunteerId === volunteerId),
  },
  actions: {
    async fetchAll(scope: CooperationListScope) {
      const list: ICooperation[] = [];
      for (let offset = 0; ; offset += LOAD_LIMIT) {
        const page = await this.service.find({ scope }, { limit: LOAD_LIMIT, offset });
        list.push(...page.list);
        if (list.length >= page.total || !page.list.length) break;
      }
      if (scope === CooperationListScope.AS_VOLUNTEER) this.asVolunteer = list;
      else this.asCoordinator = list;
      return list;
    },

    async request(coordinatorId: string) {
      const cooperation = await this.service.request(coordinatorId);
      this.upsertAsVolunteer(cooperation);
      return cooperation;
    },

    async acceptInvitation(token: string) {
      const cooperation = await this.service.acceptInvitation(token);
      this.upsertAsVolunteer(cooperation);
      return cooperation;
    },

    async terminateAsVolunteer(coordinatorId: string) {
      await this.service.terminateAsVolunteer(coordinatorId);
      this.upsertAsVolunteer(await this.service.getWithCoordinator(coordinatorId));
    },

    async approve(volunteerId: string) {
      await this.service.approve(volunteerId);
      this.upsertAsCoordinator(await this.service.getWithVolunteer(volunteerId));
    },

    async reject(volunteerId: string) {
      await this.service.reject(volunteerId);
      this.upsertAsCoordinator(await this.service.getWithVolunteer(volunteerId));
    },

    async terminateAsCoordinator(volunteerId: string) {
      await this.service.terminateAsCoordinator(volunteerId);
      this.upsertAsCoordinator(await this.service.getWithVolunteer(volunteerId));
    },

    async fetchInvitations(includeExpired = false) {
      const { list } = await this.service.findInvitations(includeExpired, { limit: LOAD_LIMIT, offset: 0 });
      this.invitations = list;
      return list;
    },

    async createInvitation(expiresAt: Date | null) {
      const invitation = await this.service.createInvitation(expiresAt);
      this.invitations = [invitation, ...this.invitations];
      return invitation;
    },

    async revokeInvitation(id: string) {
      await this.service.revokeInvitation(id);
      this.invitations = this.invitations.filter((item) => item.id !== id);
    },

    upsertAsVolunteer(cooperation: ICooperation) {
      this.asVolunteer = [...this.asVolunteer.filter((item) => item.coordinatorId !== cooperation.coordinatorId), cooperation];
    },

    upsertAsCoordinator(cooperation: ICooperation) {
      this.asCoordinator = [...this.asCoordinator.filter((item) => item.volunteerId !== cooperation.volunteerId), cooperation];
    },

    $reset() {
      this.asVolunteer = [];
      this.asCoordinator = [];
      this.invitations = [];
    },
  },
});
