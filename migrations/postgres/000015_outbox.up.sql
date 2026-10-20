CREATE TABLE IF NOT EXISTS outbox(
    id UUID DEFAULT uuidv7(),
    aggregate_id VARCHAR NOT NULL,
    aggregate_type VARCHAR NOT NULL,
    type VARCHAR NOT NULL,
    payload JSONB DEFAULT '{}'::JSONB,
    processed_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_outbox ON outbox(id);