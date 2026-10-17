-- +goose Up
CREATE TABLE volunteers
(
    id            uuid PRIMARY KEY,
    flida_issuer  text        NOT NULL,
    flida_subject text        NOT NULL,
    first_name    text        NOT NULL,
    last_name     text        NOT NULL,
    middle_name   text        NOT NULL,
    date_of_birth date        NOT NULL,
    email_address text        NOT NULL,
    phone_number  text        NOT NULL,
    avatar_path   text        NOT NULL,
    created_at    timestamptz NOT NULL,
    updated_at    timestamptz NOT NULL,
    UNIQUE (flida_issuer, flida_subject)
);

CREATE INDEX volunteers_email_address_idx ON volunteers (email_address);
CREATE INDEX volunteers_phone_number_idx ON volunteers (phone_number);

-- +goose Down
DROP TABLE volunteers;
