-- +goose Up
CREATE TABLE persons (
    id          SERIAL PRIMARY KEY,
    email       TEXT NOT NULL UNIQUE,
    nickname    TEXT NOT NULL UNIQUE,
    website     TEXT,
    bio         TEXT,
    created_on  TIMESTAMPTZ,
    optout      BOOLEAN,
    rank        INTEGER,
    last_login  TIMESTAMPTZ,
    avatar      TEXT
);

-- +goose Down
DROP TABLE persons;
