CREATE TABLE IF NOT EXISTS achievements(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    key VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    stat_id INT NOT NULL,
    target_value BIGINT NOT NULL CHECK (target_value >= 1),

    FOREIGN KEY(stat_id) REFERENCES stats(id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS achievement_progress(
    achievement_id INT NOT NULL,
    user_id UUID NOT NULL,
    current_stat_value INT NOT NULL CHECK (current_stat_value >= 0),
    updated_at TIMESTAMPTZ,

    PRIMARY KEY(user_id, achievement_id),
    FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS achievement_rewards(
    achievement_id INT PRIMARY KEY,
    reward_id VARCHAR NOT NULL,
    count INT CHECK(count >= 1),

    FOREIGN KEY(achievement_id) REFERENCES achievements(id)
    ON DELETE CASCADE,
    FOREIGN KEY(reward_id) REFERENCES rewards(id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS completed_achievements(
    user_id UUID,
    achievement_id INT NOT NULL,
    completed_on TIMESTAMPTZ NOT NULL,

    PRIMARY KEY(user_id, achievement_id, completed_on),
    FOREIGN KEY(achievement_id) REFERENCES achievements(id)
    ON DELETE CASCADE
);

