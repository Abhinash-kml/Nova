CREATE TABLE IF NOT EXISTS player_inventory(
    userid UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    itemid UUID NOT NULL REFERENCES player_items(id) ON DELETE CASCADE,
    quantity INT DEFAULT 1,

    -- composite primary key to ensure a player has only one row per unique item 
    PRIMARY KEY(userid, itemid),

    -- constraint to prevent negative item count from dev mistake
    CONSTRAINT chk_player_inventory_quantity_positive CHECK(quantity >= 0)
);

CREATE INDEX idx_player_inventory_player_id ON player_inventory(userid);