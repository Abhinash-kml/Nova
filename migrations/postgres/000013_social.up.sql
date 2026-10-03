CREATE TYPE relationship_status AS ENUM('pending', 'accepted', 'blocked');

CREATE TABLE IF NOT EXISTS friendships(
    user_one_id UUID NOT NULL,
    user_two_id UUID NOT NULL,
    status relationship_status DEFAULT 'pending',
    action_user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

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
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS conversation_participants(
    conversation_id UUID NOT NULL,
    user_id UUID NOT NULL,
    joined_at TIMESTAMP,
    last_read_message_id BIGINT,

    FOREIGN KEY(conversation_id) REFERENCES conversations(id)
    ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY(last_read_message_id) REFERENCES messages(id)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    sender_id UUID NOT NULL,
    conversation_id UUID NOT NULL,
    body VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    FOREIGN KEY(sender_id) REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY(conversation_id) REFERENCES conversations(id)
    ON DELETE CASCADE
);