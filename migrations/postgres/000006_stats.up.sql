CREATE TABLE IF NOT EXISTS stats(
    id SERIAL PRIMARY KEY UNIQUE,
    key VARCHAR(64) UNIQUE,
    name VARCHAR(64) UNIQUE,
    description TEXT,
    start_value INT,
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS player_stats_progression(
    player_id UUID REFERENCES users(id),
    stat_id INT REFERENCES stats(id),
    current_value INT,
    updated_at TIMESTAMP,

    -- Composite primary key for uniqueness accross both columns
    PRIMARY KEY(player_id, stat_id)
);

CREATE INDEX idx_player_stats_progression_player_id ON player_stats_progression(player_id);