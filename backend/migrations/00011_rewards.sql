-- +goose Up
CREATE TABLE rewards
(
    id                uuid        NOT NULL,
    coordinator_id    uuid        NOT NULL,
    title             text        NOT NULL,
    description       text        NOT NULL,
    image_path        text        NOT NULL,
    price_coins       bigint      NOT NULL,
    quantity_limit    integer     NOT NULL,
    reserved_quantity integer     NOT NULL,
    issued_quantity   integer     NOT NULL,
    deleted           boolean     NOT NULL,
    created_at        timestamptz NOT NULL,
    updated_at        timestamptz NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (coordinator_id) REFERENCES coordinators (id)
);
-- +goose Down
DROP TABLE rewards;
