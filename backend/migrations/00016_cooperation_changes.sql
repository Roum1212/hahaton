-- +goose Up
CREATE TABLE cooperation_changes
(
    sequence       BIGSERIAL PRIMARY KEY,
    coordinator_id UUID        NOT NULL,
    volunteer_id   UUID        NOT NULL,
    status         TEXT        NOT NULL,
    actor_id       UUID        NOT NULL REFERENCES volunteers (id),
    retry_after    TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (coordinator_id, volunteer_id) REFERENCES cooperations (coordinator_id, volunteer_id)
);

CREATE INDEX cooperation_changes_pair_idx ON cooperation_changes (coordinator_id, volunteer_id, sequence);

-- +goose Down
DROP TABLE cooperation_changes;
