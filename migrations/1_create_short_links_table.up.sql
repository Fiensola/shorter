CREATE TABLE short_links (
    id BIGSERIAL PRIMARY KEY,
    alias VARCHAR(100) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    click_count INT DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_reviews_alias ON short_links(alias);