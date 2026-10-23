CREATE TABLE IF NOT EXISTS leaderboards(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(64) UNIQUE,
    type VARCHAR(64),
    stat_id INT NOT NULL,
    process_interval INT,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY(stat_id) REFERENCES stats(id),
    ON DELETE RESTRICT
);