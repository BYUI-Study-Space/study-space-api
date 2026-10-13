-- +goose Up
-- Session ids are UUIDs so confirmation links can't be guessed by counting.
CREATE TABLE sessions (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id      BIGINT      NOT NULL REFERENCES locations (id),
    student_name     TEXT        NOT NULL,
    email            TEXT        NOT NULL,
    student_id       TEXT        NOT NULL,
    session_date     DATE        NOT NULL,
    start_time       TIME        NOT NULL,
    duration_minutes INTEGER     NOT NULL CHECK (duration_minutes IN (30, 60, 90, 120, 180)),
    session_type     TEXT        NOT NULL CHECK (session_type IN ('individual', 'group')),
    group_size       INTEGER     CHECK (group_size >= 2),
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK ((session_type = 'group') = (group_size IS NOT NULL))
);

CREATE INDEX sessions_location_date_idx ON sessions (location_id, session_date);

-- +goose Down
DROP TABLE sessions;
