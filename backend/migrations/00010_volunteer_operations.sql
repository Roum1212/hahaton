-- +goose Up
CREATE TABLE volunteer_operations
(
    id                      uuid        NOT NULL,
    coordinator_id          uuid        NOT NULL,
    volunteer_id            uuid        NOT NULL,
    description             text        NOT NULL,
    operation_type          integer     NOT NULL,
    source_id               text        NOT NULL,
    coin_delta              bigint      NOT NULL,
    minute_delta            bigint      NOT NULL,
    coin_balance_after      bigint      NOT NULL,
    confirmed_minutes_after bigint      NOT NULL,
    sequence                bigint      NOT NULL,
    created_at              timestamptz NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (operation_type, source_id),
    UNIQUE (volunteer_id, sequence),
    FOREIGN KEY (volunteer_id) REFERENCES wallets (volunteer_id)
);
-- +goose Down
DROP TABLE volunteer_operations;
