CREATE TABLE IF NOT EXISTS stats(
    id SERIAL PRIMARY KEY UNIQUE,
    key VARCHAR(64) UNIQUE,
    name VARCHAR(64) UNIQUE,
    description TEXT,
    start_value INT,
    created_at TIMESTAMP
    -- aggregation_type
);