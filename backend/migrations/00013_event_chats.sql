-- +goose Up
CREATE TABLE event_chats
(
    event_id        uuid PRIMARY KEY REFERENCES events (id),
    latest_sequence bigint NOT NULL
);
-- +goose Down
DROP TABLE event_chats;
