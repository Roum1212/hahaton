package repository_event_chat_postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	aggregate_event_chat "github.com/Roum1212/hahaton/internal/domain/aggregate/event-chat"
	model_actor_kind "github.com/Roum1212/hahaton/internal/domain/model/actor-kind"
	model_chat_message "github.com/Roum1212/hahaton/internal/domain/model/chat-message"
	model_chat_message_id "github.com/Roum1212/hahaton/internal/domain/model/chat-message-id"
	model_event_id "github.com/Roum1212/hahaton/internal/domain/model/event-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	repository_locking "github.com/Roum1212/hahaton/internal/infra/repository/locking"
)

const (
	idColumn      = "id"
	eventIDColumn = "event_id"

	authorDisplayNameColumn = "author_display_name"
	authorIDColumn          = "author_id"
	authorKindColumn        = "author_kind"
	encryptedTextColumn     = "encrypted_text"
	latestSequenceColumn    = "latest_sequence"
	requestIDColumn         = "request_id"
	sequenceColumn          = "sequence"

	createdAtColumn = "created_at"
)

type Crypto interface {
	Decrypt(value []byte) (string, error)
	Encrypt(text string) ([]byte, error)
}
type Database interface {
	Exec(ctx context.Context, query string, arguments ...any) (pgconn.CommandTag, error)
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
	Query(ctx context.Context, query string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, arguments ...any) pgx.Row
}
type EventChats struct {
	database Database
	crypto   Crypto
	queries  squirrel.StatementBuilderType
}

func (x *EventChats) AppendMessage(ctx context.Context, message model_chat_message.ChatMessage) error {
	encrypted, err := x.crypto.Encrypt(message.GetText())
	if err != nil {
		return err
	}

	sql, arguments, err := x.queries.Insert("chat_messages").SetMap(map[string]any{
		idColumn:                message.GetMessageID().String(),
		eventIDColumn:           message.GetEventID().String(),
		authorIDColumn:          message.GetAuthorID().String(),
		authorKindColumn:        message.GetAuthorKind().String(),
		authorDisplayNameColumn: message.GetAuthorDisplayName(),
		encryptedTextColumn:     encrypted,
		requestIDColumn:         message.GetRequestID(),
		sequenceColumn:          message.GetSequence(),
		createdAtColumn:         message.GetCreatedAt(),
	}).ToSql()
	if err != nil {
		return err
	}

	_, err = x.database.Exec(ctx, sql, arguments...)

	return err
}

func (x *EventChats) Create(ctx context.Context, chat *aggregate_event_chat.EventChat) error {
	sql, arguments, err := x.queries.Insert("event_chats").SetMap(map[string]any{
		eventIDColumn:        chat.GetEventID().String(),
		latestSequenceColumn: chat.GetLatestSequence(),
	}).ToSql()
	if err != nil {
		return err
	}

	_, err = x.database.Exec(ctx, sql, arguments...)

	return err
}

func (x *EventChats) ListMessages(
	ctx context.Context,
	eventID model_event_id.EventID,
) ([]model_chat_message.ChatMessage, error) {
	sql, arguments, err := x.queries.Select(
		idColumn,
		authorIDColumn,
		authorKindColumn,
		authorDisplayNameColumn,
		encryptedTextColumn,
		requestIDColumn,
		sequenceColumn,
		createdAtColumn,
	).From("chat_messages").Where(squirrel.Eq{
		eventIDColumn: eventID.String(),
	}).OrderBy(sequenceColumn + " ASC").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := x.database.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]model_chat_message.ChatMessage, 0)

	for rows.Next() {
		var (
			id, authorID, kind, displayName, requestID string
			encrypted                                  []byte
			sequence                                   uint64
			createdAt                                  time.Time
		)

		if err := rows.Scan(
			&id,
			&authorID,
			&kind,
			&displayName,
			&encrypted,
			&requestID,
			&sequence,
			&createdAt,
		); err != nil {
			return nil, err
		}

		messageID, err := model_chat_message_id.NewChatMessageID(id)
		if err != nil {
			return nil, fmt.Errorf("restore message id: %w", err)
		}

		userID, err := model_user_id.NewUserID(authorID)
		if err != nil {
			return nil, fmt.Errorf("restore user id: %w", err)
		}

		text, err := x.crypto.Decrypt(encrypted)
		if err != nil {
			return nil, err
		}

		authorKind, err := model_actor_kind.NewActorKind(kind)
		if err != nil {
			return nil, fmt.Errorf("restore author kind: %w", err)
		}

		message, err := model_chat_message.NewChatMessage(
			messageID,
			eventID,
			userID,
			authorKind,
			displayName,
			text,
			requestID,
			sequence,
			createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("restore message: %w", err)
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (x *EventChats) Load(
	ctx context.Context,
	eventID model_event_id.EventID,
) (*aggregate_event_chat.EventChat, error) {
	return x.load(ctx, eventID, false)
}

func (x *EventChats) Modify(
	ctx context.Context,
	eventID model_event_id.EventID,
	change func(
		ctx context.Context,
		chat *aggregate_event_chat.EventChat,
	) error,
) error {
	return x.database.InTx(ctx, func(ctx context.Context) error {
		if _, err := x.Load(ctx, eventID); err != nil {
			return err
		}

		if err := repository_locking.Event(ctx, x.database, eventID); err != nil {
			return err
		}

		chat, err := x.load(ctx, eventID, true)
		if err != nil {
			return err
		}

		if err := change(ctx, chat); err != nil {
			return err
		}

		sql, arguments, err := x.queries.Update("event_chats").Set(
			latestSequenceColumn,
			chat.GetLatestSequence(),
		).Where(
			squirrel.Eq{eventIDColumn: eventID.String()},
		).ToSql()
		if err != nil {
			return err
		}

		result, err := x.database.Exec(ctx, sql, arguments...)
		if err != nil {
			return fmt.Errorf("save event chat: %w", err)
		}

		if result.RowsAffected() == 0 {
			return aggregate_event_chat.ErrNotFound
		}

		return nil
	})
}

func (x *EventChats) load(
	ctx context.Context,
	eventID model_event_id.EventID,
	lock bool,
) (*aggregate_event_chat.EventChat, error) {
	query := x.queries.Select(latestSequenceColumn).From("event_chats").Where(squirrel.Eq{
		eventIDColumn: eventID.String(),
	})
	if lock {
		query = query.Suffix("FOR UPDATE")
	}

	sql, arguments, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var sequence uint64
	if err := x.database.QueryRow(ctx, sql, arguments...).Scan(&sequence); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, aggregate_event_chat.ErrNotFound
		default:
			return nil, err
		}
	}

	return aggregate_event_chat.NewEventChatFromFields(eventID, sequence), nil
}

func NewEventChats(database Database, crypto Crypto) *EventChats {
	return &EventChats{
		database: database,
		crypto:   crypto,
		queries:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
