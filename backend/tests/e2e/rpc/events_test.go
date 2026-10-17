//go:build e2e

package rpc

import (
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	geo_v1 "github.com/Roum1212/hahaton/genproto/geo/v1"
	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	volunteer_v1 "github.com/Roum1212/hahaton/genproto/volunteer/v1"
)

func (x *RPCTestSuite) TestEventCancellationCancelsApplications() {
	ctx := x.T().Context()
	c, a := x.coordinator(), x.login()
	x.cooperate(c, a)
	id := x.event(c, 1)
	participationID := x.participate(c, a, id)
	_, err := x.clients.event.CancelEvent(ctx, request(c, &event_v1.CancelEventRequest{Id: id}))
	x.Require().NoError(err)
	application, err := x.clients.participation.GetParticipation(ctx,
		request(a, &participation_v1.GetParticipationRequest{Id: participationID}))
	x.Require().NoError(err)
	x.Require().Equal(
		participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CANCELLED,
		application.Msg.GetParticipation().GetStatus(),
	)
	x.Require().Equal(participation_v1.ParticipationCancellationReason_PARTICIPATION_CANCELLATION_REASON_EVENT_CANCELLED,
		application.Msg.GetParticipation().GetCancellationReason())
	_, err = x.clients.participation.RestoreParticipation(ctx,
		request(a, &participation_v1.RestoreParticipationRequest{Id: participationID}))
	callError(x.T(), err, connect.CodeFailedPrecondition)
	past := x.event(c, 1)
	x.finishTime(past)
	_, err = x.clients.event.CancelEvent(ctx, request(c, &event_v1.CancelEventRequest{Id: past}))
	callError(x.T(), err, connect.CodeFailedPrecondition)
}

func (x *RPCTestSuite) TestEventCatalogAndUpdates() {
	ctx := x.T().Context()
	c, other, a := x.coordinator(), x.coordinator(), x.login()
	id := x.event(c, 2)
	x.event(c, 2)
	x.event(c, 2)

	area := &geo_v1.SearchArea{Area: &geo_v1.SearchArea_Radius{Radius: &geo_v1.RadiusArea{
		Center:       point(),
		RadiusMeters: 1000,
	}}}
	catalog, err := x.clients.event.ListEvents(ctx, request(a, &event_v1.ListEventsRequest{
		Scope:      event_v1.EventListScope_EVENT_LIST_SCOPE_CATALOG,
		SearchArea: area,
	}))
	x.Require().NoError(err)
	x.Require().Empty(catalog.Msg.GetEvents())
	_, err = x.clients.event.GetEvent(ctx, request(a, &event_v1.GetEventRequest{Id: id}))
	callError(x.T(), err, connect.CodePermissionDenied)
	x.cooperate(c, a)
	catalog, err = x.clients.event.ListEvents(ctx, request(a, &event_v1.ListEventsRequest{
		Scope:      event_v1.EventListScope_EVENT_LIST_SCOPE_CATALOG,
		SearchArea: area,
		Limit:      1,
	}))
	x.Require().NoError(err)
	x.Require().EqualValues(3, catalog.Msg.GetTotalCount())
	x.Require().Len(catalog.Msg.GetEvents(), 1)
	owned, err := x.clients.event.ListEvents(ctx, request(c, &event_v1.ListEventsRequest{
		Scope:  event_v1.EventListScope_EVENT_LIST_SCOPE_OWNED,
		Offset: 100,
	}))
	x.Require().NoError(err)
	x.Require().EqualValues(3, owned.Msg.GetTotalCount())
	x.Require().Empty(owned.Msg.GetEvents())
	_, err = x.clients.event.UpdateEvent(ctx, request(other, &event_v1.UpdateEventRequest{
		Id:    id,
		Title: proto.String("Чужое"),
	}))
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.event.UpdateEvent(ctx, request(c, &event_v1.UpdateEventRequest{
		Id:          id,
		Description: proto.String("Изменено"),
	}))
	x.Require().NoError(err)
	_, err = x.clients.event.UpdateEvent(ctx, request(c, &event_v1.UpdateEventRequest{Id: id}))
	x.Require().NoError(err)
	updated, err := x.clients.event.GetEvent(ctx, request(a, &event_v1.GetEventRequest{Id: id}))
	x.Require().NoError(err)
	x.Require().Equal("Изменено", updated.Msg.GetEvent().GetDescription())
	x.Require().Equal("Городская акция", updated.Msg.GetEvent().GetTitle())
	_, err = x.clients.event.UpdateEvent(ctx, request(c, &event_v1.UpdateEventRequest{
		Id:                    id,
		Description:           proto.String(""),
		CandidateRequirements: &event_v1.EventCandidateRequirements{},
	}))
	x.Require().NoError(err)
	cleared, err := x.clients.event.GetEvent(ctx, request(c, &event_v1.GetEventRequest{Id: id}))
	x.Require().NoError(err)
	x.Require().Empty(cleared.Msg.GetEvent().GetDescription())
	_, err = x.clients.event.UpdateEvent(ctx, request(c, &event_v1.UpdateEventRequest{
		Id:     id,
		EndsAt: timestamppb.New(updated.Msg.GetEvent().GetStartsAt().AsTime().Add(31 * time.Minute)),
	}))
	callError(x.T(), err, connect.CodeInvalidArgument)
}

func (x *RPCTestSuite) TestEventCompletionIsAtomicAndIdempotent() {
	ctx := x.T().Context()

	c, a, absent, pending := x.coordinator(), x.login(), x.login(), x.login()
	for _, volunteer := range []actor{
		a,
		absent,
		pending,
	} {
		x.cooperate(c, volunteer)
	}

	id := x.event(c, 3)
	first := x.participate(c, a, id)
	second := x.participate(c, absent, id)
	waiting, err := x.clients.participation.RequestParticipation(ctx,
		request(pending, &participation_v1.RequestParticipationRequest{EventId: id}))
	x.Require().NoError(err)

	results := []*event_v1.EventParticipationResult{
		{
			ParticipationId: first,
			Outcome:         &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 180},
		},
		{
			ParticipationId: second,
			Outcome:         &event_v1.EventParticipationResult_Absent{Absent: true},
		},
	}
	completion := &event_v1.CompleteEventRequest{
		Id:        id,
		RequestId: "finish",
		Results:   results,
	}
	_, err = x.clients.event.CompleteEvent(ctx, request(c, completion))
	callError(x.T(), err, connect.CodeFailedPrecondition)
	x.finishTime(id)
	_, err = x.clients.event.CompleteEvent(ctx, request(a, completion))
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.event.CompleteEvent(ctx, request(c, &event_v1.CompleteEventRequest{
		Id:        id,
		RequestId: "finish",
		Results:   results[:1],
	}))
	callError(x.T(), err, connect.CodeInvalidArgument)
	before, err := x.clients.volunteer.GetVolunteer(ctx, request(a, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().Zero(before.Msg.GetVolunteer().GetCoinBalance())
	completed, err := x.clients.event.CompleteEvent(ctx, request(c, completion))
	x.Require().NoError(err)
	x.Require().Equal(event_v1.EventStatus_EVENT_STATUS_COMPLETED, completed.Msg.GetEvent().GetStatus())
	_, err = x.clients.event.CompleteEvent(ctx, request(c, completion))
	x.Require().NoError(err)
	_, err = x.clients.event.CompleteEvent(ctx, request(c, &event_v1.CompleteEventRequest{
		Id:        id,
		RequestId: "different",
		Results:   results,
	}))
	callError(x.T(), err, connect.CodeFailedPrecondition)

	changed := proto.CloneOf(completion)
	changed.Results[0].Outcome = &event_v1.EventParticipationResult_ConfirmedMinutes{ConfirmedMinutes: 60}
	_, err = x.clients.event.CompleteEvent(ctx, request(c, changed))
	callError(x.T(), err, connect.CodeAlreadyExists)
	profile, err := x.clients.volunteer.GetVolunteer(ctx, request(a, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().EqualValues(30, profile.Msg.GetVolunteer().GetCoinBalance())
	x.Require().EqualValues(180, profile.Msg.GetVolunteer().GetConfirmedMinutes())
	noAward, err := x.clients.volunteer.GetVolunteer(ctx, request(absent, &volunteer_v1.GetVolunteerRequest{}))
	x.Require().NoError(err)
	x.Require().Zero(noAward.Msg.GetVolunteer().GetCoinBalance())
	closed, err := x.clients.participation.GetParticipation(ctx,
		request(pending, &participation_v1.GetParticipationRequest{Id: waiting.Msg.GetId()}))
	x.Require().NoError(err)
	x.Require().Equal(
		participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CLOSED,
		closed.Msg.GetParticipation().GetStatus(),
	)
	history, err := x.clients.volunteer.ListVolunteerOperations(ctx,
		request(a, &volunteer_v1.ListVolunteerOperationsRequest{}))
	x.Require().NoError(err)
	x.Require().Len(history.Msg.GetOperations(), 1)
	x.Require().EqualValues(30, history.Msg.GetOperations()[0].GetCoinDelta())
	_, err = x.clients.event.UpdateEvent(ctx, request(c, &event_v1.UpdateEventRequest{Id: id}))
	callError(x.T(), err, connect.CodeFailedPrecondition)
}

func (x *RPCTestSuite) TestParticipationLifecycleAndCapacity() {
	ctx := x.T().Context()
	c, other, a, b := x.coordinator(), x.coordinator(), x.login(), x.login()
	x.cooperate(c, a)
	x.cooperate(c, b)
	id := x.event(c, 1)
	first := x.participate(c, a, id)
	_, err := x.clients.event.UpdateEvent(
		ctx,
		request(c, &event_v1.UpdateEventRequest{
			Id:    id,
			Title: proto.String("Нельзя"),
		}),
	)
	callError(x.T(), err, connect.CodeFailedPrecondition)
	_, err = x.clients.participation.RequestParticipation(ctx,
		request(a, &participation_v1.RequestParticipationRequest{EventId: id}))
	callError(x.T(), err, connect.CodeAlreadyExists)
	second, err := x.clients.participation.RequestParticipation(ctx,
		request(b, &participation_v1.RequestParticipationRequest{EventId: id}))
	x.Require().NoError(err)
	_, err = x.clients.participation.ApproveParticipation(ctx,
		request(other, &participation_v1.ApproveParticipationRequest{Id: second.Msg.GetId()}))
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.participation.ApproveParticipation(ctx,
		request(c, &participation_v1.ApproveParticipationRequest{Id: second.Msg.GetId()}))
	callError(x.T(), err, connect.CodeFailedPrecondition)
	_, err = x.clients.participation.CancelParticipation(ctx,
		request(a, &participation_v1.CancelParticipationRequest{Id: first}))
	x.Require().NoError(err)
	_, err = x.clients.participation.ApproveParticipation(ctx,
		request(c, &participation_v1.ApproveParticipationRequest{Id: second.Msg.GetId()}))
	x.Require().NoError(err)
	_, err = x.clients.participation.RestoreParticipation(ctx,
		request(a, &participation_v1.RestoreParticipationRequest{Id: first}))
	x.Require().NoError(err)
	_, err = x.clients.participation.RejectParticipation(ctx, request(c, &participation_v1.RejectParticipationRequest{
		Id:     second.Msg.GetId(),
		Reason: proto.String("  Изменились условия  "),
	}))
	x.Require().NoError(err)
	rejected, err := x.clients.participation.GetParticipation(ctx,
		request(b, &participation_v1.GetParticipationRequest{Id: second.Msg.GetId()}))
	x.Require().NoError(err)
	x.Require().Equal("Изменились условия", rejected.Msg.GetParticipation().GetRejectionReason())
	_, err = x.clients.participation.GetParticipation(ctx,
		request(a, &participation_v1.GetParticipationRequest{Id: second.Msg.GetId()}))
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.participation.RestoreParticipation(ctx,
		request(b, &participation_v1.RestoreParticipationRequest{Id: second.Msg.GetId()}))
	callError(x.T(), err, connect.CodeFailedPrecondition)
	_, err = x.clients.participation.ApproveParticipation(ctx,
		request(c, &participation_v1.ApproveParticipationRequest{Id: first}))
	x.Require().NoError(err)
	history, err := x.clients.participation.ListParticipationChanges(ctx,
		request(a, &participation_v1.ListParticipationChangesRequest{ParticipationId: first}))
	x.Require().NoError(err)
	x.Require().Len(history.Msg.GetChanges(), 5)
	listed, err := x.clients.participation.ListParticipations(ctx, request(c, &participation_v1.ListParticipationsRequest{
		Scope:   participation_v1.ParticipationListScope_PARTICIPATION_LIST_SCOPE_MY_EVENTS,
		EventId: proto.String(id),
		Status:  participation_v1.ParticipationStatus_PARTICIPATION_STATUS_APPROVED.Enum(),
	}))
	x.Require().NoError(err)
	x.Require().Len(listed.Msg.GetParticipations(), 1)
	x.Require().Equal(first, listed.Msg.GetParticipations()[0].GetId())
}
