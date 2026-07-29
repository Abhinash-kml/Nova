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