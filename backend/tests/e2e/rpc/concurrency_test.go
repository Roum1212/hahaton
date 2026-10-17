//go:build e2e

package rpc

import (
	"sync"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	reward_order_v1 "github.com/Roum1212/hahaton/genproto/reward-order/v1"
	reward_v1 "github.com/Roum1212/hahaton/genproto/reward/v1"
	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
)

func (x *RPCTestSuite) TestConcurrentApprovalAndVolunteerCancellation() {
	ctx := x.T().Context()
	c, a := x.coordinator(), x.login()
	x.cooperate(c, a)

	for range 5 {
		id := x.event(c, 1)
		application, err := x.clients.participation.RequestParticipation(ctx,
			request(a, &participation_v1.RequestParticipationRequest{EventId: id}))
		x.Require().NoError(err)

		participationID := application.Msg.GetId()
		results := simultaneously(
			func() error {
				_, err := x.clients.participation.ApproveParticipation(ctx,
					request(c, &participation_v1.ApproveParticipationRequest{Id: participationID}))

				return err
			},
			func() error {
				_, err := x.clients.participation.CancelParticipation(ctx,
					request(a, &participation_v1.CancelParticipationRequest{Id: participationID}))

				return err
			},
		)
		x.Require().NoError(results[1])

		if results[0] != nil {
			callError(x.T(), results[0], connect.CodeFailedPrecondition)
		}

		canceled, err := x.clients.participation.GetParticipation(ctx,
			request(a, &participation_v1.GetParticipationRequest{Id: participationID}))
		x.Require().NoError(err)
		x.Require().Equal(participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CANCELLED,
			canceled.Msg.GetParticipation().GetStatus())
		_, err = x.clients.eventChat.GetEventChat(ctx, request(a, &event_chat_v1.GetEventChatRequest{EventId: id}))
		callError(x.T(), err, connect.CodePermissionDenied)
	}
}

func (x *RPCTestSuite) TestConcurrentCooperationApprovalAndRejection() {
	ctx := x.T().Context()

	c := x.coordinator()
	for range 5 {
		a := x.login()
		_, err := x.clients.cooperation.RequestCooperation(ctx,
			request(a, &cooperation_v1.RequestCooperationRequest{CoordinatorId: c.id}))
		x.Require().NoError(err)

		results := simultaneously(
			func() error {
				_, err := x.clients.cooperation.ApproveCooperation(ctx,
					request(c, &cooperation_v1.ApproveCooperationRequest{VolunteerId: a.id}))

				return err
			},
			func() error {
				_, err := x.clients.cooperation.RejectCooperation(ctx,
					request(c, &cooperation_v1.RejectCooperationRequest{VolunteerId: a.id}))

				return err
			},
		)
		x.assertOneWinner(results)
		cooperation, err := x.clients.cooperation.GetCooperation(ctx, request(a, &cooperation_v1.GetCooperationRequest{
			Counterpart: &cooperation_v1.GetCooperationRequest_CoordinatorId{CoordinatorId: c.id},
		}))
		x.Require().NoError(err)

		if results[0] == nil {
			x.Require().Equal(
				cooperation_v1.CooperationStatus_COOPERATION_STATUS_ACTIVE,
				cooperation.Msg.GetCooperation().GetStatus(),
			)
			x.Require().Nil(cooperation.Msg.GetCooperation().GetRetryAfter())
		} else {
			x.Require().Equal(
				cooperation_v1.CooperationStatus_COOPERATION_STATUS_REJECTED,
				cooperation.Msg.GetCooperation().GetStatus(),
			)
			x.Require().NotNil(cooperation.Msg.GetCooperation().GetRetryAfter())
		}
	}
}

func (x *RPCTestSuite) TestConcurrentEventCancellationAndApplication() {
	ctx := x.T().Context()
	c, a := x.coordinator(), x.login()
	x.cooperate(c, a)

	for range 5 {
		id := x.event(c, 1)
		results := simultaneously(
			func() error {
				_, err := x.clients.event.CancelEvent(ctx, request(c, &event_v1.CancelEventRequest{Id: id}))

				return err
			},
			func() error {
				_, err := x.clients.participation.RequestParticipation(ctx,
					request(a, &participation_v1.RequestParticipationRequest{EventId: id}))

				return err
			},
		)
		x.Require().NoError(results[0])

		if results[1] != nil {
			callError(x.T(), results[1], connect.CodeFailedPrecondition)
		}

		listed, err := x.clients.participation.ListParticipations(ctx, request(a, &participation_v1.ListParticipationsRequest{
			Scope:   participation_v1.ParticipationListScope_PARTICIPATION_LIST_SCOPE_MY_APPLICATIONS,
			EventId: proto.String(id),
		}))
		x.Require().NoError(err)

		for _, application := range listed.Msg.GetParticipations() {
			x.Require().Equal(
				participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CANCELLED,
				application.GetStatus(),
			)
		}
	}
}

func (x *RPCTestSuite) TestConcurrentIssueAndRefund() {
	ctx := x.T().Context()
	c, a := x.coordinator(), x.login()
	x.cooperate(c, a)
	x.fund(c, a)
	id := x.reward(c, 20, 1)
	created, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, &reward_order_v1.CreateRewardOrderRequest{
		RewardId:           id,
		ExpectedPriceCoins: 20,
		RequestId:          "issue-or-refund",
	}))
	x.Require().NoError(err)

	orderID := created.Msg.GetRewardOrder().GetId()
	results := simultaneously(
		func() error {
			_, err := x.clients.rewardOrder.IssueRewardOrder(ctx,
				request(c, &reward_order_v1.IssueRewardOrderRequest{Id: orderID}))

			return err
		},
		func() error {
			_, err := x.clients.rewardOrder.CancelRewardOrder(ctx,
				request(a, &reward_order_v1.CancelRewardOrderRequest{Id: orderID}))

			return err
		},
	)
	x.assertOneWinner(results)
	order, err := x.clients.rewardOrder.GetRewardOrder(
		ctx,
		request(a, &reward_order_v1.GetRewardOrderRequest{Id: orderID}),
	)
	x.Require().NoError(err)
	balance, err := x.clients.volunteer.GetVolunteer(ctx, request(a, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)

	if results[0] == nil {
		x.Require().Equal(
			reward_order_v1.RewardOrderStatus_REWARD_ORDER_STATUS_ISSUED,
			order.Msg.GetRewardOrder().GetStatus(),
		)
		x.Require().EqualValues(10, balance.Msg.GetVolunteer().GetCoinBalance())
	} else {
		x.Require().Equal(
			reward_order_v1.RewardOrderStatus_REWARD_ORDER_STATUS_CANCELLED,
			order.Msg.GetRewardOrder().GetStatus(),
		)
		x.Require().EqualValues(30, balance.Msg.GetVolunteer().GetCoinBalance())
	}
}

func (x *RPCTestSuite) TestConcurrentLastEventPlace() {
	ctx := x.T().Context()
	c, a, b := x.coordinator(), x.login(), x.login()
	x.cooperate(c, a)
	x.cooperate(c, b)
	id := x.event(c, 1)
	first, err := x.clients.participation.RequestParticipation(ctx,
		request(a, &participation_v1.RequestParticipationRequest{EventId: id}))
	x.Require().NoError(err)
	second, err := x.clients.participation.RequestParticipation(ctx,
		request(b, &participation_v1.RequestParticipationRequest{EventId: id}))
	x.Require().NoError(err)

	results := simultaneously(
		func() error {
			_, err := x.clients.participation.ApproveParticipation(ctx,
				request(c, &participation_v1.ApproveParticipationRequest{Id: first.Msg.GetId()}))

			return err
		},
		func() error {
			_, err := x.clients.participation.ApproveParticipation(ctx,
				request(c, &participation_v1.ApproveParticipationRequest{Id: second.Msg.GetId()}))

			return err
		},
	)
	x.assertOneWinner(results)
	listed, err := x.clients.participation.ListParticipations(ctx, request(c, &participation_v1.ListParticipationsRequest{
		Scope:   participation_v1.ParticipationListScope_PARTICIPATION_LIST_SCOPE_MY_EVENTS,
		EventId: proto.String(id),
		Status:  participation_v1.ParticipationStatus_PARTICIPATION_STATUS_APPROVED.Enum(),
	}))
	x.Require().NoError(err)
	x.Require().Len(listed.Msg.GetParticipations(), 1)
}

func (x *RPCTestSuite) TestConcurrentLastRewardAndWalletBalance() {
	ctx := x.T().Context()
	c, a, b := x.coordinator(), x.login(), x.login()
	x.cooperate(c, a)
	x.cooperate(c, b)
	x.fund(c, a)
	x.fund(c, b)
	id := x.reward(c, 20, 1)
	results := simultaneously(
		func() error {
			_, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, &reward_order_v1.CreateRewardOrderRequest{
				RewardId:           id,
				ExpectedPriceCoins: 20,
				RequestId:          "last",
			}))

			return err
		},
		func() error {
			_, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(b, &reward_order_v1.CreateRewardOrderRequest{
				RewardId:           id,
				ExpectedPriceCoins: 20,
				RequestId:          "last",
			}))

			return err
		},
	)
	x.assertOneWinner(results)
	reward, err := x.clients.reward.GetReward(ctx, request(c, &reward_v1.GetRewardRequest{Id: id}))
	x.Require().NoError(err)
	x.Require().Zero(reward.Msg.GetReward().GetAvailableQuantity())

	total := int64(0)

	for _, volunteer := range []actor{
		a,
		b,
	} {
		profile, err := x.clients.volunteer.GetVolunteer(ctx, request(volunteer, &volunteer_v1.GetVolunteerRequest{}))
		x.Require().NoError(err)

		total += profile.Msg.GetVolunteer().GetCoinBalance()
	}

	x.Require().EqualValues(40, total)
	firstReward, secondReward := x.reward(c, 20, 1), x.reward(c, 20, 1)
	// The loser of the previous purchase still has 30 coins.
	buyer := a
	if results[0] == nil {
		buyer = b
	}

	results = simultaneously(
		func() error {
			_, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(buyer, &reward_order_v1.CreateRewardOrderRequest{
				RewardId:           firstReward,
				ExpectedPriceCoins: 20,
				RequestId:          "wallet-first",
			}))

			return err
		},
		func() error {
			_, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(buyer, &reward_order_v1.CreateRewardOrderRequest{
				RewardId:           secondReward,
				ExpectedPriceCoins: 20,
				RequestId:          "wallet-second",
			}))

			return err
		},
	)
	x.assertOneWinner(results)
	balance, err := x.clients.volunteer.GetVolunteer(ctx, request(buyer, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().EqualValues(10, balance.Msg.GetVolunteer().GetCoinBalance())
}

func (x *RPCTestSuite) TestConcurrentMultiVolunteerCompletion() {
	ctx := x.T().Context()
	for range 5 {
		c, other, a, b := x.coordinator(), x.coordinator(), x.login(), x.login()
		for _, coordinator := range []actor{c, other} {
			for _, volunteer := range []actor{a, b} {
				x.cooperate(coordinator, volunteer)
			}
		}

		first, second := x.event(c, 2), x.event(other, 2)
		firstA, firstB := x.participate(c, a, first), x.participate(c, b, first)
		secondB, secondA := x.participate(other, b, second), x.participate(other, a, second)
		x.finishTime(first)
		x.finishTime(second)

		firstResult := &event_v1.CompleteEventRequest{
			Id: first, RequestId: first,
			Results: []*event_v1.EventParticipationResult{
				{ParticipationId: firstA, Outcome: &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 180}},
				{ParticipationId: firstB, Outcome: &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 180}},
			},
		}
		secondResult := &event_v1.CompleteEventRequest{
			Id: second, RequestId: second,
			Results: []*event_v1.EventParticipationResult{
				{ParticipationId: secondB, Outcome: &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 180}},
				{ParticipationId: secondA, Outcome: &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 180}},
			},
		}

		results := simultaneously(
			func() error {
				_, err := x.clients.event.CompleteEvent(ctx, request(c, firstResult))

				return err
			},
			func() error {
				_, err := x.clients.event.CompleteEvent(ctx, request(other, secondResult))

				return err
			},
		)
		for _, err := range results {
			x.Require().NoError(err)
		}

		for _, volunteer := range []actor{a, b} {
			profile, err := x.clients.volunteer.GetVolunteer(ctx, request(volunteer, &volunteer_v1.GetVolunteerRequest{}))
			x.Require().NoError(err)
			x.Require().EqualValues(60, profile.Msg.GetVolunteer().GetCoinBalance())
			x.Require().EqualValues(360, profile.Msg.GetVolunteer().GetConfirmedMinutes())
		}
	}
}

func (x *RPCTestSuite) TestConcurrentRepeatedCompletionAndPurchase() {
	ctx := x.T().Context()
	c, a := x.coordinator(), x.login()
	x.cooperate(c, a)
	id := x.event(c, 1)
	participationID := x.participate(c, a, id)
	x.finishTime(id)
	complete := &event_v1.CompleteEventRequest{
		Id:        id,
		RequestId: "repeated",
		Results: []*event_v1.EventParticipationResult{{
			ParticipationId: participationID,
			Outcome:         &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 180},
		}},
	}

	finish := func() error {
		_, err := x.clients.event.CompleteEvent(ctx, request(c, complete))

		return err
	}
	for _, err := range simultaneously(finish, finish) {
		x.Require().NoError(err)
	}

	balance, err := x.clients.volunteer.GetVolunteer(ctx, request(a, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().EqualValues(30, balance.Msg.GetVolunteer().GetCoinBalance())
	rewardID := x.reward(c, 20, 2)
	orders := make([]string, 2)

	actions := make([]func() error, 2)
	for i := range actions {
		actions[i] = func() error {
			response, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, &reward_order_v1.CreateRewardOrderRequest{
				RewardId:           rewardID,
				ExpectedPriceCoins: 20,
				RequestId:          "same-purchase",
			}))
			if err == nil {
				orders[i] = response.Msg.GetRewardOrder().GetId()
			}

			return err
		}
	}

	for _, err := range simultaneously(actions...) {
		x.Require().NoError(err)
	}

	x.Require().NotEmpty(orders[0])
	x.Require().Equal(orders[0], orders[1])
	history, err := x.clients.volunteer.ListVolunteerOperations(ctx,
		request(a, &volunteer_v1.ListVolunteerOperationsRequest{}))
	x.Require().NoError(err)
	x.Require().Len(history.Msg.GetOperations(), 2)
	x.Require().EqualValues(10, history.Msg.GetOperations()[0].GetCoinBalanceAfter())
}

func (x *RPCTestSuite) TestConcurrentTerminationAndCompletion() {
	ctx := x.T().Context()
	for range 5 {
		c, a := x.coordinator(), x.login()
		x.cooperate(c, a)
		id := x.event(c, 1)
		participationID := x.participate(c, a, id)
		x.finishTime(id)

		results := simultaneously(
			func() error {
				_, err := x.clients.cooperation.TerminateCooperation(ctx, request(c, &cooperation_v1.TerminateCooperationRequest{
					Counterpart: &cooperation_v1.TerminateCooperationRequest_VolunteerId{VolunteerId: a.id},
				}))

				return err
			},
			func() error {
				_, err := x.clients.event.CompleteEvent(ctx, request(c, &event_v1.CompleteEventRequest{
					Id:        id,
					RequestId: "complete-or-terminate",
					Results: []*event_v1.EventParticipationResult{{
						ParticipationId: participationID,
						Outcome:         &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 180},
					}},
				}))

				return err
			},
		)
		x.Require().NoError(results[0])
		profile, err := x.clients.volunteer.GetVolunteer(ctx, request(a, &volunteer_v1.GetVolunteerRequest{}))
		x.Require().NoError(err)
		participation, err := x.clients.participation.GetParticipation(ctx,
			request(a, &participation_v1.GetParticipationRequest{Id: participationID}))
		x.Require().NoError(err)

		if results[1] == nil {
			x.Require().EqualValues(30, profile.Msg.GetVolunteer().GetCoinBalance())
			x.Require().EqualValues(180, profile.Msg.GetVolunteer().GetConfirmedMinutes())
			x.Require().Equal(participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CONFIRMED,
				participation.Msg.GetParticipation().GetStatus())
		} else {
			callError(x.T(), results[1], connect.CodeInvalidArgument)
			x.Require().Zero(profile.Msg.GetVolunteer().GetCoinBalance())
			x.Require().Zero(profile.Msg.GetVolunteer().GetConfirmedMinutes())
			x.Require().Equal(participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CANCELLED,
				participation.Msg.GetParticipation().GetStatus())
		}
	}
}

func (x *RPCTestSuite) TestConcurrentTerminationAndVolunteerActions() {
	ctx := x.T().Context()
	for range 5 {
		c, a := x.coordinator(), x.login()
		x.cooperate(c, a)
		x.fund(c, a)
		eventID := x.event(c, 1)
		rewardID := x.reward(c, 20, 1)
		results := simultaneously(
			func() error {
				_, err := x.clients.cooperation.TerminateCooperation(ctx, request(c, &cooperation_v1.TerminateCooperationRequest{
					Counterpart: &cooperation_v1.TerminateCooperationRequest_VolunteerId{VolunteerId: a.id},
				}))

				return err
			},
			func() error {
				_, err := x.clients.participation.RequestParticipation(ctx,
					request(a, &participation_v1.RequestParticipationRequest{EventId: eventID}))

				return err
			},
			func() error {
				_, err := x.clients.rewardOrder.CreateRewardOrder(ctx, request(a, &reward_order_v1.CreateRewardOrderRequest{
					RewardId:           rewardID,
					ExpectedPriceCoins: 20,
					RequestId:          "terminate",
				}))

				return err
			},
		)
		x.Require().NoError(results[0])

		for _, err := range results[1:] {
			if err != nil {
				callError(x.T(), err, connect.CodeFailedPrecondition)
			}
		}

		applications, err := x.clients.participation.ListParticipations(ctx, request(
			a,
			&participation_v1.ListParticipationsRequest{
				Scope:   participation_v1.ParticipationListScope_PARTICIPATION_LIST_SCOPE_MY_APPLICATIONS,
				EventId: proto.String(eventID),
			},
		))
		x.Require().NoError(err)

		for _, application := range applications.Msg.GetParticipations() {
			x.Require().Equal(
				participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CANCELLED,
				application.GetStatus(),
			)
		}

		orders, err := x.clients.rewardOrder.ListRewardOrders(ctx, request(a, &reward_order_v1.ListRewardOrdersRequest{
			Scope: reward_order_v1.RewardOrderListScope_REWARD_ORDER_LIST_SCOPE_MY_PURCHASES,
		}))
		x.Require().NoError(err)

		for _, order := range orders.Msg.GetRewardOrders() {
			x.Require().Equal(reward_order_v1.RewardOrderStatus_REWARD_ORDER_STATUS_CANCELLED, order.GetStatus())
		}

		balance, err := x.clients.volunteer.GetVolunteer(ctx, request(a, &volunteer_v1.GetVolunteerRequest{}))
		x.Require().NoError(err)
		x.Require().EqualValues(30, balance.Msg.GetVolunteer().GetCoinBalance())
		_, err = x.clients.eventChat.GetEventChat(
			ctx,
			request(a, &event_chat_v1.GetEventChatRequest{EventId: eventID}),
		)
		callError(x.T(), err, connect.CodePermissionDenied)
	}
}

func (x *RPCTestSuite) assertOneWinner(results []error) {
	successes := 0

	for _, err := range results {
		if err == nil {
			successes++

			continue
		}

		callError(x.T(), err, connect.CodeFailedPrecondition)
	}

	x.Require().Equal(1, successes)
}

func simultaneously(actions ...func() error) []error {
	start := make(chan struct{})
	results := make([]error, len(actions))

	var workers sync.WaitGroup
	for i, action := range actions {
		workers.Go(func() {
			<-start

			results[i] = action()
		})
	}

	close(start)
	workers.Wait()

	return results
}
