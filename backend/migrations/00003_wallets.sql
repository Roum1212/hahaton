-- +goose Up
CREATE TABLE wallets
(
    volunteer_id      uuid PRIMARY KEY REFERENCES volunteers (id),
    coin_balance      bigint      NOT NULL,
    confirmed_minutes bigint      NOT NULL,
    last_sequence     bigint      NOT NULL,
    created_at        timestamptz NOT NULL,
    updated_at        timestamptz NOT NULL
);

-- +goose Down
DROP TABLE wallets;
