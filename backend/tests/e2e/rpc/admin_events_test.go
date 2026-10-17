//go:build e2e

package rpc

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	admin_v1 "github.com/Roum1212/hahaton/genproto/admin/v1"
	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	reward_order_v1 "github.com/Roum1212/hahaton/genproto/reward-order/v1"
	updates_v1 "github.com/Roum1212/hahaton/genproto/updates/v1"
	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
)

const missingForceCompletionID = "00000000-0000-7000-8000-000000000001"

func (x *RPCTestSuite) TestForceCompleteEventAppliesAllResults() {
	ctx := x.T().Context()
	c, volunteer, other, pending, stranger := x.coordinator(), x.login(), x.login(), x.login(), x.login()

	for _, value := range []actor{volunteer, other, pending} {
		x.cooperate(c, value)
	}

	id := x.event(c, 3)
	first := x.participate(c, volunteer, id)
	second := x.participate(c, other, id)
	waiting, err := x.clients.participation.RequestParticipation(ctx,
		request(pending, &participation_v1.RequestParticipationRequest{EventId: id}))
	x.Require().NoError(err)

	before, err := x.clients.event.GetEvent(ctx, request(c, &event_v1.GetEventRequest{Id: id}))
	x.Require().NoError(err)
	x.Require().True(before.Msg.GetEvent().GetStartsAt().AsTime().After(time.Now()))

	message := &admin_v1.ForceCompleteEventRequest{Id: id}
	_, err = x.clients.event.CompleteEvent(ctx, request(c, &event_v1.CompleteEventRequest{
		Id: id, RequestId: "normal-finish",
		Results: []*event_v1.EventParticipationResult{
			{ParticipationId: first, Outcome: &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 180}},
			{ParticipationId: second, Outcome: &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 180}},
		},
	}))
	callError(x.T(), err, connect.CodeFailedPrecondition)

	_, err = x.clients.eventChat.SendEventChatMessage(ctx, request(c, &event_chat_v1.SendEventChatMessageRequest{
		EventId: id, Text: "История сохраняется", RequestId: "before-completion",
	}))
	x.Require().NoError(err)

	streamCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	chat, err := x.clients.eventChat.SubscribeEventChatMessages(streamCtx,
		request(volunteer, &event_chat_v1.SubscribeEventChatMessagesRequest{EventId: id}))
	x.Require().NoError(err)
	x.T().Cleanup(func() { x.Require().NoError(chat.Close()) })
	x.Require().True(chat.Receive(), "%v", chat.Err())

	ownerUpdates, err := x.clients.updates.SubscribeUpdates(streamCtx, request(c, &updates_v1.SubscribeUpdatesRequest{}))
	x.Require().NoError(err)
	x.T().Cleanup(func() { x.Require().NoError(ownerUpdates.Close()) })
	x.update(ownerUpdates, updates_v1.UpdateKind_UPDATE_KIND_READY)
	volunteerUpdates, err := x.clients.updates.SubscribeUpdates(streamCtx,
		request(volunteer, &updates_v1.SubscribeUpdatesRequest{}))
	x.Require().NoError(err)
	x.T().Cleanup(func() { x.Require().NoError(volunteerUpdates.Close()) })
	x.update(volunteerUpdates, updates_v1.UpdateKind_UPDATE_KIND_READY)

	strangerCtx, stopStranger := context.WithTimeout(ctx, noUpdateTimeout)
	defer stopStranger()

	strangerUpdates, err := x.clients.updates.SubscribeUpdates(strangerCtx,
		request(stranger, &updates_v1.SubscribeUpdatesRequest{}))
	x.Require().NoError(err)
	x.T().Cleanup(func() { x.Require().NoError(strangerUpdates.Close()) })
	x.update(strangerUpdates, updates_v1.UpdateKind_UPDATE_KIND_READY)

	r := connect.NewRequest(message)
	r.Header().Set("Authorization", "Bearer "+adminAPIKey)
	completed, err := x.clients.admin.ForceCompleteEvent(ctx, r)
	x.Require().NoError(err)
	x.Require().Equal(event_v1.EventStatus_EVENT_STATUS_COMPLETED, completed.Msg.GetEvent().GetStatus())
	x.Require().True(proto.Equal(before.Msg.GetEvent().GetStartsAt(), completed.Msg.GetEvent().GetStartsAt()))
	x.Require().True(proto.Equal(before.Msg.GetEvent().GetEndsAt(), completed.Msg.GetEvent().GetEndsAt()))
	x.changed(ownerUpdates, volunteerUpdates)
	x.Require().False(strangerUpdates.Receive())
	callError(x.T(), strangerUpdates.Err(), connect.CodeDeadlineExceeded)
	x.Require().False(chat.Receive())
	x.Require().NoError(chat.Err())
	_, err = x.clients.admin.ForceCompleteEvent(ctx, r)
	x.Require().NoError(err)
	profile, err := x.clients.volunteer.GetVolunteer(ctx, request(volunteer, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().EqualValues(30, profile.Msg.GetVolunteer().GetCoinBalance())
	x.Require().EqualValues(180, profile.Msg.GetVolunteer().GetConfirmedMinutes())
	operations, err := x.clients.volunteer.ListVolunteerOperations(ctx,
		request(volunteer, &volunteer_v1.ListVolunteerOperationsRequest{}))
	x.Require().NoError(err)
	x.Require().Len(operations.Msg.GetOperations(), 1)
	x.Require().EqualValues(30, operations.Msg.GetOperations()[0].GetCoinDelta())

	for _, value := range []struct {
		actor  actor
		id     string
		status participation_v1.ParticipationStatus
	}{
		{volunteer, first, participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CONFIRMED},
		{other, second, participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CONFIRMED},
		{pending, waiting.Msg.GetId(), participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CONFIRMED},
	} {
		participation, err := x.clients.participation.GetParticipation(ctx,
			request(value.actor, &participation_v1.GetParticipationRequest{Id: value.id}))
		x.Require().NoError(err)
		x.Require().Equal(value.status, participation.Msg.GetParticipation().GetStatus())
		changes, err := x.clients.participation.ListParticipationChanges(ctx,
			request(value.actor, &participation_v1.ListParticipationChangesRequest{ParticipationId: value.id}))
		x.Require().NoError(err)

		last := changes.Msg.GetChanges()[len(changes.Msg.GetChanges())-1]
		x.Require().Equal(value.status, last.GetStatus())

		x.Require().Equal(c.id, last.GetCoordinatorId())
		x.Require().EqualValues(180, participation.Msg.GetParticipation().GetAwardedMinutes())
		x.Require().EqualValues(30, participation.Msg.GetParticipation().GetAwardedCoins())

		profile, err := x.clients.volunteer.GetVolunteer(ctx,
			request(value.actor, &volunteer_v1.GetVolunteerRequest{}))
		x.Require().NoError(err)
		x.Require().EqualValues(30, profile.Msg.GetVolunteer().GetCoinBalance())
		x.Require().EqualValues(180, profile.Msg.GetVolunteer().GetConfirmedMinutes())
		operations, err := x.clients.volunteer.ListVolunteerOperations(ctx,
			request(value.actor, &volunteer_v1.ListVolunteerOperationsRequest{}))
		x.Require().NoError(err)
		x.Require().Len(operations.Msg.GetOperations(), 1)
	}

	readOnly, err := x.clients.eventChat.GetEventChat(ctx,
		request(volunteer, &event_chat_v1.GetEventChatRequest{EventId: id}))
	x.Require().NoError(err)
	x.Require().True(readOnly.Msg.GetChat().GetReadOnly())
	_, err = x.clients.eventChat.SendEventChatMessage(ctx, request(c, &event_chat_v1.SendEventChatMessageRequest{
		EventId: id, Text: "После завершения", RequestId: "after-completion",
	}))
	callError(x.T(), err, connect.CodeFailedPrecondition)
	_, err = x.clients.event.UpdateEvent(ctx, request(c, &event_v1.UpdateEventRequest{Id: id}))
	callError(x.T(), err, connect.CodeFailedPrecondition)
	_, err = x.clients.participation.RequestParticipation(ctx,
		request(pending, &participation_v1.RequestParticipationRequest{EventId: id}))
	callError(x.T(), err, connect.CodeFailedPrecondition)

	rewardID := x.reward(c, 20, 1)
	_, err = x.clients.rewardOrder.CreateRewardOrder(ctx, request(volunteer, &reward_order_v1.CreateRewardOrderRequest{
		RewardId: rewardID, ExpectedPriceCoins: 20, RequestId: "spend-forced-accrual",
	}))
	x.Require().NoError(err)
	profile, err = x.clients.volunteer.GetVolunteer(ctx, request(volunteer, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().EqualValues(10, profile.Msg.GetVolunteer().GetCoinBalance())
}

func (x *RPCTestSuite) TestForceCompleteEventRequiresAdminToken() {
	ctx := x.T().Context()
	c := x.coordinator()
	id := x.event(c, 1)
	r := request(c, &admin_v1.ForceCompleteEventRequest{Id: id})
	_, err := x.clients.admin.ForceCompleteEvent(ctx, r)
	callError(x.T(), err, connect.CodeUnauthenticated)
	r.Header().Set("Authorization", "Bearer incorrect")
	_, err = x.clients.admin.ForceCompleteEvent(ctx, r)
	callError(x.T(), err, connect.CodeUnauthenticated)
	r.Header().Del("Cookie")
	r.Header().Set("Authorization", "Bearer "+adminAPIKey)
	_, err = x.clients.admin.ForceCompleteEvent(ctx, r)
	x.Require().NoError(err)
}

func (x *RPCTestSuite) TestForceCompleteEventValidatesState() {
	ctx := x.T().Context()
	c, volunteer := x.coordinator(), x.login()
	x.cooperate(c, volunteer)
	id := x.event(c, 1)
	x.participate(c, volunteer, id)

	for _, id := range []string{"", "invalid"} {
		r := connect.NewRequest(&admin_v1.ForceCompleteEventRequest{Id: id})
		r.Header().Set("Authorization", "Bearer "+adminAPIKey)
		_, err := x.clients.admin.ForceCompleteEvent(ctx, r)
		callError(x.T(), err, connect.CodeInvalidArgument)
	}

	_, err := x.clients.event.CancelEvent(ctx, request(c, &event_v1.CancelEventRequest{Id: id}))
	x.Require().NoError(err)

	r := connect.NewRequest(&admin_v1.ForceCompleteEventRequest{Id: id})
	r.Header().Set("Authorization", "Bearer "+adminAPIKey)
	_, err = x.clients.admin.ForceCompleteEvent(ctx, r)
	callError(x.T(), err, connect.CodeFailedPrecondition)

	r.Msg.Id = missingForceCompletionID
	_, err = x.clients.admin.ForceCompleteEvent(ctx, r)
	callError(x.T(), err, connect.CodeNotFound)

	id = x.event(c, 1)
	participationID := x.participate(c, volunteer, id)
	x.finishTime(id)
	_, err = x.clients.event.CompleteEvent(ctx, request(c, &event_v1.CompleteEventRequest{
		Id: id, RequestId: "regular-completion",
		Results: []*event_v1.EventParticipationResult{
			{
				ParticipationId: participationID,
				Outcome:         &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 60},
			},
		},
	}))
	x.Require().NoError(err)

	r.Msg.Id = id
	_, err = x.clients.admin.ForceCompleteEvent(ctx, r)
	callError(x.T(), err, connect.CodeFailedPrecondition)
}

func (x *RPCTestSuite) TestForceCompleteEventSkipsInactiveAndOtherEvents() {
	ctx := x.T().Context()
	c, canceled, rejected, pending := x.coordinator(), x.login(), x.login(), x.login()

	for _, value := range []actor{canceled, rejected, pending} {
		x.cooperate(c, value)
	}

	id := x.event(c, 2)
	canceledID := x.participate(c, canceled, id)
	rejectedID := x.participate(c, rejected, id)
	_, err := x.clients.participation.CancelParticipation(ctx,
		request(canceled, &participation_v1.CancelParticipationRequest{Id: canceledID}))
	x.Require().NoError(err)
	_, err = x.clients.participation.RejectParticipation(ctx,
		request(c, &participation_v1.RejectParticipationRequest{Id: rejectedID}))
	x.Require().NoError(err)

	otherEventID := x.event(c, 1)
	waiting, err := x.clients.participation.RequestParticipation(ctx,
		request(pending, &participation_v1.RequestParticipationRequest{EventId: otherEventID}))
	x.Require().NoError(err)

	r := connect.NewRequest(&admin_v1.ForceCompleteEventRequest{Id: id})
	r.Header().Set("Authorization", "Bearer "+adminAPIKey)
	_, err = x.clients.admin.ForceCompleteEvent(ctx, r)
	x.Require().NoError(err)

	for _, value := range []struct {
		actor  actor
		id     string
		status participation_v1.ParticipationStatus
	}{
		{canceled, canceledID, participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CANCELLED},
		{rejected, rejectedID, participation_v1.ParticipationStatus_PARTICIPATION_STATUS_REJECTED},
		{pending, waiting.Msg.GetId(), participation_v1.ParticipationStatus_PARTICIPATION_STATUS_PENDING},
	} {
		participation, err := x.clients.participation.GetParticipation(ctx,
			request(value.actor, &participation_v1.GetParticipationRequest{Id: value.id}))
		x.Require().NoError(err)
		x.Require().Equal(value.status, participation.Msg.GetParticipation().GetStatus())
		profile, err := x.clients.volunteer.GetVolunteer(ctx,
			request(value.actor, &volunteer_v1.GetVolunteerRequest{}))
		x.Require().NoError(err)
		x.Require().Zero(profile.Msg.GetVolunteer().GetCoinBalance())
		x.Require().Zero(profile.Msg.GetVolunteer().GetConfirmedMinutes())
	}
}

func (x *RPCTestSuite) TestForceCompleteEventConcurrentRetries() {
	ctx := x.T().Context()
	c, volunteer := x.coordinator(), x.login()
	x.cooperate(c, volunteer)
	id := x.event(c, 1)
	_, err := x.clients.participation.RequestParticipation(ctx,
		request(volunteer, &participation_v1.RequestParticipationRequest{EventId: id}))
	x.Require().NoError(err)

	complete := func() error {
		r := connect.NewRequest(&admin_v1.ForceCompleteEventRequest{Id: id})
		r.Header().Set("Authorization", "Bearer "+adminAPIKey)
		_, err := x.clients.admin.ForceCompleteEvent(ctx, r)

		return err
	}
	for _, err := range simultaneously(complete, complete) {
		x.Require().NoError(err)
	}

	profile, err := x.clients.volunteer.GetVolunteer(ctx, request(volunteer, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().EqualValues(30, profile.Msg.GetVolunteer().GetCoinBalance())
	x.Require().EqualValues(180, profile.Msg.GetVolunteer().GetConfirmedMinutes())
	operations, err := x.clients.volunteer.ListVolunteerOperations(ctx,
		request(volunteer, &volunteer_v1.ListVolunteerOperationsRequest{}))
	x.Require().NoError(err)
	x.Require().Len(operations.Msg.GetOperations(), 1)
}
