-- +goose Up
CREATE TABLE IF NOT EXISTS links (
    id  UUID PRIMARY KEY               DEFAULT gen_random_uuid(),
    original_link TEXT UNIQUE NOT NULL,
    short_link VARCHAR(10) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)

CREATE INDEX idx_short_links ON links (short_link)

-- +goose Down
DELETE TABLE IF EXISTS links
