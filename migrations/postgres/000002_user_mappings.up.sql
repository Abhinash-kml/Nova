CREATE TABLE IF NOT EXISTS user_mappings(
    user_id UUID,
    provider_id VARCHAR UNIQUE NOT NULL,
    provider VARCHAR NOT NULL,

    CONSTRAINT fkey_user_mappings FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_composite_mappings ON user_mappings(provider_id, provider);