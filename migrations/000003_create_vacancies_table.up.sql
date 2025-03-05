CREATE TABLE vacancies
(
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    employment_type VARCHAR(50)  NOT NULL,
    work_format     VARCHAR(50)  NOT NULL,
    experience      VARCHAR(50)  NOT NULL,
    salary_min      INTEGER      NULL,
    salary_max      INTEGER      NULL,
    location        VARCHAR(255) NULL,
    category        VARCHAR(100) NULL,
    skills          TEXT[]       NULL,
    description     TEXT         NOT NULL,
    created_by      INTEGER      NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE CASCADE
);
