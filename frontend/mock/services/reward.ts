import { Code, type ConnectRouter } from '@connectrpc/connect';
import { RewardOrderCancellationReason, RewardOrderStatus } from '../../genproto/reward-order/v1/reward_order_service_pb';
import { RewardListScope, RewardService } from '../../genproto/reward/v1/reward_service_pb';
import { activePartnerIds, hasActiveCooperation, requireCoordinator, requireUser } from '../auth';
import { compareStrings, db, fail, newId, paginate, type RewardRec } from '../db';
import { cancelRewardOrder, issuedQuantity, reservedQuantity } from '../domain';
import { toReward } from '../mappers';

export const requireOwnReward = (id: string, coordinatorId: string): RewardRec => {
  const reward = db.rewards.get(id);
  if (!reward) return fail(Code.NotFound, 'Поощрение не найдено');
  if (reward.coordinatorId !== coordinatorId) return fail(Code.PermissionDenied, 'Чужое поощрение');
  return reward;
};

export const registerRewardService = (router: ConnectRouter) =>
  router.service(RewardService, {
    createReward(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const title = req.title.trim();
      if (!title) return fail(Code.InvalidArgument, 'Укажите название');
      if (req.priceCoins <= 0n) return fail(Code.InvalidArgument, 'Цена должна быть больше нуля');
      if (req.quantityLimit <= 0) return fail(Code.InvalidArgument, 'Количество должно быть больше нуля');
      const reward: RewardRec = {
        id: newId(),
        coordinatorId: coordinator.id,
        title,
        description: req.description,
        priceCoins: Number(req.priceCoins),
        quantityLimit: req.quantityLimit,
        deleted: false,
      };
      db.rewards.set(reward.id, reward);
      return { id: reward.id };
    },
    getReward(req, ctx) {
      const user = requireUser(ctx);
      const reward = db.rewards.get(req.id);
      if (!reward || reward.deleted) return fail(Code.NotFound, 'Поощрение не найдено');
      if (reward.coordinatorId !== user.id && !hasActiveCooperation(reward.coordinatorId, user.id)) {
        return fail(Code.PermissionDenied, 'Нет доступа к поощрению');
      }
      return { reward: toReward(reward) };
    },
    listRewards(req, ctx) {
      const user = requireUser(ctx);
      let owners: Set<string>;
      if (req.scope === RewardListScope.CATALOG) owners = activePartnerIds(user.id);
      else if (req.scope === RewardListScope.OWNED) owners = new Set([requireCoordinator(ctx).id]);
      else return fail(Code.InvalidArgument, 'Не указана область выборки');
      if (req.coordinatorId !== undefined && !req.coordinatorId) return fail(Code.InvalidArgument, 'Пустой фильтр координатора');
      const items = [...db.rewards.values()]
        .filter((item) => !item.deleted && owners.has(item.coordinatorId))
        .filter((item) => !req.coordinatorId || item.coordinatorId === req.coordinatorId)
        .sort((a, b) => compareStrings(a.title, b.title) || compareStrings(a.id, b.id));
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return { rewards: page.map(toReward), totalCount };
    },
    updateReward(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const reward = requireOwnReward(req.id, coordinator.id);
      if (reward.deleted) return fail(Code.NotFound, 'Поощрение не найдено');
      if (req.title !== undefined && !req.title.trim()) return fail(Code.InvalidArgument, 'Название не может быть пустым');
      if (req.priceCoins !== undefined && req.priceCoins <= 0n) return fail(Code.InvalidArgument, 'Цена должна быть больше нуля');
      if (req.quantityLimit !== undefined) {
        if (req.quantityLimit < 0) return fail(Code.InvalidArgument, 'Количество не может быть отрицательным');
        if (req.quantityLimit < issuedQuantity(reward.id) + reservedQuantity(reward.id)) {
          return fail(Code.FailedPrecondition, 'Количество меньше выданных и зарезервированных единиц');
        }
      }
      if (req.title !== undefined) reward.title = req.title.trim();
      if (req.description !== undefined) reward.description = req.description;
      if (req.priceCoins !== undefined) reward.priceCoins = Number(req.priceCoins);
      if (req.quantityLimit !== undefined) reward.quantityLimit = req.quantityLimit;
      return {};
    },
    deleteReward(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const reward = requireOwnReward(req.id, coordinator.id);
      if (reward.deleted) return {};
      reward.deleted = true;
      db.rewardOrders.forEach((order) => {
        if (order.rewardId === reward.id && order.status === RewardOrderStatus.PENDING) {
          cancelRewardOrder(order, RewardOrderCancellationReason.REWARD_DELETED);
        }
      });
      return {};
    },
  });
