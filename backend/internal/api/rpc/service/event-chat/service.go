package rpc_service_event_chat

import (
	"context"

	"connectrpc.com/connect"

	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	event_chat_v1connect "github.com/Roum1212/hahaton/genproto/event-chat/v1/event_chat_v1connect"
	rpc_service_event_chat_rpc_get_event_chat "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/rpc/get-event-chat"
	rpc_service_event_chat_rpc_list_event_chat_messages "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/rpc/list-event-chat-messages"
	rpc_service_event_chat_rpc_list_my_event_chats "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/rpc/list-my-event-chats"
	rpc_service_event_chat_rpc_send_event_chat_message "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/rpc/send-event-chat-message"
	rpc_service_event_chat_rpc_subscribe_event_chat_messages "github.com/Roum1212/hahaton/internal/api/rpc/service/event-chat/rpc/subscribe-event-chat-messages"
)

var _ event_chat_v1connect.EventChatServiceHandler = Service{}

type Service struct {
	getEventChatRPC               rpc_service_event_chat_rpc_get_event_chat.RPC
	listEventChatMessagesRPC      rpc_service_event_chat_rpc_list_event_chat_messages.RPC
	listMyEventChatsRPC           rpc_service_event_chat_rpc_list_my_event_chats.RPC
	sendEventChatMessageRPC       rpc_service_event_chat_rpc_send_event_chat_message.RPC
	subscribeEventChatMessagesRPC rpc_service_event_chat_rpc_subscribe_event_chat_messages.RPC
}

func (x Service) GetEventChat(
	ctx context.Context,
	r *connect.Request[event_chat_v1.GetEventChatRequest],
) (*connect.Response[event_chat_v1.GetEventChatResponse], error) {
	return x.getEventChatRPC.GetEventChat(ctx, r)
}

func (x Service) ListEventChatMessages(
	ctx context.Context,
	r *connect.Request[event_chat_v1.ListEventChatMessagesRequest],
) (*connect.Response[event_chat_v1.ListEventChatMessagesResponse], error) {
	return x.listEventChatMessagesRPC.ListEventChatMessages(ctx, r)
}

func (x Service) ListMyEventChats(
	ctx context.Context,
	r *connect.Request[event_chat_v1.ListMyEventChatsRequest],
) (*connect.Response[event_chat_v1.ListMyEventChatsResponse], error) {
	return x.listMyEventChatsRPC.ListMyEventChats(ctx, r)
}

func (x Service) SendEventChatMessage(
	ctx context.Context,
	r *connect.Request[event_chat_v1.SendEventChatMessageRequest],
) (*connect.Response[event_chat_v1.SendEventChatMessageResponse], error) {
	return x.sendEventChatMessageRPC.SendEventChatMessage(ctx, r)
}

func (x Service) SubscribeEventChatMessages(
	ctx context.Context,
	r *connect.Request[event_chat_v1.SubscribeEventChatMessagesRequest],
	stream *connect.ServerStream[event_chat_v1.SubscribeEventChatMessagesResponse],
) error {
	return x.subscribeEventChatMessagesRPC.SubscribeEventChatMessages(ctx, r, stream)
}

func NewService(
	getEventChatRPC rpc_service_event_chat_rpc_get_event_chat.RPC,
	listEventChatMessagesRPC rpc_service_event_chat_rpc_list_event_chat_messages.RPC,
	listMyEventChatsRPC rpc_service_event_chat_rpc_list_my_event_chats.RPC,
	sendEventChatMessageRPC rpc_service_event_chat_rpc_send_event_chat_message.RPC,
	subscribeEventChatMessagesRPC rpc_service_event_chat_rpc_subscribe_event_chat_messages.RPC,
) Service {
	return Service{
		getEventChatRPC:               getEventChatRPC,
		listEventChatMessagesRPC:      listEventChatMessagesRPC,
		listMyEventChatsRPC:           listMyEventChatsRPC,
		sendEventChatMessageRPC:       sendEventChatMessageRPC,
		subscribeEventChatMessagesRPC: subscribeEventChatMessagesRPC,
	}
}
