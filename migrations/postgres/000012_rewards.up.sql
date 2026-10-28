CREATE TABLE IF NOT EXISTS rewards(
    id VARCHAR PRIMARY KEY,
    name VARCHAR(100),
    description VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS reward_items(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    reward_id VARCHAR NOT NULL,
    item_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity >= 1),

    FOREIGN KEY(reward_id) REFERENCES rewards(id)
    ON DELETE CASCADE,
    FOREIGN KEY(item_id) REFERENCES player_items(id)
    ON DELETE CASCADE
);