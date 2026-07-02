CREATE TABLE IF NOT EXISTS clans (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR NOT NULL,
    tag VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    leader_id UUID NOT NULL REFERENCES users(id),
    coleader_ids UUID[] DEFAULT '{}',
    level INT DEFAULT 1,
    member_ids UUID[] DEFAULT '{}',
    max_members INT, 
    islocked BOOL DEFAULT false,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP
);