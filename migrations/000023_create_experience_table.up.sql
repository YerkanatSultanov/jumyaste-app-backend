CREATE TABLE work_experience
(
    id              SERIAL PRIMARY KEY,
    resume_id       INTEGER REFERENCES resume (id) ON DELETE CASCADE,
    company_name    TEXT,
    position        TEXT,
    start_date      TEXT,
    end_date        TEXT,
    location        TEXT,
    employment_type TEXT,
    description     TEXT
);
