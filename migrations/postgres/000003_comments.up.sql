CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
    post_id UUID NOT NULL REFERENCES posts(id),
    author_id UUID NOT NULL REFERENCES users(id),
    body VARCHAR,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP
);