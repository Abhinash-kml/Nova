CREATE TABLE IF NOT EXISTS leaderboards(
    id UUID PRIMARY KEY UNIQUE,
    name VARCHAR(64) UNIQUE,
    type VARCHAR(64),
    stat INT REFERENCES stats(id),
    process_interval INT,
    created_by UUID,
    created_at TIMESTAMP
);