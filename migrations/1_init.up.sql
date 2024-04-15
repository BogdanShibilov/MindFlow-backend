CREATE TABLE IF NOT EXISTS users
(
    uuid uuid DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL,
    roles VARCHAR[] DEFAULT '{}',
    PRIMARY KEY (uuid)
);
CREATE INDEX IF NOT EXISTS idx_email on users (email);

CREATE TABLE IF NOT EXISTS user_details
(
    user_uuid uuid PRIMARY KEY,
    name VARCHAR(255),
    phone_number VARCHAR(20),
    professional_field VARCHAR(255),
    experience_description TEXT,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS expert_info
(
    user_uuid UUID PRIMARY KEY,
    charge_per_hour INTEGER NOT NULL,
    expertise_at_description TEXT NOT NULL,
    submitted_at TIMESTAMP DEFAULT now(),
    is_approved BOOLEAN DEFAULT false,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS enrollments
(
    uuid uuid DEFAULT gen_random_uuid(),
    mentor_uuid uuid NOT NULL,
    mentee_uuid uuid NOT NULL,
    is_approved BOOLEAN DEFAULT false,
    mentee_questions TEXT,
    PRIMARY KEY (uuid),
    FOREIGN KEY (mentor_uuid) REFERENCES users(uuid) ON DELETE CASCADE,
    FOREIGN KEY (mentee_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS meetings
(
    uuid uuid DEFAULT gen_random_uuid(),
    enrollment_uuid uuid NOT NULL,
    link TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    PRIMARY KEY (uuid),
    FOREIGN KEY (enrollment_uuid) REFERENCES enrollments(uuid) ON DELETE CASCADE
);