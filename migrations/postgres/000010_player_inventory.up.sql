CREATE TABLE IF NOT EXISTS player_inventory(
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES player_items(id) ON DELETE CASCADE,
    quantity INT DEFAULT 1,
    source VARCHAR(20),
    is_equipped BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP NOT NULL,

    -- composite primary key to ensure a player has only one row per unique item 
    PRIMARY KEY(user_id, item_id),

    -- constraint to prevent negative item count from dev mistake
    CONSTRAINT chk_player_inventory_quantity_positive CHECK(quantity >= 0)
);

CREATE INDEX idx_player_inventory_player_id ON player_inventory(user_id);