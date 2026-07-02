CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
    title VARCHAR,
    body VARCHAR,
    author_id UUID NOT NULL REFERENCES users(id),
    likes INT,
    comments INT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP
);