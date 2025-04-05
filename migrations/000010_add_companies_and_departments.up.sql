CREATE TABLE companies
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255)   NOT NULL,
    owner_id   INTEGER UNIQUE REFERENCES users (id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE departments
(
    id         SERIAL PRIMARY KEY,
    company_id INTEGER      NOT NULL,
    name       VARCHAR(255) NOT NULL,
    hr_count   INTEGER DEFAULT 0,
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

CREATE TABLE invitations
(
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(200) NOT NULL UNIQUE,
    company_id INTEGER      NOT NULL REFERENCES companies (id) ON DELETE CASCADE,
    dep_id     INTEGER      NOT NULL REFERENCES departments (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
