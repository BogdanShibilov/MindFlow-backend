CREATE TABLE IF NOT EXISTS users
(
    uuid uuid DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL,
    roles VARCHAR[] NOT NULL,
    PRIMARY KEY (uuid)
);
CREATE INDEX IF NOT EXISTS idx_email on users (email);

CREATE TABLE IF NOT EXISTS user_details
(
    user_uuid uuid PRIMARY KEY,
    name VARCHAR(255),
    surname VARCHAR(255),
    phone_number VARCHAR(20),
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS expert_info
(
    user_uuid UUID PRIMARY KEY,
    position VARCHAR(255),
    charge_per_hour INTEGER,
    experience_description TEXT,
    expertise_at_description TEXT,
    submitted_at TIMESTAMP,
    approved_at TIMESTAMP,
    is_approved BOOLEAN,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);