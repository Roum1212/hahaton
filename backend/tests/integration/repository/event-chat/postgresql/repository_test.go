//go:build integration

package repository_event_chat_postgresql_test

import (
	"context"
	"errors"

	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	aggregate_event_chat_command "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat/command"
	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_chat_message "github.com/Roum1212/hahaton/internal/domain/model/chat-message"
	model_chat_message_id "github.com/Roum1212/hahaton/internal/domain/model/chat-message-id"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	repository_fixture "github.com/Roum1212/hahaton/tests/integration/repository/fixture"
)

func (x *RepositoryTestSuite) TestRepository_CreateAndLoad() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	got, err := x.repository.Load(ctx, value.GetEventID())
	x.Require().NoError(err)
	assertEventChat(x.T(), value, got)
	x.Require().Error(x.repository.Create(ctx, value))
}

func (x *RepositoryTestSuite) TestRepository_ListMessages_OrderAndEncryption() {
	ctx := x.T().Context()
	chat, other := x.newValue(), x.newValue()
	x.Require().NoError(x.repository.Create(ctx, chat))
	x.Require().NoError(x.repository.Create(ctx, other))
	messages, err := x.repository.ListMessages(ctx, chat.GetEventID())
	x.Require().NoError(err)
	x.Require().Empty(messages)
	authorID := repository_fixture.CreateVolunteer(x.T(), x.database).GetUserID()
	first := x.newMessage(chat.GetEventID(), authorID, model_actor_kind.Coordinator, 1)
	second := x.newMessage(chat.GetEventID(), authorID, model_actor_kind.Volunteer, 2)
	foreign := x.newMessage(other.GetEventID(), authorID, model_actor_kind.Volunteer, 1)
	x.Require().NoError(x.repository.AppendMessage(ctx, second))
	x.Require().NoError(x.repository.AppendMessage(ctx, first))
	x.Require().NoError(x.repository.AppendMessage(ctx, foreign))
	messages, err = x.repository.ListMessages(ctx, chat.GetEventID())
	x.Require().NoError(err)
	x.Require().Len(messages, 2)
	assertMessage(x, first, messages[0])
	assertMessage(x, second, messages[1])

	var encrypted []byte

	err = x.pool.QueryRow(ctx, "SELECT encrypted_text FROM chat_messages WHERE id = $1", first.GetMessageID().String()).
		Scan(&encrypted)
	x.Require().NoError(err)
	x.Require().NotEmpty(encrypted)
	x.Require().NotEqual([]byte(first.GetText()), encrypted)
	x.Require().Error(x.repository.AppendMessage(ctx, first))
}

func (x *RepositoryTestSuite) TestRepository_Load_NotFound() {
	got, err := x.repository.Load(x.T().Context(), x.newValue().GetEventID())
	x.Require().ErrorIs(err, aggregate_event_chat.ErrNotFound)
	x.Require().Nil(got)
}

func (x *RepositoryTestSuite) TestRepository_Modify() {
	ctx := x.T().Context()
	value := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, value))
	err := x.repository.Modify(ctx, value.GetEventID(), func(
		ctx context.Context,
		current *aggregate_event_chat.EventChat,
	) error {
		return current.Append(aggregate_event_chat_command.NewAppend(1))
	})
	x.Require().NoError(err)
	got, err := x.repository.Load(ctx, value.GetEventID())
	x.Require().NoError(err)
	x.Require().Equal(uint64(1), got.GetLatestSequence())
}

func (x *RepositoryTestSuite) TestRepository_Modify_NotFound() {
	err := x.repository.Modify(x.T().Context(), x.newValue().GetEventID(), func(
		ctx context.Context,
		value *aggregate_event_chat.EventChat,
	) error {
		return nil
	})
	x.Require().ErrorIs(err, aggregate_event_chat.ErrNotFound)
}

func (x *RepositoryTestSuite) TestRepository_Modify_RollbackMessage() {
	ctx := x.T().Context()
	chat := x.newValue()
	x.Require().NoError(x.repository.Create(ctx, chat))
	author := repository_fixture.CreateVolunteer(x.T(), x.database).GetUserID()
	message := x.newMessage(chat.GetEventID(), author, model_actor_kind.Volunteer, 1)
	stop := errors.New("cancel message")
	err := x.repository.Modify(ctx, chat.GetEventID(), func(
		ctx context.Context,
		current *aggregate_event_chat.EventChat,
	) error {
		if err := current.Append(aggregate_event_chat_command.NewAppend(1)); err != nil {
			return err
		}

		if err := x.repository.AppendMessage(ctx, message); err != nil {
			return err
		}

		return stop
	})
	x.Require().ErrorIs(err, stop)
	got, err := x.repository.Load(ctx, chat.GetEventID())
	x.Require().NoError(err)
	assertEventChat(x.T(), chat, got)
	messages, err := x.repository.ListMessages(ctx, chat.GetEventID())
	x.Require().NoError(err)
	x.Require().Empty(messages)
}

func (x *RepositoryTestSuite) newMessage(
	eventID model_event_id.EventID,
	authorID model_user_id.UserID,
	kind model_actor_kind.ActorKind,
	sequence uint64,
) model_chat_message.ChatMessage {
	id, err := model_chat_message_id.GenerateChatMessageID()
	x.Require().NoError(err)
	message, err := model_chat_message.NewChatMessage(
		id, eventID, authorID, kind, "Иван Иванов", "Привет, волонтёры!", id.String(), sequence, repository_fixture.Now(),
	)
	x.Require().NoError(err)

	return message
}

func assertMessage(x *RepositoryTestSuite, expected, actual model_chat_message.ChatMessage) {
	x.T().Helper()

	x.Require().Equal(expected.GetMessageID(), actual.GetMessageID())
	x.Require().Equal(expected.GetEventID(), actual.GetEventID())
	x.Require().Equal(expected.GetAuthorID(), actual.GetAuthorID())
	x.Require().Equal(expected.GetAuthorKind(), actual.GetAuthorKind())
	x.Require().Equal(expected.GetAuthorDisplayName(), actual.GetAuthorDisplayName())
	x.Require().Equal(expected.GetText(), actual.GetText())
	x.Require().Equal(expected.GetRequestID(), actual.GetRequestID())
	x.Require().Equal(expected.GetSequence(), actual.GetSequence())
	x.Require().WithinDuration(expected.GetCreatedAt(), actual.GetCreatedAt(), 0)
}
