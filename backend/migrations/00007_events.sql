-- +goose Up
CREATE TABLE events
(
    id                    uuid             NOT NULL,
    coordinator_id        uuid             NOT NULL,
    title                 text             NOT NULL,
    description           text             NOT NULL,
    image_path            text             NOT NULL,
    location              text             NOT NULL,
    latitude              double precision NOT NULL,
    longitude             double precision NOT NULL,
    starts_at             timestamptz      NOT NULL,
    ends_at               timestamptz      NOT NULL,
    requirements          integer[] NOT NULL,
    volunteer_limit       integer          NOT NULL,
    status                text             NOT NULL,
    completion_request_id text,
    result_fingerprint    text             NOT NULL,
    created_at            timestamptz      NOT NULL,
    updated_at            timestamptz      NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (coordinator_id, completion_request_id),
    FOREIGN KEY (coordinator_id) REFERENCES coordinators (id)
);
-- +goose Down
DROP TABLE events;
