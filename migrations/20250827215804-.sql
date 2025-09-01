-- +migrate Up
CREATE TABLE vessels
(
    id         VARCHAR(50)      PRIMARY KEY,
    name       VARCHAR(255)     NOT NULL,
    capacity   INTEGER          NOT NULL CHECK (capacity > 0),
    latitude   DOUBLE PRECISION NOT NULL CHECK (latitude >= -90 AND latitude <= 90),
    longitude  DOUBLE PRECISION NOT NULL CHECK (longitude >= -180 AND longitude <= 180),
    created_at TIMESTAMPTZ      NOT NULL,
    updated_at TIMESTAMPTZ      NOT NULL,
    deleted_at TIMESTAMPTZ
);
-- +migrate Down
DROP TABLE IF EXISTS vessels;
