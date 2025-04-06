CREATE TABLE resume
(
    id               SERIAL PRIMARY KEY,
    user_id          INTEGER      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    full_name        VARCHAR(255) NOT NULL,
    desired_position VARCHAR(255) NOT NULL,
    skills           TEXT[]       NULL,
    city             VARCHAR(255) NULL,
    about            TEXT         NULL,
    parsed_data      JSONB        NULL,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);