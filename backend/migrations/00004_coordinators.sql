-- +goose Up
CREATE TABLE coordinators
(
    id               uuid        NOT NULL,
    description      text        NOT NULL,
    avatar_path      text        NOT NULL,
    cooldown_seconds bigint      NOT NULL,
    created_at       timestamptz NOT NULL,
    updated_at       timestamptz NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (id) REFERENCES volunteers (id)
);
-- +goose Down
DROP TABLE coordinators;
