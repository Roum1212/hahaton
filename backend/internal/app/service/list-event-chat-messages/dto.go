package service_list_event_chat_messages

import (
	"fmt"

	model_chat_message "github.com/Roum1212/hahaton/internal/domain/model/chat-message"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	userID  model_user_id.UserID
	eventID model_event_id.EventID

	snapshotSequence utils_value.Optional[uint64]

	page model_page.Page
}

func (x Request) GetEventID() model_event_id.EventID {
	return x.eventID
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetSnapshotSequence() utils_value.Optional[uint64] {
	return x.snapshotSequence
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.eventID.Validate(); err != nil {
		return fmt.Errorf("event id: %w", err)
	}

	if err := x.page.Validate(); err != nil {
		return fmt.Errorf("page: %w", err)
	}

	return nil
}

type RequestOption func(x *Request)

func WithSnapshotSequence(value uint64) RequestOption {
	return func(x *Request) {
		x.snapshotSequence = utils_value.NewOptional(value)
	}
}

func NewRequest(
	userID model_user_id.UserID,
	eventID model_event_id.EventID,
	page model_page.Page,
	opts ...RequestOption,
) Request {
	x := Request{
		userID:  userID,
		eventID: eventID,
		page:    page,
	}
	for _, option := range opts {
		option(&x)
	}

	return x
}

type Response struct {
	messages             []model_chat_message.ChatMessage
	snapshotSequence     uint64
	visibleAfterSequence uint64
	totalCount           int64
}

func (x Response) GetMessages() []model_chat_message.ChatMessage {
	return x.messages
}

func (x Response) GetSnapshotSequence() uint64 {
	return x.snapshotSequence
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func (x Response) GetVisibleAfterSequence() uint64 {
	return x.visibleAfterSequence
}

func NewResponse(
	messages []model_chat_message.ChatMessage,
	snapshotSequence uint64,
	visibleAfterSequence uint64,
	totalCount int64,
) Response {
	return Response{
		messages:             messages,
		snapshotSequence:     snapshotSequence,
		visibleAfterSequence: visibleAfterSequence,
		totalCount:           totalCount,
	}
}
