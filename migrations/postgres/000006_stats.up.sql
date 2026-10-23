CREATE TABLE IF NOT EXISTS stats(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    key VARCHAR(64) UNIQUE,
    name VARCHAR(64) UNIQUE,
    description TEXT,
    start_value INT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS player_stats_progression(
    player_id UUID REFERENCES users(id),
    stat_id INT REFERENCES stats(id),
    current_value BIGINT,
    updated_at TIMESTAMPTZ,

    -- Composite primary key for uniqueness accross both columns
    PRIMARY KEY(player_id, stat_id),

    FOREIGN KEY(player_id) REFERENCES users(id)
    ON DELETE CASCADE,

    FOREIGN KEY(stat_id) REFERENCES stats(id)
    ON DELETE CASCADE
);