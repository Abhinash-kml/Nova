CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT uuidv7(),
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
    disabled_at TIMESTAMP,

    PRIMARY KEY(id, username)
);

-- INSERT a system user
INSERT INTO users(id, username, displayname, email) VALUES ('01a1498d-fb97-721b-90f4-b134122fcd23', 'system', 'system', 'system@system.com')

CREATE TABLE IF NOT EXISTS user_name_audit(
    id UUID,
    display_name VARCHAR NOT NULL,
    changed_on TIMESTAMP NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY(id) REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_user_name_audit ON user_name_audit(id);

CREATE TABLE IF NOT EXISTS users_banned(
    id UUID,
    reason VARCHAR NOT NULL,
    ban_type VARCHAR NOT NULL,
    banned_at TIMESTAMPTZ NOT NULL,
    banned_till TIMESTAMPTZ NOT NULL,
    banned_by UUID NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY(id) REFERENCES users(id)
    ON DELETE CASCADE
);

-- CREATE TRIGGER FUNCTION FOR AUDITING DISPLAY NAME CHANGE
CREATE OR REPLACE FUNCTION log_displayname()
RETURNS TRIGGER
LANGUAGE PLPGSQL
AS $$
    BEGIN
        INSERT INTO user_name_audit(id, display_name, changed_on)
        VALUES (NEW.id, NEW.displayname, CURRENT_TIMESTAMP);

        RETURN NEW;
    END;
$$;

-- CREATE TRIGGER FOR AUDITING DISPLAY NAME CHANGE
CREATE TRIGGER trigger_log_displayname
AFTER UPDATE OF displayname
OR INSERT
ON users
FOR EACH ROW
EXECUTE FUNCTION log_displayname();
    