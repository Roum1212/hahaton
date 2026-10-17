-- +goose Up
CREATE TABLE cities
(
    id              TEXT PRIMARY KEY,
    name            TEXT  NOT NULL,
    region_name     TEXT  NOT NULL,
    latitude        DOUBLE PRECISION,
    longitude       DOUBLE PRECISION,
    polygons        JSONB NOT NULL,
    catalog_version TEXT  NOT NULL
);

-- +goose Down
DROP TABLE cities;
