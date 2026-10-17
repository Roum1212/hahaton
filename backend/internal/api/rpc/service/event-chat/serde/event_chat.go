package rpc_service_event_chat_serde

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	event_chat_v1 "github.com/Roum1212/hahaton/genproto/event-chat/v1"
	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_chat_message "github.com/Roum1212/hahaton/internal/domain/model/chat-message"
	model_event_status "github.com/Roum1212/hahaton/internal/domain/model/event-status"
)

func AggregateToPBEventChat(
	event *aggregate_event.Event,
	chat *aggregate_event_chat.EventChat,
	boundary uint64,
) *event_chat_v1.EventChat {
	return &event_chat_v1.EventChat{
		EventId:              event.GetEventID().String(),
		EventTitle:           event.GetTitle().String(),
		ReadOnly:             event.GetStatus() != model_event_status.Published,
		LatestSequence:       chat.GetLatestSequence(),
		VisibleAfterSequence: boundary,
	}
}

func MessageToPB(value model_chat_message.ChatMessage) *event_chat_v1.EventChatMessage {
	result := &event_chat_v1.EventChatMessage{
		Id:                value.GetMessageID().String(),
		EventId:           value.GetEventID().String(),
		AuthorDisplayName: value.GetAuthorDisplayName(),
		Text:              value.GetText(),
		Sequence:          value.GetSequence(),
		CreatedAt:         timestamppb.New(value.GetCreatedAt()),
	}
	switch value.GetAuthorKind() {
	case model_actor_kind.Coordinator:
		result.Author = &event_chat_v1.EventChatMessage_CoordinatorId{CoordinatorId: value.GetAuthorID().String()}
	case model_actor_kind.Volunteer:
		result.Author = &event_chat_v1.EventChatMessage_VolunteerId{VolunteerId: value.GetAuthorID().String()}
	case model_actor_kind.System:
	}

	return result
}
