CREATE TABLE IF NOT EXISTS users
(
    id              SERIAL PRIMARY KEY,
    email           VARCHAR(200) NOT NULL UNIQUE,
    password        VARCHAR(200) NOT NULL,
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    profile_picture VARCHAR      NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE password_resets
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    reset_code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP  NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE email_verifications
(
    id                SERIAL PRIMARY KEY,
    email             VARCHAR(200) NOT NULL UNIQUE,
    verification_code VARCHAR(6)   NOT NULL,
    expires_at        TIMESTAMP    NOT NULL,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

