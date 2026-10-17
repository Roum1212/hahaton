//go:build e2e

package rpc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
	updates_v1 "github.com/Roum1212/hahaton/genproto/updates/v1"
	updates_v1connect "github.com/Roum1212/hahaton/genproto/updates/v1/updates_v1connect"
)

const noUpdateTimeout = 500 * time.Millisecond

func (x *RPCTestSuite) TestSharedUpdatesThroughBufferingProxy() {
	ctx := x.T().Context()
	coordinator, volunteer := x.coordinator(), x.login()
	proxy := x.bufferingProxy()
	client := updates_v1connect.NewUpdatesServiceClient(&http.Client{Timeout: requestTimeout}, proxy)

	streamCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	coordinatorStream, err := client.SubscribeUpdates(streamCtx,
		request(coordinator, &updates_v1.SubscribeUpdatesRequest{}))
	x.Require().NoError(err)

	x.T().Cleanup(func() { x.Require().NoError(coordinatorStream.Close()) })

	volunteerStream, err := client.SubscribeUpdates(streamCtx, request(volunteer, &updates_v1.SubscribeUpdatesRequest{}))
	x.Require().NoError(err)

	x.T().Cleanup(func() { x.Require().NoError(volunteerStream.Close()) })

	x.update(coordinatorStream, updates_v1.UpdateKind_UPDATE_KIND_READY)
	x.update(volunteerStream, updates_v1.UpdateKind_UPDATE_KIND_READY)

	_, err = x.clients.cooperation.RequestCooperation(ctx, request(volunteer,
		&cooperation_v1.RequestCooperationRequest{CoordinatorId: coordinator.id}))
	x.Require().NoError(err)
	x.changed(coordinatorStream, volunteerStream)
	_, err = x.clients.cooperation.ApproveCooperation(ctx, request(coordinator,
		&cooperation_v1.ApproveCooperationRequest{VolunteerId: volunteer.id}))
	x.Require().NoError(err)
	x.changed(coordinatorStream, volunteerStream)

	eventID := x.event(coordinator, 1)
	x.changed(coordinatorStream, volunteerStream)
	participation, err := x.clients.participation.RequestParticipation(ctx, request(volunteer,
		&participation_v1.RequestParticipationRequest{EventId: eventID}))
	x.Require().NoError(err)
	x.changed(coordinatorStream, volunteerStream)
	_, err = x.clients.participation.ApproveParticipation(ctx, request(coordinator,
		&participation_v1.ApproveParticipationRequest{Id: participation.Msg.GetId()}))
	x.Require().NoError(err)
	x.changed(coordinatorStream, volunteerStream)

	_, err = x.clients.eventChat.SendEventChatMessage(ctx, request(volunteer,
		&event_chat_v1.SendEventChatMessageRequest{EventId: eventID, Text: "Общий стрим", RequestId: "shared-stream"}))
	x.Require().NoError(err)
	x.changed(coordinatorStream, volunteerStream)
	messages, err := x.clients.eventChat.ListEventChatMessages(ctx, request(coordinator,
		&event_chat_v1.ListEventChatMessagesRequest{EventId: eventID}))
	x.Require().NoError(err)
	x.Require().Len(messages.Msg.GetMessages(), 1)
	x.Require().Equal("Общий стрим", messages.Msg.GetMessages()[0].GetText())

	_, err = x.clients.participation.RejectParticipation(ctx, request(coordinator,
		&participation_v1.RejectParticipationRequest{Id: participation.Msg.GetId()}))
	x.Require().NoError(err)
	x.changed(coordinatorStream, volunteerStream)
	_, err = x.clients.eventChat.ListEventChatMessages(ctx, request(volunteer,
		&event_chat_v1.ListEventChatMessagesRequest{EventId: eventID}))
	callError(x.T(), err, connect.CodePermissionDenied)

	_, err = x.clients.session.CloseSession(ctx, request(volunteer, &emptypb.Empty{}))
	x.Require().NoError(err)
	x.Require().False(volunteerStream.Receive())
	callError(x.T(), volunteerStream.Err(), connect.CodeUnauthenticated)
}

func (x *RPCTestSuite) TestUpdatesRequireAuthenticatedSession() {
	stream, err := x.clients.updates.SubscribeUpdates(x.T().Context(),
		connect.NewRequest(&updates_v1.SubscribeUpdatesRequest{}))
	if err == nil {
		x.Require().False(stream.Receive())
		err = stream.Err()
		x.Require().NoError(stream.Close())
	}

	callError(x.T(), err, connect.CodeUnauthenticated)
}

func (x *RPCTestSuite) TestUpdatesAreScopedAndFailedWritesDoNotNotify() {
	ctx := x.T().Context()
	coordinator, volunteer, stranger := x.coordinator(), x.login(), x.login()
	x.cooperate(coordinator, volunteer)
	eventID := x.event(coordinator, 1)
	participationID := x.participate(coordinator, volunteer, eventID)

	streamCtx, cancel := context.WithTimeout(ctx, noUpdateTimeout)
	defer cancel()

	stream, err := x.clients.updates.SubscribeUpdates(streamCtx, request(stranger, &updates_v1.SubscribeUpdatesRequest{}))
	x.Require().NoError(err)

	x.T().Cleanup(func() { x.Require().NoError(stream.Close()) })

	x.update(stream, updates_v1.UpdateKind_UPDATE_KIND_READY)

	_, err = x.clients.participation.RejectParticipation(ctx,
		request(stranger, &participation_v1.RejectParticipationRequest{Id: participationID}))
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.eventChat.SendEventChatMessage(ctx, request(coordinator,
		&event_chat_v1.SendEventChatMessageRequest{EventId: eventID, Text: "Закрытый чат", RequestId: "scoped"}))
	x.Require().NoError(err)
	x.Require().False(stream.Receive())
	callError(x.T(), stream.Err(), connect.CodeDeadlineExceeded)
}

func (x *RPCTestSuite) TestUpdatesReconnectResyncsStoredMessages() {
	ctx := x.T().Context()
	coordinator, volunteer := x.coordinator(), x.login()
	x.cooperate(coordinator, volunteer)
	eventID := x.event(coordinator, 1)
	x.participate(coordinator, volunteer, eventID)

	firstCtx, stopFirst := context.WithCancel(ctx)
	first, err := x.clients.updates.SubscribeUpdates(firstCtx, request(volunteer, &updates_v1.SubscribeUpdatesRequest{}))
	x.Require().NoError(err)
	x.update(first, updates_v1.UpdateKind_UPDATE_KIND_READY)
	stopFirst()
	x.Require().NoError(first.Close())

	for i := range 20 {
		_, err := x.clients.eventChat.SendEventChatMessage(ctx, request(coordinator,
			&event_chat_v1.SendEventChatMessageRequest{
				EventId: eventID, Text: "Сохранённое сообщение", RequestId: fmt.Sprintf("gap-%d", i),
			}))
		x.Require().NoError(err)
	}

	secondCtx, stopSecond := context.WithTimeout(ctx, requestTimeout)
	defer stopSecond()

	second, err := x.clients.updates.SubscribeUpdates(secondCtx, request(volunteer, &updates_v1.SubscribeUpdatesRequest{}))
	x.Require().NoError(err)

	x.T().Cleanup(func() { x.Require().NoError(second.Close()) })

	x.update(second, updates_v1.UpdateKind_UPDATE_KIND_READY)
	history, err := x.clients.eventChat.ListEventChatMessages(ctx, request(volunteer,
		&event_chat_v1.ListEventChatMessagesRequest{EventId: eventID, Limit: 100}))
	x.Require().NoError(err)
	x.Require().Len(history.Msg.GetMessages(), 20)
}

func (x *RPCTestSuite) TestUpdatesDetectExpiredSessionWhileIdle() {
	ctx := x.T().Context()
	volunteer := x.login()
	client := updates_v1connect.NewUpdatesServiceClient(&http.Client{
		Transport: x.httpClient.Transport,
		Timeout:   15 * time.Second,
	}, x.baseURL)

	streamCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	stream, err := client.SubscribeUpdates(streamCtx, request(volunteer, &updates_v1.SubscribeUpdatesRequest{}))
	x.Require().NoError(err)
	x.T().Cleanup(func() { x.Require().NoError(stream.Close()) })
	x.update(stream, updates_v1.UpdateKind_UPDATE_KIND_READY)
	_, err = x.pool.Exec(ctx,
		"UPDATE sessions SET expires_at = now() - interval '1 second' WHERE user_id = $1", volunteer.id,
	)
	x.Require().NoError(err)
	x.Require().False(stream.Receive())
	callError(x.T(), stream.Err(), connect.CodeUnauthenticated)
}

func (x *RPCTestSuite) TestClosingOneTabKeepsOtherSubscription() {
	ctx := x.T().Context()
	volunteer := x.login()

	firstCtx, closeFirst := context.WithCancel(ctx)
	defer closeFirst()

	first, err := x.clients.updates.SubscribeUpdates(firstCtx, request(volunteer, &updates_v1.SubscribeUpdatesRequest{}))
	x.Require().NoError(err)
	x.update(first, updates_v1.UpdateKind_UPDATE_KIND_READY)

	secondCtx, closeSecond := context.WithTimeout(ctx, requestTimeout)
	defer closeSecond()

	second, err := x.clients.updates.SubscribeUpdates(secondCtx, request(volunteer, &updates_v1.SubscribeUpdatesRequest{}))
	x.Require().NoError(err)
	x.T().Cleanup(func() { x.Require().NoError(second.Close()) })
	x.update(second, updates_v1.UpdateKind_UPDATE_KIND_READY)
	closeFirst()
	x.Require().NoError(first.Close())
	coordinator := x.coordinator()
	_, err = x.clients.cooperation.RequestCooperation(ctx, request(volunteer,
		&cooperation_v1.RequestCooperationRequest{CoordinatorId: coordinator.id}))
	x.Require().NoError(err)
	x.update(second, updates_v1.UpdateKind_UPDATE_KIND_CHANGED)
}

func (x *RPCTestSuite) changed(streams ...*connect.ServerStreamForClient[updates_v1.SubscribeUpdatesResponse]) {
	for _, stream := range streams {
		x.update(stream, updates_v1.UpdateKind_UPDATE_KIND_CHANGED)
	}
}

func (x *RPCTestSuite) update(
	stream *connect.ServerStreamForClient[updates_v1.SubscribeUpdatesResponse],
	kind updates_v1.UpdateKind,
) {
	x.Require().True(stream.Receive(), "%v", stream.Err())
	x.Require().Equal(kind, stream.Msg().GetKind())
}
