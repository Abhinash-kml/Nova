CREATE TABLE IF NOT EXISTS achievements(
    id SERIAL UNIQUE NOT NULL,
    key VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(64) NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS achievement_criteria(
    achievement_id INT REFERENCES achievements(id) NOT NULL,
    stat_id INT REFERENCES stats(id) NOT NULL,
    target_value INT CHECK (target_value < 0)
);

CREATE TABLE IF NOT EXISTS completed_achievements(
    user_id UUID REFERENCES users(id) NOT NULL,
    achievement_id REFERENCES achievements(id) NOT NULL,
    completed_on TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_completed_achievements_userid ON completed_achievements(id);