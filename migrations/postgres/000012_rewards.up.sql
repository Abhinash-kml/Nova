CREATE TABLE IF NOT EXISTS rewards(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    key VARCHAR(32) UNIQUE NOT NULL,
    description VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS reward_items(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    reward_id INT REFERENCES rewards(id) ON DELETE CASCADE NOT NULL,
    item_id UUID REFERENCES player_items(id) ON DELETE CASCADE NOT NULL,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity >= 1)
);

CREATE INDEX idx_reward_items_rewardid ON reward_items(reward_id);