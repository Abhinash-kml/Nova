CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY UNIQUE DEFAULT uuidv7(),
    username VARCHAR UNIQUE NOT NULL,
    displayname VARCHAR UNIQUE NOT NULL,
    email VARCHAR UNIQUE,
    country VARCHAR,
    state VARCHAR,
    avatar_url VARCHAR,
    lang_tag VARCHAR,
    timezone VARCHAR,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP,
    verified_at TIMESTAMP,
    disabled_at TIMESTAMP
);