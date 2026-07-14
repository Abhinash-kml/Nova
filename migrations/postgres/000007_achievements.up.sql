CREATE TABLE IF NOT EXISTS achievements(
    id SERIAL,
    key VARCHAR(64),
    name VARCHAR(64),
    description TEXT
);