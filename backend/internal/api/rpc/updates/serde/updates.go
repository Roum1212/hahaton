package rpc_updates_serde

import (
	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	coordinator_v1 "github.com/Roum1212/hahaton/genproto/coordinator/v1"
	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	reward_order_v1 "github.com/Roum1212/hahaton/genproto/reward-order/v1"
	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
	service_get_update_recipients "github.com/Roum1212/hahaton/internal/app/service/get-update-recipients"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

func RequestFromRPC(
	userID model_user_id.UserID,
	r any,
	response any,
) (service_get_update_recipients.Request, error) {
	counterpart, err := optionalID(counterpartID(userID, r, response), model_user_id.NewUserID)
	if err != nil {
		return service_get_update_recipients.Request{}, err
	}

	var rawEventID, rawParticipationID, rawRewardID, rawOrderID string
	if value, ok := r.(interface{ GetEventId() string }); ok {
		rawEventID = value.GetEventId()
	}

	if value, ok := r.(interface{ GetRewardId() string }); ok {
		rawRewardID = value.GetRewardId()
	}

	switch value := r.(type) {
	case *event_v1.CancelEventRequest:
		rawEventID = value.GetId()
	case *event_v1.CompleteEventRequest:
		rawEventID = value.GetId()
	case *event_v1.UpdateEventRequest:
		rawEventID = value.GetId()
	case *participation_v1.ApproveParticipationRequest:
		rawParticipationID = value.GetId()
	case *participation_v1.CancelParticipationRequest:
		rawParticipationID = value.GetId()
	case *participation_v1.RejectParticipationRequest:
		rawParticipationID = value.GetId()
	case *participation_v1.RestoreParticipationRequest:
		rawParticipationID = value.GetId()
	case *reward_v1.DeleteRewardRequest:
		rawRewardID = value.GetId()
	case *reward_v1.UpdateRewardRequest:
		rawRewardID = value.GetId()
	case *reward_order_v1.CancelRewardOrderRequest:
		rawOrderID = value.GetId()
	case *reward_order_v1.IssueRewardOrderRequest:
		rawOrderID = value.GetId()
	}

	eventID, err := optionalID(rawEventID, model_event_id.NewEventID)
	if err != nil {
		return service_get_update_recipients.Request{}, err
	}

	participationID, err := optionalID(rawParticipationID, model_participation_id.NewParticipationID)
	if err != nil {
		return service_get_update_recipients.Request{}, err
	}

	rewardID, err := optionalID(rawRewardID, model_reward_id.NewRewardID)
	if err != nil {
		return service_get_update_recipients.Request{}, err
	}

	orderID, err := optionalID(rawOrderID, model_reward_order_id.NewRewardOrderID)
	if err != nil {
		return service_get_update_recipients.Request{}, err
	}

	return service_get_update_recipients.NewRequest(userID, counterpart, eventID, participationID, rewardID, orderID), nil
}

func counterpartID(
	userID model_user_id.UserID,
	r any,
	response any,
) string {
	if value, ok := response.(interface {
		GetCooperation() *cooperation_v1.Cooperation
	}); ok {
		cooperation := value.GetCooperation()
		if cooperation.GetCoordinatorId() == userID.String() {
			return cooperation.GetVolunteerId()
		}

		return cooperation.GetCoordinatorId()
	}

	if value, ok := response.(interface {
		GetCoordinator() *coordinator_v1.Coordinator
	}); ok {
		return value.GetCoordinator().GetId()
	}

	if value, ok := r.(interface{ GetCoordinatorId() string }); ok && value.GetCoordinatorId() != "" {
		return value.GetCoordinatorId()
	}

	if value, ok := r.(interface{ GetVolunteerId() string }); ok && value.GetVolunteerId() != "" {
		return value.GetVolunteerId()
	}

	return ""
}

func optionalID[T any](value string, parse func(value string) (T, error)) (T, error) {
	if value == "" {
		var zero T

		return zero, nil
	}

	return parse(value)
}
