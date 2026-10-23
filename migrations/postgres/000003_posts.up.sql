CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    title VARCHAR,
    body VARCHAR,
    author_id UUID,
    likes BIGINT,
    comments BIGINT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ,

    FOREIGN KEY(author_id) REFERENCES users(id)
    ON DELETE CASCADE
);