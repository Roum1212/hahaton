//go:build integration

package integration

import (
	"context"
	"fmt"
	"net/http/httptest"
	"time"

	"connectrpc.com/connect"

	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	event_chat_v1connect "github.com/Roum1212/hahaton/genproto/event-chat/v1/event_chat_v1connect"
)

func checkChat(x *fixture) {
	c, a := x.coordinator("chat-owner"), x.login("chat-volunteer")
	x.cooperate(c, a)
	id := x.event(c, "Переписка", 1)
	before := x.call(c, "event_chat.v1.EventChatService/SendEventChatMessage", map[string]any{
		"event_id":   id,
		"text":       "До принятия",
		"request_id": "before",
	})
	require(x.t, object(x.t, before, "message")["coordinator_id"] == c.id, "coordinator author kind missing")
	participationID := x.participate(c, a, id)
	initial := x.call(a, "event_chat.v1.EventChatService/ListEventChatMessages", map[string]any{"event_id": id})
	require(x.t, len(list(x.t, initial, "messages")) == 0, "initial admission revealed old messages")
	sent := x.call(a, "event_chat.v1.EventChatService/SendEventChatMessage", map[string]any{
		"event_id":   id,
		"text":       "После принятия",
		"request_id": "mine",
	})
	messageID := object(x.t, sent, "message")["id"]
	repeated := x.call(a, "event_chat.v1.EventChatService/SendEventChatMessage", map[string]any{
		"event_id":   id,
		"text":       "После принятия",
		"request_id": "mine",
	})
	require(x.t, object(x.t, repeated, "message")["id"] == messageID, "message duplicated on retry")
	x.fail(a, "event_chat.v1.EventChatService/SendEventChatMessage", map[string]any{
		"event_id":   id,
		"text":       "Другой текст",
		"request_id": "mine",
	}, "already_exists")
	snapshot := x.call(a, "event_chat.v1.EventChatService/ListEventChatMessages", map[string]any{
		"event_id": id,
		"limit":    1,
	})
	x.call(c, "event_chat.v1.EventChatService/SendEventChatMessage", map[string]any{
		"event_id":   id,
		"text":       "Новое сообщение",
		"request_id": "new",
	})
	stable := x.call(a, "event_chat.v1.EventChatService/ListEventChatMessages", map[string]any{
		"event_id":          id,
		"snapshot_sequence": snapshot["snapshot_sequence"],
	})
	require(x.t, len(list(x.t, stable, "messages")) == 1, "fixed snapshot shifted after new message")
	x.call(a, "participation.v1.ParticipationService/CancelParticipation", map[string]any{"id": participationID})
	x.fail(a, "event_chat.v1.EventChatService/GetEventChat", map[string]any{"event_id": id}, "permission_denied")
	x.call(a, "participation.v1.ParticipationService/RestoreParticipation", map[string]any{"id": participationID})
	x.call(c, "participation.v1.ParticipationService/ApproveParticipation", map[string]any{"id": participationID})
	x.fail(a, "event_chat.v1.EventChatService/SendEventChatMessage", map[string]any{
		"event_id":   id,
		"text":       "После принятия",
		"request_id": "mine",
	}, "failed_precondition")
	x.fail(a, "event_chat.v1.EventChatService/SendEventChatMessage", map[string]any{
		"event_id":   id,
		"text":       "Другой текст",
		"request_id": "mine",
	}, "failed_precondition")
	history := x.call(a, "event_chat.v1.EventChatService/ListEventChatMessages", map[string]any{"event_id": id})
	require(x.t, len(list(x.t, history, "messages")) == 0, "readmission revealed old messages")
	x.call(a, "event_chat.v1.EventChatService/SendEventChatMessage", map[string]any{
		"event_id":   id,
		"text":       "Новый период",
		"request_id": "new-period",
	})
	chats := x.call(a, "event_chat.v1.EventChatService/ListMyEventChats", nil)
	require(x.t, len(list(x.t, chats, "chats")) == 1, "my chats list missing approved participation")

	server := httptest.NewServer(x.handler)
	defer server.Close()

	client := event_chat_v1connect.NewEventChatServiceClient(server.Client(), server.URL)
	x.seen["event_chat.v1.EventChatService/SubscribeEventChatMessages"] = true

	ctx, cancel := context.WithTimeout(x.t.Context(), 5*time.Second)
	defer cancel()

	request := connect.NewRequest(&event_chat_v1.SubscribeEventChatMessagesRequest{EventId: id})
	request.Header().Set("Cookie", a.cookie.String())
	stream, err := client.SubscribeEventChatMessages(ctx, request)
	must(x.t, err)
	require(x.t, stream.Receive(), "stream failed to deliver new period: "+fmt.Sprint(stream.Err()))
	require(x.t, stream.Msg().GetMessage().GetText() == "Новый период", "stream leaked old period")
	x.call(c, "participation.v1.ParticipationService/RejectParticipation", map[string]any{"id": participationID})
	require(
		x.t,
		!stream.Receive() && connect.CodeOf(stream.Err()) == connect.CodePermissionDenied,
		"idle stream was not revoked",
	)
	x.finishTime(id)
	x.call(c, "event.v1.EventService/CompleteEvent", map[string]any{
		"id":         id,
		"request_id": "close-chat",
	})
	readOnly := object(
		x.t,
		x.call(c, "event_chat.v1.EventChatService/GetEventChat", map[string]any{"event_id": id}),
		"chat",
	)
	require(x.t, readOnly["read_only"] == true, "completed chat is writable")
	x.fail(c, "event_chat.v1.EventChatService/SendEventChatMessage", map[string]any{
		"event_id":   id,
		"text":       "Поздно",
		"request_id": "late",
	}, "failed_precondition")
	request = connect.NewRequest(&event_chat_v1.SubscribeEventChatMessagesRequest{EventId: id})
	request.Header().Set("Cookie", c.cookie.String())
	completed, err := client.SubscribeEventChatMessages(ctx, request)
	must(x.t, err)

	count := 0
	for completed.Receive() {
		count++
	}

	must(x.t, completed.Err())
	require(x.t, count == 4, "completed stream lost persisted messages")

	var plaintext int

	query := "SELECT count(*) FROM chat_messages WHERE position($1::bytea in encrypted_text)>0"
	must(x.t, x.pool.QueryRow(x.t.Context(), query, []byte("Новый период")).Scan(&plaintext))
	require(x.t, plaintext == 0, "chat text stored as plaintext")
}
