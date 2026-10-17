package rpc_service_reward_order_serde

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	reward_order_v1 "github.com/Roum1212/hahaton/genproto/reward-order/v1"
	aggregate_reward_order "github.com/Roum1212/hahaton/internal/domain/aggregate/reward-order"
	model_reward_order_status "github.com/Roum1212/hahaton/internal/domain/model/reward-order-status"
)

func AggregateToPBRewardOrder(value *aggregate_reward_order.RewardOrder) *reward_order_v1.RewardOrder {
	statuses := map[model_reward_order_status.RewardOrderStatus]reward_order_v1.RewardOrderStatus{
		model_reward_order_status.Pending:  reward_order_v1.RewardOrderStatus_REWARD_ORDER_STATUS_PENDING,
		model_reward_order_status.Issued:   reward_order_v1.RewardOrderStatus_REWARD_ORDER_STATUS_ISSUED,
		model_reward_order_status.Canceled: reward_order_v1.RewardOrderStatus_REWARD_ORDER_STATUS_CANCELLED,
	}

	return &reward_order_v1.RewardOrder{
		Id:                 value.GetRewardOrderID().String(),
		RewardId:           value.GetRewardID().String(),
		CoordinatorId:      value.GetCoordinatorID().String(),
		VolunteerId:        value.GetVolunteerID().String(),
		RewardTitle:        value.GetRewardTitle(),
		PriceCoins:         value.GetPriceCoins(),
		CancellationReason: CancellationReasonToPB(value.GetCancellationReason()),
		Status:             statuses[value.GetStatus()],
		CreatedAt:          timestamppb.New(value.GetCreatedAt()),
		UpdatedAt:          timestamppb.New(value.GetUpdatedAt()),
	}
}

func CancellationReasonToPB(value string) reward_order_v1.RewardOrderCancellationReason {
	switch value {
	case "cooperation_terminated":
		return reward_order_v1.RewardOrderCancellationReason_REWARD_ORDER_CANCELLATION_REASON_COOPERATION_TERMINATED
	case "coordinator_canceled":
		return reward_order_v1.RewardOrderCancellationReason_REWARD_ORDER_CANCELLATION_REASON_COORDINATOR_CANCELLED
	case "volunteer_canceled":
		return reward_order_v1.RewardOrderCancellationReason_REWARD_ORDER_CANCELLATION_REASON_VOLUNTEER_CANCELLED
	case "reward_deleted":
		return reward_order_v1.RewardOrderCancellationReason_REWARD_ORDER_CANCELLATION_REASON_REWARD_DELETED
	default:
		return reward_order_v1.RewardOrderCancellationReason_REWARD_ORDER_CANCELLATION_REASON_UNSPECIFIED
	}
}

func StatusFromPB(value reward_order_v1.RewardOrderStatus) (model_reward_order_status.RewardOrderStatus, error) {
	switch value {
	case reward_order_v1.RewardOrderStatus_REWARD_ORDER_STATUS_PENDING:
		return model_reward_order_status.NewRewardOrderStatus(model_reward_order_status.Pending.String())
	case reward_order_v1.RewardOrderStatus_REWARD_ORDER_STATUS_ISSUED:
		return model_reward_order_status.NewRewardOrderStatus(model_reward_order_status.Issued.String())
	case reward_order_v1.RewardOrderStatus_REWARD_ORDER_STATUS_CANCELLED:
		return model_reward_order_status.NewRewardOrderStatus(model_reward_order_status.Canceled.String())
	case reward_order_v1.RewardOrderStatus_REWARD_ORDER_STATUS_UNSPECIFIED:
		return "", errors.New("unspecified reward-order status")
	default:
		return "", fmt.Errorf("unknown reward-order status: %d", value)
	}
}
