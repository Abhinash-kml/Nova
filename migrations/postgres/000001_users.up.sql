CREATE TABLE IF NOT EXISTS users (
    id UUID UNIQUE DEFAULT uuidv7(),
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
INSERT INTO users(id, username, displayname, email, country, state, avatar_url, lang_tag, timezone, created_at, updated_at, verified_at, disabled_at) 
VALUES 
('01a1498d-fb97-721b-90f4-b134122fcd23', 'system0', 'system', 'system@system.com', '-', '-', '-', '-', '-', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL);

-- INSERT some seed users
INSERT INTO users(id, username, displayname, email, country, state, avatar_url, lang_tag, timezone, created_at, updated_at, verified_at, disabled_at) 
VALUES 
('2c71b0d3-0171-4da2-99f4-3aebb8b93db4', 'systemname1', 'TestUser1', 'system1@system.com', '-', '-', '-', '-', '-', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
('5cc59547-508d-4e74-bb36-675eac84e10b', 'systemname2', 'TestUser2', 'system2@system.com', '-', '-', '-', '-', '-', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
('dcfd8b0c-fb72-48ed-9cde-ca4285d3f598', 'systemname3', 'TestUser3', 'system3@system.com', '-', '-', '-', '-', '-', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
('29d6e32a-8b82-4de7-9f2d-8a3a8900d2ec', 'systemname4', 'TestUser4', 'system4@system.com', '-', '-', '-', '-', '-', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
('e21de38b-d2a1-42ae-8960-72d725c41cfd', 'systemname5', 'TestUser5', 'system5@system.com', '-', '-', '-', '-', '-', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
('40ef63e4-7e79-4e4e-bd7f-49a9ae1ce84a', 'systemname6', 'TestUser6', 'system6@system.com', '-', '-', '-', '-', '-', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
('dc1fd6a4-cdd9-4a30-ac8c-ebef1297f593', 'systemname7', 'TestUser7', 'system7@system.com', '-', '-', '-', '-', '-', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL);

-- Name change audit log
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
    