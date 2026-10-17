//go:build e2e

package rpc

import (
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	reward_order_v1 "github.com/Roum1212/hahaton/genproto/reward-order/v1"
	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
)

func (x *RPCTestSuite) TestRewardDeletionAndCooperationRefund() {
	ctx := x.T().Context()
	c, a := x.coordinator(), x.login()
	x.cooperate(c, a)
	x.fund(c, a)
	id := x.reward(c, 20, 1)
	bought, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, &reward_order_v1.CreateRewardOrderRequest{
		RewardId:           id,
		ExpectedPriceCoins: 20,
		RequestId:          "deleted",
	}))
	x.Require().NoError(err)

	for range 2 {
		_, err = x.clients.reward.DeleteReward(ctx, request(c, &reward_v1.DeleteRewardRequest{Id: id}))
		x.Require().NoError(err)
	}

	_, err = x.clients.reward.GetReward(ctx, request(a, &reward_v1.GetRewardRequest{Id: id}))
	callError(x.T(), err, connect.CodeNotFound)
	canceled, err := x.clients.rewardOrder.GetRewardOrder(ctx, request(a,
		&reward_order_v1.GetRewardOrderRequest{Id: bought.Msg.GetRewardOrder().GetId()}))
	x.Require().NoError(err)
	x.Require().Equal(reward_order_v1.RewardOrderCancellationReason_REWARD_ORDER_CANCELLATION_REASON_REWARD_DELETED,
		canceled.Msg.GetRewardOrder().GetCancellationReason())
	next := x.reward(c, 20, 1)
	pending, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, &reward_order_v1.CreateRewardOrderRequest{
		RewardId:           next,
		ExpectedPriceCoins: 20,
		RequestId:          "termination",
	}))
	x.Require().NoError(err)
	_, err = x.clients.cooperation.TerminateCooperation(ctx, request(a, &cooperation_v1.TerminateCooperationRequest{
		Counterpart: &cooperation_v1.TerminateCooperationRequest_CoordinatorId{CoordinatorId: c.id},
	}))
	x.Require().NoError(err)
	terminated, err := x.clients.rewardOrder.GetRewardOrder(ctx, request(a,
		&reward_order_v1.GetRewardOrderRequest{Id: pending.Msg.GetRewardOrder().GetId()}))
	x.Require().NoError(err)
	x.Require().Equal(
		reward_order_v1.RewardOrderCancellationReason_REWARD_ORDER_CANCELLATION_REASON_COOPERATION_TERMINATED,
		terminated.Msg.GetRewardOrder().GetCancellationReason(),
	)
	balance, err := x.clients.volunteer.GetVolunteer(ctx, request(a, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().EqualValues(30, balance.Msg.GetVolunteer().GetCoinBalance())
	x.Require().EqualValues(180, balance.Msg.GetVolunteer().GetConfirmedMinutes())
	x.cooperate(c, a)
	available, err := x.clients.reward.GetReward(ctx, request(a, &reward_v1.GetRewardRequest{Id: next}))
	x.Require().NoError(err)
	x.Require().EqualValues(1, available.Msg.GetReward().GetAvailableQuantity())
	listed, err := x.clients.rewardOrder.ListRewardOrders(ctx, request(a, &reward_order_v1.ListRewardOrdersRequest{
		Scope: reward_order_v1.RewardOrderListScope_REWARD_ORDER_LIST_SCOPE_MY_PURCHASES,
	}))
	x.Require().NoError(err)
	x.Require().Len(listed.Msg.GetRewardOrders(), 2)
}

func (x *RPCTestSuite) TestRewardPurchaseRefundAndIssue() {
	ctx := x.T().Context()
	c, partner, a, outsider := x.coordinator(), x.coordinator(), x.login(), x.login()
	x.cooperate(c, a)
	x.cooperate(partner, a)
	x.fund(c, a)
	id := x.reward(partner, 20, 2)
	x.reward(partner, 5, 1)
	x.reward(partner, 30, 1)
	_, err := x.clients.reward.GetReward(ctx, request(outsider, &reward_v1.GetRewardRequest{Id: id}))
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.reward.UpdateReward(
		ctx,
		request(c, &reward_v1.UpdateRewardRequest{
			Id:    id,
			Title: proto.String("Чужое"),
		}),
	)
	callError(x.T(), err, connect.CodePermissionDenied)
	catalog, err := x.clients.reward.ListRewards(ctx, request(a, &reward_v1.ListRewardsRequest{
		Scope:         reward_v1.RewardListScope_REWARD_LIST_SCOPE_CATALOG,
		CoordinatorId: proto.String(partner.id),
	}))
	x.Require().NoError(err)
	x.Require().Len(catalog.Msg.GetRewards(), 3)
	_, err = x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, &reward_order_v1.CreateRewardOrderRequest{
		RewardId:           id,
		ExpectedPriceCoins: 19,
		RequestId:          "purchase",
	}))
	callError(x.T(), err, connect.CodeAborted)

	purchase := &reward_order_v1.CreateRewardOrderRequest{
		RewardId:           id,
		ExpectedPriceCoins: 20,
		RequestId:          "purchase",
	}
	bought, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, purchase))
	x.Require().NoError(err)

	orderID := bought.Msg.GetRewardOrder().GetId()
	profile, err := x.clients.volunteer.GetVolunteer(ctx, request(a, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().EqualValues(10, profile.Msg.GetVolunteer().GetCoinBalance())
	repeated, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, purchase))
	x.Require().NoError(err)
	x.Require().Equal(orderID, repeated.Msg.GetRewardOrder().GetId())
	_, err = x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, &reward_order_v1.CreateRewardOrderRequest{
		RewardId:           id,
		ExpectedPriceCoins: 20,
		RequestId:          "overdraft",
	}))
	callError(x.T(), err, connect.CodeFailedPrecondition)
	_, err = x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, &reward_order_v1.CreateRewardOrderRequest{
		RewardId:           id,
		ExpectedPriceCoins: 21,
		RequestId:          "purchase",
	}))
	callError(x.T(), err, connect.CodeAlreadyExists)
	_, err = x.clients.rewardOrder.GetRewardOrder(
		ctx,
		request(outsider, &reward_order_v1.GetRewardOrderRequest{Id: orderID}),
	)
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.reward.UpdateReward(ctx, request(partner, &reward_v1.UpdateRewardRequest{
		Id:         id,
		PriceCoins: proto.Int64(25),
		Title:      proto.String("Новое название"),
	}))
	x.Require().NoError(err)
	snapshot, err := x.clients.rewardOrder.GetRewardOrder(
		ctx,
		request(a, &reward_order_v1.GetRewardOrderRequest{Id: orderID}),
	)
	x.Require().NoError(err)
	x.Require().Equal("Сувенир", snapshot.Msg.GetRewardOrder().GetRewardTitle())
	x.Require().EqualValues(20, snapshot.Msg.GetRewardOrder().GetPriceCoins())

	for range 2 {
		_, err = x.clients.rewardOrder.CancelRewardOrder(
			ctx,
			request(a, &reward_order_v1.CancelRewardOrderRequest{Id: orderID}),
		)
		x.Require().NoError(err)
	}

	refunded, err := x.clients.volunteer.GetVolunteer(ctx, request(a, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().EqualValues(30, refunded.Msg.GetVolunteer().GetCoinBalance())
	issued, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, &reward_order_v1.CreateRewardOrderRequest{
		RewardId:           id,
		ExpectedPriceCoins: 25,
		RequestId:          "issued",
	}))
	x.Require().NoError(err)

	issuedID := issued.Msg.GetRewardOrder().GetId()
	_, err = x.clients.rewardOrder.IssueRewardOrder(
		ctx,
		request(a, &reward_order_v1.IssueRewardOrderRequest{Id: issuedID}),
	)
	callError(x.T(), err, connect.CodePermissionDenied)

	for range 2 {
		_, err = x.clients.rewardOrder.IssueRewardOrder(
			ctx,
			request(partner, &reward_order_v1.IssueRewardOrderRequest{Id: issuedID}),
		)
		x.Require().NoError(err)
	}

	_, err = x.clients.rewardOrder.CancelRewardOrder(
		ctx,
		request(a, &reward_order_v1.CancelRewardOrderRequest{Id: issuedID}),
	)
	callError(x.T(), err, connect.CodeFailedPrecondition)
	incoming, err := x.clients.rewardOrder.ListRewardOrders(ctx, request(
		partner,
		&reward_order_v1.ListRewardOrdersRequest{
			Scope:       reward_order_v1.RewardOrderListScope_REWARD_ORDER_LIST_SCOPE_INCOMING,
			VolunteerId: proto.String(a.id),
		},
	))
	x.Require().NoError(err)
	x.Require().Len(incoming.Msg.GetRewardOrders(), 2)
	history, err := x.clients.volunteer.ListVolunteerOperations(ctx,
		request(a, &volunteer_v1.ListVolunteerOperationsRequest{}))
	x.Require().NoError(err)
	x.Require().Len(history.Msg.GetOperations(), 4)

	deltas := make([]int64, 0, 4)
	for _, operation := range history.Msg.GetOperations() {
		deltas = append(deltas, operation.GetCoinDelta())
	}

	x.Require().Equal([]int64{-25, 20, -20, 30}, deltas)
}

func (x *RPCTestSuite) fund(coordinator, volunteer actor) {
	id := x.event(coordinator, 1)
	participationID := x.participate(coordinator, volunteer, id)
	x.finishTime(id)
	_, err := x.clients.event.CompleteEvent(x.T().Context(), request(coordinator, &event_v1.CompleteEventRequest{
		Id:        id,
		RequestId: id,
		Results: []*event_v1.EventParticipationResult{{
			ParticipationId: participationID,
			Outcome:         &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 180},
		}},
	}))
	x.Require().NoError(err)
}

func (x *RPCTestSuite) reward(coordinator actor, price int64, quantity int32) string {
	response, err := x.clients.reward.CreateReward(x.T().Context(), request(coordinator, &reward_v1.CreateRewardRequest{
		Title:         "Сувенир",
		PriceCoins:    price,
		QuantityLimit: quantity,
	}))
	x.Require().NoError(err)

	return response.Msg.GetId()
}
