CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    post_id UUID,
    author_id UUID,
    body VARCHAR,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP,

    FOREIGN KEY(post_id) REFERENCES posts(id)
    ON DELETE CASCADE,

    FOREIGN KEY(author_id) REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_comments_postid ON comments(post_id);
CREATE INDEX idx_comments_userid ON comments(author_id);