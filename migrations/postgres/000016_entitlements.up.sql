CREATE TABLE IF NOT EXISTS entitlements(
    entitlement_id UUID DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    item_id UUID NOT NULL,
    source TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    obtained_on TIMESTAMP NOT NULL DEAFULT CURRENT_TIMESTAMP,
    valid_till TIMESTAMP,

    PRIMARY KEY(user_id, item_id),

    FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE CASCADE,

    FOREIGN KEY(item_id) REFERENCES player_items(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_entitlements_id ON entitlements(id);