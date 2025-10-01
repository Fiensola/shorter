CREATE TABLE enriched_clicks (
    id BIGSERIAL PRIMARY KEY,
    alias VARCHAR(100) NOT NULL,
    ip VARCHAR(45),
    country VARCHAR(100),
    city VARCHAR(100),
    device_type VARCHAR(20),
    os VARCHAR(50),
    browser VARCHAR(50),
    referer VARCHAR(500),
    timestamp TIMESTAMPTZ DEFAULT NOW()
);