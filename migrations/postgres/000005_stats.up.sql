CREATE TABLE IF NOT EXISTS stats(
    id SERIAL,
    key VARCHAR(64) UNIQUE,
    name VARCHAR(64) UNIQUE,
    description TEXT,
    start_value INT
    -- aggregation_type
);