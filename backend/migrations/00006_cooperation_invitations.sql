-- +goose Up
CREATE TABLE cooperation_invitations
(
    id             uuid        NOT NULL,
    coordinator_id uuid        NOT NULL,
    token          text        NOT NULL,
    created_at     timestamptz NOT NULL,
    expires_at     timestamptz NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (token),
    FOREIGN KEY (coordinator_id) REFERENCES coordinators (id)
);
-- +goose Down
DROP TABLE cooperation_invitations;
