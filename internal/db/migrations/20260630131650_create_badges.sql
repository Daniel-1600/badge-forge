-- +goose Up
CREATE TABLE badges (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    image       TEXT NOT NULL,
    description TEXT NOT NULL,
    criteria    TEXT NOT NULL,
    issuer_id   INTEGER,
    created_on  TIMESTAMPTZ,
    tags        TEXT,
    stl         TEXT,
    CONSTRAINT fk_badges_issuer FOREIGN KEY (issuer_id) REFERENCES persons(id)
);

-- +goose Down
DROP TABLE badges;