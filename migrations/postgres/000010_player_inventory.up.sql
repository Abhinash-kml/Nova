CREATE TABLE IF NOT EXISTS player_inventory(
    user_id UUID NOT NULL,
    item_id UUID NOT NULL,
    quantity INT DEFAULT 1,
    is_equipped BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP NOT NULL,

    -- composite primary key to ensure a player has only one row per unique item 
    PRIMARY KEY(user_id, item_id),
    FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY(item_id) REFERENCES player_items(id)
    ON DELETE CASCADE,

    -- constraint to prevent negative item count from dev mistake
    CONSTRAINT chk_player_inventory_quantity_positive CHECK(quantity >= 0)
);