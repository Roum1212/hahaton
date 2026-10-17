-- +goose Up
CREATE TABLE participations
(
    id                     uuid        NOT NULL,
    event_id               uuid        NOT NULL,
    coordinator_id         uuid        NOT NULL,
    volunteer_id           uuid        NOT NULL,
    status                 text        NOT NULL,
    cancellation_reason    text        NOT NULL,
    rejection_reason       text        NOT NULL,
    awarded_minutes        bigint      NOT NULL,
    awarded_coins          bigint      NOT NULL,
    accrual_operation_id   text        NOT NULL,
    has_chat_access        boolean     NOT NULL,
    visible_after_sequence bigint      NOT NULL,
    created_at             timestamptz NOT NULL,
    updated_at             timestamptz NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (event_id, volunteer_id),
    FOREIGN KEY (event_id) REFERENCES events (id),
    FOREIGN KEY (volunteer_id) REFERENCES volunteers (id)
);
-- +goose Down
DROP TABLE participations;
