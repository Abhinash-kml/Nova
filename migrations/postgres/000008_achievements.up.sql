CREATE TABLE IF NOT EXISTS achievements(
    id SERIAL UNIQUE NOT NULL,
    key VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(64) NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS achievement_criteria(
    id SERIAL PRIMARY KEY,
    achievement_id INT REFERENCES achievements(id) ON DELETE CASCADE NOT NULL,
    stat_id INT REFERENCES stats(id) ON DELETE CASCADE NOT NULL,
    target_value INT CHECK (target_value >= 1),
    criteria_text TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS achievement_progress(
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) NOT NULL,
    criterion_id INT REFERENCES achievement_criteria(id) ON DELETE CASCADE NOT NULL,
    current_value INT NOT NULL CHECK (current_value >= 0),
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS achievement_rewards(
    id SERIAL PRIMARY KEY,
    achievement_id INT REFERENCES achievements(id) ON DELETE CASCADE NOT NULL,
    reward_id INT REFERENCES rewards(id),
    count INT CHECK(count >= 1)
);

CREATE TABLE IF NOT EXISTS completed_achievements(
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    achievement_id INT REFERENCES achievements(id) ON DELETE CASCADE NOT NULL,
    completed_on TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_completed_achievements_userid ON completed_achievements(user_id);

