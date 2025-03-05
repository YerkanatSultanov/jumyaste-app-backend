CREATE TABLE IF NOT EXISTS users
(
    id              SERIAL PRIMARY KEY,
    email           VARCHAR(200) NOT NULL UNIQUE,
    password        VARCHAR(200) NOT NULL,
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    profile_picture VARCHAR      NULL,
    role_id         INTEGER NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_users_roles FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
);


CREATE TABLE password_resets
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    reset_code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP  NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS roles
(
    id        SERIAL PRIMARY KEY,
    role_name VARCHAR(100) NOT NULL UNIQUE
);
