import { create } from '@bufbuild/protobuf';
import { Code, type ConnectRouter } from '@connectrpc/connect';
import { CoordinatorService, CoordinatorSettingsSchema } from '../../genproto/coordinator/v1/coordinator_service_pb';
import { requireCoordinator, requireUser } from '../auth';
import { compareStrings, db, fail, paginate, requireFound } from '../db';
import { resolveImagePath } from '../images';
import { toCoordinator } from '../mappers';

export const registerCoordinatorService = (router: ConnectRouter) =>
  router.service(CoordinatorService, {
    getCoordinator(req, ctx) {
      const user = requireUser(ctx);
      if (req.id !== undefined && !req.id) return fail(Code.InvalidArgument, 'Пустой id');
      const coordinator = requireFound(db.coordinators.get(req.id ?? user.id), 'Координатор не найден');
      return { coordinator: toCoordinator(coordinator) };
    },
    listCoordinators(req, ctx) {
      requireUser(ctx);
      const items = [...db.coordinators.values()]
        .map(toCoordinator)
        .sort((a, b) => compareStrings(a.displayName, b.displayName) || compareStrings(a.id, b.id));
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return { coordinators: page, totalCount };
    },
    updateCoordinator(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      if (req.avatarPath !== undefined) coordinator.avatarPath = resolveImagePath(req.avatarPath, coordinator.id);
      if (req.description !== undefined) coordinator.description = req.description;
      return {};
    },
    getCoordinatorSettings(_, ctx) {
      const coordinator = requireCoordinator(ctx);
      return {
        settings: create(CoordinatorSettingsSchema, { cooperationRequestCooldownSeconds: BigInt(coordinator.cooldownSeconds) }),
      };
    },
    updateCoordinatorSettings(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      if (req.cooperationRequestCooldownSeconds !== undefined) {
        if (req.cooperationRequestCooldownSeconds < 0n) return fail(Code.InvalidArgument, 'Период ожидания не может быть отрицательным');
        coordinator.cooldownSeconds = Number(req.cooperationRequestCooldownSeconds);
      }
      return {};
    },
  });
