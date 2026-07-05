-- +goose Up
CREATE TABLE IF NOT EXISTS links (
    id            UUID PRIMARY KEY       DEFAULT gen_random_uuid(),
    original_link TEXT        NOT NULL UNIQUE,
    short_link    VARCHAR(10) NOT NULL UNIQUE,
    created_at    TIMESTAMPTZ NOT NULL   DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS links;
