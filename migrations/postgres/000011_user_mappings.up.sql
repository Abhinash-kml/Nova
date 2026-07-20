CREATE TABLE IF NOT EXISTS user_mappings(
    user_id UUID REFERENCES users(id),
    provider_id VARCHAR UNIQUE NOT NULL,
    provider VARCHAR NOT NULL
);

CREATE INDEX idx_composite_mappings ON user_mappings(provider_id, provider);