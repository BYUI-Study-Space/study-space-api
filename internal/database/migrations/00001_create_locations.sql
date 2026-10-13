-- +goose Up
CREATE TABLE locations (
    id             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name           TEXT        NOT NULL,
    building       TEXT        NOT NULL,
    room           TEXT        NOT NULL DEFAULT '',
    description    TEXT        NOT NULL DEFAULT '',
    image_url      TEXT        NOT NULL DEFAULT '',
    has_wifi       BOOLEAN     NOT NULL DEFAULT TRUE,
    group_friendly BOOLEAN     NOT NULL DEFAULT FALSE,
    noise_level    TEXT        NOT NULL CHECK (noise_level IN ('quiet', 'moderate', 'collaborative')),
    capacity       INTEGER     NOT NULL CHECK (capacity > 0),
    open_time      TIME        NOT NULL,
    close_time     TIME        NOT NULL CHECK (close_time > open_time),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE locations;
