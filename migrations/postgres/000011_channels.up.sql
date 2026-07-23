CREATE TABLE IF NOT EXISTS channels(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(64) UNIQUE,
    persistant BOOLEAN DEFAULT false,
    process_interval INT,
    created_by UUID,
    created_at TIMESTAMP,
    updated_by UUID,
    updated_at TIMESTAMP
);