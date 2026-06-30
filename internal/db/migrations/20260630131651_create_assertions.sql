--- +goose Up
CREATE TABLE assertions (
    id          TEXT PRIMARY KEY,
    badge_id    TEXT,
    person_id   INTEGER,
    salt        TEXT,
    issued_on   TIMESTAMPTZ,
    recipient   TEXT,
    issued_for  TEXT,
    CONSTRAINT fk_assertions_badge FOREIGN KEY (badge_id) REFERENCES badges(id),
    CONSTRAINT fk_assertions_person FOREIGN KEY (person_id) REFERENCES persons(id)
);

-- +goose Down
DROP TABLE assertions;
