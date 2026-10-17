package service_list_my_event_chats

import (
	"fmt"

	aggregate_event "github.com/Roum1212/hahaton/internal/domain/aggregate/event"
	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID model_user_id.UserID
	page   model_page.Page
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.page.Validate(); err != nil {
		return fmt.Errorf("page: %w", err)
	}

	return nil
}

func NewRequest(userID model_user_id.UserID, page model_page.Page) Request {
	return Request{
		userID: userID,
		page:   page,
	}
}

type Response struct {
	events     []*aggregate_event.Event
	chats      map[model_event_id.EventID]*aggregate_event_chat.EventChat
	boundaries map[model_event_id.EventID]uint64
	totalCount int64
}

func (x Response) GetBoundaries() map[model_event_id.EventID]uint64 {
	return x.boundaries
}

func (x Response) GetChats() map[model_event_id.EventID]*aggregate_event_chat.EventChat {
	return x.chats
}

func (x Response) GetEvents() []*aggregate_event.Event {
	return x.events
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func NewResponse(
	events []*aggregate_event.Event,
	chats map[model_event_id.EventID]*aggregate_event_chat.EventChat,
	boundaries map[model_event_id.EventID]uint64,
	totalCount int64,
) Response {
	return Response{
		events:     events,
		chats:      chats,
		boundaries: boundaries,
		totalCount: totalCount,
	}
}
