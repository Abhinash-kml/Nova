CREATE TYPE relationship_status AS ENUM('pending', 'accepted', 'blocked');

CREATE TABLE IF NOT EXISTS friendships(
    user_one_id UUID NOT NULL,
    user_two_id UUID NOT NULL,
    status relationship_status DEFAULT 'pending',
    action_user_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(user_one_id, user_two_id),
    FOREIGN KEY(user_one_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(user_two_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(action_user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT chk_ordered_users CHECK (user_one_id < user_two_id),
    CONSTRAINT chk_no_self_friend CHECK (user_one_id <> user_two_id)
);

CREATE INDEX idx_friendship_user_two ON friendships(user_two_id, status);
CREATE INDEX idx_friendship_user_one ON friendships(user_one_id, status);


CREATE TABLE IF NOT EXISTS conversations(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    type VARCHAR NOT NULL,
    name VARCHAR UNIQUE,
    last_message_id BIGINT,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sender_id UUID NOT NULL,
    conversation_id UUID NOT NULL,
    body VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    FOREIGN KEY(sender_id) REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY(conversation_id) REFERENCES conversations(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_messages_user_conversation ON messages(conversation_id, sender_id);

CREATE TABLE IF NOT EXISTS conversation_participants(
    conversation_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_read_message_id BIGINT,

    FOREIGN KEY(conversation_id) REFERENCES conversations(id)
    ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY(last_read_message_id) REFERENCES messages(id)
    ON DELETE CASCADE
);