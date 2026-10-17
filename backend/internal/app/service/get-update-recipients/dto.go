package service_get_update_recipients

import (
	"errors"

	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_reward_order_id "github.com/Roum1212/hahaton/internal/domain/model/reward-order-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID          model_user_id.UserID
	counterpartID   model_user_id.UserID
	eventID         model_event_id.EventID
	participationID model_participation_id.ParticipationID
	rewardID        model_reward_id.RewardID
	rewardOrderID   model_reward_order_id.RewardOrderID
}

func (x Request) GetCounterpartID() model_user_id.UserID { return x.counterpartID }

func (x Request) GetEventID() model_event_id.EventID { return x.eventID }

func (x Request) GetParticipationID() model_participation_id.ParticipationID {
	return x.participationID
}

func (x Request) GetRewardID() model_reward_id.RewardID { return x.rewardID }

func (x Request) GetRewardOrderID() model_reward_order_id.RewardOrderID { return x.rewardOrderID }

func (x Request) GetUserID() model_user_id.UserID { return x.userID }

func (x Request) Validate() error {
	if x.userID.IsZero() && x.counterpartID.IsZero() && x.eventID.IsZero() {
		return errors.New("update recipient is required")
	}

	return nil
}

func NewRequest(
	userID model_user_id.UserID,
	counterpartID model_user_id.UserID,
	eventID model_event_id.EventID,
	participationID model_participation_id.ParticipationID,
	rewardID model_reward_id.RewardID,
	rewardOrderID model_reward_order_id.RewardOrderID,
) Request {
	return Request{
		userID:          userID,
		counterpartID:   counterpartID,
		eventID:         eventID,
		participationID: participationID,
		rewardID:        rewardID,
		rewardOrderID:   rewardOrderID,
	}
}

type Response struct {
	recipients []model_user_id.UserID
}

func (x Response) GetRecipients() []model_user_id.UserID { return x.recipients }

func NewResponse(recipients []model_user_id.UserID) Response { return Response{recipients: recipients} }
