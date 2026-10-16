import { Code, type ConnectRouter } from '@connectrpc/connect';
import {
  RewardOrderCancellationReason,
  RewardOrderListScope,
  RewardOrderService,
  RewardOrderStatus,
} from '../../genproto/reward-order/v1/reward_order_service_pb';
import { VolunteerOperationType } from '../../genproto/volunteer/v1/volunteer_service_pb';
import { hasActiveCooperation, requireCoordinator, requireUser, touch } from '../auth';
import { compareStrings, db, fail, newId, now, paginate, requireFound, type RewardOrderRec } from '../db';
import { cancelRewardOrder, issuedQuantity, recordOperation, reservedQuantity } from '../domain';
import { toRewardOrder } from '../mappers';

const requireOrder = (id: string) => requireFound(db.rewardOrders.get(id), 'Заказ не найден');

export const registerRewardOrderService = (router: ConnectRouter) =>
  router.service(RewardOrderService, {
    createRewardOrder(req, ctx) {
      const user = requireUser(ctx);
      if (!req.requestId) return fail(Code.InvalidArgument, 'Нужен request_id');
      if (req.expectedPriceCoins <= 0n) return fail(Code.InvalidArgument, 'Ожидаемая цена должна быть больше нуля');
      const repeated = [...db.rewardOrders.values()].find((item) => item.volunteerId === user.id && item.requestId === req.requestId);
      if (repeated) {
        if (repeated.rewardId !== req.rewardId || BigInt(repeated.priceCoins) !== req.expectedPriceCoins) {
          return fail(Code.AlreadyExists, 'request_id уже использован для другой покупки');
        }
        return { rewardOrder: toRewardOrder(repeated) };
      }
      const reward = db.rewards.get(req.rewardId);
      if (!reward || reward.deleted) return fail(Code.NotFound, 'Поощрение не найдено');
      if (reward.coordinatorId === user.id) return fail(Code.FailedPrecondition, 'Нельзя купить своё поощрение');
      if (!hasActiveCooperation(reward.coordinatorId, user.id)) return fail(Code.FailedPrecondition, 'Нужно активное сотрудничество');
      if (BigInt(reward.priceCoins) !== req.expectedPriceCoins) return fail(Code.Aborted, 'Цена поощрения изменилась');
      if (reward.quantityLimit - issuedQuantity(reward.id) - reservedQuantity(reward.id) <= 0) {
        return fail(Code.FailedPrecondition, 'Поощрение закончилось');
      }
      if (user.coinBalance < reward.priceCoins) return fail(Code.FailedPrecondition, 'Недостаточно коинов');
      const at = now();
      const order: RewardOrderRec = {
        id: newId(),
        coordinatorId: reward.coordinatorId,
        rewardId: reward.id,
        volunteerId: user.id,
        priceCoins: reward.priceCoins,
        rewardTitle: reward.title,
        status: RewardOrderStatus.PENDING,
        cancellationReason: RewardOrderCancellationReason.UNSPECIFIED,
        requestId: req.requestId,
        createdAt: at,
        updatedAt: at,
      };
      recordOperation({
        volunteerId: user.id,
        coordinatorId: reward.coordinatorId,
        type: VolunteerOperationType.REWARD_PURCHASE,
        source: { case: 'rewardOrderId', value: order.id },
        description: `Покупка «${reward.title}»`,
        coinDelta: -reward.priceCoins,
        minuteDelta: 0,
        createdAt: at,
      });
      db.rewardOrders.set(order.id, order);
      return { rewardOrder: toRewardOrder(order) };
    },
    getRewardOrder(req, ctx) {
      const user = requireUser(ctx);
      const order = requireOrder(req.id);
      if (order.volunteerId !== user.id && order.coordinatorId !== user.id) return fail(Code.PermissionDenied, 'Чужой заказ');
      return { rewardOrder: toRewardOrder(order) };
    },
    listRewardOrders(req, ctx) {
      const user = requireUser(ctx);
      let filter: (item: RewardOrderRec) => boolean;
      if (req.scope === RewardOrderListScope.MY_PURCHASES) filter = (item) => item.volunteerId === user.id;
      else if (req.scope === RewardOrderListScope.INCOMING) {
        const coordinator = requireCoordinator(ctx);
        filter = (item) => item.coordinatorId === coordinator.id;
      } else return fail(Code.InvalidArgument, 'Не указана область выборки');
      for (const value of [req.coordinatorId, req.rewardId, req.volunteerId]) {
        if (value !== undefined && !value) return fail(Code.InvalidArgument, 'Пустой фильтр');
      }
      if (req.status === RewardOrderStatus.UNSPECIFIED) return fail(Code.InvalidArgument, 'Некорректный статус');
      const items = [...db.rewardOrders.values()]
        .filter(filter)
        .filter((item) => !req.coordinatorId || item.coordinatorId === req.coordinatorId)
        .filter((item) => !req.rewardId || item.rewardId === req.rewardId)
        .filter((item) => !req.volunteerId || item.volunteerId === req.volunteerId)
        .filter((item) => req.status === undefined || item.status === req.status)
        .sort((a, b) => b.createdAt.getTime() - a.createdAt.getTime() || compareStrings(b.id, a.id));
      const { page, totalCount } = paginate(items, req.limit, req.offset);
      return { rewardOrders: page.map(toRewardOrder), totalCount };
    },
    issueRewardOrder(req, ctx) {
      const coordinator = requireCoordinator(ctx);
      const order = requireOrder(req.id);
      if (order.coordinatorId !== coordinator.id) return fail(Code.PermissionDenied, 'Чужой заказ');
      if (order.status === RewardOrderStatus.ISSUED) return {};
      if (order.status === RewardOrderStatus.CANCELLED) return fail(Code.FailedPrecondition, 'Отменённый заказ нельзя выдать');
      if (!hasActiveCooperation(coordinator.id, order.volunteerId)) return fail(Code.FailedPrecondition, 'Нет активного сотрудничества');
      order.status = RewardOrderStatus.ISSUED;
      touch(order);
      return {};
    },
    cancelRewardOrder(req, ctx) {
      const user = requireUser(ctx);
      const order = requireOrder(req.id);
      if (order.volunteerId !== user.id && order.coordinatorId !== user.id) return fail(Code.PermissionDenied, 'Чужой заказ');
      if (order.coordinatorId === user.id) requireCoordinator(ctx);
      const reason =
        order.coordinatorId === user.id
          ? RewardOrderCancellationReason.COORDINATOR_CANCELLED
          : RewardOrderCancellationReason.VOLUNTEER_CANCELLED;
      cancelRewardOrder(order, reason);
      return {};
    },
  });
