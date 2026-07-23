CREATE TABLE IF NOT EXISTS player_items(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    rarity INT NOT NULL DEFAULT 1,
    is_stackable BOOLEAN NOT NULL DEFAULT TRUE,
    meta JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ
);

CREATE INDEX idx_player_items_meta ON player_items USING gin(meta);