CREATE TABLE IF NOT EXISTS achievements(
    id SERIAL UNIQUE,
    key VARCHAR(64) UNIQUE,
    name VARCHAR(64),
    description TEXT
);