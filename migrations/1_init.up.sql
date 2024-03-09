CREATE TABLE IF NOT EXISTS users
(
    uuid uuid DEFAULT gen_random_uuid(),
    email VARCHAR NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL,
    PRIMARY KEY (uuid)
);
CREATE INDEX IF NOT EXISTS idx_email on users (email);