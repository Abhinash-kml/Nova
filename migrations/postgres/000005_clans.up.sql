CREATE TABLE IF NOT EXISTS clans (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR NOT NULL,
    tag VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    leader_id UUID UNIQUE NOT NULL,
    level INT DEFAULT 1,
    max_members INT, 
    is_locked BOOL DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    
    FOREIGN KEY(leader_id) REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE TYPE CLAN_MEMBER_ROLE AS ENUM (
    'LEADER',
    'CO-LEADER',
    'ELITE',
    'MEMBER'
);

CREATE TABLE IF NOT EXISTS clan_members (
    clan_id UUID,
    user_id UUID,
    role CLAN_MEMBER_ROLE DEFAULT 'MEMBER',
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(clan_id, user_id, role)
);

CREATE TABLE IF NOT EXISTS clan_operations (
    clan_id UUID,
    initiator_id UUID,
    target_id UUID,

    PRIMARY KEY(clan_id, initiator_id, target_id),
    FOREIGN KEY(clan_id) REFERENCES clans(id)
    ON DELETE CASCADE,
    FOREIGN KEY(initiator_id) REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY(target_id) REFERENCES users(id)
    ON DELETE CASCADE
);




