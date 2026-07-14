CREATE TABLE IF NOT EXISTS player_progress(
    player_id UUID REFERENCES users(id),
    stat_id INT REFERENCES stats(id),
    current_value INT,
    updated_at TIMESTAMP
);