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

INSERT INTO player_items(id, name, description) VALUES ('a5a7c125-cf90-4d1f-8a32-46883ebd506a', 'Fire Weapon', 'Very Fire Weapon'),
('b5a76e53-cfa1-4a9f-9749-4436280f5519', 'Frost Weapon', 'Very Cold Weapon'),
('1d9fe0fb-c42c-4686-b224-52653727150f', 'HealthPack', 'This will staisfy you')