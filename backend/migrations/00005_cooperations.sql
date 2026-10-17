-- +goose Up
CREATE TABLE cooperations
(
    coordinator_id uuid        NOT NULL,
    volunteer_id   uuid        NOT NULL,
    status         text        NOT NULL,
    retry_after    timestamptz NOT NULL,
    created_at     timestamptz NOT NULL,
    updated_at     timestamptz NOT NULL,
    PRIMARY KEY (coordinator_id, volunteer_id),
    UNIQUE (coordinator_id, volunteer_id),
    FOREIGN KEY (coordinator_id) REFERENCES coordinators (id),
    FOREIGN KEY (volunteer_id) REFERENCES volunteers (id)
);
-- +goose Down
DROP TABLE cooperations;
