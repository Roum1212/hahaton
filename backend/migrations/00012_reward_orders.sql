-- +goose Up
CREATE TABLE reward_orders
(
    id                  uuid        NOT NULL,
    reward_id           uuid        NOT NULL,
    coordinator_id      uuid        NOT NULL,
    volunteer_id        uuid        NOT NULL,
    reward_title        text        NOT NULL,
    price_coins         bigint      NOT NULL,
    request_id          text        NOT NULL,
    status              text        NOT NULL,
    cancellation_reason text        NOT NULL,
    created_at          timestamptz NOT NULL,
    updated_at          timestamptz NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (volunteer_id, request_id),
    FOREIGN KEY (reward_id) REFERENCES rewards (id),
    FOREIGN KEY (volunteer_id) REFERENCES volunteers (id)
);
-- +goose Down
DROP TABLE reward_orders;
