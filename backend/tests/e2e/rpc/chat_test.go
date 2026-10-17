//go:build e2e

package rpc

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	cooperation_v1 "github.com/Roum1212/hahaton/genproto/cooperation/v1"
	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	event_v1 "github.com/Roum1212/hahaton/genproto/event/v1"
	participation_v1 "github.com/Roum1212/hahaton/genproto/participation/v1"
)

func (x *RPCTestSuite) TestChatHistoryStreamingAndRevocation() {
	ctx := x.T().Context()
	c, a := x.coordinator(), x.login()
	x.cooperate(c, a)
	id := x.event(c, 1)
	old, err := x.clients.eventChat.SendEventChatMessage(ctx, request(c, &event_chat_v1.SendEventChatMessageRequest{
		EventId:   id,
		Text:      "До принятия",
		RequestId: "old",
	}))
	x.Require().NoError(err)
	x.Require().Equal(c.id, old.Msg.GetMessage().GetCoordinatorId())
	participationID := x.participate(c, a, id)
	initial, err := x.clients.eventChat.ListEventChatMessages(ctx,
		request(a, &event_chat_v1.ListEventChatMessagesRequest{EventId: id}))
	x.Require().NoError(err)
	x.Require().Empty(initial.Msg.GetMessages())

	send := &event_chat_v1.SendEventChatMessageRequest{
		EventId:   id,
		Text:      "После принятия",
		RequestId: "first",
	}
	sent, err := x.clients.eventChat.SendEventChatMessage(ctx, request(a, send))
	x.Require().NoError(err)
	x.Require().Equal(a.id, sent.Msg.GetMessage().GetVolunteerId())
	repeated, err := x.clients.eventChat.SendEventChatMessage(ctx, request(a, send))
	x.Require().NoError(err)
	x.Require().Equal(sent.Msg.GetMessage().GetId(), repeated.Msg.GetMessage().GetId())
	_, err = x.clients.eventChat.SendEventChatMessage(ctx, request(a, &event_chat_v1.SendEventChatMessageRequest{
		EventId:   id,
		Text:      "Другой текст",
		RequestId: "first",
	}))
	callError(x.T(), err, connect.CodeAlreadyExists)
	snapshot, err := x.clients.eventChat.ListEventChatMessages(ctx,
		request(a, &event_chat_v1.ListEventChatMessagesRequest{
			EventId: id,
			Limit:   1,
		}))
	x.Require().NoError(err)
	_, err = x.clients.eventChat.SendEventChatMessage(ctx, request(c, &event_chat_v1.SendEventChatMessageRequest{
		EventId:   id,
		Text:      "Свежее сообщение",
		RequestId: "new",
	}))
	x.Require().NoError(err)
	stable, err := x.clients.eventChat.ListEventChatMessages(ctx, request(a, &event_chat_v1.ListEventChatMessagesRequest{
		EventId:          id,
		SnapshotSequence: proto.Uint64(snapshot.Msg.GetSnapshotSequence()),
	}))
	x.Require().NoError(err)
	x.Require().Len(stable.Msg.GetMessages(), 1)
	_, err = x.clients.participation.CancelParticipation(ctx,
		request(a, &participation_v1.CancelParticipationRequest{Id: participationID}))
	x.Require().NoError(err)
	_, err = x.clients.eventChat.GetEventChat(ctx, request(a, &event_chat_v1.GetEventChatRequest{EventId: id}))
	callError(x.T(), err, connect.CodePermissionDenied)
	_, err = x.clients.participation.RestoreParticipation(ctx,
		request(a, &participation_v1.RestoreParticipationRequest{Id: participationID}))
	x.Require().NoError(err)
	_, err = x.clients.participation.ApproveParticipation(ctx,
		request(c, &participation_v1.ApproveParticipationRequest{Id: participationID}))
	x.Require().NoError(err)
	reset, err := x.clients.eventChat.ListEventChatMessages(ctx,
		request(a, &event_chat_v1.ListEventChatMessagesRequest{EventId: id}))
	x.Require().NoError(err)
	x.Require().Empty(reset.Msg.GetMessages())
	_, err = x.clients.eventChat.SendEventChatMessage(ctx, request(a, &event_chat_v1.SendEventChatMessageRequest{
		EventId:   id,
		Text:      "Новый период",
		RequestId: "period",
	}))
	x.Require().NoError(err)
	chats, err := x.clients.eventChat.ListMyEventChats(ctx, request(a, &event_chat_v1.ListMyEventChatsRequest{}))
	x.Require().NoError(err)
	x.Require().Len(chats.Msg.GetChats(), 1)

	streamCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	stream, err := x.clients.eventChat.SubscribeEventChatMessages(streamCtx,
		request(a, &event_chat_v1.SubscribeEventChatMessagesRequest{EventId: id}))
	x.Require().NoError(err)
	x.Require().True(stream.Receive(), "%v", stream.Err())
	x.Require().Equal("Новый период", stream.Msg().GetMessage().GetText())
	_, err = x.clients.eventChat.SendEventChatMessage(ctx, request(c, &event_chat_v1.SendEventChatMessageRequest{
		EventId:   id,
		Text:      "В эфире",
		RequestId: "live",
	}))
	x.Require().NoError(err)
	x.Require().True(stream.Receive(), "%v", stream.Err())
	x.Require().Equal("В эфире", stream.Msg().GetMessage().GetText())
	_, err = x.clients.participation.RejectParticipation(ctx,
		request(c, &participation_v1.RejectParticipationRequest{Id: participationID}))
	x.Require().NoError(err)
	x.Require().False(stream.Receive())
	callError(x.T(), stream.Err(), connect.CodePermissionDenied)
	x.Require().NoError(stream.Close())
	x.finishTime(id)
	_, err = x.clients.event.CompleteEvent(
		ctx,
		request(c, &event_v1.CompleteEventRequest{
			Id:        id,
			RequestId: "closed",
		}),
	)
	x.Require().NoError(err)
	chat, err := x.clients.eventChat.GetEventChat(ctx, request(c, &event_chat_v1.GetEventChatRequest{EventId: id}))
	x.Require().NoError(err)
	x.Require().True(chat.Msg.GetChat().GetReadOnly())
	_, err = x.clients.eventChat.SendEventChatMessage(ctx, request(c, &event_chat_v1.SendEventChatMessageRequest{
		EventId:   id,
		Text:      "Поздно",
		RequestId: "late",
	}))
	callError(x.T(), err, connect.CodeFailedPrecondition)
	completed, err := x.clients.eventChat.SubscribeEventChatMessages(streamCtx,
		request(c, &event_chat_v1.SubscribeEventChatMessagesRequest{EventId: id}))
	x.Require().NoError(err)

	count := 0
	for completed.Receive() {
		count++
	}

	x.Require().NoError(completed.Err())
	x.Require().Equal(5, count)
	x.Require().NoError(completed.Close())
}

func (x *RPCTestSuite) TestCooperationTerminationRevokesChatAndParticipation() {
	ctx := x.T().Context()
	c, a := x.coordinator(), x.login()
	x.cooperate(c, a)
	id := x.event(c, 1)
	participationID := x.participate(c, a, id)
	_, err := x.clients.cooperation.TerminateCooperation(ctx, request(c, &cooperation_v1.TerminateCooperationRequest{
		Counterpart: &cooperation_v1.TerminateCooperationRequest_VolunteerId{VolunteerId: a.id},
	}))
	x.Require().NoError(err)
	application, err := x.clients.participation.GetParticipation(ctx,
		request(a, &participation_v1.GetParticipationRequest{Id: participationID}))
	x.Require().NoError(err)
	x.Require().Equal(
		participation_v1.ParticipationStatus_PARTICIPATION_STATUS_CANCELLED,
		application.Msg.GetParticipation().GetStatus(),
	)
	_, err = x.clients.eventChat.GetEventChat(ctx, request(a, &event_chat_v1.GetEventChatRequest{EventId: id}))
	callError(x.T(), err, connect.CodePermissionDenied)
	x.cooperate(c, a)
	_, err = x.clients.eventChat.GetEventChat(ctx, request(a, &event_chat_v1.GetEventChatRequest{EventId: id}))
	callError(x.T(), err, connect.CodePermissionDenied)
}
