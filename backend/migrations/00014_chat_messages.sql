-- +goose Up
CREATE TABLE chat_messages
(
    id                  uuid PRIMARY KEY,
    event_id            uuid        NOT NULL REFERENCES event_chats (event_id),
    author_id           uuid        NOT NULL REFERENCES volunteers (id),
    author_kind         text        NOT NULL,
    author_display_name text        NOT NULL,
    encrypted_text      bytea       NOT NULL,
    request_id          text        NOT NULL,
    sequence            bigint      NOT NULL,
    created_at          timestamptz NOT NULL,
    UNIQUE (event_id, sequence),
    UNIQUE (event_id, author_kind, author_id, request_id)
);
-- +goose Down
DROP TABLE chat_messages;
