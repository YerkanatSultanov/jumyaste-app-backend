CREATE TABLE resume
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    full_name   VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    skills      TEXT[]       NULL,
    experience  TEXT         NULL,
    parsed_data JSONB        NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
