-- +goose Up
CREATE TABLE participation_changes
(
    id                  uuid        NOT NULL,
    participation_id    uuid        NOT NULL,
    previous_status     text        NOT NULL,
    status              text        NOT NULL,
    actor_id            text        NOT NULL,
    actor_kind          text        NOT NULL,
    cancellation_reason text        NOT NULL,
    rejection_reason    text        NOT NULL,
    sequence            bigserial   NOT NULL,
    created_at          timestamptz NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (participation_id) REFERENCES participations (id)
);
-- +goose Down
DROP TABLE participation_changes;
