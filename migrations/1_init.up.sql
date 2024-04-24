CREATE TABLE IF NOT EXISTS users
(
    uuid uuid DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL,
    roles VARCHAR[] DEFAULT '{}',
    PRIMARY KEY (uuid)
);
CREATE INDEX IF NOT EXISTS idx_username on users (username);

CREATE TABLE IF NOT EXISTS staff
(
    user_uuid uuid PRIMARY KEY,
    permissions INTEGER[] NOT NULL,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_profiles
(
    user_uuid uuid PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    professional_field VARCHAR(255),
    experience_description TEXT,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS expert_information
(
    user_uuid uuid PRIMARY KEY,
    price INTEGER NOT NULL,
    help_description TEXT NOT NULL,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS expert_application
(
    user_uuid uuid PRIMARY KEY,
    status INTEGER DEFAULT 0,
    submitted_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS consultation
(
    uuid uuid DEFAULT gen_random_uuid(),
    expert_uuid uuid NOT NULL,
    mentee_uuid uuid NOT NULL,
    PRIMARY KEY (uuid),
    FOREIGN KEY (expert_uuid) REFERENCES users(uuid) ON DELETE CASCADE,
    FOREIGN KEY (mentee_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS consultation_application
(
    consultation_uuid uuid PRIMARY KEY,
    status INTEGER DEFAULT 0,
    mentee_questions TEXT NOT NULL,
    submitted_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (consultation_uuid) REFERENCES consultation(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS consultation_meeting
(
    uuid uuid DEFAULT gen_random_uuid(),
    consultation_uuid uuid NOT NULL,
    start_time TIMESTAMP NOT NULL,
    link TEXT NOT NULL,
    PRIMARY KEY (uuid),
    FOREIGN KEY (consultation_uuid) REFERENCES consultation(uuid) ON DELETE CASCADE
);